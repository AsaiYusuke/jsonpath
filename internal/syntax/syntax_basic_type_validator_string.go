package syntax

type syntaxBasicStringTypeValidator struct {
}

func (c *syntaxBasicStringTypeValidator) validate(values []any) bool {
	var foundValue bool
	for index := range values {
		switch values[index].(type) {
		case string:
			foundValue = true
		case struct{}:
		default:
			values[index] = emptyEntity
		}
	}
	return foundValue
}
