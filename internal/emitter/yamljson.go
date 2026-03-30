package emitter

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"

	"github.com/notwillk/nara/internal/config"
	"gopkg.in/yaml.v3"
)

func WriteYAML(path string, objects []map[string]any, cfg *config.Config) error {
	clean := cleanObjects(objects, cfg)
	b, err := yaml.Marshal(clean)
	if err != nil {
		return err
	}
	return writeFile(path, b)
}

func WriteJSON(path string, objects []map[string]any, cfg *config.Config) error {
	clean := cleanObjects(objects, cfg)
	b, err := json.MarshalIndent(clean, "", "  ")
	if err != nil {
		return err
	}
	b = append(b, '\n')
	return writeFile(path, b)
}

func writeFile(path string, b []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o644)
}

func cleanObjects(in []map[string]any, cfg *config.Config) []map[string]any {
	out := make([]map[string]any, 0, len(in))
	for _, m := range in {
		out = append(out, cleanMap(m, cfg))
	}
	return out
}

func cleanMap(in map[string]any, cfg *config.Config) map[string]any {
	include := map[string]bool{}
	for _, k := range cfg.Meta.IncludeKeys {
		include[k] = true
	}
	refKey := cfg.Meta.Ref
	idKey := cfg.Meta.ID
	schemaKey := cfg.Meta.Schema
	if refKey == "" {
		refKey = "$ref"
	}
	if idKey == "" {
		idKey = "$id"
	}
	if schemaKey == "" {
		schemaKey = "$schema"
	}

	keys := make([]string, 0, len(in))
	for k := range in {
		keys = append(keys, k)
	}
	slices.Sort(keys)

	out := map[string]any{}
	for _, k := range keys {
		if (k == refKey || k == idKey || k == schemaKey) && !include[k] {
			continue
		}
		out[k] = cleanValue(in[k], cfg)
	}
	return out
}

func cleanValue(v any, cfg *config.Config) any {
	switch t := v.(type) {
	case map[string]any:
		return cleanMap(t, cfg)
	case []any:
		out := make([]any, len(t))
		for i := range t {
			out[i] = cleanValue(t[i], cfg)
		}
		return out
	default:
		return t
	}
}

