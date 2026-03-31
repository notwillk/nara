package lint

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunValidWorkspace(t *testing.T) {
	dir := t.TempDir()
	writeWorkspaceConfig(t, dir)
	mustWriteFile(t, filepath.Join(dir, "schemas", "note.cue"), "name: string\nfriend?: _\n")
	mustWriteFile(t, filepath.Join(dir, "examples", "friend.note.yaml"), "$id: friend\n$schema: note\nname: Friend\n")
	mustWriteFile(t, filepath.Join(dir, "examples", "hello.note.yaml"), "$id: hello\n$schema: note\nname: Hello\nfriend:\n  $ref: friend\n")

	if err := Run(filepath.Join(dir, "nara.yaml")); err != nil {
		t.Fatal(err)
	}
}

func TestRunMissingRefFails(t *testing.T) {
	dir := t.TempDir()
	writeWorkspaceConfig(t, dir)
	mustWriteFile(t, filepath.Join(dir, "schemas", "note.cue"), "name: string\nfriend?: _\n")
	mustWriteFile(t, filepath.Join(dir, "examples", "hello.note.yaml"), "$id: hello\n$schema: note\nname: Hello\nfriend:\n  $ref: missing\n")

	if err := Run(filepath.Join(dir, "nara.yaml")); err == nil {
		t.Fatal("expected lint to fail for missing ref")
	}
}

func writeWorkspaceConfig(t *testing.T, dir string) {
	t.Helper()
	configBody := `version: 1
paths:
  "/": "."
meta:
  ref: "$ref"
  id: "$id"
  schema: "$schema"
  includeKeys: []
schemas:
  sources:
    - schemas/*.cue
  inferFromFilename: true
resolution:
  filenamePattern: <id>.<schema>
  extensions:
    - .yaml
`
	mustWriteFile(t, filepath.Join(dir, "nara.yaml"), configBody)
}

func mustWriteFile(t *testing.T, path string, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}
