package jsonpath

type syntaxIndexSubscript struct {
	*syntaxBasicSubscript

	number    int
	isOmitted bool
}

func (i *syntaxIndexSubscript) getIndexes(src []interface{}) []int {
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
