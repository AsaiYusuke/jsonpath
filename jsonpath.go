package jsonpath

import (
	"regexp"
)

// Retrieve returns the retrieved JSON using the given JSONPath.
func Retrieve(jsonPath string, src interface{}) ([]interface{}, error) {
	jsonPathFunc, err := Parse(jsonPath)
	if err != nil {
		return nil, err
	}
	return jsonPathFunc(src)
}

// Parse returns the parser function using the given JSONPath.
func Parse(jsonPath string) (func(src interface{}) ([]interface{}, error), error) {
	unescapeRegex, _ := regexp.Compile(`\\(.)`)

	parser := pegJSONPathParser{
		Buffer: jsonPath,
		jsonPathParser: jsonPathParser{
			resultPtr:     &[]interface{}{},
			unescapeRegex: unescapeRegex,
		},
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
