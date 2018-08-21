package mapping

import (
	"sort"
	"strings"
)

func ensureBefore(mappings []Mapping, mappingID string, beforeID string) []Mapping {
	indexEl := -1
	indexBefore := -1

	for i, n := range mappings {
		if n.MappingID == mappingID {
			indexEl = i
		}

		if n.MappingID == beforeID {
			indexBefore = i
		}
	}

	// We only need to move things around if not in the right order
	if indexBefore < indexEl {
		above := make([]Mapping, indexBefore)
		copy(above, mappings[:indexBefore])

		between := make([]Mapping, indexEl-indexBefore)
		copy(between, mappings[indexBefore:indexEl])

		after := make([]Mapping, len(mappings)-indexEl-1)
		copy(after, mappings[indexEl+1:])

		mappings = append(above, mappings[indexEl])
		mappings = append(mappings, between...)
		mappings = append(mappings, after...)
	}

	return mappings
}

func sortMappings(mappings []Mapping) []Mapping {
	sortBySpecificity(mappings)

	// Move accordingly to before attribute
	for _, mapping := range mappings {
		if mapping.Before != "" {
			mappings = ensureBefore(mappings, mapping.MappingID, mapping.Before)
		}
	}

	return mappings
}

func sortBySpecificity(mappings []Mapping) {
	sort.Slice(mappings, func(i, j int) bool {
		pathI := strings.ToLower(mappings[i].From)
		pathJ := strings.ToLower(mappings[j].From)

		if pathI == "" {
			pathI = strings.ToLower(mappings[i].Regexp)
		}

		if pathJ == "" {
			pathJ = strings.ToLower(mappings[j].Regexp)
		}

		if len(pathI) == len(pathJ) {
			return strings.Compare(pathI, pathJ) < 0
		}

		return len(pathI) > len(pathJ)
	})
}
