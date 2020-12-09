package jsonpath

type syntaxSlicePositiveStep struct {
	*syntaxBasicSubscript

	start syntaxIndex
	end   syntaxIndex
	step  syntaxIndex
}

func (s syntaxSlicePositiveStep) getIndexes(src []interface{}) []int {
	srcLength := len(src)
	loopStart := s.getLoopStart(srcLength)
	loopEnd := s.getLoopEnd(srcLength)

	result := make([]int, 0)

	for i := loopStart; i < loopEnd; i += s.step.number {
		if i < 0 || i >= srcLength {
			break
		}
		result = append(result, i)
	}

	return result
}

func (s syntaxSlicePositiveStep) getLoopStart(srcLength int) int {
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

func (s syntaxSlicePositiveStep) getLoopEnd(srcLength int) int {
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
