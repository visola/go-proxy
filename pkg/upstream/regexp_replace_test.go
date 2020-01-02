package upstream

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegexpReplaceNilRegexp(t *testing.T) {
	var re *regexp.Regexp
	toMatch := "/some/path"
	toReplaceIn := "http://notreal.com/another"

	replaced := replaceRegexp(toMatch, toReplaceIn, re)
	assert.Equal(t, toReplaceIn, replaced)
}

func TestRegexpReplaceNoMatch(t *testing.T) {
	re := regexp.MustCompile("nothing")
	toMatch := "/some/path"
	toReplaceIn := "http://notreal.com/another"

	replaced := replaceRegexp(toMatch, toReplaceIn, re)
	assert.Equal(t, toReplaceIn, replaced)
}

func TestRegexpReplaceReplacesMatches(t *testing.T) {
	re := regexp.MustCompile("/(.+)/(.+)")
	toMatch := "/some/path"
	toReplaceIn := "http://notreal.com/another/$2/$1"

	replaced := replaceRegexp(toMatch, toReplaceIn, re)
	assert.Equal(t, "http://notreal.com/another/path/some", replaced)
}
