package jsonpath

import "regexp"

type syntaxCompareRegex struct {
	*syntaxBasicStringComparator

	regex *regexp.Regexp
}

func (r syntaxCompareRegex) comparator(left, _ interface{}) bool {
	return r.regex.MatchString(left.(string))
}
