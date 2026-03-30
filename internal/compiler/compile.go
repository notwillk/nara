package compiler

import (
	"path/filepath"

	"github.com/notwillk/nara/internal/config"
	"github.com/notwillk/nara/internal/cuevalidate"
	"github.com/notwillk/nara/internal/errors"
	"github.com/notwillk/nara/internal/graph"
	"github.com/notwillk/nara/internal/schema"
)

type Result struct {
	Graph   *graph.Graph
	Entries []map[string]any
}

func Compile(configPath string, entryGlobs []string) (*Result, error) {
	cfg, err := config.Load(configPath)
	if err != nil {
		return nil, err
	}
	baseDir := filepath.Dir(configPath)

	discovered, err := schema.Discover(cfg, baseDir)
	if err != nil {
		return nil, err
	}

	validator, err := cuevalidate.New(discovered)
	if err != nil {
		return nil, err
	}

	entryPaths, err := graph.NormalizeEntryPaths(entryGlobs)
	if err != nil {
		return nil, errors.New(errors.CategoryConfig, configPath, "invalid entry glob")
	}
	if len(entryPaths) == 0 {
		return nil, errors.New(errors.CategoryResolution, configPath, "no entry files matched")
	}

	g, err := graph.Build(entryPaths, cfg, baseDir)
	if err != nil {
		return nil, err
	}

	for p, ent := range g.Entities {
		payload := g.Expanded[p]
		if err := validator.ValidateEntity(ent.Schema, payload, p); err != nil {
			return nil, err
		}
	}

	return &Result{
		Graph:   g,
		Entries: graph.EntryObjects(g),
	}, nil
}

