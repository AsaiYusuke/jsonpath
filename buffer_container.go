package jsonpath

import (
	"sort"
	"sync"
)

type bufferContainer struct {
	result []interface{}
}

var bufferContainerSortSliceSyncPool = &sync.Pool{
	New: func() interface{} { return new(sort.StringSlice) },
}

func (b *bufferContainer) getSortSlice(length int) *sort.StringSlice {
	sortKeys := bufferContainerSortSliceSyncPool.Get().(*sort.StringSlice)

	if cap(*sortKeys) < length {
		*sortKeys = make(sort.StringSlice, length)
	}
	*sortKeys = (*sortKeys)[:length]
	return sortKeys
}

func (b *bufferContainer) putSortSlice(sortKeys *sort.StringSlice) {
	if sortKeys != nil {
		bufferContainerSortSliceSyncPool.Put(sortKeys)
	}
}
