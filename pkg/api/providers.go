package api

// SupportedBenchmark represents simplified benchmark reference for provider list
type SupportedBenchmark struct {
	ID string `json:"id"`
}

// Provider represents provider specification
type ProviderResource struct {
	ID                  string               `json:"id"`
	Label               string               `json:"label"`
	SupportedBenchmarks []SupportedBenchmark `json:"supported_benchmarks,omitempty"`
}

// ProviderResourceList represents response for listing providers
type ProviderResourceList struct {
	TotalCount int                `json:"total_count"`
	Items      []ProviderResource `json:"items"`
}
