package jsonpath

import "fmt"

// ErrorTypeUnmatched represents the error that the node type specified in the JSONPath did not exist in the JSON object.
type ErrorTypeUnmatched struct {
	expectedType string
	foundType    string
	path         string
}

func (e ErrorTypeUnmatched) Error() string {
	return fmt.Sprintf(`type unmatched (expected=%s, found=%s, path=%s)`, e.expectedType, e.foundType, e.path)
}
