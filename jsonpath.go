package jsonpath

import (
	"regexp"
)

// Retrieve returns the retrieved JSON using the given JSONPath.
func Retrieve(jsonPath string, src interface{}, config ...Config) ([]interface{}, error) {
	jsonPathFunc, err := Parse(jsonPath, config...)
	if err != nil {
		return nil, err
	}
	return jsonPathFunc(src)
}

// Parse returns the parser function using the given JSONPath.
func Parse(jsonPath string, config ...Config) (func(src interface{}) ([]interface{}, error), error) {
	parser := pegJSONPathParser{
		Buffer: jsonPath,
		jsonPathParser: jsonPathParser{
			resultPtr: &[]interface{}{},
		},
	}

	parser.unescapeRegex, _ = regexp.Compile(`\\(.)`)

	if len(config) > 0 {
		parser.filterFunctions = config[0].filterFunctions
		parser.aggregateFunctions = config[0].aggregateFunctions
	}

	parser.Init()
	parser.Parse()
	parser.Execute()

	if parser.thisError != nil {
		return nil, parser.thisError
	}

	return func(src interface{}) ([]interface{}, error) {
		result := make([]interface{}, 0)
		parser.srcJSON = &src
		parser.resultPtr = &result
		if err := parser.root.retrieve(src); err != nil {
			return nil, err
		}
		return result, nil
	}, nil
}
