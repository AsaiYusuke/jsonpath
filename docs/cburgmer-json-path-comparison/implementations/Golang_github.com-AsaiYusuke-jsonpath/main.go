package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/AsaiYusuke/jsonpath"
)

func main() {
	defer func() {
		e := recover()
		if e != nil {
			fmt.Fprintln(os.Stderr, e)
			os.Exit(1)
		}
	}()

	selector := os.Args[1]

	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	var jsonData interface{}
	json.Unmarshal([]byte(data), &jsonData)

	result, err := jsonpath.Retrieve(selector, jsonData)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		switch err.(type) {
		case jsonpath.ErrorInvalidArgument:
			os.Exit(2)
		case jsonpath.ErrorInvalidSyntax:
			os.Exit(2)
		case jsonpath.ErrorNotSupported:
			os.Exit(2)
		case jsonpath.ErrorMemberNotExist:
			os.Exit(3)
		case jsonpath.ErrorIndexOutOfRange:
			os.Exit(3)
		case jsonpath.ErrorTypeUnmatched:
			os.Exit(3)
		case jsonpath.ErrorNoneMatched:
			os.Exit(3)
		}
		os.Exit(1)
	}
	jsonResult, err := json.Marshal(result)
	fmt.Println(string(jsonResult))
}
