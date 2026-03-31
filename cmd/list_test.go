package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestListEntriesOutputsRelativeRows(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "nara.yaml"), `version: 1
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
`)
	writeFile(t, filepath.Join(dir, "schemas", "note.cue"), "title: string\n")
	writeFile(t, filepath.Join(dir, "examples", "hello.note.yaml"), "$id: hello\n$schema: note\ntitle: Hello\n")

	oldConfigPath := configPath
	configPath = filepath.Join(dir, "nara.yaml")
	defer func() { configPath = oldConfigPath }()

	buf := &bytes.Buffer{}
	command := newListCmd()
	command.SetOut(buf)
	command.SetErr(buf)
	command.SetArgs([]string{"entries"})
	if err := command.Execute(); err != nil {
		t.Fatal(err)
	}

	want := "hello\tnote\texamples/hello.note.yaml\n"
	if buf.String() != want {
		t.Fatalf("unexpected output\nwant: %q\ngot:  %q", want, buf.String())
	}
}

func writeFile(t *testing.T, path string, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}
