package errors

import "testing"

func TestErrorFormatting(t *testing.T) {
	err := Wrap(CategoryConfig, "nara.yaml", 12, "paths./", "invalid path")
	got := err.Error()
	want := "config error: nara.yaml:12:paths./: invalid path"
	if got != want {
		t.Fatalf("unexpected error format\nwant: %s\ngot:  %s", want, got)
	}
}

