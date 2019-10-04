package mapping

import (
	"errors"
	"fmt"
	"net/http"
	"path"
	"regexp"
	"strings"

	"github.com/visola/variables/variables"
)

// Mapping represents a mapping that can be active or not
type Mapping struct {
	Active    bool      `json:"active"`
	Before    string    `json:"before"`
	From      string    `json:"from"`
	Inject    Injection `json:"injection"`
	MappingID string    `json:"mappingID"`
	Origin    string    `json:"origin"`
	Proxy     bool      `json:"proxy"`
	Regexp    string    `json:"regexp"`
	Tags      []string  `json:"tags"`
	To        string    `json:"to"`
}

// Injection represents parameters that can be injected into proxied requests
type Injection struct {
	Headers map[string]string `json:"headers"`
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
// return a match result with the new path where to go. Otherwise it will return nil.
func (mapping *Mapping) Match(req *http.Request) *MatchResult {
	if mapping.From != "" {
		return mapping.matchWithFrom(req)
	}

	if mapping.Regexp != "" {
		return mapping.matchWithRegexp(req)
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

// WithReplacedVariables creates a new mapping with all its variables replaced by the
// values passed in the context
func (mapping Mapping) WithReplacedVariables(context map[string]string) Mapping {
	result := Mapping{
		Active:    mapping.Active,
		Before:    mapping.Before,
		From:      mapping.From,
		MappingID: mapping.MappingID,
		Origin:    mapping.Origin,
		Proxy:     mapping.Proxy,
		Regexp:    mapping.Regexp,
		Tags:      mapping.Tags,
		To:        mapping.To,
	}

	result.From = variables.ReplaceVariables(result.From, context)
	result.To = variables.ReplaceVariables(result.To, context)

	inject := Injection{Headers: make(map[string]string)}
	result.Inject = inject
	for k, v := range mapping.Inject.Headers {
		inject.Headers[k] = variables.ReplaceVariables(v, context)
	}

	return result
}

func (mapping *Mapping) matchForProxyWithFrom(req *http.Request) *MatchResult {
	newPath := req.URL.Path[len(mapping.From):]

	if strings.HasPrefix(newPath, "/") {
		newPath = mapping.To + newPath
	} else {
		newPath = mapping.To + "/" + newPath
	}

	return &MatchResult{
		Mapping: *mapping,
		NewPath: newPath,
		Parts:   []string{req.URL.Path},
	}
}

func (mapping *Mapping) matchForStaticWithFrom(req *http.Request) *MatchResult {
	return &MatchResult{
		Mapping: *mapping,
		NewPath: path.Join(mapping.To, req.URL.Path[len(mapping.From):]),
		Parts:   []string{req.URL.Path},
	}
}

func (mapping *Mapping) matchWithFrom(req *http.Request) *MatchResult {
	if !strings.HasPrefix(req.URL.Path, mapping.From) {
		return nil
	}

	if mapping.Proxy {
		return mapping.matchForProxyWithFrom(req)
	}

	return mapping.matchForStaticWithFrom(req)
}

func (mapping *Mapping) matchWithRegexp(req *http.Request) *MatchResult {
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

	return nil
}
