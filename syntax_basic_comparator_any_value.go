package jsonpath

import "encoding/json"

type syntaxBasicAnyValueComparator struct {
}

func (c *syntaxBasicAnyValueComparator) typeCast(value interface{}) (interface{}, bool) {
	if number, ok := value.(json.Number); ok {
		if floatNumber, err := number.Float64(); err == nil {
			return floatNumber, true
		}
	}
	return value, true
}
