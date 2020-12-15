package jsonpath

type syntaxSlicePositiveStepSubscript struct {
	*syntaxBasicSubscript

	start *syntaxIndexSubscript
	end   *syntaxIndexSubscript
	step  *syntaxIndexSubscript
}

func (s *syntaxSlicePositiveStepSubscript) getIndexes(src []interface{}) []int {
	srcLength := len(src)
	loopStart := s.getLoopStart(srcLength)
	loopEnd := s.getLoopEnd(srcLength)

	index, result := 0, make([]int, srcLength)
	for i := loopStart; i < loopEnd; i += s.step.number {
		if i < 0 || i >= srcLength {
			break
		}
		result[index] = i
		index++
	}

	return result[:index]
}

func (s *syntaxSlicePositiveStepSubscript) getLoopStart(srcLength int) int {
	loopStart := s.start.number
	if s.start.isOmitted {
		loopStart = 0
	}
	if loopStart < 0 {
		loopStart = loopStart + srcLength
		if loopStart < 0 {
			loopStart = 0
		}
	}
	return loopStart
}

func (s *syntaxSlicePositiveStepSubscript) getLoopEnd(srcLength int) int {
	loopEnd := s.end.number
	if s.end.isOmitted {
		loopEnd = srcLength
	}
	if loopEnd < 0 {
		loopEnd = loopEnd + srcLength
		if loopEnd < 0 {
			loopEnd = 0
		}
	}
	return loopEnd
}
