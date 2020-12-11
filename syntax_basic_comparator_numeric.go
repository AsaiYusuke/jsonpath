package jsonpath

import "encoding/json"

type syntaxBasicNumericComparator struct {
}

func (c *syntaxBasicNumericComparator) typeCast(value interface{}) (interface{}, bool) {
	switch value.(type) {
	case float64:
		return value, true
	case json.Number:
		number := value.(json.Number)
		if numberFloat, err := number.Float64(); err == nil {
			return numberFloat, true
		}
	}
	return 0, false
}
