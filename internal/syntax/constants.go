package syntax

const (
	msgErrorInvalidSyntaxUnrecognizedInput string = `unrecognized input`
	msgErrorInvalidSyntaxTwoCurrentNode    string = `comparison between two current nodes is prohibited`
	msgErrorInvalidSyntaxFilterValueGroup  string = `JSONPath that returns a value group is prohibited`

	msgTypeNull          string = `null`
	msgTypeObject        string = `object`
	msgTypeArray         string = `array`
	msgTypeObjectOrArray string = `object/array`
)

var emptyEntity = struct{}{}
var emptyList = []interface{}{emptyEntity}
var fullList = []interface{}{true}
