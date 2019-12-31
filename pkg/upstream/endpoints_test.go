package upstream

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEndpointsSorting(t *testing.T) {
	e1 := &StaticEndpoint{
		BaseEndpoint: BaseEndpoint{
			From: "/first/second/third",
		},
	}

	e2 := &StaticEndpoint{
		BaseEndpoint: BaseEndpoint{
			Regexp: "/first/(.+)/forth",
		},
	}

	e3 := &StaticEndpoint{
		BaseEndpoint: BaseEndpoint{
			From: "/",
		},
	}

	e4 := &StaticEndpoint{
		BaseEndpoint: BaseEndpoint{
			Regexp: "/another/(.+)",
		},
	}

	arr := Endpoints{e2, e3, e1, e4}
	sort.Sort(arr)

	assert.Equal(t, e1, arr[0])
	assert.Equal(t, e2, arr[1])
	assert.Equal(t, e4, arr[2])
	assert.Equal(t, e3, arr[3])
}
