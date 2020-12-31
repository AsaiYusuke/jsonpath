package jsonpath

import "encoding/json"

type syntaxBasicNumericComparator struct {
}

func (c *syntaxBasicNumericComparator) typeCast(values map[int]interface{}) {
	for index, value := range values {
		switch value.(type) {
		case float64:
			continue
		case json.Number:
			number := value.(json.Number)
			if floatNumber, err := number.Float64(); err == nil {
				values[index] = floatNumber
			}
		default:
			delete(values, index)
		}
	}
}
