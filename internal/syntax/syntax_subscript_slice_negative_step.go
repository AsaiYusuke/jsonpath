package syntax

type syntaxSliceNegativeStepSubscript struct {
	*syntaxBasicSubscript

	start *syntaxIndexSubscript
	end   *syntaxIndexSubscript
	step  *syntaxIndexSubscript
}

func (s *syntaxSliceNegativeStepSubscript) forEachIndex(srcLength int, handleIndex func(index int)) {
	loopStart, loopEnd := s.getLoopStart(srcLength), s.getLoopEnd(srcLength)
	for i := loopStart; i > loopEnd; i += s.step.number {
		handleIndex(i)
	}
}

func (s *syntaxSliceNegativeStepSubscript) getLoopStart(srcLength int) int {
	if s.start.isOmitted {
		return s.getNormalizedValue(srcLength-1, srcLength)
	}
	return s.getNormalizedValue(s.start.number, srcLength)
}

func (s *syntaxSliceNegativeStepSubscript) getLoopEnd(srcLength int) int {
	if s.end.isOmitted {
		return s.getNormalizedValue(-srcLength-1, srcLength)
	}
	return s.getNormalizedValue(s.end.number, srcLength)
}

func (s *syntaxSliceNegativeStepSubscript) getNormalizedValue(value int, srcLength int) int {
	if value > srcLength-1 {
		return srcLength - 1
	}
	if value < -srcLength-1 {
		return -1
	}
	if value < 0 {
		return value + srcLength
	}
	return value
}
