package jsonpath

import "encoding/json"

type syntaxBasicNumericComparator struct {
}

func (c *syntaxBasicNumericComparator) typeCast(values []interface{}) {
	for index, value := range values {
		switch typedValue := value.(type) {
		case float64:
			continue
		case json.Number:
			if floatNumber, err := typedValue.Float64(); err == nil {
				values[index] = floatNumber
			}
		default:
			values[index] = struct{}{}
		}
	}
}
