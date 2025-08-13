package syntax

import (
	"sort"
	"sync"
)

var sortSliceSyncPool = &sync.Pool{
	New: func() any { return new(sort.StringSlice) },
}

var resultSyncPool = &sync.Pool{
	New: func() any { return new(bufferContainer) },
}

}

func getSortedKeys(srcMap map[string]interface{}) *sort.StringSlice {
	length := len(srcMap)
	sortKeys := sortSliceSyncPool.Get().(*sort.StringSlice)
	if cap(*sortKeys) < length {
		*sortKeys = make(sort.StringSlice, length)
	}
	*sortKeys = (*sortKeys)[:length]
	index := 0
	for key := range srcMap {
		(*sortKeys)[index] = key
		index++
	}
	if length > 1 {
		sortKeys.Sort()
	}
	return sortKeys
}

func putSortSlice(sortKeys *sort.StringSlice) {
	if sortKeys != nil {
		sortSliceSyncPool.Put(sortKeys)
	}
}

func getContainer() *bufferContainer {
	return resultSyncPool.Get().(*bufferContainer)
}

func putContainer(container *bufferContainer) {
	container.result = container.result[:0]
	resultSyncPool.Put(container)
}
