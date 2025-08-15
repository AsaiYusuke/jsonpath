package syntax

type syntaxIndexSubscript struct {
	*syntaxBasicSubscript

	number    int
	isOmitted bool
}

// Dummy to satisfy syntaxSubscript; not used in normal paths.
func (i *syntaxIndexSubscript) forEachIndex(_ int, _ func(_ int)) {
	// no-op
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
