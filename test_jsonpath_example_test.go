package jsonpath_test

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/AsaiYusuke/jsonpath"
	"github.com/AsaiYusuke/jsonpath/errors"
)

func Example() {
	jsonPath, srcJSON := `$.key`, `{"key":"value"}`
	var src any
	json.Unmarshal([]byte(srcJSON), &src)
	output, _ := jsonpath.Retrieve(jsonPath, src)
	outputJSON, _ := json.Marshal(output)
	fmt.Println(string(outputJSON))
	// Output:
	// ["value"]
}

func ExampleRetrieve() {
	jsonPath, srcJSON := `$.key`, `{"key":"value"}`
	var src any
	json.Unmarshal([]byte(srcJSON), &src)
	output, err := jsonpath.Retrieve(jsonPath, src)
	if err != nil {
		fmt.Printf(`type: %v, value: %v`, reflect.TypeOf(err), err)
		return
	}
	outputJSON, _ := json.Marshal(output)
	fmt.Println(string(outputJSON))
	// Output:
	// ["value"]
}

func ExampleParse() {
	jsonPath := `$.key`
	srcJSON1 := `{"key":"value1"}`
	srcJSON2 := `{"key":"value2"}`
	jsonPathParser, err := jsonpath.Parse(jsonPath)
	if err != nil {
		fmt.Printf(`type: %v, value: %v`, reflect.TypeOf(err), err)
		return
	}
	var src1, src2 any
	json.Unmarshal([]byte(srcJSON1), &src1)
	json.Unmarshal([]byte(srcJSON2), &src2)
	output1, err := jsonPathParser(src1)
	if err != nil {
		fmt.Printf(`type: %v, value: %v`, reflect.TypeOf(err), err)
		return
	}
	output2, err := jsonPathParser(src2)
	if err != nil {
		fmt.Printf(`type: %v, value: %v`, reflect.TypeOf(err), err)
		return
	}
	outputJSON1, _ := json.Marshal(output1)
	outputJSON2, _ := json.Marshal(output2)
	fmt.Println(string(outputJSON1))
	fmt.Println(string(outputJSON2))
	// Output:
	// ["value1"]
	// ["value2"]
}

func ExampleErrorInvalidSyntax() {
	jsonPath, srcJSON := `$.`, `{}`
	var src any
	json.Unmarshal([]byte(srcJSON), &src)
	output, err := jsonpath.Retrieve(jsonPath, src)
	switch err.(type) {
	case errors.ErrorInvalidSyntax:
		fmt.Printf(`type: %v, value: %v`, reflect.TypeOf(err), err)
		return
	}
	outputJSON, _ := json.Marshal(output)
	fmt.Println(string(outputJSON))
	// Output:
	// type: errors.ErrorInvalidSyntax, value: invalid syntax (position=1, reason=unrecognized input, near=.)
}

func ExampleErrorInvalidArgument() {
	jsonPath, srcJSON := `$[?(1.0.0>0)]`, `{}`
	var src any
	json.Unmarshal([]byte(srcJSON), &src)
	output, err := jsonpath.Retrieve(jsonPath, src)
	switch err.(type) {
	case errors.ErrorInvalidArgument:
		fmt.Printf(`type: %v, value: %v`, reflect.TypeOf(err), err)
		return
	}
	outputJSON, _ := json.Marshal(output)
	fmt.Println(string(outputJSON))
	// Output:
	// type: errors.ErrorInvalidArgument, value: invalid argument (argument=1.0.0, error=strconv.ParseFloat: parsing "1.0.0": invalid syntax)
}

func ExampleErrorFunctionNotFound() {
	jsonPath, srcJSON := `$.unknown()`, `{}`
	var src any
	json.Unmarshal([]byte(srcJSON), &src)
	output, err := jsonpath.Retrieve(jsonPath, src)
	switch err.(type) {
	case errors.ErrorFunctionNotFound:
		fmt.Printf(`type: %v, value: %v`, reflect.TypeOf(err), err)
		return
	}
	outputJSON, _ := json.Marshal(output)
	fmt.Println(string(outputJSON))
	// Output:
	// type: errors.ErrorFunctionNotFound, value: function not found (path=.unknown())
}

func ExampleErrorNotSupported() {
	jsonPath, srcJSON := `$[(command)]`, `{}`
	var src any
	json.Unmarshal([]byte(srcJSON), &src)
	output, err := jsonpath.Retrieve(jsonPath, src)
	switch err.(type) {
	case errors.ErrorNotSupported:
		fmt.Printf(`type: %v, value: %v`, reflect.TypeOf(err), err)
		return
	}
	outputJSON, _ := json.Marshal(output)
	fmt.Println(string(outputJSON))
	// Output:
	// type: errors.ErrorNotSupported, value: not supported (path=[(command)], feature=script)
}

func ExampleErrorMemberNotExist() {
	jsonPath, srcJSON := `$.none`, `{}`
	var src any
	json.Unmarshal([]byte(srcJSON), &src)
	output, err := jsonpath.Retrieve(jsonPath, src)
	switch err.(type) {
	case errors.ErrorMemberNotExist:
		fmt.Printf(`type: %v, value: %v`, reflect.TypeOf(err), err)
		return
	}
	outputJSON, _ := json.Marshal(output)
	fmt.Println(string(outputJSON))
	// Output:
	// type: errors.ErrorMemberNotExist, value: member did not exist (path=.none)
}

func ExampleErrorTypeUnmatched() {
	jsonPath, srcJSON := `$.a`, `[]`
	var src any
	json.Unmarshal([]byte(srcJSON), &src)
	output, err := jsonpath.Retrieve(jsonPath, src)
	switch err.(type) {
	case errors.ErrorTypeUnmatched:
		fmt.Printf(`type: %v, value: %v`, reflect.TypeOf(err), err)
		return
	}
	outputJSON, _ := json.Marshal(output)
	fmt.Println(string(outputJSON))
	// Output:
	// type: errors.ErrorTypeUnmatched, value: type unmatched (path=.a, expected=object, found=[]interface {})
}

func ExampleErrorFunctionFailed() {
	config := jsonpath.Config{}
	config.SetFilterFunction(`invalid`, func(param any) (any, error) {
		return nil, fmt.Errorf(`invalid function executed`)
	})
	jsonPath, srcJSON := `$.invalid()`, `{}`
	var src any
	json.Unmarshal([]byte(srcJSON), &src)
	output, err := jsonpath.Retrieve(jsonPath, src, config)
	switch err.(type) {
	case errors.ErrorFunctionFailed:
		fmt.Printf(`type: %v, value: %v`, reflect.TypeOf(err), err)
		return
	}
	outputJSON, _ := json.Marshal(output)
	fmt.Println(string(outputJSON))
	// Output:
	// type: errors.ErrorFunctionFailed, value: function failed (path=.invalid(), error=invalid function executed)
}

func ExampleConfig_SetFilterFunction() {
	config := jsonpath.Config{}
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
	config := jsonpath.Config{}
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
	config := jsonpath.Config{}
	config.SetAccessorMode()
	jsonPath, srcJSON := `$.a`, `{"a":1,"b":0}`
	var src any
	json.Unmarshal([]byte(srcJSON), &src)
	output, err := jsonpath.Retrieve(jsonPath, src, config)
	if err != nil {
		fmt.Printf(`type: %v, value: %v`, reflect.TypeOf(err), err)
		return
	}
	accessor := output[0].(jsonpath.Accessor)
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
