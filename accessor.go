package jsonpath

// Accessor represents the accessor to the result nodes of JSONPath.
type Accessor struct {
	Get func() interface{}
	Set func(interface{})
}
