package internal

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

func downloadPlaylistIfRecent(url string, lastSequence int) (int, error) {
	BackendFs = afero.NewOsFs()
	fs := &FileStorage{}

	var h Hls
	err := h.fetchPlaylist(url)
	if err != nil {
		return 0, errors.Wrapf(err, "fetchPlaylist")
	}

	if int(h.Mp.SeqNo) != lastSequence {
		err = h.fetchAndSaveAll(fs)
		if err != nil {
			return 0, errors.Wrapf(err, "fetchAndSaveAll")
		}
		h.blockTillExpires()
		return int(h.Mp.SeqNo), nil
	}

	return lastSequence, nil
}

func LoopPlayList(url string) error {
	lastSequence := 0
	for {
		var err error
		lastSequence, err = downloadPlaylistIfRecent(url, lastSequence)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
}
