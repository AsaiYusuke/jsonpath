package jsonpath

import (
	"math"
)

type syntaxSlice struct {
	*syntaxBasicSubscript

	start syntaxIndex
	end   syntaxIndex
	step  syntaxIndex
}

func (s syntaxSlice) getIndexes(src []interface{}) []int {
	loopStep := s.step.number
	if s.step.isOmitted || s.step.number == 0 {
		loopStep = 1
	}

	direction := 1
	if loopStep < 0 {
		direction = -1
	}

	srcLength := len(src)

	loopStart := s.start.number
	if direction > 0 {
		if s.start.isOmitted {
			loopStart = 0
		}
		if loopStart < 0 {
			loopStart = int(math.Max(float64(loopStart+srcLength), 0))
		}
	} else {
		if s.start.isOmitted {
			loopStart = srcLength - 1
		}
		if loopStart > srcLength-1 {
			loopStart = int(math.Min(float64(loopStart-srcLength), float64(srcLength-1)))
		}
	}

	loopEnd := s.end.number
	if direction > 0 {
		if s.end.isOmitted {
			loopEnd = srcLength
		}
		if loopEnd < 0 {
			loopEnd = int(math.Max(float64(loopEnd+srcLength), 0))
		}
	} else {
		if s.end.isOmitted {
			loopEnd = -1
		}
		if loopEnd >= srcLength-1 {
			loopEnd = int(math.Min(float64(loopEnd-srcLength), float64(srcLength)))
		}
	}

	result := make([]int, 0)

	for i := loopStart; i*direction < loopEnd*direction; i += loopStep {
		if i < 0 || i >= srcLength {
			break
		}
		result = append(result, i)
	}

	return result
}
