package configschema

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/notwillk/nara/internal/config"
)

func TestFormatGeneratesSchemaAndNormalizesPath(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "nara.yaml")
	body := `version: 1
$schema: wrong/path.json
updateSchemaOnFormat: true
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
	mustWriteFile(t, filepath.Join(dir, "schemas", "note.cue"), "title: string\n")
	mustWriteFile(t, configPath, body)

	if err := Format(configPath); err != nil {
		t.Fatal(err)
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.SchemaPath != GeneratedSchemaRelPath() {
		t.Fatalf("unexpected schema path: %q", cfg.SchemaPath)
	}

	schemaBytes, err := os.ReadFile(filepath.Join(dir, ".schemas", GeneratedFile))
	if err != nil {
		t.Fatal(err)
	}
	var doc map[string]any
	if err := json.Unmarshal(schemaBytes, &doc); err != nil {
		t.Fatal(err)
	}

	configBytes, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(configBytes), GeneratedSchemaRelPath()) {
		t.Fatalf("config did not contain normalized schema path: %s", string(configBytes))
	}
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
