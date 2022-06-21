package internal

import (
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"log"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type Storage interface {
	SavePlaylist(h *Hls) error
	SaveSegment(h *Hls, url string, segment []byte) error
	ReadPlaylistNear(desiredTimestamp int) ([]byte, error)
}

type FileStorage struct {}

var BackendFs afero.Fs

func (d *FileStorage) SavePlaylist(h *Hls) error {
	fileName, err := getPlaylistFilename(h.playlistUrl,
		h.downloadTime)
	if err != nil {
		return errors.Wrapf(err, "getPlaylistFilename")
	}

	_ = BackendFs.MkdirAll(filepath.Dir(fileName), os.ModeDir|os.ModePerm)
	out, err := BackendFs.Create(fileName)
	if err != nil {
		return errors.Wrapf(err, "Create")
	}
	defer out.Close()

	log.Printf("saved %s\n", fileName)

	_, err = h.Mp.Encode().WriteTo(out)
	if err != nil {
		return errors.Wrapf(err, "WriteTo")
	}
	return nil
}

func (d *FileStorage) SaveSegment(h *Hls, url string, segment []byte) error {
	fileName, err := getSegmentFilename(url)
	if _, err := BackendFs.Stat(fileName); err == nil {
		//log.Printf("file %s already exists, exiting\n", fileName)
		return nil
	}
	_ = BackendFs.MkdirAll(filepath.Dir(fileName), os.ModeDir|os.ModePerm)

	log.Printf("saved %s\n", fileName)

	out, err := BackendFs.Create(fileName)
	if err != nil {
		return errors.Wrapf(err, "Create")
	}
	defer out.Close()

	_, err = out.Write(segment)
	if err != nil {
		return errors.Wrapf(err, "Write")
	}
	return nil
}

func (d *FileStorage) ReadPlaylistNear(desiredTimestamp int) ([]byte, error) {
	fileName, err := GetClosestPlaylistFilename(desiredTimestamp)

	if err != nil {
		return nil, errors.Wrapf(err, "GetClosestPlaylistFilename")
	}

	fileBytes, err := afero.ReadFile(BackendFs, fileName)
	if err != nil {
		return nil, errors.Wrapf(err, "ReadFile")
	}

	return fileBytes, nil
}

const SEPARATOR = "---"

func getPlaylistFilename(playListUrl string, unixTime int64) (string, error) {
	u, err := url.Parse(playListUrl)
	if err != nil {
		return "", errors.Wrapf(err, "Parse")
	}
	fileName := filepath.Join(OutPath, "m3u8", path.Base(u.Path)+SEPARATOR+strconv.FormatInt(unixTime, 10))
	return fileName, nil
}

func GetClosestPlaylistFilename(approxTime int) (string, error) {
	fs, err := BackendFs.Open(filepath.Join(OutPath, "m3u8"))
	if err != nil {
		return "", errors.Wrapf(err, "Open")
	}
	fileNames, err := fs.Readdirnames(0)
	if err != nil {
		return "", errors.Wrapf(err, "Readdir")
	}
	sort.Strings(fileNames)

	for _, fn := range fileNames {
		fields := strings.Split(fn, SEPARATOR)
		whenDownloaded, _ := strconv.Atoi(fields[1])
		if whenDownloaded > approxTime {
			difference := whenDownloaded-approxTime
			log.Printf("getClosestPlaylist returns %s (wanted: %d got: %d difference: %d)\n",
				fn,
				approxTime,
				whenDownloaded,
				difference)

			if difference > 60 {
				// really should be an error, but do not want to tear everything down
				log.Printf("difference %v is too large\n", difference)
			}

			return filepath.Join(filepath.Join(OutPath, "m3u8"), fn), nil
		}
	}
	return "", nil
}

func getSegmentFilename(segmentUrl string) (string, error) {
	u, err := url.Parse(segmentUrl)
	if err != nil {
		return "", errors.Wrapf(err, "Parse")
	}
	fileName := filepath.Join(OutPath, "ts", u.Path)
	return fileName, nil
}
