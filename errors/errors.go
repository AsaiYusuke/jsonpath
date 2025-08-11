package errors

import "fmt"

// ErrorInvalidSyntax represents the error that have syntax error in the JSONPath.
type ErrorInvalidSyntax struct {
	Position int
	Reason   string
	Near     string
}

func (e ErrorInvalidSyntax) Error() string {
	return fmt.Sprintf(`invalid syntax (position=%d, reason=%s, near=%s)`, e.Position, e.Reason, e.Near)
}

func NewErrorInvalidSyntax(position int, reason string, near string) ErrorInvalidSyntax {
	return ErrorInvalidSyntax{
		Position: position,
		Reason:   reason,
		Near:     near,
	}
}

// ErrorInvalidArgument represents the error that argument specified in the JSONPath is treated as the invalid error in Go syntax.
type ErrorInvalidArgument struct {
	ArgumentName string
	Err          error
}

func (e ErrorInvalidArgument) Error() string {
	return fmt.Sprintf(`invalid argument (argument=%s, error=%s)`, e.ArgumentName, e.Err)
}

func NewErrorInvalidArgument(argument string, err error) ErrorInvalidArgument {
	return ErrorInvalidArgument{
		ArgumentName: argument,
		Err:          err,
	}
}

// ErrorNotSupported represents the error that the unsupported syntaxes specified in the JSONPath.
type ErrorNotSupported struct {
	Feature string
	Path    string
}

func (e ErrorNotSupported) Error() string {
	return fmt.Sprintf(`not supported (path=%s, feature=%s)`, e.Path, e.Feature)
}

func NewErrorNotSupported(feature string, path string) ErrorNotSupported {
	return ErrorNotSupported{
		Feature: feature,
		Path:    path,
	}
}

// ErrorFunctionNotFound represents the error that the function specified in the JSONPath is not found.
type ErrorFunctionNotFound struct {
	Function string
}

func (e ErrorFunctionNotFound) Error() string {
	return fmt.Sprintf(`function not found (path=%s)`, e.Function)
}

func NewErrorFunctionNotFound(function string) ErrorFunctionNotFound {
	return ErrorFunctionNotFound{
		Function: function,
	}
}

// ErrorTypeUnmatched represents the error that the node type specified in the JSONPath did not exist in the JSON object.
type ErrorTypeUnmatched struct {
	*ErrorBasicRuntime
	ExpectedType string
	FoundType    string
}

func (e ErrorTypeUnmatched) Error() string {
	return fmt.Sprintf(`type unmatched (path=%s, expected=%s, found=%s)`, e.ErrorBasicRuntime.GetPath(), e.ExpectedType, e.FoundType)
}

func NewErrorTypeUnmatched(path string, remainingPathLen int, expected string, found string) ErrorTypeUnmatched {
	return ErrorTypeUnmatched{
		ErrorBasicRuntime: &ErrorBasicRuntime{path: path, remainingPathLen: remainingPathLen},
		ExpectedType:      expected,
		FoundType:         found,
	}
}

// ErrorMemberNotExist represents the error that the member specified in the JSONPath did not exist in the JSON object.
type ErrorMemberNotExist struct {
	*ErrorBasicRuntime
}

func (e ErrorMemberNotExist) Error() string {
	return fmt.Sprintf(`member did not exist (path=%s)`, e.ErrorBasicRuntime.GetPath())
}

func NewErrorMemberNotExist(path string, remainingPathLen int) ErrorMemberNotExist {
	return ErrorMemberNotExist{
		ErrorBasicRuntime: &ErrorBasicRuntime{path: path, remainingPathLen: remainingPathLen},
	}
}

// ErrorFunctionFailed represents the error that function execution failed in the JSONPath.
type ErrorFunctionFailed struct {
	*ErrorBasicRuntime
	Err error
}

func (e ErrorFunctionFailed) Error() string {
	return fmt.Sprintf(`function failed (path=%s, error=%s)`, e.ErrorBasicRuntime.GetPath(), e.Err)
}

func NewErrorFunctionFailed(path string, remainingPathLen int, errorString string) ErrorFunctionFailed {
	return ErrorFunctionFailed{
		ErrorBasicRuntime: &ErrorBasicRuntime{path: path, remainingPathLen: remainingPathLen},
		Err:               fmt.Errorf(`%s`, errorString),
	}
}
