package mapping

import (
	"log"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatch(t *testing.T) {
	t.Run("Returns nil if does't match", testNoMatch)
	t.Run("Matches Path to proxy", testMatchProxyPath)
	t.Run("Matches Path to proxy and doesn't add slash at the end", testMatchDoesntAddSlashAtTheEnd)
	t.Run("Matches Path to static", testMatchStaticPath)
	t.Run("Matches Regexp", testMatchRegexp)

	t.Run("Test header injection mapping", testHeaderInjection)
	t.Run("Valid mapping returns no error", testValidMapping)
	t.Run("Fail if invalid Regexp", testInvalidRegexp)
	t.Run("Fail if missing information", testMissingInfoMappings)
}

func testMatchProxyPath(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "https://localhost:12312/some/path", nil)

	mapping := Mapping{
		From:  "/some",
		To:    "http://nowhere.com/another",
		Proxy: true,
	}

	matchResult := mapping.Match(req)
	assert.NotNil(t, matchResult, "Should match path")
	if matchResult == nil {
		log.Fatal("Should've matched path")
	}
	assert.Equal(t, "http://nowhere.com/another/path", matchResult.NewPath, "Should return new path")
	assert.Equal(t, mapping, matchResult.Mapping, "Should set self as the mapping in result")
	assert.Equal(t, []string{"/some/path"}, matchResult.Parts, "Should set path as part")
}

func testMatchDoesntAddSlashAtTheEnd(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "https://localhost:12312/some", nil)

	mapping := Mapping{
		From:  "/some",
		To:    "http://nowhere.com/another",
		Proxy: true,
	}

	matchResult := mapping.Match(req)
	assert.NotNil(t, matchResult, "Should match path")
	if matchResult == nil {
		log.Fatal("Should've matched path")
	}
	assert.Equal(t, "http://nowhere.com/another", matchResult.NewPath, "Should return new path")
	assert.Equal(t, mapping, matchResult.Mapping, "Should set self as the mapping in result")
	assert.Equal(t, []string{"/some"}, matchResult.Parts, "Should set path as part")
}

func testMatchStaticPath(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "https://localhost:12312/some/path", nil)

	mapping := Mapping{
		From:  "/some",
		To:    "/home/user/another",
		Proxy: false,
	}

	matchResult := mapping.Match(req)
	assert.NotNil(t, matchResult, "Should match path")
	if matchResult == nil {
		log.Fatal("Should've matched path")
	}
	assert.Equal(t, "/home/user/another/path", matchResult.NewPath, "Should return new path")
	assert.Equal(t, mapping, matchResult.Mapping, "Should set self as the mapping in result")
	assert.Equal(t, []string{"/some/path"}, matchResult.Parts, "Should set path as part")
}

func testMatchRegexp(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "https://localhost:12312/some/here/path", nil)

	mapping := Mapping{
		Regexp: "/some/([A-Za-z]+)/([A-Za-z]+)",
		To:     "http://nowhere.com/another/$2/$1",
		Proxy:  true,
	}

	matchResult := mapping.Match(req)
	assert.NotNil(t, matchResult, "Should match path")
	if matchResult == nil {
		log.Fatal("Should've matched path")
	}
	assert.Equal(t, "http://nowhere.com/another/path/here", matchResult.NewPath, "Should replace matched part in new path")
	assert.Equal(t, mapping, matchResult.Mapping, "Should set self as the mapping in result")
	assert.Equal(t, []string{"/some/here/path", "here", "path"}, matchResult.Parts, "Should set path as part")
}

func testNoMatch(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "https://localhost:12312/will/not/match", nil)

	mapping := Mapping{
		From:   "/some",
		Regexp: "/some/([A-Za-z]+)/([A-Za-z]+)",
		To:     "/another/$2/$1",
	}

	matchResult := mapping.Match(req)
	assert.Nil(t, matchResult, "Should return nil if no match")
}

func testValidMapping(t *testing.T) {
	mapping := Mapping{
		From: "/some/path",
		To:   "http://some.server.com/another/path",
	}

	validationError := mapping.Validate()
	assert.Nil(t, validationError, "Should not return error when From and To are correctly set")

	mapping = Mapping{
		Regexp: "/some/([A-Za-z]+)/([A-Za-z]+)",
		To:     "http://some.server.com/another/path",
	}

	validationError = mapping.Validate()
	assert.Nil(t, validationError, "Should not return error when valid Regexp and To are correctly set")
}

func testInvalidRegexp(t *testing.T) {
	mapping := Mapping{
		Regexp: "/some/([A-Za-z]+", // Missing closing bracket
		To:     "http://some.server.com/another/path",
	}

	validationError := mapping.Validate()
	assert.NotNil(t, validationError, "Should return error when invalid Regexp")
}

func testMissingInfoMappings(t *testing.T) {
	var validationError error

	missingFromRegexp := Mapping{
		To: "http://some.server.com/another/path",
	}

	validationError = missingFromRegexp.Validate()
	assert.NotNil(t, validationError, "Should return error when missing From and Regexp")

	missingToWithRegexp := Mapping{
		Regexp: "/some/([A-Za-z]+)",
	}

	validationError = missingToWithRegexp.Validate()
	assert.NotNil(t, validationError, "Should return error when missing To and has Regexp")

	missingToWithFrom := Mapping{
		From: "/some/path",
	}

	validationError = missingToWithFrom.Validate()
	assert.NotNil(t, validationError, "Should return error when missing To and has From")
}

func testHeaderInjection(t *testing.T) {
	var validationError error

	headerInjection := Mapping{
		From:  "/some/path",
		To:    "http://some.server.com/another/path",
		Proxy: true,
		Inject: Injection{
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		},
	}

	validationError = headerInjection.Validate()
	assert.Nil(t, validationError, "Should not return error if valid injection")

	headerInjectionNonProxy := Mapping{
		From:  "/some/path",
		To:    "http://some.server.com/another/path",
		Proxy: false,
		Inject: Injection{
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		},
	}

	validationError = headerInjectionNonProxy.Validate()
	assert.NotNil(t, validationError, "Should return error if injecting in non-proxy")
}
