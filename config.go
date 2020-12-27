package jsonpath

// Config represents the configuration parameters.
type Config struct {
	filterFunctions    map[string]func(interface{}) (interface{}, error)
	aggregateFunctions map[string]func([]interface{}) (interface{}, error)
	accessorMode       bool
}

// SetFilterFunction sets the custom function.
func (c *Config) SetFilterFunction(id string, function func(interface{}) (interface{}, error)) {
	if c.filterFunctions == nil {
		c.filterFunctions = map[string]func(interface{}) (interface{}, error){}
	}
	c.filterFunctions[id] = function
}

// SetAggregateFunction sets the custom function.
func (c *Config) SetAggregateFunction(id string, function func([]interface{}) (interface{}, error)) {
	if c.aggregateFunctions == nil {
		c.aggregateFunctions = map[string]func([]interface{}) (interface{}, error){}
	}
	c.aggregateFunctions[id] = function
}

// SetAccessorMode sets a collection of accessors to the result.
func (c *Config) SetAccessorMode() {
	c.accessorMode = true
}
