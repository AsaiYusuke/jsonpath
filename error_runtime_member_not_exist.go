package jsonpath

import "fmt"

// ErrorMemberNotExist represents the error that the member specified in the JSONPath did not exist in the JSON object.
type ErrorMemberNotExist struct {
	*errorBasicRuntime
}

func (e ErrorMemberNotExist) Error() string {
	return fmt.Sprintf(`member did not exist (path=%s)`, e.node.text)
}
