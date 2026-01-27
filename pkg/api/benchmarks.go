package api

// Benchmark represents benchmark specification
type BenchmarkResource struct {
	Resource
	Label       string   `json:"label"`
	Description string   `json:"description,omitempty"`
	Category    string   `json:"category,omitempty"`
	ProviderID  string   `json:"provider_id"`
	Tags        []string `json:"tags,omitempty"`
}

// BenchmarkResourceList represents list of benchmarks
type BenchmarkResourceList struct {
	TotalCount int                 `json:"total_count"`
	Items      []BenchmarkResource `json:"items"`
}
