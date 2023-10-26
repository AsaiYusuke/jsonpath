package jsonpath

type syntaxBasicAnyValueTypeValidator struct {
}

func (c *syntaxBasicAnyValueTypeValidator) validate(values []interface{}) bool {
	var foundValue bool
	for index := range values {
		if values[index] != emptyEntity {
			foundValue = true
		}
	}
	return foundValue
}
