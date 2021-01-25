package jsonpath

import "fmt"

// ErrorInvalidSyntax represents the error that have syntax error in the JSONPath.
type ErrorInvalidSyntax struct {
	position int
	reason   string
	near     string
}

func (e ErrorInvalidSyntax) Error() string {
	return fmt.Sprintf(`invalid syntax (position=%d, reason=%s, near=%s)`, e.position, e.reason, e.near)
}
