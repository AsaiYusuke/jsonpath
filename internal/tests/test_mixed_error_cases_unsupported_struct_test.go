package tests

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/AsaiYusuke/jsonpath"
)

type UnsupportedStructChild struct {
	B string
	C int
}

type UnsupportedStructParent struct {
	A UnsupportedStructChild
}

func TestError_UnsupportedStruct(t *testing.T) {
	inputJSON := UnsupportedStructParent{A: UnsupportedStructChild{B: `test`, C: 123}}
	jsonPath := `$.A.B`
	expectedError := createErrorTypeUnmatched(`.A`, `object`, `tests.UnsupportedStructParent`)
	_, err := jsonpath.Retrieve(jsonPath, inputJSON)

	if reflect.TypeOf(expectedError) != reflect.TypeOf(err) ||
		fmt.Sprintf(`%s`, expectedError) != fmt.Sprintf(`%s`, err) {
		t.Errorf("expected error<%s> != actual error<%s>\n",
			expectedError, err)
	}
}
