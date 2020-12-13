package jsonpath

import (
	"regexp"
	"strconv"
)

type jsonPathParser struct {
	root          syntaxNode
	srcJSON       *interface{}
	resultPtr     *[]interface{}
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

func (j *jsonPathParser) updateResultPtr(checkNode syntaxNode, result **[]interface{}) {
	for checkNode != nil {
		checkNode.setResultPtr(result)
		checkNode = checkNode.getNext()
	}
}

func (j *jsonPathParser) syntaxErr(pos int, reason string, buffer string) {
	j.thisError = ErrorInvalidSyntax{pos, reason, buffer[pos:]}
}

func (j *jsonPathParser) hasErr() bool {
	return j.thisError != nil
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
