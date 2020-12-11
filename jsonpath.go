package jsonpath

import (
	"regexp"
	"strconv"
)

type jsonPathParser struct {
	root          syntaxNode
	params        []interface{}
	thisError     error
	unescapeRegex *regexp.Regexp
}

func (j *jsonPathParser) push(param interface{}) {
	j.params = append(j.params, param)
}

func (j *jsonPathParser) pop() interface{} {
	var param interface{}
	param, j.params = j.params[len(j.params)-1], j.params[:len(j.params)-1]
	return param
}

func (j *jsonPathParser) toInt(text string) int {
	value, err := strconv.Atoi(text)
	if err != nil {
		j.thisError = ErrorInvalidArgument{text, err}
		return 0
	}
	return value
}

func (j *jsonPathParser) toFloat(text string) float64 {
	value, err := strconv.ParseFloat(text, 64)
	if err != nil {
		j.thisError = ErrorInvalidArgument{text, err}
		return 0
	}
	return value
}

func (j *jsonPathParser) unescape(text string) string {
	return j.unescapeRegex.ReplaceAllStringFunc(text, func(block string) string {
		varBlockSet := j.unescapeRegex.FindStringSubmatch(block)
		return varBlockSet[1]
	})
}

func (j *jsonPathParser) syntaxErr(pos int, reason string, buffer string) {
	j.thisError = ErrorInvalidSyntax{pos, reason, buffer[pos:]}
}

func (j *jsonPathParser) hasErr() bool {
	return j.thisError != nil
}

func (j *jsonPathParser) parse(jsonPath string) error {
	parser := parser{Buffer: jsonPath}

	regex, _ := regexp.Compile(`\\(.)`)
	parser.unescapeRegex = regex

	parser.Init()
	parser.Parse()
	parser.Execute()

	if parser.thisError != nil {
		return parser.thisError
	}

	j.root = parser.root
	return nil
}

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
	jsonpath := jsonPathParser{}
	if err := jsonpath.parse(jsonPath); err != nil {
		return nil, err
	}
	return func(src interface{}) ([]interface{}, error) {
		result := make([]interface{}, 0)
		if err := jsonpath.root.retrieve(src, src, &result); err != nil {
			return nil, err
		}
		return result, nil
	}, nil
}
