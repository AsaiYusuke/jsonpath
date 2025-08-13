package syntax

type syntaxIndexSubscript struct {
	*syntaxBasicSubscript

	number    int
	isOmitted bool
}

func (i *syntaxIndexSubscript) getIndexes(srcLength int) []int {
	return []int{}
}

func (i *syntaxIndexSubscript) getIndex(srcLength int) int {
	if i.number < -srcLength || i.number >= srcLength {
		return -1
	}

	if i.number < 0 {
		return i.number + srcLength
	}

	return i.number
}
