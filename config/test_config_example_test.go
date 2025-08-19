package config_test

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/AsaiYusuke/jsonpath/v2"
	"github.com/AsaiYusuke/jsonpath/v2/config"
	"github.com/AsaiYusuke/jsonpath/v2/errors"
)

func ExampleConfig_SetFilterFunction() {
	config := config.Config{}
	config.SetFilterFunction(`twice`, func(param any) (any, error) {
		if floatParam, ok := param.(float64); ok {
			return floatParam * 2, nil
		}
		return nil, fmt.Errorf(`type error`)
	})
	jsonPath, srcJSON := `$[*].twice()`, `[1,3]`
	var src any
	json.Unmarshal([]byte(srcJSON), &src)
	output, err := jsonpath.Retrieve(jsonPath, src, config)
	switch err.(type) {
	case errors.ErrorFunctionNotFound:
		fmt.Printf(`type: %v, value: %v`, reflect.TypeOf(err), err)
		return
	}
	outputJSON, _ := json.Marshal(output)
	fmt.Println(string(outputJSON))
	// Output:
	// [2,6]
}

func ExampleConfig_SetAggregateFunction() {
	config := config.Config{}
	config.SetAggregateFunction(`max`, func(params []any) (any, error) {
		var result float64
		for _, param := range params {
			if floatParam, ok := param.(float64); ok {
				if result < floatParam {
					result = floatParam
				}
				continue
			}
			return nil, fmt.Errorf(`type error`)
		}
		return result, nil
	})
	jsonPath, srcJSON := `$[*].max()`, `[1,3]`
	var src any
	json.Unmarshal([]byte(srcJSON), &src)
	output, err := jsonpath.Retrieve(jsonPath, src, config)
	if err != nil {
		fmt.Printf(`type: %v, value: %v`, reflect.TypeOf(err), err)
		return
	}
	outputJSON, _ := json.Marshal(output)
	fmt.Println(string(outputJSON))
	// Output:
	// [3]
}

func ExampleConfig_SetAccessorMode() {
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
