package entity

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/notwillk/nara/internal/config"
	"github.com/notwillk/nara/internal/errors"
	"github.com/notwillk/nara/internal/naming"
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

	metaID := cfg.Meta.ID
	metaSchema := cfg.Meta.Schema
	if metaID == "" {
		metaID = "$id"
	}
	if metaSchema == "" {
		metaSchema = "$schema"
	}

	rawID := stringValue(data[metaID])
	rawSchema := stringValue(data[metaSchema])
	inferredID, inferredSchema, _ := naming.InferFromPath(path, cfg.Resolution.FilenamePattern)

	id := rawID
	if id == "" {
		id = inferredID
	}
	schemaName := rawSchema
	if schemaName == "" {
		schemaName = inferredSchema
	}
	if id == "" || schemaName == "" {
		return nil, errors.New(errors.CategoryValidation, path, "unable to infer id/schema")
	}
	if inferredID != "" && rawID != "" && rawID != inferredID {
		return nil, errors.New(errors.CategoryValidation, path, "entity $id does not match filename id")
	}
	if inferredSchema != "" && rawSchema != "" && schemaName != inferredSchema {
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
	id, schema, _ = naming.InferFromPath(path, naming.TokenID+"."+naming.TokenSchema)
	return id, schema
}

func stringValue(v any) string {
	s, _ := v.(string)
	return s
}
