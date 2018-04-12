package mapping

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
)

// Injection represents parameters that can be injected into proxied requests
type Injection struct {
	Headers map[string]string
}

// Mapping represents a mapping configuration loaded from some file.
type Mapping struct {
	From      string    `json:"from"`
	Inject    Injection `json:"inject"`
	MappingID string    `json:"mappingID"`
	Origin    string    `json:"origin"`
	Proxy     bool      `json:"proxy"`
	Regexp    string    `json:"regexp"`
	To        string    `json:"to"`
}

// fromYAMLMapping creates a new mapping from a mapping loaded from a YAML file
func fromYAMLMapping(loaded mapping, origin string, proxy bool) Mapping {
	newMapping := Mapping{
		From:   loaded.From,
		Inject: loaded.Inject,
		Origin: origin,
		Proxy:  proxy,
		Regexp: loaded.Regexp,
		To:     loaded.To,
	}

	newMapping.MappingID = generateID(newMapping)
	return newMapping
}

func generateID(mapping Mapping) string {
	hasher := sha1.New()
	isProxy := "0"
	if mapping.Proxy {
		isProxy = "1"
	}
	headers := fmt.Sprintf("%s", mapping.Inject.Headers)
	hasher.Write([]byte(mapping.Origin + mapping.From + mapping.To + mapping.Regexp + isProxy + headers))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}
