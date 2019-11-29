package upstream

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatches(t *testing.T) {
	t.Run("Doesn't match wrong path", testDoesntMatch)
	t.Run("Matches by from", testMatchesByFrom)
	t.Run("Matches by regexp", testMatchesByRegexp)
}

func testDoesntMatch(t *testing.T) {
	m := BaseEndpoint{}
	req := httptest.NewRequest(http.MethodGet, "http://nowhere.com/another", strings.NewReader("Some Body"))
	assert.Equal(t, false, m.Matches(*req))
}

func testMatchesByFrom(t *testing.T) {
	m := BaseEndpoint{
		From: "/test",
	}

	req := httptest.NewRequest(http.MethodGet, "http://nowhere.com/test", strings.NewReader("Some Body"))
	assert.Equal(t, true, m.Matches(*req))
}

func testMatchesByRegexp(t *testing.T) {
	m := BaseEndpoint{
		Regexp: "/.*/more/(.*)",
	}

	req := httptest.NewRequest(http.MethodGet, "http://nowhere.com/test/more/complex", strings.NewReader("Some Body"))
	assert.Equal(t, true, m.Matches(*req))
}
