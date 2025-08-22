package tests

import (
	"testing"

	syntax "github.com/AsaiYusuke/jsonpath/v2/internal/syntax"
)

func TestFilterTest_SlicesGrow(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$[?(@)]`,
			inputJSON:    `{"01":1,"02":2,"03":3,"04":4,"05":5,"06":6,"07":7,"08":8,"09":9,"10":10,"11":11}`,
			expectedJSON: `[1,2,3,4,5,6,7,8,9,10,11]`,
		},
	}

	runTestCasesSerial(t, "TestFilterTest_SlicesGrow", tests, func() {
		syntax.ResetNodeSliceSyncPool()
	})
}
