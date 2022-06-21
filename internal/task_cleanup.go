package internal

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

func DeleteOldFiles() {
	dir := filepath.Join(OutPath, "m3u8")
	wipeDir(dir)

	dir2 := filepath.Join(OutPath, "ts")
	wipeDir(dir2)
}

func wipeDir(dir string) {
	err := filepath.Walk(dir,
		func(path string, fileInfo os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if fileInfo.IsDir() {
				return nil
			}
			// not here: check for nasty stuff like symlinks
			if fileInfo.ModTime().AddDate(0, 0, MaxAgeFilesDays).Before(time.Now()) {
				log.Printf("deleting file %s (too old)\n", path)
				err = os.Remove(path)
				if err != nil {
					log.Printf("error: Remove %s %v\n", path, err)
				}
			}
			return nil
		})
	if err != nil {
		log.Println(err)
	}
}
