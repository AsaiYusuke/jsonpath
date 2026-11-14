package syntax

import "sync"

const (
	maxQueryScratchBuffers  = 16
	maxQueryScratchCapacity = 1 << 16
)

// Query expressions can keep multiple intermediate slices live until evaluation ends.
type syntaxQueryScratch struct {
	buffers [][]any
	next    int
}

var syntaxQueryScratchSyncPool = &sync.Pool{
	New: func() any { return &syntaxQueryScratch{} },
}

func getSyntaxQueryScratch() *syntaxQueryScratch {
	return syntaxQueryScratchSyncPool.Get().(*syntaxQueryScratch)
}

func (s *syntaxQueryScratch) getSlice(length int) []any {
	if s.next == len(s.buffers) {
		s.buffers = append(s.buffers, nil)
	}

	buffer := s.buffers[s.next]
	if cap(buffer) < length {
		buffer = make([]any, length)
	} else {
		buffer = buffer[:length]
	}
	s.buffers[s.next] = buffer
	s.next++
	return buffer
}

func (s *syntaxQueryScratch) reset() {
	for index, buffer := range s.buffers {
		clear(buffer)
		if cap(buffer) > maxQueryScratchCapacity {
			s.buffers[index] = nil
		} else {
			s.buffers[index] = buffer[:0]
		}
	}

	if len(s.buffers) > maxQueryScratchBuffers {
		clear(s.buffers[maxQueryScratchBuffers:])
		s.buffers = s.buffers[:maxQueryScratchBuffers]
	}
	s.next = 0
}

func putSyntaxQueryScratch(s *syntaxQueryScratch) {
	s.reset()
	syntaxQueryScratchSyncPool.Put(s)
}
