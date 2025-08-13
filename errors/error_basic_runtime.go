package errors

// ErrorBasicRuntime is a basic runtime error structure
type ErrorBasicRuntime struct {
	path             string
	remainingPathLen int
}

func (e *ErrorBasicRuntime) GetPath() string {
	return e.path
}

func (e *ErrorBasicRuntime) GetRemainingPathLen() int {
	return e.remainingPathLen
}
