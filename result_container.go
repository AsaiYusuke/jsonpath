package jsonpath

type resultContainer struct {
	result []interface{}
}

func (c *resultContainer) checkResult() {
	if c.result == nil {
		c.result = make([]interface{}, 0)
	}
}

func (c *resultContainer) append(entry interface{}) {
	c.checkResult()
	c.result = append(c.result, entry)
}

func (c *resultContainer) hasResult() bool {
	c.checkResult()
	return len(c.result) > 0
}

func (c *resultContainer) getResult() []interface{} {
	c.checkResult()
	return c.result
}
