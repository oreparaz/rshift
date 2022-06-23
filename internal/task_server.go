package internal

import (
	"compress/gzip"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/goji/httpauth"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

var storage Storage

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func TimeshiftHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	ts, err := strconv.Atoi(vars["timeShift"])
	if err != nil {
		log.Printf("Atoi returned error converting %s\n", vars["timeShift"])
		return
	}
	startTime := time.Now()

	fileBytes, err := storage.ReadPlaylistNear(int(time.Now().Unix()) - ts)
	endTime := time.Now()

	log.Println("ReadPlaylistNear took: ", endTime.Sub(startTime))

	w.Header().Set("Content-Type", "application/x-mpegURL")
	w.WriteHeader(http.StatusOK)
	w.Write(fileBytes)
}

func gzipHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			h.ServeHTTP(w, r)
			return
		}
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		h.ServeHTTP(gzipResponseWriter{Writer: gz, ResponseWriter: w}, r)
	})
}

func MainServer() {
	storage = &FileStorage{}

	// bits and pieces from https://github.com/jeongmin/m3u8-reader/blob/master/main.go
	corsObj := handlers.AllowedOrigins([]string{"*"})

	// this is a hack.
	u, err := url.Parse(M3u8DownloadUrl)
	if err != nil { panic(err) }
	// when M3u3DownloadUrl is http://a/b/c/d.m3u8, urlDirectory is /b/c/
	urlDirectory := path.Dir(u.Path)

	r := mux.NewRouter()
	r.HandleFunc("/timeshift/{timeShift}.m3u8", TimeshiftHandler)
	r.PathPrefix("/timeshift/").Handler(
		http.StripPrefix("/timeshift/",
			http.FileServer(http.Dir(path.Join(OutPath, path.Join("ts/", urlDirectory))))))

	handler := gzipHandler(handlers.CORS(corsObj)(r))
	if os.Getenv("RSHIFT_USERNAME") != "" {
		auth := httpauth.SimpleBasicAuth(os.Getenv("RSHIFT_USERNAME"), os.Getenv("RSHIFT_PASSWORD"))
		handler = auth(handler)
	} else {
		log.Println("warning: proceeding without any authentication")
	}
	http.ListenAndServe(":8080", handler)
}
