package config

// Config represents the configuration parameters.
type Config struct {
	FilterFunctions    map[string]func(any) (any, error)
	AggregateFunctions map[string]func([]any) (any, error)
	AccessorMode       bool
}

// SetFilterFunction sets the custom function.
func (c *Config) SetFilterFunction(id string, function func(any) (any, error)) {
	if c.FilterFunctions == nil {
		c.FilterFunctions = map[string]func(any) (any, error){}
	}
	c.FilterFunctions[id] = function
}

// SetAggregateFunction sets the custom function.
func (c *Config) SetAggregateFunction(id string, function func([]any) (any, error)) {
	if c.AggregateFunctions == nil {
		c.AggregateFunctions = map[string]func([]any) (any, error){}
	}
	c.AggregateFunctions[id] = function
}

// SetAccessorMode sets a collection of accessors to the result.
func (c *Config) SetAccessorMode() {
	c.AccessorMode = true
}
