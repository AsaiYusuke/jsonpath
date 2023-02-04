package jsonpath

type syntaxLogicalNot struct {
	query syntaxQuery
}

func (l *syntaxLogicalNot) compute(
	root interface{}, currentList []interface{}) []interface{} {

	computedList := l.query.compute(root, currentList)
	if len(computedList) == 1 {
		if computedList[0] == struct{}{} {
			return fullList
		}
		return emptyList
	}

	var hasValue bool
	for index := range computedList {
		if computedList[index] == struct{}{} {
			computedList[index] = true
			hasValue = true
		} else {
			computedList[index] = struct{}{}
		}
	}
	if hasValue {
		return computedList
	}
	return emptyList
}
