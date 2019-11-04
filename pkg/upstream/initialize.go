package upstream

import (
	"io/ioutil"
	"log"
	"path/filepath"
)

var filesToLoad = make(chan string)

func Initialize(baseDir string) {
	log.Printf("Initializing upstreams...")
	for i := 0; i < 5; i++ {
		go processFilesToLoad()
	}
	go findFilesInConfiguratonDirectory(baseDir)
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
