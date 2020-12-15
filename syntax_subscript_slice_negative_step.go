package jsonpath

type syntaxSliceNegativeStepSubscript struct {
	*syntaxBasicSubscript

	start *syntaxIndexSubscript
	end   *syntaxIndexSubscript
	step  *syntaxIndexSubscript
}

func (s *syntaxSliceNegativeStepSubscript) getIndexes(src []interface{}) []int {
	srcLength := len(src)
	loopStart := s.getLoopStart(srcLength)
	loopEnd := s.getLoopEnd(srcLength)

	index, result := 0, make([]int, srcLength)
	for i := loopStart; i > loopEnd; i += s.step.number {
		if i < 0 || i >= srcLength {
			break
		}
		result[index] = i
		index++
	}

	return result[:index]
}

func (s *syntaxSliceNegativeStepSubscript) getLoopStart(srcLength int) int {
	loopStart := s.start.number
	if s.start.isOmitted {
		loopStart = srcLength - 1
	}
	if loopStart > srcLength-1 {
		loopStart = loopStart - srcLength
		if loopStart > srcLength-1 {
			loopStart = srcLength - 1
		}
	}
	return loopStart
}

func (s *syntaxSliceNegativeStepSubscript) getLoopEnd(srcLength int) int {
	loopEnd := s.end.number
	if s.end.isOmitted {
		loopEnd = -1
	}
	if loopEnd >= srcLength-1 {
		loopEnd = loopEnd - srcLength
		if loopEnd > srcLength {
			loopEnd = srcLength
		}
	}
	return loopEnd
}
