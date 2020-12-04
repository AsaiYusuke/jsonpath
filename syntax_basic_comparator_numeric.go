package jsonpath

import "encoding/json"

type syntaxBasicNumericComparator struct {
}

func (c *syntaxBasicNumericComparator) typeCast(value interface{}) (interface{}, bool) {
	switch value.(type) {
	case float64:
		return value, true
	case json.Number:
		leftNumber := value.(json.Number)
		if leftNumberFloat, err := leftNumber.Float64(); err == nil {
			return leftNumberFloat, true
		}
	}
	return 0, false
}
