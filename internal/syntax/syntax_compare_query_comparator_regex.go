package syntax

import "regexp"

type syntaxCompareRegex struct {
	*syntaxStringTypeValidator

	regex *regexp.Regexp
}

func (r *syntaxCompareRegex) comparator(left []any, _ any) bool {
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
