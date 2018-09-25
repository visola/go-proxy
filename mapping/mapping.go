package mapping

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/visola/variables/variables"
)

// Mapping represents a mapping that can be active or not
type Mapping struct {
	Active    bool                 `json:"active"`
	Before    string               `json:"before"`
	From      string               `json:"from"`
	Inject    Injection            `json:"injection"`
	MappingID string               `json:"mappingID"`
	Origin    string               `json:"origin"`
	Proxy     bool                 `json:"proxy"`
	Regexp    string               `json:"regexp"`
	Tags      []string             `json:"tags"`
	To        string               `json:"to"`
	Variables []variables.Variable `json:"variables"`
}

// Injection represents parameters that can be injected into proxied requests
type Injection struct {
	Headers map[string]string
}

// MatchResult stores the result for a mapping that matched a request
type MatchResult struct {
	Mapping Mapping
	NewPath string
	Parts   []string
}

// GetVariables returns all the variables set in a mapping
func (mapping *Mapping) GetVariables() []variables.Variable {
	result := make([]variables.Variable, 0)
	result = append(result, variables.FindVariables(mapping.From)...)
	result = append(result, variables.FindVariables(mapping.To)...)

	for _, v := range mapping.Inject.Headers {
		result = append(result, variables.FindVariables(v)...)
	}

	return result
}

// Match tests if the mapping matches the specific request. If it does, it will
// return a match result, otherwise it will return nil.
func (mapping *Mapping) Match(req *http.Request) *MatchResult {
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

// Validate makes sure that the a mapping is correctly setup
func (mapping *Mapping) Validate() error {
	if mapping.From == "" && mapping.Regexp == "" {
		return errors.New("Either `from` or `regexp` need to be present")
	}

	if mapping.To == "" {
		return errors.New("Missing value for `to`")
	}

	if mapping.Regexp != "" {
		_, err := regexp.Compile(mapping.Regexp)
		if err != nil {
			return fmt.Errorf("Error compiling regexp: '%s'\n%s", mapping.Regexp, err)
		}
	}

	if len(mapping.Inject.Headers) > 0 && !mapping.Proxy {
		return errors.New("Inject is only available for proxy mappings")
	}

	return nil
}
