package jsonpath

import (
	"encoding/json"
	"testing"
)

func BenchmarkTypeCast(b *testing.B) {
	node := &syntaxBasicNumericComparator{}

	for i := 0; i < b.N; i++ {
		for i := 0; i < 10000; i++ {
			node.typeCast([]interface{}{json.Number(`1`)})
		}
	}
}

// BenchmarkTypeCast-4   	    2227	    505221 ns/op	   80000 B/op	   10000 allocs/op
