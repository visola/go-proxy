package testutil

import (
	"io/ioutil"
	"os"
	"testing"
)

// CreateTempFile creates a temporary file with the specified content
func CreateTempFile(t *testing.T, pattern string, content string) *os.File {
	tempFile, err := ioutil.TempFile("", pattern)
	if err != nil {
		t.Fatalf("Error while creating temp file: %s", err)
	}

	if _, err := tempFile.Write([]byte(content)); err != nil {
		t.Fatalf("Error while writing content to temp file: %s", err)
	}

	return tempFile
}
