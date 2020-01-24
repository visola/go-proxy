package configuration

// Origin is where the upstream was loaded from
type Origin struct {
	File     string `json:"file"`
	LoadedAt int64  `json:"loadedAt"`
}