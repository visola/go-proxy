package proxy

import (
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/visola/go-proxy/config"
)

func serveStaticFile(req *http.Request, w http.ResponseWriter, mapping config.Mapping) {
	oldPath := req.URL.Path
	newPath := path.Join(mapping.To, oldPath[len(mapping.From):])

	file, err := os.Open(newPath)

	if err == os.ErrNotExist {
		notFound(req, w, newPath)
		return
	}

	if err != nil {
		internalError(req, w, err)
		return
	}

	defer file.Close()

	log.Printf("Serving '%s' for '%s', from '%s'", newPath, req.URL.Path, mapping.Origin)

	buffer := make([]byte, 512)
	loopCount := 0
	for {
		bytesRead, readError := file.Read(buffer)

		if readError != nil && readError != io.EOF {
			internalError(req, w, readError)
		}

		if bytesRead == 0 {
			break
		}

		loopCount++

		contentType := ""
		if loopCount == 1 {
			contentType = mime.TypeByExtension(filepath.Ext(file.Name()))
			if contentType == "" {
				contentType = http.DetectContentType(buffer)
			}
		}

		w.Header().Set("Content-Type", contentType)
		w.Write(buffer[:bytesRead])
	}
}
