package upstream

import "strings"

// Endpoints represents a sortable array of endpoints
type Endpoints []Endpoint

func (a Endpoints) Len() int {
	return len(a)
}

func (a Endpoints) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a Endpoints) Less(i, j int) bool {
	e1 := a[i]
	e2 := a[j]

	path1 := e1.Path()
	path2 := e2.Path()
	pathParts1 := strings.Split(path1, "/")
	pathParts2 := strings.Split(path2, "/")

	if len(pathParts1) != len(pathParts2) {
		return len(pathParts1) > len(pathParts2)
	}

	for i, pathPart1 := range pathParts1 {
		if len(pathPart1) != len(pathParts2[i]) {
			return len(pathPart1) > len(pathParts2[i])
		}
	}

	return path1 < path2
}
