package jsonpath

type syntaxBasicStringComparator struct {
}

func (c *syntaxBasicStringComparator) typeCast(values []interface{}) bool {
	var foundValue bool
	for index := range values {
		if _, ok := values[index].(string); ok {
			foundValue = true
		} else {
			values[index] = emptyEntity
		}
	}
	return foundValue
}
