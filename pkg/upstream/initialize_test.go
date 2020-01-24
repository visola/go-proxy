package upstream

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/visola/go-proxy/pkg/configuration"
	"github.com/visola/go-proxy/pkg/testutil"
)

func TestRefreshStaleUpstreams(t *testing.T) {
	// Ensure no upstreams
	upstreams = make(map[string]Upstream, 0)

	fileName := "test.yml"
	fileContent := `
upstreams:
  backend:
`

	tempFile := testutil.CreateTempFile(t, fileName, fileContent)
	defer os.Remove(tempFile.Name())

	origin := configuration.Origin{
		File:     tempFile.Name(),
		LoadedAt: 10,
	}

	AddUpstreams([]Upstream{
		Upstream{
			Name:   "test",
			Origin: origin,
		},
	})

	refreshStaleUpstreams()
	assert.Equal(t, 2, len(upstreams))

	// TODO - rewrite this to avoid accessing globals
	first := upstreams["test"]
	assert.Equal(t, "test", first.Name)
	assert.Equal(t, tempFile.Name(), first.Origin.File)

	second := upstreams["backend"]
	assert.Equal(t, "backend", second.Name)
	assert.Equal(t, tempFile.Name(), second.Origin.File)

	refreshStaleUpstreams()
	assert.Equal(t, 2, len(upstreams))
}
