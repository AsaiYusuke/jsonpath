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

var nodeSliceSyncPool = &sync.Pool{
	New: func() any {
		slice := make([]any, 0, 10)
		return &slice
	},
}

func getSortedKeys(srcMap map[string]any) (*sort.StringSlice, int) {
	mapLength := len(srcMap)
	sortKeys := sortSliceSyncPool.Get().(*sort.StringSlice)
	if cap(*sortKeys) < mapLength {
		*sortKeys = make(sort.StringSlice, mapLength)
	}
	*sortKeys = (*sortKeys)[:mapLength]
	index := 0
	for key := range srcMap {
		(*sortKeys)[index] = key
		index++
	}
	if mapLength > 1 {
		sort.Sort(sortKeys)
	}
	return sortKeys, mapLength
}

func getSortedRecursiveKeys(srcMap map[string]any) (*sort.StringSlice, int) {
	mapLength := len(srcMap)
	sortKeys := sortSliceSyncPool.Get().(*sort.StringSlice)
	if cap(*sortKeys) < mapLength {
		*sortKeys = make(sort.StringSlice, mapLength)
	}
	*sortKeys = (*sortKeys)[:0]
	for key, value := range srcMap {
		switch value.(type) {
		case map[string]any, []any:
			*sortKeys = append(*sortKeys, key)
		}
	}
	keyLength := len(*sortKeys)
	if keyLength > 1 {
		sort.Sort(sortKeys)
	}
	return sortKeys, keyLength
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

func getNodeSlice() *[]any { return nodeSliceSyncPool.Get().(*[]any) }

func putNodeSlice(nodes *[]any) {
	*nodes = (*nodes)[:0]
	nodeSliceSyncPool.Put(nodes)
}

func ResetNodeSliceSyncPool() {
	nodeSliceSyncPool = &sync.Pool{
		New: func() any {
			slice := make([]any, 0, 10)
			return &slice
		},
	}
}
