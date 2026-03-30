package entity

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/notwillk/nara/internal/config"
	"github.com/notwillk/nara/internal/errors"
	"gopkg.in/yaml.v3"
)

type Entity struct {
	ID     string
	Schema string
	Path   string
	Data   map[string]any
}

// Load reads a YAML/JSON entity file and extracts id/schema metadata.
func Load(path string, cfg *config.Config) (*Entity, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.New(errors.CategoryResolution, path, "unable to read entity file")
	}

	data := map[string]any{}
	switch strings.ToLower(filepath.Ext(path)) {
	case ".json":
		if err := json.Unmarshal(b, &data); err != nil {
			return nil, errors.New(errors.CategoryValidation, path, "invalid JSON")
		}
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(b, &data); err != nil {
			return nil, errors.New(errors.CategoryValidation, path, "invalid YAML")
		}
	default:
		return nil, errors.New(errors.CategoryValidation, path, "unsupported entity extension")
	}

	inferredID, inferredSchema := InferFromFilename(path)
	metaID := cfg.Meta.ID
	metaSchema := cfg.Meta.Schema
	if metaID == "" {
		metaID = "$id"
	}
	if metaSchema == "" {
		metaSchema = "$schema"
	}

	id := stringValue(data[metaID])
	schemaName := stringValue(data[metaSchema])
	if id == "" {
		id = inferredID
	}
	if schemaName == "" {
		schemaName = inferredSchema
	}
	if id == "" || schemaName == "" {
		return nil, errors.New(errors.CategoryValidation, path, "unable to infer id/schema")
	}

	// If $schema exists, it must match inferred schema when inferable.
	if inferredSchema != "" && stringValue(data[metaSchema]) != "" && schemaName != inferredSchema {
		return nil, errors.New(errors.CategoryValidation, path, "entity $schema does not match filename schema")
	}

	return &Entity{
		ID:     id,
		Schema: schemaName,
		Path:   path,
		Data:   data,
	}, nil
}

func InferFromFilename(path string) (id string, schema string) {
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)
	parts := strings.Split(name, ".")
	if len(parts) < 2 {
		return "", ""
	}
	return strings.Join(parts[:len(parts)-1], "."), parts[len(parts)-1]
}

func stringValue(v any) string {
	s, _ := v.(string)
	return s
}

