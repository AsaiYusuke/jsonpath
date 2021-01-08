package jsonpath_test

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/AsaiYusuke/jsonpath"
)

func ExampleRetrieve() {
	jsonPath, srcJSON := `$.key`, `{"key":"value"}`
	var src interface{}
	json.Unmarshal([]byte(srcJSON), &src)
	output, err := jsonpath.Retrieve(jsonPath, src)
	if err != nil {
		fmt.Printf(`%v, %v`, reflect.TypeOf(err), err)
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
		fmt.Printf(`%v, %v`, reflect.TypeOf(err), err)
		return
	}
	var src1, src2 interface{}
	json.Unmarshal([]byte(srcJSON1), &src1)
	json.Unmarshal([]byte(srcJSON2), &src2)
	output1, err := jsonPathParser(src1)
	if err != nil {
		fmt.Printf(`%v, %v`, reflect.TypeOf(err), err)
		return
	}
	output2, err := jsonPathParser(src2)
	if err != nil {
		fmt.Printf(`%v, %v`, reflect.TypeOf(err), err)
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
	var src interface{}
	json.Unmarshal([]byte(srcJSON), &src)
	output, err := jsonpath.Retrieve(jsonPath, src)
	if err != nil {
		fmt.Printf(`%v, %v`, reflect.TypeOf(err), err)
		return
	}
	outputJSON, _ := json.Marshal(output)
	fmt.Println(string(outputJSON))
	// Output:
	// jsonpath.ErrorInvalidSyntax, invalid syntax (position=1, reason=unrecognized input, near=.)
}

func ExampleErrorInvalidArgument() {
	jsonPath, srcJSON := `$[?(1.0.0>0)]`, `{}`
	var src interface{}
	json.Unmarshal([]byte(srcJSON), &src)
	output, err := jsonpath.Retrieve(jsonPath, src)
	if err != nil {
		fmt.Printf(`%v, %v`, reflect.TypeOf(err), err)
		return
	}
	outputJSON, _ := json.Marshal(output)
	fmt.Println(string(outputJSON))
	// Output:
	// jsonpath.ErrorInvalidArgument, invalid argument (argument=1.0.0, error=strconv.ParseFloat: parsing "1.0.0": invalid syntax)
}

func ExampleErrorFunctionNotFound() {
	jsonPath, srcJSON := `$.unknown()`, `{}`
	var src interface{}
	json.Unmarshal([]byte(srcJSON), &src)
	output, err := jsonpath.Retrieve(jsonPath, src)
	if err != nil {
		fmt.Printf(`%v, %v`, reflect.TypeOf(err), err)
		return
	}
	outputJSON, _ := json.Marshal(output)
	fmt.Println(string(outputJSON))
	// Output:
	// jsonpath.ErrorFunctionNotFound, function not found (function=.unknown())
}

func ExampleErrorNotSupported() {
	jsonPath, srcJSON := `$[(command)]`, `{}`
	var src interface{}
	json.Unmarshal([]byte(srcJSON), &src)
	output, err := jsonpath.Retrieve(jsonPath, src)
	if err != nil {
		fmt.Printf(`%v, %v`, reflect.TypeOf(err), err)
		return
	}
	outputJSON, _ := json.Marshal(output)
	fmt.Println(string(outputJSON))
	// Output:
	// jsonpath.ErrorNotSupported, not supported (feature=script, path=[(command)])
}

func ExampleErrorMemberNotExist() {
	jsonPath, srcJSON := `$.none`, `{}`
	var src interface{}
	json.Unmarshal([]byte(srcJSON), &src)
	output, err := jsonpath.Retrieve(jsonPath, src)
	if err != nil {
		fmt.Printf(`%v, %v`, reflect.TypeOf(err), err)
		return
	}
	outputJSON, _ := json.Marshal(output)
	fmt.Println(string(outputJSON))
	// Output:
	// jsonpath.ErrorMemberNotExist, member did not exist (path=.none)
}

func ExampleErrorIndexOutOfRange() {
	jsonPath, srcJSON := `$[1]`, `["a"]`
	var src interface{}
	json.Unmarshal([]byte(srcJSON), &src)
	output, err := jsonpath.Retrieve(jsonPath, src)
	if err != nil {
		fmt.Printf(`%v, %v`, reflect.TypeOf(err), err)
		return
	}
	outputJSON, _ := json.Marshal(output)
	fmt.Println(string(outputJSON))
	// Output:
	// jsonpath.ErrorIndexOutOfRange, index out of range (path=[1])
}

func ExampleErrorTypeUnmatched() {
	jsonPath, srcJSON := `$.a`, `[]`
	var src interface{}
	json.Unmarshal([]byte(srcJSON), &src)
	output, err := jsonpath.Retrieve(jsonPath, src)
	if err != nil {
		fmt.Printf(`%v, %v`, reflect.TypeOf(err), err)
		return
	}
	outputJSON, _ := json.Marshal(output)
	fmt.Println(string(outputJSON))
	// Output:
	// jsonpath.ErrorTypeUnmatched, type unmatched (expected=object, found=[]interface {}, path=.a)
}

func ExampleErrorNoneMatched() {
	jsonPath, srcJSON := `$[1,2]`, `["a"]`
	var src interface{}
	json.Unmarshal([]byte(srcJSON), &src)
	output, err := jsonpath.Retrieve(jsonPath, src)
	if err != nil {
		fmt.Printf(`%v, %v`, reflect.TypeOf(err), err)
		return
	}
	outputJSON, _ := json.Marshal(output)
	fmt.Println(string(outputJSON))
	// Output:
	// jsonpath.ErrorNoneMatched, none matched (path=[1,2])
}

func ExampleErrorFunctionFailed() {
	config := jsonpath.Config{}
	config.SetFilterFunction(`invalid`, func(param interface{}) (interface{}, error) {
		return nil, fmt.Errorf(`invalid function executed`)
	})
	jsonPath, srcJSON := `$.invalid()`, `{}`
	var src interface{}
	json.Unmarshal([]byte(srcJSON), &src)
	output, err := jsonpath.Retrieve(jsonPath, src, config)
	if err != nil {
		fmt.Printf(`%v, %v`, reflect.TypeOf(err), err)
		return
	}
	outputJSON, _ := json.Marshal(output)
	fmt.Println(string(outputJSON))
	// Output:
	// jsonpath.ErrorFunctionFailed, function failed (function=.invalid(), error=invalid function executed)
}

func ExampleConfig_SetFilterFunction() {
	config := jsonpath.Config{}
	config.SetFilterFunction(`twice`, func(param interface{}) (interface{}, error) {
		if floatParam, ok := param.(float64); ok {
			return floatParam * 2, nil
		}
		return nil, fmt.Errorf(`type error`)
	})
	jsonPath, srcJSON := `$[*].twice()`, `[1,3]`
	var src interface{}
	json.Unmarshal([]byte(srcJSON), &src)
	output, err := jsonpath.Retrieve(jsonPath, src, config)
	if err != nil {
		fmt.Printf(`%v, %v`, reflect.TypeOf(err), err)
		return
	}
	outputJSON, _ := json.Marshal(output)
	fmt.Println(string(outputJSON))
	// Output:
	// [2,6]
}

func ExampleConfig_SetAggregateFunction() {
	config := jsonpath.Config{}
	config.SetAggregateFunction(`max`, func(params []interface{}) (interface{}, error) {
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
	var src interface{}
	json.Unmarshal([]byte(srcJSON), &src)
	output, err := jsonpath.Retrieve(jsonPath, src, config)
	if err != nil {
		fmt.Printf(`%v, %v`, reflect.TypeOf(err), err)
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
	var src interface{}
	json.Unmarshal([]byte(srcJSON), &src)
	output, err := jsonpath.Retrieve(jsonPath, src, config)
	if err != nil {
		fmt.Printf(`%v, %v`, reflect.TypeOf(err), err)
		return
	}
	accessor := output[0].(jsonpath.Accessor)
	srcMap := src.(map[string]interface{})

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
