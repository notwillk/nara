package workspace

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/notwillk/nara/internal/config"
)

func TestDiscoverEntriesFiltersToKnownSchemas(t *testing.T) {
	dir := t.TempDir()
	mustWriteFile(t, filepath.Join(dir, "examples", "hello.note.yaml"), "$id: hello\n$schema: note\ntitle: hello\n")
	mustWriteFile(t, filepath.Join(dir, ".schemas", "nara.config.schema.json"), "{}\n")
	mustWriteFile(t, filepath.Join(dir, "test.checksy.yaml"), "rules: []\n")

	cfg := testConfig()
	entries, err := DiscoverEntries(cfg, dir, map[string]bool{"note": true})
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].ID != "hello" || entries[0].Schema != "note" {
		t.Fatalf("unexpected entry: %+v", entries[0])
	}
}

func TestEntriesForGlobsRejectUnknownSchema(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "examples", "hello.note.yaml")
	mustWriteFile(t, path, "$id: hello\n$schema: note\ntitle: hello\n")

	cfg := testConfig()
	_, err := EntriesForGlobs([]string{path}, cfg, map[string]bool{"pizza": true})
	if err == nil {
		t.Fatal("expected unknown schema to be rejected")
	}
}

func testConfig() *config.Config {
	return &config.Config{
		Paths: map[string]string{"/": "."},
		Resolution: config.Resolution{
			FilenamePattern: "<id>.<schema>",
			Extensions:      []string{".yaml", ".json"},
		},
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
