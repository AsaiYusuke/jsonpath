package jsonpath

import (
	"fmt"
	"reflect"
	"testing"
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
	jsonpath := `$.A.B`
	expectedError := createErrorTypeUnmatched(`.A`, `object`, `jsonpath.UnsupportedStructParent`)
	_, err := Retrieve(jsonpath, inputJSON)

	if reflect.TypeOf(expectedError) != reflect.TypeOf(err) ||
		fmt.Sprintf(`%s`, expectedError) != fmt.Sprintf(`%s`, err) {
		t.Errorf("expected error<%s> != actual error<%s>\n",
			expectedError, err)
	}
}
