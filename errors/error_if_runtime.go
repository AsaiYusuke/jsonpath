package errors

type ErrorRuntime interface {
	GetPath() string
	GetRemainingPathLen() int
	Error() string
}
