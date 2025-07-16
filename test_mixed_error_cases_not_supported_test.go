package jsonpath

import "testing"

func TestRetrieve_notSupported_script_command(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[(command)]`,
		inputJSON:   `{}`,
		expectedErr: ErrorNotSupported{feature: `script`, path: `[(command)]`},
	}
	runTestCase(t, testCase, "TestRetrieve_notSupported_script_command")
}

func TestRetrieve_notSupported_script_command_with_spaces(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[( command )]`,
		inputJSON:   `{}`,
		expectedErr: ErrorNotSupported{feature: `script`, path: `[(command)]`},
	}
	runTestCase(t, testCase, "TestRetrieve_notSupported_script_command_with_spaces")
}
