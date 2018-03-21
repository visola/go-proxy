package proxy

import (
	"net/http"

	"github.com/visola/go-proxy/config"
)

func matchConfiguration(req *http.Request, configurations []config.DynamicMapping) *config.MatchResult {
	for _, config := range configurations {
		if !config.Active {
			continue
		}

		matchResult := config.Match(req)
		if matchResult != nil {
			return matchResult
		}
	}

	return nil
}
