package jsonpath

type syntaxBasicAnyValueTypeValidator struct {
}

func (c *syntaxBasicAnyValueTypeValidator) validate(values []interface{}) bool {
	for index := range values {
		if values[index] != emptyEntity {
			return true
		}
	}
	return false
}
