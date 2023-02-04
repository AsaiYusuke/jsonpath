package jsonpath

type syntaxLogicalNot struct {
	query     syntaxQuery
	emptyList []interface{}
	fullList  []interface{}
}

func (l *syntaxLogicalNot) compute(
	root interface{}, currentList []interface{}) []interface{} {

	computedList := l.query.compute(root, currentList)
	if len(computedList) == 1 {
		if computedList[0] == struct{}{} {
			return l.fullList
		}
		return l.emptyList
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
	return l.emptyList
}
