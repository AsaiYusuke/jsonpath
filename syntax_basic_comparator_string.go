package jsonpath

type syntaxBasicStringComparator struct {
}

func (c *syntaxBasicStringComparator) typeCast(values map[int]interface{}) {
	for index := range values {
		if _, ok := values[index].(string); !ok {
			delete(values, index)
		}
	}
}
