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
	return fmt.Sprintf(`not supported (feature=%s, path=%s)`, e.Feature, e.Path)
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
	return fmt.Sprintf(`function not found (function=%s)`, e.Function)
}

func NewErrorFunctionNotFound(function string) ErrorFunctionNotFound {
	return ErrorFunctionNotFound{
		Function: function,
	}
}

// ErrorTypeUnmatched represents the error that the node type specified in the JSONPath did not exist in the JSON object.
type ErrorTypeUnmatched struct {
	Path         string
	ExpectedType string
	FoundType    string
}

func (e ErrorTypeUnmatched) Error() string {
	return fmt.Sprintf(`type unmatched (expected=%s, found=%s, path=%s)`, e.ExpectedType, e.FoundType, e.Path)
}

func NewErrorTypeUnmatched(path string, expected string, found string) ErrorTypeUnmatched {
	return ErrorTypeUnmatched{
		Path:         path,
		ExpectedType: expected,
		FoundType:    found,
	}
}

// ErrorMemberNotExist represents the error that the member specified in the JSONPath did not exist in the JSON object.
type ErrorMemberNotExist struct {
	Path string
}

func (e ErrorMemberNotExist) Error() string {
	return fmt.Sprintf(`member did not exist (path=%s)`, e.Path)
}

func NewErrorMemberNotExist(path string) ErrorMemberNotExist {
	return ErrorMemberNotExist{
		Path: path,
	}
}

// ErrorFunctionFailed represents the error that function execution failed in the JSONPath.
type ErrorFunctionFailed struct {
	FunctionName string
	Err          error
}

func (e ErrorFunctionFailed) Error() string {
	return fmt.Sprintf(`function failed (function=%s, error=%s)`, e.FunctionName, e.Err)
}

func NewErrorFunctionFailed(functionName string, errorString string) ErrorFunctionFailed {
	return ErrorFunctionFailed{
		FunctionName: functionName,
		Err:          fmt.Errorf(`%s`, errorString),
	}
}
