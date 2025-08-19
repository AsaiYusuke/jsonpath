package jsonpath

import (
	"github.com/AsaiYusuke/jsonpath/v2/config"
	"github.com/AsaiYusuke/jsonpath/v2/internal/syntax"
)

func Retrieve(jsonPath string, src any, config ...config.Config) ([]any, error) {
	return syntax.Retrieve(jsonPath, src, config...)
}

func Parse(jsonPath string, config ...config.Config) (func(src any) ([]any, error), error) {
	return syntax.Parse(jsonPath, config...)
}
