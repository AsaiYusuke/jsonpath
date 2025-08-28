package syntax

type syntaxIndexSubscript struct {
	*syntaxBasicSubscript

	number    int
	isOmitted bool
}

func (i *syntaxIndexSubscript) getIndex(srcLength int) int {
	if i.number >= srcLength {
		return -1
	}

	if i.number < 0 {
		return i.number + srcLength
	}

	return i.number
}

func (i *syntaxIndexSubscript) count(srcLength int) int {
	if i.getIndex(srcLength) >= 0 {
		return 1
	}
	return 0
}

func (i *syntaxIndexSubscript) indexAt(srcLength int, ordinal int) int {
	return i.getIndex(srcLength)
}
