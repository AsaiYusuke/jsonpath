package jsonpath

// Accessor represents the accessor to the result elements of JSONPath.
type Accessor struct {
	Get func() interface{}
	Set func(interface{})
}
