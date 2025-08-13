package jsonpath

import (
	"github.com/AsaiYusuke/jsonpath/internal/syntax"
)

func Retrieve(jsonPath string, src any, config ...Config) ([]any, error) {
	return syntax.Retrieve(jsonPath, src, config...)
}

func Parse(jsonPath string, config ...Config) (func(src any) ([]any, error), error) {
	return syntax.Parse(jsonPath, config...)
}
