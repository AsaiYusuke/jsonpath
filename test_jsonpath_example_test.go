package jsonpath_test

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/AsaiYusuke/jsonpath/v2"
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

func ExampleParse_reuseBuffer() {
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
	buf := make([]any, 0, 4)
	output1, err := jsonPathParser(src1, &buf)
	if err != nil {
		fmt.Printf(`type: %v, value: %v`, reflect.TypeOf(err), err)
		return
	}
	output2, err := jsonPathParser(src2, &buf)
	if err != nil {
		fmt.Printf(`type: %v, value: %v`, reflect.TypeOf(err), err)
		return
	}
	outputJSON1, _ := json.Marshal(output1)
	outputJSON2, _ := json.Marshal(output2)
	bufJSON, _ := json.Marshal(buf)
	fmt.Println(string(outputJSON1))
	fmt.Println(string(outputJSON2))
	fmt.Println(string(bufJSON))
	// Output:
	// ["value2"]
	// ["value2"]
	// ["value2"]
}
