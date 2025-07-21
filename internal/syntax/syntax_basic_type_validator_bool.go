package syntax

type syntaxBasicBoolTypeValidator struct {
}

func (c *syntaxBasicBoolTypeValidator) validate(values []interface{}) bool {
	var foundValue bool
	for index := range values {
		switch values[index].(type) {
		case bool:
			foundValue = true
		case struct{}:
		default:
			values[index] = emptyEntity
		}
	}
	return foundValue
}
