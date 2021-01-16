package jsonpath

import "encoding/json"

type syntaxBasicNumericComparator struct {
}

func (c *syntaxBasicNumericComparator) typeCast(values []interface{}) bool {
	var foundValue bool
	for index := range values {
		switch typedValue := values[index].(type) {
		case float64:
			foundValue = true
		case json.Number:
			foundValue = true
			if floatNumber, err := typedValue.Float64(); err == nil {
				values[index] = floatNumber
			}
		default:
			values[index] = struct{}{}
		}
	}
	return foundValue
}
