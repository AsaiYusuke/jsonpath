package jsonpath

type syntaxBasicNilTypeValidator struct {
}

func (c *syntaxBasicNilTypeValidator) validate(values []interface{}) bool {
	var foundValue bool
	for index := range values {
		switch values[index].(type) {
		case nil:
			foundValue = true
		case struct{}:
		default:
			values[index] = emptyEntity
		}
	}
	return foundValue
}
