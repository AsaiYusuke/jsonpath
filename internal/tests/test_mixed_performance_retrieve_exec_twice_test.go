package tests

import (
	"encoding/json"
	"testing"

	"github.com/AsaiYusuke/jsonpath/v2"
)

func TestRetrieveExecTwice(t *testing.T) {
	jsonpath1 := `$.a`
	srcJSON1 := `{"a":123}`
	expectedOutput1 := "[123]"
	jsonpath2 := `$[1].b`
	srcJSON2 := `[123,{"b":456}]`
	expectedOutput2 := "[456]"

	var src1 any
	if err := json.Unmarshal([]byte(srcJSON1), &src1); err != nil {
		t.Errorf("%s", err)
		return
	}
	var src2 any
	if err := json.Unmarshal([]byte(srcJSON2), &src2); err != nil {
		t.Errorf("%s", err)
		return
	}

	actualObject1, err := jsonpath.Retrieve(jsonpath1, src1)
	if err != nil {
		t.Errorf("expected error<nil> != actual error<%s>\n", err)
		return
	}
	actualObject2, err := jsonpath.Retrieve(jsonpath2, src2)
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
