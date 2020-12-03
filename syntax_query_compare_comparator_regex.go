package jsonpath

import "regexp"

type syntaxCompareRegex struct {
	regex *regexp.Regexp
}

func (r syntaxCompareRegex) comparator(left, _ interface{}) bool {
	leftValue, ok := left.(string)
	return ok && r.regex.Match([]byte(leftValue))
}
