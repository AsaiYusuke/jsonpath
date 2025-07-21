package tests

import "testing"

func TestRetrieve_notSupported_script_command(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[(command)]`,
		inputJSON:   `{}`,
		expectedErr: createErrorNotSupported(`script`, `[(command)]`),
	}
	runTestCase(t, testCase, "TestRetrieve_notSupported_script_command")
}

func TestRetrieve_notSupported_script_command_with_spaces(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[( command )]`,
		inputJSON:   `{}`,
		expectedErr: createErrorNotSupported(`script`, `[(command)]`),
	}
	runTestCase(t, testCase, "TestRetrieve_notSupported_script_command_with_spaces")
}
