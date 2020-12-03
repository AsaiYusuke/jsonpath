package jsonpath

import "fmt"

// ErrorMemberNotExist represents the error that the object member specified in the JSONPath did not exist in the JSON object.
type ErrorMemberNotExist struct {
	path string
}

func (e ErrorMemberNotExist) Error() string {
	return fmt.Sprintf(`member did not exist (path=%s)`, e.path)
}
