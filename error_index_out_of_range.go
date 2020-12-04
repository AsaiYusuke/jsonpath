package jsonpath

import "fmt"

// ErrorIndexOutOfRange represents the error that the array indexes specified in the JSONPath are out of range.
type ErrorIndexOutOfRange struct {
	path string
}

func (e ErrorIndexOutOfRange) Error() string {
	return fmt.Sprintf(`index out of range (path=%s)`, e.path)
}