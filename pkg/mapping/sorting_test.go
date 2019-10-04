package mapping

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSorting(t *testing.T) {
	mappingZero := Mapping{
		From:      "/some",
		MappingID: "mappingZero",
		To:        "doesnt_matter",
	}

	mappingOne := Mapping{
		From:      "/some/path",
		MappingID: "mappingOne",
		To:        "doesnt_matter",
	}

	mappingTwo := Mapping{
		MappingID: "mappingTwo",
		Regexp:    "/static/js(.*\\.chunk\\.js)",
		To:        "doesnt_matter",
	}

	mappingThree := Mapping{
		MappingID: "mappingThree",
		Regexp:    "/js(.*\\.chunk\\.js)",
		To:        "doesnt_matter",
	}

	mappingFour := Mapping{
		MappingID: "mappingFour",
		Regexp:    "/same/path",
		To:        "doesnt_matter",
	}

	mappingZero.Before = mappingFour.MappingID
	mappingThree.Before = mappingTwo.MappingID

	mappings := []Mapping{mappingZero, mappingOne, mappingTwo, mappingThree, mappingFour}
	expectedBySpecificity := []Mapping{mappingTwo, mappingThree, mappingFour, mappingOne, mappingZero}

	sortBySpecificity(mappings)

	assert.Equal(t, len(expectedBySpecificity), len(mappings), "Should have the same size")
	assert.Equal(t, expectedBySpecificity, mappings, "Should match specificity")

	expectedSorting := []Mapping{mappingThree, mappingTwo, mappingZero, mappingFour, mappingOne}
	mappings = sortMappings(mappings)

	assert.Equal(t, len(expectedSorting), len(mappings), "Should have the same size")
	assert.Equal(t, expectedSorting, mappings, "Should match final sorting")
}
