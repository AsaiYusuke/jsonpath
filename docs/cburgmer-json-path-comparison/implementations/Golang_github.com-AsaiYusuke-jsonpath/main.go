package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/AsaiYusuke/jsonpath/v2"
	"github.com/AsaiYusuke/jsonpath/v2/errors"
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
	err = json.Unmarshal([]byte(data), &jsonData)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	result, err := jsonpath.Retrieve(selector, jsonData)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		switch err.(type) {
		case errors.ErrorInvalidArgument,
			errors.ErrorInvalidSyntax,
			errors.ErrorNotSupported,
			errors.ErrorFunctionNotFound:
			os.Exit(2)
		case errors.ErrorMemberNotExist,
			errors.ErrorTypeUnmatched,
			errors.ErrorFunctionFailed:
			os.Exit(3)
		}
		os.Exit(1)
	}
	jsonResult, err := json.Marshal(result)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println(string(jsonResult))
}
