package syntax

// Config represents the configuration parameters.
type Config struct {
	filterFunctions    map[string]func(any) (any, error)
	aggregateFunctions map[string]func([]any) (any, error)
	accessorMode       bool
}

// SetFilterFunction sets the custom function.
func (c *Config) SetFilterFunction(id string, function func(any) (any, error)) {
	if c.filterFunctions == nil {
		c.filterFunctions = map[string]func(any) (any, error){}
	}
	c.filterFunctions[id] = function
}

// SetAggregateFunction sets the custom function.
func (c *Config) SetAggregateFunction(id string, function func([]any) (any, error)) {
	if c.aggregateFunctions == nil {
		c.aggregateFunctions = map[string]func([]any) (any, error){}
	}
	c.aggregateFunctions[id] = function
}

// SetAccessorMode sets a collection of accessors to the result.
func (c *Config) SetAccessorMode() {
	c.accessorMode = true
}
