package tests

import (
	"encoding/json"
	"testing"

	"github.com/AsaiYusuke/jsonpath"
)

func TestParserFuncExecTwice(t *testing.T) {
	jsonPath := `$.a`
	srcJSON1 := `{"a":1}`
	expectedOutput1 := "[1]"
	srcJSON2 := `{"a":2}`
	expectedOutput2 := "[2]"

	var src1 interface{}
	if err := json.Unmarshal([]byte(srcJSON1), &src1); err != nil {
		t.Errorf("%s", err)
		return
	}
	var src2 interface{}
	if err := json.Unmarshal([]byte(srcJSON2), &src2); err != nil {
		t.Errorf("%s", err)
		return
	}

	parserFunc, err := jsonpath.Parse(jsonPath)
	if err != nil {
		t.Errorf("expected error<nil> != actual error<%s>\n", err)
		return
	}

	actualObject1, err := parserFunc(src1)
	if err != nil {
		t.Errorf("expected error<nil> != actual error<%s>\n", err)
		return
	}
	actualObject2, err := parserFunc(src2)
	if err != nil {
		t.Errorf("expected error<nil> != actual error<%s>\n", err)
		return
	}

	actualOutputJSON1, err := json.Marshal(actualObject1)
	if err != nil {
		t.Errorf("%s", err)
		return
	}
	actualOutputJSON2, err := json.Marshal(actualObject2)
	if err != nil {
		t.Errorf("%s", err)
		return
	}

	if string(actualOutputJSON1) != string(expectedOutput1) || string(actualOutputJSON2) != string(expectedOutput2) {
		t.Errorf("actualOutputJSON1<%s> != expectedOutput1<%s> || actualOutputJSON2<%s> != expectedOutput2<%s>\n",
			string(actualOutputJSON1), string(expectedOutput1), string(actualOutputJSON2), string(expectedOutput2))
		return
	}
}
