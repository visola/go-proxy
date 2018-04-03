package mapping

import (
	"crypto/sha1"
	"encoding/base64"
)

// Mapping represents a mapping configuration loaded from some file.
type Mapping struct {
	From      string `json:"from"`
	MappingID string `json:"mappingID"`
	Origin    string `json:"origin"`
	Proxy     bool   `json:"proxy"`
	Regexp    string `json:"regexp"`
	To        string `json:"to"`
}

// fromYAMLMapping creates a new mapping from a mapping loaded from a YAML file
func fromYAMLMapping(loaded mapping, origin string, proxy bool) Mapping {
	newMapping := Mapping{
		From:   loaded.From,
		To:     loaded.To,
		Origin: origin,
		Proxy:  proxy,
		Regexp: loaded.Regexp,
	}

	newMapping.MappingID = generateID(newMapping)
	return newMapping
}

func generateID(mapping Mapping) string {
	hasher := sha1.New()
	hasher.Write([]byte(mapping.Origin + mapping.From + mapping.To + mapping.Regexp))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}
