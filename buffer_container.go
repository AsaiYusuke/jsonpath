package jsonpath

import (
	"sort"
	"sync"
)

type bufferContainer struct {
	result   []interface{}
	sortKeys *sort.StringSlice
}

var bufferContainerSortSliceSyncPool = &sync.Pool{
	New: func() interface{} { return new(sort.StringSlice) },
}

func (b *bufferContainer) expandSortSlice(length int) {
	if b.sortKeys == nil {
		b.sortKeys = bufferContainerSortSliceSyncPool.Get().(*sort.StringSlice)
	}
	if cap(*b.sortKeys) < length {
		*b.sortKeys = make(sort.StringSlice, length)
	}
	*b.sortKeys = (*b.sortKeys)[:length]
}

func (b *bufferContainer) putSortSlice() {
	if b.sortKeys != nil {
		bufferContainerSortSliceSyncPool.Put(b.sortKeys)
	}
}
