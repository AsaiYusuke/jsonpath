package jsonpath

type syntaxIndex struct {
	*syntaxBasicSubscript

	number    int
	isOmitted bool
}

func (i *syntaxIndex) getIndexes(src []interface{}) []int {
	index := i.number
	srcLength := len(src)

	if index < 0 {
		index += srcLength
	}

	if index < 0 || index >= srcLength {
		return []int{}
	}
	return []int{index}
}
