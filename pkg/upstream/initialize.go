package upstream

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

var filesToLoad = make(chan string)

// Initialize loads all upstreams and initalize a routine to refresh the read files
func Initialize(baseDir string) {
	log.Printf("Initializing upstreams...")
	for i := 0; i < 5; i++ {
		go processFilesToLoad()
	}

	go findFilesInConfiguratonDirectory(baseDir)
	go func() {
		for {
			time.Sleep(5 * time.Second)
			refreshStaleUpstreams()
		}
	}()
}

func findFilesInConfiguratonDirectory(baseDir string) {
	log.Printf("Reading configuration directory: %s", baseDir)
	files, err := ioutil.ReadDir(baseDir)
	if err != nil {
		log.Fatalf("Error while reading directory: %s, %v", baseDir, err)
	}

	for _, file := range files {
		if !file.IsDir() {
			ext := filepath.Ext(file.Name())
			if ext == ".yml" || ext == ".yaml" {
				filesToLoad <- filepath.Join(baseDir, file.Name())
			}
		}
	}
}

func processFilesToLoad() {
	fileToLoad := <-filesToLoad
	upstreams, err := loadFromFile(fileToLoad)
	log.Printf("Found %d upstreams in file: %s", len(upstreams), fileToLoad)
	if err != nil {
		log.Fatalf("Error while loading upstream from %s, %v", fileToLoad, err)
	}

	AddUpstreams(upstreams)
}

func refreshStaleUpstreams() {
	newUpstreams := make(map[string]Upstream)
	for filePath, upstreamsInFile := range UpstreamsPerFile() {
		fileStat, fileErr := os.Stat(filePath)
		if fileErr != nil {
			log.Fatalf("Error while checking file status: %s, %v", filePath, fileErr)
		}

		if upstreamsInFile[0].Origin.LoadedAt < fileStat.ModTime().Unix() {
			log.Printf("Found stale upstreams from file: %s", filePath)
			loadedUpstreams, loadErr := loadFromFile(filePath)
			if loadErr != nil {
				log.Fatalf("Error while refresing upstreams from %s, %v", filePath, loadErr)
			}
			upstreamsInFile = loadedUpstreams
		}

		for _, toAdd := range upstreamsInFile {
			newUpstreams[toAdd.Name] = toAdd
		}
	}

	upstreamsMutex.Lock()
	defer upstreamsMutex.Unlock()
	upstreams = newUpstreams
}
