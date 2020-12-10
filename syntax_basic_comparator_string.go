package jsonpath

type syntaxBasicStringComparator struct {
}

func (c *syntaxBasicStringComparator) typeCast(value interface{}) (interface{}, bool) {
	_, ok := value.(string)
	return value, ok
}
