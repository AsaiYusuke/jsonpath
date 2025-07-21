package syntax

import "regexp"

type syntaxCompareRegex struct {
	*syntaxBasicStringTypeValidator

	regex *regexp.Regexp
}

func (r *syntaxCompareRegex) comparator(left []interface{}, _ interface{}) bool {
	var hasValue bool
	for leftIndex := range left {
		if left[leftIndex] == emptyEntity {
			continue
		}
		if r.regex.MatchString(left[leftIndex].(string)) {
			hasValue = true
		} else {
			left[leftIndex] = emptyEntity
		}
	}
	return hasValue
}
