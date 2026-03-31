package emitter

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	"github.com/notwillk/nara/internal/config"
	"github.com/notwillk/nara/internal/entity"
	"github.com/notwillk/nara/internal/graph"
	_ "modernc.org/sqlite"
)

func TestWriteSQLite(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "graph.db")
	g := &graph.Graph{
		Entities: map[string]*entity.Entity{
			"/tmp/hello": {ID: "hello", Schema: "note"},
			"/tmp/world": {ID: "world", Schema: "note"},
		},
		Expanded: map[string]map[string]any{
			"/tmp/hello": {"$id": "hello", "$schema": "note", "title": "Hello"},
			"/tmp/world": {"$id": "world", "$schema": "note", "title": "World"},
		},
		Edges: []graph.Edge{{FromID: "hello", Field: "$.friend", ToID: "world"}},
	}
	cfg := &config.Config{Meta: config.Meta{ID: "$id", Schema: "$schema"}}

	if err := WriteSQLite(path, g, cfg); err != nil {
		t.Fatal(err)
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM entities`).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Fatalf("expected 2 entities, got %d", count)
	}

	var schemaName string
	var jsonPayload string
	if err := db.QueryRow(`SELECT schema, json FROM entities WHERE id = 'hello'`).Scan(&schemaName, &jsonPayload); err != nil {
		t.Fatal(err)
	}
	if schemaName != "note" {
		t.Fatalf("unexpected schema: %q", schemaName)
	}
	if jsonPayload != `{"title":"Hello"}` {
		t.Fatalf("unexpected payload: %s", jsonPayload)
	}

	if err := db.QueryRow(`SELECT COUNT(*) FROM edges`).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("expected 1 edge, got %d", count)
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatal(err)
	}
}
