package jsonpath

import (
	"github.com/AsaiYusuke/jsonpath/internal/syntax"
)

func Retrieve(jsonPath string, src interface{}, config ...Config) ([]interface{}, error) {
	return syntax.Retrieve(jsonPath, src, config...)
}

func Parse(jsonPath string, config ...Config) (func(src interface{}) ([]interface{}, error), error) {
	return syntax.Parse(jsonPath, config...)
}
