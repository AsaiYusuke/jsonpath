package syntax

import "testing"

func TestSyntaxIndexSubscript_forEachIndex_DoesNotCallCallback(t *testing.T) {
	s := &syntaxIndexSubscript{number: 0}
	called := false
	s.forEachIndex(1, func(index int) { called = true })
	if called {
		t.Fatalf("callback should not be called for dummy getIndexes")
	}
}
