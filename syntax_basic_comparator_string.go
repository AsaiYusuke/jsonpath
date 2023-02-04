package jsonpath

type syntaxBasicStringComparator struct {
}

func (c *syntaxBasicStringComparator) typeCast(values []interface{}) bool {
	var foundValue bool
	for index := range values {
		switch values[index].(type) {
		case string:
			foundValue = true
		case struct{}:
		default:
			values[index] = emptyEntity
		}
	}
	return foundValue
}
