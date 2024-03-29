package upstream

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/visola/go-proxy/pkg/configuration"
)

func TestUpstreamsPerFile(t *testing.T) {
	// Ensure no upstreams
	upstreams = make(map[string]Upstream, 0)

	firstOrigin := configuration.Origin{
		File:     "/some/place/first",
		LoadedAt: 10,
	}

	AddUpstreams([]Upstream{
		Upstream{
			Name:   "Upstream 1",
			Origin: firstOrigin,
		},
		Upstream{
			Name:   "Upstream 2",
			Origin: firstOrigin,
		},
	})

	secondOrigin := configuration.Origin{
		File:     "/another/file/second",
		LoadedAt: 20,
	}

	AddUpstreams([]Upstream{
		Upstream{
			Name:   "Upstream 3",
			Origin: secondOrigin,
		},
		Upstream{
			Name:   "Upstream 4",
			Origin: secondOrigin,
		},
		Upstream{
			Name:   "Upstream 5",
			Origin: secondOrigin,
		},
	})

	groupedByFile := UpstreamsPerFile()

	assert.Equal(t, 2, len(groupedByFile), "Should group by two origins")

	assert.Equal(t, 2, len(groupedByFile[firstOrigin.File]))
	assert.Equal(t, 3, len(groupedByFile[secondOrigin.File]))
}
