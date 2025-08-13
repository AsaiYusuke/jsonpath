package syntax

type syntaxIndexSubscript struct {
	*syntaxBasicSubscript

	number    int
	isOmitted bool
}

func (i *syntaxIndexSubscript) getIndexes(srcLength int) []int {
	index := i.number

	if index < 0 {
		index += srcLength
	}

	if index < 0 || index >= srcLength {
		return []int{}
	}

	return []int{index}
}
