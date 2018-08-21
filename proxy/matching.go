package proxy

import (
	"net/http"

	"github.com/Everbridge/go-proxy/mapping"
)

func matchMapping(req *http.Request, mappings []mapping.Mapping) *mapping.MatchResult {
	for _, mapping := range mappings {
		if !mapping.Active {
			continue
		}

		matchResult := mapping.Match(req)
		if matchResult != nil {
			return matchResult
		}
	}

	return nil
}
