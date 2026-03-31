package scaffold

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/notwillk/nara/internal/config"
)

func TestInitCreatesProjectScaffold(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "nara.yaml")
	if err := Init(configPath); err != nil {
		t.Fatal(err)
	}

	paths := []string{
		configPath,
		filepath.Join(dir, "schemas", "note.cue"),
		filepath.Join(dir, "examples", "hello.note.yaml"),
		filepath.Join(dir, ".schemas", "nara.config.schema.json"),
	}
	for _, path := range paths {
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("expected %s to exist: %v", path, err)
		}
	}

	if _, err := config.Load(configPath); err != nil {
		t.Fatal(err)
	}
}
