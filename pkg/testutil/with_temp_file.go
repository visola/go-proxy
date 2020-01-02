package testutil

import (
	"io/ioutil"
	"os"
	"testing"
)

// WithTempFile wraps a test function with a temporary file
func WithTempFile(t *testing.T, pattern string, content string, callback func(*os.File)) {
	tempFile, err := ioutil.TempFile("", pattern)
	if err != nil {
		t.Fatalf("Error while creating temp file: %s", err)
	}

	if _, err := tempFile.Write([]byte(content)); err != nil {
		t.Fatalf("Error while writing content to temp file: %s", err)
	}

	defer os.Remove(tempFile.Name())
	callback(tempFile)
}
