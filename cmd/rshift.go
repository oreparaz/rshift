package main

import (
	"flag"
	"log"
	"rshift/internal"
	"sync"
	"time"
)

func workerDownload(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		err := internal.LoopPlayList(internal.M3u8DownloadUrl)
		log.Printf("download error: %v, restarting...\n", err)
	}
}

func workerServer(wg *sync.WaitGroup) {
	defer wg.Done()
	internal.MainServer()
}

func workerCleanup(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		internal.DeleteOldFiles()
		time.Sleep(time.Hour)
	}
}



func main() {
	flag.StringVar(&internal.M3u8DownloadUrl, "download-m3u8-url", "https://example.com/file.m3u8", "URL with M3U8 playlist")
	flag.StringVar(&internal.OutPath, "output-path", "/mnt/disks/sdb/out", "path to store cached files")
	flag.Parse()

	var wg sync.WaitGroup
	wg.Add(1)
	go workerDownload(&wg)
	wg.Add(1)
	go workerServer(&wg)

	go workerCleanup(&wg)
	wg.Add(1)

	wg.Wait()
	log.Println("Main: Completed")
}
