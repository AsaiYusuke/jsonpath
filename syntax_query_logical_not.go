package jsonpath

type syntaxLogicalNot struct {
	query syntaxQuery
}

func (l *syntaxLogicalNot) compute(
	root interface{}, currentList []interface{}) []interface{} {

	computedList := l.query.compute(root, currentList)
	var hasValue bool
	for index := range computedList {
		if _, ok := computedList[index].(struct{}); ok {
			computedList[index] = true
			hasValue = true
		} else {
			computedList[index] = struct{}{}
		}
	}
	if hasValue {
		return computedList
	}
	return []interface{}{struct{}{}}
}
