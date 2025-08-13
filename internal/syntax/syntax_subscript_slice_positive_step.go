package syntax

type syntaxSlicePositiveStepSubscript struct {
	*syntaxBasicSubscript

	start *syntaxIndexSubscript
	end   *syntaxIndexSubscript
	step  *syntaxIndexSubscript
}

func (s *syntaxSlicePositiveStepSubscript) getIndexes(srcLength int) []int {
	loopStart := s.getLoopStart(srcLength)
	loopEnd := s.getLoopEnd(srcLength)

	index, result := 0, make([]int, srcLength)
	if s.step.number > 0 {
		for i := loopStart; i < loopEnd; i += s.step.number {
			result[index] = i
			index++
		}
	}

	return result[:index]
}

func (s *syntaxSlicePositiveStepSubscript) getLoopStart(srcLength int) int {
	loopStart := s.start.number
	if s.start.isOmitted {
		loopStart = 0
	}
	return s.getNormalizedValue(loopStart, srcLength)
}

func (s *syntaxSlicePositiveStepSubscript) getLoopEnd(srcLength int) int {
	loopEnd := s.end.number
	if s.end.isOmitted {
		loopEnd = srcLength
	}
	return s.getNormalizedValue(loopEnd, srcLength)
}

func (s *syntaxSlicePositiveStepSubscript) getNormalizedValue(value int, srcLength int) int {
	if value < 0 {
		value += srcLength
		if value < 0 {
			value = 0
		}
	}
	if value > srcLength {
		value = srcLength
	}
	return value
}
