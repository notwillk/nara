package graph

import (
	"fmt"
	"path/filepath"
	"slices"
	"strings"

	"github.com/notwillk/nara/internal/config"
	"github.com/notwillk/nara/internal/entity"
	"github.com/notwillk/nara/internal/errors"
	"github.com/notwillk/nara/internal/resolver"
)

type Edge struct {
	FromID string
	Field  string
	ToID   string
}

type Graph struct {
	Entities map[string]*entity.Entity       // key: absolute file path
	Expanded map[string]map[string]any       // key: absolute file path
	Edges    []Edge
	Entries  []string                        // absolute entry paths
}

// Build recursively resolves refs and builds a DAG.
func Build(entryPaths []string, cfg *config.Config, baseDir string) (*Graph, error) {
	g := &Graph{
		Entities: map[string]*entity.Entity{},
		Expanded: map[string]map[string]any{},
		Edges:    []Edge{},
		Entries:  []string{},
	}

	inStack := map[string]bool{}
	for _, p := range entryPaths {
		abs, _ := filepath.Abs(p)
		if err := visit(abs, g, cfg, baseDir, inStack); err != nil {
			return nil, err
		}
		g.Entries = append(g.Entries, abs)
	}
	return g, nil
}

func visit(path string, g *Graph, cfg *config.Config, baseDir string, inStack map[string]bool) error {
	if inStack[path] {
		return errors.New(errors.CategoryResolution, path, "illegal cycle detected")
	}
	if _, ok := g.Entities[path]; ok {
		return nil
	}
	inStack[path] = true
	defer delete(inStack, path)

	ent, err := entity.Load(path, cfg)
	if err != nil {
		return err
	}
	g.Entities[path] = ent

	expanded, edges, err := expandValue(ent, ent.Data, "$", cfg, baseDir, g, inStack)
	if err != nil {
		return err
	}
	rootMap, ok := expanded.(map[string]any)
	if !ok {
		return errors.New(errors.CategoryValidation, path, "entity root must be object")
	}
	g.Expanded[path] = rootMap
	g.Edges = append(g.Edges, edges...)
	return nil
}

func expandValue(owner *entity.Entity, value any, fieldPath string, cfg *config.Config, baseDir string, g *Graph, inStack map[string]bool) (any, []Edge, error) {
	metaRef := cfg.Meta.Ref
	if metaRef == "" {
		metaRef = "$ref"
	}

	switch t := value.(type) {
	case map[string]any:
		if refRaw, ok := t[metaRef]; ok {
			ref, ok := refRaw.(string)
			if !ok {
				return nil, nil, errors.New(errors.CategoryResolution, owner.Path, "ref must be string")
			}

			targetPath, err := resolver.Resolve(owner.Path, owner.Schema, ref, cfg, baseDir)
			if err != nil {
				return nil, nil, err
			}
			targetAbs, _ := filepath.Abs(targetPath)
			if err := visit(targetAbs, g, cfg, baseDir, inStack); err != nil {
				return nil, nil, err
			}
			targetEnt := g.Entities[targetAbs]
			targetExpanded := deepCopyMap(g.Expanded[targetAbs])

			// Shallow override (excluding $ref key).
			override := map[string]any{}
			for k, v := range t {
				if k == metaRef {
					continue
				}
				override[k] = v
			}
			merged := shallowMerge(targetExpanded, override)

			edges := []Edge{{FromID: owner.ID, Field: fieldPath, ToID: targetEnt.ID}}
			return merged, edges, nil
		}

		out := map[string]any{}
		edges := []Edge{}
		keys := make([]string, 0, len(t))
		for k := range t {
			keys = append(keys, k)
		}
		slices.Sort(keys)
		for _, k := range keys {
			v := t[k]
			childPath := fmt.Sprintf("%s.%s", fieldPath, k)
			ev, ee, err := expandValue(owner, v, childPath, cfg, baseDir, g, inStack)
			if err != nil {
				return nil, nil, err
			}
			out[k] = ev
			edges = append(edges, ee...)
		}
		return out, edges, nil
	case []any:
		out := make([]any, 0, len(t))
		edges := []Edge{}
		for i, v := range t {
			childPath := fmt.Sprintf("%s[%d]", fieldPath, i)
			ev, ee, err := expandValue(owner, v, childPath, cfg, baseDir, g, inStack)
			if err != nil {
				return nil, nil, err
			}
			out = append(out, ev)
			edges = append(edges, ee...)
		}
		return out, edges, nil
	default:
		return t, nil, nil
	}
}

func shallowMerge(base map[string]any, override map[string]any) map[string]any {
	out := deepCopyMap(base)
	for k, v := range override {
		out[k] = v
	}
	return out
}

func deepCopyMap(in map[string]any) map[string]any {
	out := make(map[string]any, len(in))
	for k, v := range in {
		out[k] = deepCopyAny(v)
	}
	return out
}

func deepCopyAny(v any) any {
	switch t := v.(type) {
	case map[string]any:
		return deepCopyMap(t)
	case []any:
		out := make([]any, len(t))
		for i := range t {
			out[i] = deepCopyAny(t[i])
		}
		return out
	default:
		return t
	}
}

func EntryObjects(g *Graph) []map[string]any {
	out := make([]map[string]any, 0, len(g.Entries))
	for _, p := range g.Entries {
		if m, ok := g.Expanded[p]; ok {
			out = append(out, deepCopyMap(stripPathKeys(m)))
		}
	}
	return out
}

func stripPathKeys(m map[string]any) map[string]any {
	out := map[string]any{}
	for k, v := range m {
		out[k] = v
	}
	return out
}

func NormalizeEntryPaths(globs []string) ([]string, error) {
	seen := map[string]bool{}
	out := []string{}
	for _, g := range globs {
		matches, err := filepath.Glob(g)
		if err != nil {
			return nil, err
		}
		for _, m := range matches {
			ext := strings.ToLower(filepath.Ext(m))
			if ext != ".yaml" && ext != ".yml" && ext != ".json" {
				continue
			}
			abs, _ := filepath.Abs(m)
			if !seen[abs] {
				seen[abs] = true
				out = append(out, abs)
			}
		}
	}
	slices.Sort(out)
	return out, nil
}

