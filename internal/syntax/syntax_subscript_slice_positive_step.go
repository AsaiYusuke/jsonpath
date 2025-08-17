package syntax

type syntaxSlicePositiveStepSubscript struct {
	*syntaxBasicSubscript

	start *syntaxIndexSubscript
	end   *syntaxIndexSubscript
	step  *syntaxIndexSubscript
}

func (s *syntaxSlicePositiveStepSubscript) forEachIndex(srcLength int, handleIndex func(index int)) {
	if s.step.number == 0 {
		return
	}

	loopStart, loopEnd := s.getLoopStart(srcLength), s.getLoopEnd(srcLength)
	for i := loopStart; i < loopEnd; i += s.step.number {
		handleIndex(i)
	}
}

func (s *syntaxSlicePositiveStepSubscript) getLoopStart(srcLength int) int {
	if s.start.isOmitted {
		return s.getNormalizedValue(0, srcLength)
	}
	return s.getNormalizedValue(s.start.number, srcLength)
}

func (s *syntaxSlicePositiveStepSubscript) getLoopEnd(srcLength int) int {
	if s.end.isOmitted {
		return s.getNormalizedValue(srcLength, srcLength)
	}
	return s.getNormalizedValue(s.end.number, srcLength)
}

func (s *syntaxSlicePositiveStepSubscript) getNormalizedValue(value int, srcLength int) int {
	if value > srcLength {
		return srcLength
	}
	if value < -srcLength {
		return 0
	}
	if value < 0 {
		return value + srcLength
	}
	return value
}
