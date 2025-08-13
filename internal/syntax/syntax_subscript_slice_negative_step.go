package syntax

type syntaxSliceNegativeStepSubscript struct {
	*syntaxBasicSubscript

	start *syntaxIndexSubscript
	end   *syntaxIndexSubscript
	step  *syntaxIndexSubscript
}

func (s *syntaxSliceNegativeStepSubscript) getIndexes(srcLength int) []int {
	loopStart := s.getLoopStart(srcLength)
	loopEnd := s.getLoopEnd(srcLength)

	index, result := 0, make([]int, srcLength)
	if s.step.number < 0 {
		for i := loopStart; i > loopEnd; i += s.step.number {
			result[index] = i
			index++
		}
	}

	return result[:index]
}

func (s *syntaxSliceNegativeStepSubscript) getLoopStart(srcLength int) int {
	loopStart := s.start.number
	if s.start.isOmitted {
		loopStart = srcLength - 1
	}
	return s.getNormalizedValue(loopStart, srcLength)
}

func (s *syntaxSliceNegativeStepSubscript) getLoopEnd(srcLength int) int {
	loopEnd := s.end.number
	if s.end.isOmitted {
		loopEnd = -srcLength - 1
	}
	return s.getNormalizedValue(loopEnd, srcLength)
}

func (s *syntaxSliceNegativeStepSubscript) getNormalizedValue(value int, srcLength int) int {
	if value < 0 {
		value += srcLength
		if value < -1 {
			value = -1
		}
	}
	if value > srcLength-1 {
		value = srcLength - 1
	}
	return value
}
