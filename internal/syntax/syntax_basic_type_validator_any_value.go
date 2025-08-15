package syntax

type syntaxBasicAnyValueTypeValidator struct {
}

func (c *syntaxBasicAnyValueTypeValidator) validate(values []any) bool {
	return true
}
