package jsonpath

import "fmt"

// ErrorFunctionFailed represents the error that the function specified in the JSONPath failed.
type ErrorFunctionFailed struct {
	function string
	err      error
}

func (e ErrorFunctionFailed) Error() string {
	return fmt.Sprintf(`function failed (function=%s, error=%s)`, e.function, e.err)
}
