package internal

import (
	"bytes"
	"github.com/grafov/m3u8"
	"github.com/pkg/errors"
	"strings"
	"time"
)

type Hls struct {
	Mp           *m3u8.MediaPlaylist
	playlistUrl  string
	downloadTime int64
	segmentUrls  map[int]string
}

func (h *Hls) fetchPlaylist(playlistUrl string) error {
	var err error
	h.Mp, err = m3u8.NewMediaPlaylist(100000, 100000)
	if err != nil {
		return errors.Wrapf(err, "NewMediaPlaylist")
	}

	body, err := download(playlistUrl)
	if err != nil {
		return errors.Wrapf(err, "download")
	}

	if err = h.Mp.Decode(*bytes.NewBuffer(body), false); err != nil {
		return errors.Wrapf(err, "DecodeFrom")
	}
	h.downloadTime = time.Now().Unix()
	h.playlistUrl = playlistUrl

	return nil
}

func (h *Hls) savePlaylist(storage Storage) error {
	return storage.SavePlaylist(h)
}

func (h *Hls) parseSegments() error {
	startAt := 0
	endAt := 2000

	// https://golang.hotexamples.com/examples/github.com.grafov.m3u8/-/NewMediaPlaylist/golang-newmediaplaylist-function-examples.html
	playlistHasMoreItems := func(z int) bool { return (h.Mp.Segments[z] != nil && (z < endAt)) }

	segmentUrls := make(map[int]string, 0)

	segmentCounter := 0
	for i := startAt; playlistHasMoreItems(i); i++ {
		segmentUrl := ""
		relativePath := h.Mp.Segments[i].URI

		if comps := strings.Split(h.playlistUrl, "/"); len(comps) > 0 {
			suffix := comps[len(comps)-1]
			baseUrl := strings.TrimSuffix(h.playlistUrl, suffix)
			segmentUrl = baseUrl + relativePath
		} else {
			segmentUrl = relativePath
		}

		segmentUrls[int(h.Mp.SeqNo)+segmentCounter] = segmentUrl
		segmentCounter++
	}

	h.segmentUrls = segmentUrls
	return nil
}

func (h *Hls) fetchAndSaveSegments(storage Storage) error {
	for _, v := range h.segmentUrls {
		b, err := download(v)
		if err != nil {
			return errors.Wrapf(err, "Download")
		}

		err = storage.SaveSegment(h, v, b)
	}
	return nil
}

func (h *Hls) blockTillExpires() {
	secondsToSleep := int(h.Mp.TargetDuration) - 1
	if secondsToSleep == 0 {
		secondsToSleep = 1
	}
	time.Sleep(time.Duration(int64(time.Second) * int64(secondsToSleep)))
}

func (h *Hls) fetchAndSaveAll(storage Storage) error {
	err := h.savePlaylist(storage)
	if err != nil {
		return errors.Wrapf(err, "savePlaylist")
	}

	err = h.parseSegments()
	if err != nil {
		return errors.Wrapf(err, "parseSegments")
	}

	err = h.fetchAndSaveSegments(storage)
	if err != nil {
		return errors.Wrapf(err, "fetchAndSaveSegments")
	}
	return nil
}
