package syntax

type syntaxBasicAnyValueTypeValidator struct {
}

func (c *syntaxBasicAnyValueTypeValidator) validate(values []any) bool {
	for index := range values {
		if values[index] != emptyEntity {
			return true
		}
	}
	return false
}
