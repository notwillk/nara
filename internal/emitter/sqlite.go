package emitter

import (
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"
	"slices"

	"github.com/notwillk/nara/internal/config"
	"github.com/notwillk/nara/internal/graph"
	_ "modernc.org/sqlite"
)

func WriteSQLite(path string, g *graph.Graph, cfg *config.Config) error {
	if g == nil {
		return os.ErrInvalid
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(`
		CREATE TABLE entities (
			id TEXT PRIMARY KEY,
			schema TEXT NOT NULL,
			json TEXT NOT NULL
		);
		CREATE TABLE edges (
			from_id TEXT NOT NULL,
			field TEXT NOT NULL,
			to_id TEXT NOT NULL
		);
	`); err != nil {
		return err
	}

	paths := make([]string, 0, len(g.Entities))
	for path := range g.Entities {
		paths = append(paths, path)
	}
	slices.Sort(paths)

	for _, path := range paths {
		ent := g.Entities[path]
		payload := cleanMap(g.Expanded[path], cfg)
		b, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		if _, err := tx.Exec(`INSERT INTO entities (id, schema, json) VALUES (?, ?, ?)`, ent.ID, ent.Schema, string(b)); err != nil {
			return err
		}
	}

	edges := append([]graph.Edge(nil), g.Edges...)
	slices.SortFunc(edges, func(a graph.Edge, b graph.Edge) int {
		if a.FromID != b.FromID {
			if a.FromID < b.FromID {
				return -1
			}
			return 1
		}
		if a.Field != b.Field {
			if a.Field < b.Field {
				return -1
			}
			return 1
		}
		if a.ToID < b.ToID {
			return -1
		}
		if a.ToID > b.ToID {
			return 1
		}
		return 0
	})
	for _, edge := range edges {
		if _, err := tx.Exec(`INSERT INTO edges (from_id, field, to_id) VALUES (?, ?, ?)`, edge.FromID, edge.Field, edge.ToID); err != nil {
			return err
		}
	}

	return tx.Commit()
}
