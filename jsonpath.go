package jsonpath

import (
	"fmt"
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
	if len(j.params) < 1 {
		j.thisError = fmt.Errorf(`internal error (empty queue)`)
		return nil
	}
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
		if len(varBlockSet) != 2 {
			return block
		}
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

	regex, err := regexp.Compile(`\\(.)`)
	if err != nil {
		return err
	}
	parser.unescapeRegex = regex

	parser.Init()
	err = parser.Parse()
	if err != nil {
		return err
	}

	parser.Execute()
	if parser.thisError != nil {
		return parser.thisError
	}

	if len(parser.params) > 0 {
		return fmt.Errorf(`internal error (%v)`, parser.params)
	}

	j.root = parser.root
	return nil
}

func (j *jsonPathParser) retrieve(src interface{}) ([]interface{}, error) {
	result := resultContainer{}
	if err := j.root.retrieve(src, src, &result); err != nil {
		return nil, err
	}

	return result.getResult(), nil
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
		return jsonpath.retrieve(src)
	}, nil
}
