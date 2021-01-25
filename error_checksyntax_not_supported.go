package jsonpath

import "fmt"

// ErrorNotSupported represents the error that the unsupported syntaxes specified in the JSONPath.
type ErrorNotSupported struct {
	feature string
	path    string
}

func (e ErrorNotSupported) Error() string {
	return fmt.Sprintf(`not supported (feature=%s, path=%s)`, e.feature, e.path)
}
