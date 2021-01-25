package jsonpath

import "fmt"

// ErrorTypeUnmatched represents the error that the node type specified in the JSONPath did not exist in the JSON object.
type ErrorTypeUnmatched struct {
	*errorBasicRuntime

	expectedType string
	foundType    string
}

func (e ErrorTypeUnmatched) Error() string {
	return fmt.Sprintf(`type unmatched (expected=%s, found=%s, path=%s)`, e.expectedType, e.foundType, e.node.text)
}
