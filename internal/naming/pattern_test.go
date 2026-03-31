package naming

import "testing"

func TestInferFromPath(t *testing.T) {
	id, schema, ok := InferFromPath("examples/hello.world.note.yaml", TokenID+"."+TokenSchema)
	if !ok {
		t.Fatal("expected filename match")
	}
	if id != "hello.world" || schema != "note" {
		t.Fatalf("unexpected inference: id=%q schema=%q", id, schema)
	}
}

func TestInferFromPathSchemaFirstPattern(t *testing.T) {
	id, schema, ok := InferFromPath("examples/note--alpha.beta.yaml", TokenSchema+"--"+TokenID)
	if !ok {
		t.Fatal("expected filename match")
	}
	if id != "alpha.beta" || schema != "note" {
		t.Fatalf("unexpected inference: id=%q schema=%q", id, schema)
	}
}

func TestValidateFilenamePattern(t *testing.T) {
	cases := []string{
		"",
		TokenID + TokenSchema,
		TokenID + "/" + TokenSchema,
		TokenID + `\\` + TokenSchema,
	}
	for _, pattern := range cases {
		if err := ValidateFilenamePattern(pattern); err == nil {
			t.Fatalf("expected pattern %q to fail", pattern)
		}
	}
}
