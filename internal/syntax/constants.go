package syntax

import "reflect"

const (
	msgErrorInvalidSyntaxUnrecognizedInput string = `unrecognized input`
	msgErrorInvalidSyntaxTwoCurrentNode    string = `comparison between two current nodes is prohibited`
	msgErrorInvalidSyntaxFilterValueGroup  string = `JSONPath that returns a value group is prohibited`

	msgTypeNull          string = `null`
	msgTypeObject        string = `object`
	msgTypeArray         string = `array`
	msgTypeObjectOrArray string = `object/array`
)

type emptyEntityIdentifier struct{}
type fullEntityIdentifier struct{}

var emptyEntity = emptyEntityIdentifier{}
var emptyList = []any{emptyEntity}

var fullEntity = fullEntityIdentifier{}
var fullList = []any{fullEntity}

var literalParamTypes = map[reflect.Type]struct{}{
	reflect.TypeFor[syntaxQueryParamLiteral]():      {},
	reflect.TypeFor[syntaxQueryParamRootNode]():     {},
	reflect.TypeFor[syntaxQueryParamRootNodePath](): {},
}

func isLiteralParam(v any) bool {
	t := reflect.TypeOf(v)
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	_, isLiteral := literalParamTypes[t]
	return isLiteral
}
