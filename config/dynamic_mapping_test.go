package config

import (
	"log"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatch(t *testing.T) {
	t.Run("Returns nil if does't match", testNoMatch)
	t.Run("Matches Path", testMatchPath)
	t.Run("Matches Regexp", testMatchRegexp)
}

func testMatchPath(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "https://localhost:12312/some/path", nil)

	mapping := DynamicMapping{
		From: "/some",
		To:   "/another",
	}

	matchResult := mapping.Match(req)
	assert.NotNil(t, matchResult, "Should match path")
	if matchResult == nil {
		log.Fatal("Should've matched path")
	}
	assert.Equal(t, "/some/path", matchResult.NewPath, "Should return matched path in new path")
	assert.Equal(t, mapping, matchResult.Mapping, "Should set self as the mapping in result")
	assert.Equal(t, []string{"/some/path"}, matchResult.Parts, "Should set path as part")
}

func testMatchRegexp(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "https://localhost:12312/some/here/path", nil)

	mapping := DynamicMapping{
		Regexp: "/some/([A-Za-z]+)/([A-Za-z]+)",
		To:     "/another/$2/$1",
	}

	matchResult := mapping.Match(req)
	assert.NotNil(t, matchResult, "Should match path")
	if matchResult == nil {
		log.Fatal("Should've matched path")
	}
	assert.Equal(t, "/another/path/here", matchResult.NewPath, "Should replace matched part in new path")
	assert.Equal(t, mapping, matchResult.Mapping, "Should set self as the mapping in result")
	assert.Equal(t, []string{"/some/here/path", "here", "path"}, matchResult.Parts, "Should set path as part")
}

func testNoMatch(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "https://localhost:12312/will/not/match", nil)

	mapping := DynamicMapping{
		From:   "/some",
		Regexp: "/some/([A-Za-z]+)/([A-Za-z]+)",
		To:     "/another/$2/$1",
	}

	matchResult := mapping.Match(req)
	assert.Nil(t, matchResult, "Should return nil if no match")
}
