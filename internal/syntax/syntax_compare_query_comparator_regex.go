package syntax

import "regexp"

type syntaxCompareRegex struct {
	regex *regexp.Regexp
}

func (r *syntaxCompareRegex) comparator(left []any, _ any) bool {
	var hasValue bool
	for leftIndex := range left {
		if left[leftIndex] == emptyEntity {
			continue
		}
		switch leftValue := left[leftIndex].(type) {
		case string:
			if r.regex.MatchString(leftValue) {
				hasValue = true
			} else {
				left[leftIndex] = emptyEntity
			}
		default:
			left[leftIndex] = emptyEntity
		}
	}
	return hasValue
}
