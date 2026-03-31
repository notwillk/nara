package compiler

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestFixturesCompile(t *testing.T) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("unable to locate test file")
	}

	repoRoot := filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
	fixturesRoot := filepath.Join(repoRoot, "fixtures")

	entries, err := os.ReadDir(fixturesRoot)
	if err != nil {
		t.Fatal(err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		t.Run(entry.Name(), func(t *testing.T) {
			configPath := filepath.Join(fixturesRoot, entry.Name(), "nara.yaml")
			entryGlob := filepath.Join(fixturesRoot, entry.Name(), "entries", "*.*")

			result, err := Compile(configPath, []string{entryGlob})
			if err != nil {
				t.Fatalf("compile failed: %v", err)
			}
			if len(result.Entries) == 0 {
				t.Fatal("expected compiled entries")
			}
		})
	}
}
