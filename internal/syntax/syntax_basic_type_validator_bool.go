package syntax

type syntaxBasicBoolTypeValidator struct {
}

func (c *syntaxBasicBoolTypeValidator) validate(values []any) bool {
	var foundValue bool
	for index := range values {
		switch values[index].(type) {
		case bool:
			foundValue = true
		default:
			values[index] = emptyEntity
		}
	}
	return foundValue
}
