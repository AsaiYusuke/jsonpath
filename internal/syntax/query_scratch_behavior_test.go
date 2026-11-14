package syntax

import "testing"

func TestFilterScratchDoesNotMutateSource(t *testing.T) {
	src := []any{
		map[string]any{"a": float64(1)},
		map[string]any{"a": float64(2)},
	}

	parserFunc, err := Parse(`$[?(@.a==1)]`)
	if err != nil {
		t.Fatal(err)
	}

	buf := make([]any, 0, len(src))
	args := []*[]any{&buf}
	for range 2 {
		got, err := parserFunc(src, args...)
		if err != nil {
			t.Fatal(err)
		}
		if len(got) != 1 || got[0].(map[string]any)["a"] != float64(1) {
			t.Fatalf("got %#v, want the first source element", got)
		}
	}

	if src[0].(map[string]any)["a"] != float64(1) || src[1].(map[string]any)["a"] != float64(2) {
		t.Fatalf("filter modified source: %#v", src)
	}
}
