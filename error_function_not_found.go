package jsonpath

import "fmt"

// ErrorFunctionNotFound represents the error that the function specified in the JSONPath is not found.
type ErrorFunctionNotFound struct {
	function string
}

func (e ErrorFunctionNotFound) Error() string {
	return fmt.Sprintf(`function not found (function=%s)`, e.function)
}
