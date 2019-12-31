package upstream

import (
	"math/rand"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEndpointsSorting(t *testing.T) {
	expectedPaths := []string{
		"/another/(.+)/forth",
		"/first/second/on",
		"/first/(.+)/on",
		"/another/second",
		"/another/(.+)",
		"/single",
		"/(.+)",
		"/",
	}

	arr := make(Endpoints, len(expectedPaths))
	for i, p := range expectedPaths {
		if i%2 == 0 {
			arr[i] = &StaticEndpoint{
				BaseEndpoint: BaseEndpoint{From: p},
			}
		} else {
			arr[i] = &StaticEndpoint{
				BaseEndpoint: BaseEndpoint{Regexp: p},
			}
		}
	}

	rand.Seed(10)
	rand.Shuffle(len(arr), func(i, j int) {
		arr[i], arr[j] = arr[j], arr[i]
	})

	sort.Sort(arr)

	for i := 0; i < len(expectedPaths); i++ {
		assert.Equal(t, expectedPaths[i], arr[i].Path())
	}
}
