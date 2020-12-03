package jsonpath

import "fmt"

// ErrorNoneMatched represents the error that the child paths specified in the JSONPath result in empty output.
type ErrorNoneMatched struct {
	path string
}

func (e ErrorNoneMatched) Error() string {
	return fmt.Sprintf(`none matched (path=%s)`, e.path)
}
