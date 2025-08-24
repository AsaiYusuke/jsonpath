package syntax

import (
	"regexp"
	"sync"

	"github.com/AsaiYusuke/jsonpath/v2/config"
)

var parseMutex sync.Mutex
var parser = pegJSONPathParser[uint32]{}
var unescapeRegex = regexp.MustCompile(`\\(.)`)

// Retrieve returns the retrieved JSON using the given JSONPath.
func Retrieve(jsonPath string, src any, config ...config.Config) ([]any, error) {
	jsonPathFunc, err := Parse(jsonPath, config...)
	if err != nil {
		return nil, err
	}
	return jsonPathFunc(src)
}

// Parse returns the parser function using the given JSONPath.
func Parse(jsonPath string, config ...config.Config) (f func(src any, dst ...*[]any) ([]any, error), err error) {
	parseMutex.Lock()
	defer func() {
		if exception := recover(); exception != nil {
			if _err, ok := exception.(error); ok {
				f = nil
				err = _err
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
		parser.jsonPathParser.filterFunctions = config[0].FilterFunctions
		parser.jsonPathParser.aggregateFunctions = config[0].AggregateFunctions
		parser.jsonPathParser.accessorMode = config[0].AccessorMode
	}

	parser.Parse()
	parser.Execute()

	root := parser.jsonPathParser.root
	return func(src any, dst ...*[]any) ([]any, error) {
		var buf *[]any
		usePool := true
		if len(dst) > 0 && dst[0] != nil {
			*(dst[0]) = (*(dst[0]))[:0]
			buf = dst[0]
			usePool = false
		} else {
			buf = getNodeSlice()
		}

		if err := root.retrieve(src, src, buf); err != nil {
			if usePool {
				putNodeSlice(buf)
			}
			return nil, err
		}

		if usePool {
			res := *buf
			out := make([]any, len(res))
			copy(out, res)
			putNodeSlice(buf)
			return out, nil
		}

		return *buf, nil
	}, nil
}
