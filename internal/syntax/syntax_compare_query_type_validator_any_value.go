package syntax

type syntaxAnyValueTypeValidator struct {
}

func (c *syntaxAnyValueTypeValidator) validate(values []any) bool {
	return true
}
