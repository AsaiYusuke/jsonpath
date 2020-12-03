package jsonpath

type syntaxIndex struct {
	*syntaxBasicSubscript

	number    int
	isOmitted bool
}

func (i syntaxIndex) getIndexes(src []interface{}) []int {
	index := i.number
	srcLength := len(src)

	if index < 0 {
		index = index + srcLength
	}

	if index < 0 || index >= srcLength {
		return make([]int, 0)
	}
	return []int{index}
}
