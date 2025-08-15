package syntax

type syntaxNilTypeValidator struct {
}

func (c *syntaxNilTypeValidator) validate(values []any) bool {
	var foundValue bool
	for index := range values {
		switch values[index].(type) {
		case nil:
			foundValue = true
		default:
			values[index] = emptyEntity
		}
	}
	return foundValue
}
