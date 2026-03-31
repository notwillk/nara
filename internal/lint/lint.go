package lint

import (
	"path/filepath"

	"github.com/notwillk/nara/internal/config"
	"github.com/notwillk/nara/internal/cuevalidate"
	"github.com/notwillk/nara/internal/graph"
	"github.com/notwillk/nara/internal/schema"
	"github.com/notwillk/nara/internal/workspace"
)

func Run(configPath string) error {
	cfg, err := config.Load(configPath)
	if err != nil {
		return err
	}

	baseDir := filepath.Dir(configPath)
	discovered, err := schema.Discover(cfg, baseDir)
	if err != nil {
		return err
	}
	validator, err := cuevalidate.New(discovered)
	if err != nil {
		return err
	}

	allowedSchemas := schemaNames(discovered)
	entries, err := workspace.DiscoverEntries(cfg, baseDir, allowedSchemas)
	if err != nil {
		return err
	}
	if len(entries) == 0 {
		return nil
	}

	paths := make([]string, 0, len(entries))
	for _, entry := range entries {
		paths = append(paths, entry.Path)
	}

	g, err := graph.Build(paths, cfg, baseDir)
	if err != nil {
		return err
	}

	for path, ent := range g.Entities {
		if err := validator.ValidateEntity(ent.Schema, g.Expanded[path], path); err != nil {
			return err
		}
	}
	return nil
}

func schemaNames(discovered []schema.Schema) map[string]bool {
	out := make(map[string]bool, len(discovered))
	for _, item := range discovered {
		out[item.Name] = true
	}
	return out
}
