package mapping

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"
	"sync"
)

var callbackLock sync.Mutex

func walkDir(pathToDir string, wg *sync.WaitGroup, callback func(string, bool)) {
	defer wg.Done()

	if strings.HasSuffix(pathToDir, "/.git") {
		return
	}

	files, readDirErr := ioutil.ReadDir(pathToDir)
	if readDirErr != nil {
		fmt.Printf("Error while reading directory: %s\n%s\n", pathToDir, readDirErr.Error())
		return
	}

	for _, fileInfo := range files {
		pathToFile := path.Join(pathToDir, fileInfo.Name())
		if fileInfo.IsDir() {
			wg.Add(1)
			go walkDir(pathToFile, wg, callback)
		}
		callbackLock.Lock()
		callback(pathToFile, fileInfo.IsDir())
		callbackLock.Unlock()
	}
}
