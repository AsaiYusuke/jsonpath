package jsonpath

type syntaxBasicStringComparator struct {
}

func (c *syntaxBasicStringComparator) typeCast(values []interface{}) {
	for index := range values {
		if _, ok := values[index].(string); !ok {
			values[index] = struct{}{}
		}
	}
}
