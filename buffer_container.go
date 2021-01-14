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

func (b *bufferContainer) getSortedKeys(srcMap map[string]interface{}) *sort.StringSlice {
	length := len(srcMap)
	sortKeys := bufferContainerSortSliceSyncPool.Get().(*sort.StringSlice)
	if cap(*sortKeys) < length {
		*sortKeys = make(sort.StringSlice, length)
	}
	*sortKeys = (*sortKeys)[:length]
	index := 0
	for key := range srcMap {
		(*sortKeys)[index] = key
		index++
	}
	if len(*sortKeys) > 1 {
		sortKeys.Sort()
	}
	return sortKeys
}

func (b *bufferContainer) putSortSlice(sortKeys *sort.StringSlice) {
	if sortKeys != nil {
		bufferContainerSortSliceSyncPool.Put(sortKeys)
	}
}
