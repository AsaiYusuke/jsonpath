package config_test

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/AsaiYusuke/jsonpath/v2"
	"github.com/AsaiYusuke/jsonpath/v2/config"
)

func ExampleAccessor() {
	cfg := config.Config{}
	cfg.SetAccessorMode()
	jsonPath, srcJSON := `$.a`, `{"a":1,"b":0}`
	var src any
	json.Unmarshal([]byte(srcJSON), &src)
	output, err := jsonpath.Retrieve(jsonPath, src, cfg)
	if err != nil {
		fmt.Printf(`type: %v, value: %v`, reflect.TypeOf(err), err)
		return
	}
	accessor := output[0].(config.Accessor)
	srcMap := src.(map[string]any)

	fmt.Printf("Get : %v\n", accessor.Get())

	accessor.Set(2)
	fmt.Printf("Set -> Src : %v\n", srcMap[`a`])

	accessor.Set(3)
	fmt.Printf("Set -> Get : %v\n", accessor.Get())

	srcMap[`a`] = 4
	fmt.Printf("Src -> Get : %v\n", accessor.Get())

	// Output:
	// Get : 1
	// Set -> Src : 2
	// Set -> Get : 3
	// Src -> Get : 4
}
