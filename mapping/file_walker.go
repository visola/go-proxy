package mapping

import (
	"fmt"
	"io/ioutil"
	"path"
	"sync"
)

func walkDirectories(directories []string, callback func(string, bool)) {
	var wg sync.WaitGroup
	var callbackLock sync.Mutex

	wg.Add(len(directories))
	directoriesToVisit := make(chan string, 1000000)
	for _, dir := range directories {
		directoriesToVisit <- dir
	}

	for i := 0; i < 5; i++ {
		go readDirectories(directoriesToVisit, &wg, &callbackLock, callback)
	}

	wg.Wait()
}

func readDirectories(directoriesToVisit chan string, wg *sync.WaitGroup, callbackLock *sync.Mutex, callback func(string, bool)) {
	for dirToRead := range directoriesToVisit {
		readDirectory(dirToRead, directoriesToVisit, wg, callbackLock, callback)
	}
}

func readDirectory(dirToRead string, directoriesToVisit chan string, wg *sync.WaitGroup, callbackLock *sync.Mutex, callback func(string, bool)) {
	defer wg.Done()

	files, readDirErr := ioutil.ReadDir(dirToRead)
	if readDirErr != nil {
		fmt.Printf("Error while reading directory: %s\n%s\n", dirToRead, readDirErr.Error())
		return
	}

	for _, file := range files {
		pathToFile := path.Join(dirToRead, file.Name())

		if file.IsDir() {
			wg.Add(1)
			directoriesToVisit <- pathToFile
		}

		callbackLock.Lock()
		callback(pathToFile, file.IsDir())
		callbackLock.Unlock()
	}
}
