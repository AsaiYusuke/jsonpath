package main

import (
	"encoding/json"
	"fmt"

	"github.com/AsaiYusuke/jsonpath"
)

func main() {
	jsonPath, srcJSON := `$.key`, `{"key":"value"}`
	var src interface{}
	json.Unmarshal([]byte(srcJSON), &src)
	output, _ := jsonpath.Retrieve(jsonPath, src)
	outputJSON, _ := json.Marshal(output)
	fmt.Println(string(outputJSON))
	// Output:
	// ["value"]
}
