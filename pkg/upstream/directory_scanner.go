package upstream

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

func isGoProxyUpstreamsFile(pathToFile string) bool {
	return strings.HasSuffix(pathToFile, "go-proxy.yaml") || strings.HasSuffix(pathToFile, ".go-proxy.yml")
}

// ScanFilesInDirectories scans a set of directories to search for upstreams
func ScanFilesInDirectories(dirsToScan []string) []Upstream {
	if len(dirsToScan) == 0 {
		return make([]Upstream, 0)
	}

	result := make([]Upstream, 0)
	totalFiles := 0
	filesFound := 0
	for len(dirsToScan) > 0 {
		dirToRead := dirsToScan[0]
		dirsToScan = dirsToScan[1:]

		files, readDirErr := ioutil.ReadDir(dirToRead)
		if readDirErr != nil {
			log.Printf("Error while reading directory: %s\n%s\n", dirToRead, readDirErr.Error())
			continue
		}

		for _, file := range files {
			totalFiles++
			pathToFile := path.Join(dirToRead, file.Name())

			if file.IsDir() {
				if !shouldIgnoreDirectory(file) {
					dirsToScan = append(dirsToScan, pathToFile)
				}
				continue
			}

			if isGoProxyUpstreamsFile(pathToFile) {
				filesFound++
				loadedUpstreams, err := loadFromFile(pathToFile)
				if err != nil {
					log.Printf("Error while loading upstream from file: %s\n%s\n", pathToFile, err.Error())
					continue
				}

				result = append(result, loadedUpstreams...)
			}
		}
	}

	return result
}

func shouldIgnoreDirectory(dir os.FileInfo) bool {
	name := dir.Name()
	return strings.HasPrefix(name, ".") ||
		name == "build" ||
		name == "target" ||
		name == "dist" ||
		name == "node_modules"
}
