package config

import (
	"crypto/sha1"
	"encoding/base64"
)

// DynamicMapping represents a mapping that can be active or not
type DynamicMapping struct {
	Active    bool   `json:"active"`
	From      string `json:"from"`
	MappingID string `json:"mappingID"`
	Origin    string `json:"origin"`
	Proxy     bool   `json:"proxy"`
	To        string `json:"to"`
}

// Mapping represents a mapping configuration loaded from some file.
type Mapping struct {
	From      string `json:"from"`
	MappingID string `json:"mappingID"`
	Origin    string `json:"origin"`
	Proxy     bool   `json:"proxy"`
	To        string `json:"to"`
}

// fromYAMLMapping creates a new mapping from a mapping loaded from a YAML config
func fromYAMLMapping(loaded mapping, origin string, proxy bool) Mapping {
	newMapping := Mapping{
		From:   loaded.From,
		To:     loaded.To,
		Origin: origin,
		Proxy:  proxy,
	}

	newMapping.MappingID = generateID(newMapping)
	return newMapping
}

func generateID(mapping Mapping) string {
	hasher := sha1.New()
	hasher.Write([]byte(mapping.Origin + mapping.From + mapping.To))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}
