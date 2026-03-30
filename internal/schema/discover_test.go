package schema

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/notwillk/nara/internal/config"
)

func TestDuplicateSchemaNames(t *testing.T) {
	dir := t.TempDir()
	a := filepath.Join(dir, "a")
	b := filepath.Join(dir, "b")
	if err := os.MkdirAll(a, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(b, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(a, "pizza.cue"), []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(b, "pizza.cue"), []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg := &config.Config{
		Schemas: config.Schemas{
			Sources: []string{"./a/*.cue", "./b/*.cue"},
		},
	}
	_, err := Discover(cfg, dir)
	if err == nil {
		t.Fatal("expected duplicate schema name error")
	}
}

