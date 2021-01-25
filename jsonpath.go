package jsonpath

import (
	"regexp"
	"sync"
)

var parseMutex sync.Mutex
var parser = pegJSONPathParser{}
var unescapeRegex = regexp.MustCompile(`\\(.)`)

// Retrieve returns the retrieved JSON using the given JSONPath.
func Retrieve(jsonPath string, src interface{}, config ...Config) ([]interface{}, error) {
	jsonPathFunc, err := Parse(jsonPath, config...)
	if err != nil {
		return nil, err
	}
	return jsonPathFunc(src)
}

// Parse returns the parser function using the given JSONPath.
func Parse(jsonPath string, config ...Config) (f func(src interface{}) ([]interface{}, error), err error) {
	parseMutex.Lock()
	defer func() {
		if exception := recover(); exception != nil {
			switch typedException := exception.(type) {
			case error:
				err = typedException
			default:
				panic(typedException)
			}
		}
		parser.jsonPathParser = jsonPathParser{}
		parseMutex.Unlock()
	}()

	parser.Buffer = jsonPath

	if parser.parse == nil {
		parser.Init()
	} else {
		parser.Reset()
	}

	parser.jsonPathParser.unescapeRegex = unescapeRegex

	if len(config) > 0 {
		parser.jsonPathParser.filterFunctions = config[0].filterFunctions
		parser.jsonPathParser.aggregateFunctions = config[0].aggregateFunctions
		parser.jsonPathParser.accessorMode = config[0].accessorMode
	}

	parser.Parse()
	parser.Execute()

	root := parser.jsonPathParser.root
	return func(src interface{}) ([]interface{}, error) {
		container := bufferContainer{}

		err := root.retrieve(src, src, &container)
		if err != nil {
			return container.result, err.(error)
		}
		return container.result, nil

	}, nil
}
