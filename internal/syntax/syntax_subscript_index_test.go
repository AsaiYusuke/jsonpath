package syntax

import "testing"

func TestSyntaxIndexSubscript_getIndexes_DummyAlwaysNil(t *testing.T) {
	s := &syntaxIndexSubscript{number: 0}
	got := s.getIndexes(1)
	if got != nil {
		t.Fatalf("expected nil, got %#v", got)
	}
}
