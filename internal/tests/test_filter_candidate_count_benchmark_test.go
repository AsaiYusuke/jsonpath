package tests

import (
	"fmt"
	"testing"

	"github.com/AsaiYusuke/jsonpath/v2"
)

func BenchmarkFilterComparisonCandidates(b *testing.B) {
	for _, size := range []int{1, 100, 10000} {
		b.Run(fmt.Sprint(size), func(b *testing.B) {
			src := make([]any, size)
			for index := range src {
				value := float64(2)
				if index == 0 {
					value = 1
				}
				src[index] = map[string]any{"a": value}
			}

			parserFunc, err := jsonpath.Parse(`$[?(@.a==1)]`)
			if err != nil {
				b.Fatal(err)
			}

			result := make([]any, 0, size)
			args := []*[]any{&result}
			b.ReportAllocs()
			b.ResetTimer()
			for b.Loop() {
				if _, err := parserFunc(src, args...); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
