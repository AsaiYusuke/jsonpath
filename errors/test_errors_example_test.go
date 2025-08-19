package errors_test

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/AsaiYusuke/jsonpath/v2"
	"github.com/AsaiYusuke/jsonpath/v2/config"
	"github.com/AsaiYusuke/jsonpath/v2/errors"
)

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
	config := config.Config{}
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
