package jsonpath

import "fmt"

// ErrorInvalidArgument represents the error that argument specified in the JSONPath is treated as the invalid error in Go syntax.
type ErrorInvalidArgument struct {
	argument string
	err      error
}

func (e ErrorInvalidArgument) Error() string {
	return fmt.Sprintf(`invalid argument (argument=%s, error=%s)`, e.argument, e.err)
}
