package config

// Accessor represents the accessor to the result nodes of JSONPath.
type Accessor struct {
	Get func() any
	Set func(any)
}
