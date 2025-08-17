package syntax

import "testing"

func TestSyntaxQueryParamRootNode_isValueGroupParameter_Coverage(t *testing.T) {
	s := &syntaxQueryParamRootNode{}
	s.isValueGroupParameter()
}
