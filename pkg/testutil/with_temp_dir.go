package testutil

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// WithTempDir wraps a test function with a temporary directory to be used during the test
func WithTempDir(t *testing.T, test func(*testing.T, string)) func(*testing.T) {
	dir, err := ioutil.TempDir("", "goproxytest")
	if err != nil {
		assert.FailNow(t, "Error while creating temp dir", err)
	}

	return func(t *testing.T) {
		defer os.RemoveAll(dir)
		test(t, dir)
	}
}
