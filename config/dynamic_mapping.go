package config

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

// DynamicMapping represents a mapping that can be active or not
type DynamicMapping struct {
	Active    bool   `json:"active"`
	From      string `json:"from"`
	MappingID string `json:"mappingID"`
	Origin    string `json:"origin"`
	Proxy     bool   `json:"proxy"`
	Regexp    string `json:"regexp"`
	To        string `json:"to"`
}

// MatchResult stores the result for a mapping that matched a request
type MatchResult struct {
	Mapping DynamicMapping
	NewPath string
	Parts   []string
}

// Match tests if the mapping matches the specific request. If it does it will
// return a match result, otherwise it will return nil.
func (mapping *DynamicMapping) Match(req *http.Request) *MatchResult {
	if mapping.From != "" && strings.HasPrefix(req.URL.Path, mapping.From) {
		return &MatchResult{
			Mapping: *mapping,
			NewPath: req.URL.Path,
			Parts:   []string{req.URL.Path},
		}
	}

	if mapping.Regexp != "" {
		r, err := regexp.Compile(mapping.Regexp)
		if err != nil {
			// Assume validation errors will be exposed some other way
			return nil
		}

		matched := r.FindStringSubmatch(req.URL.Path)
		if len(matched) > 0 {
			newPath := mapping.To
			for index, part := range matched[1:] {
				newPath = strings.Replace(newPath, fmt.Sprintf("$%d", index+1), part, -1)
			}
			return &MatchResult{
				Mapping: *mapping,
				NewPath: newPath,
				Parts:   matched,
			}
		}
	}

	return nil
}
