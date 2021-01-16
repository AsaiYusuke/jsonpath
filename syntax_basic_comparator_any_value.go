package jsonpath

import "encoding/json"

type syntaxBasicAnyValueComparator struct {
}

func (c *syntaxBasicAnyValueComparator) typeCast(values []interface{}) bool {
	for index := range values {
		if number, ok := values[index].(json.Number); ok {
			if floatNumber, err := number.Float64(); err == nil {
				values[index] = floatNumber
			}
		}
	}
	return len(values) > 0
}
