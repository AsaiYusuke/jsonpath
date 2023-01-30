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
			values[index], _ = typedValue.Float64()
		default:
			values[index] = struct{}{}
		}
	}
	return foundValue
}
