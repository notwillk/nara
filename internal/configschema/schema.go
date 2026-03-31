package configschema

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/notwillk/nara/internal/config"
	"github.com/notwillk/nara/internal/errors"
	"gopkg.in/yaml.v3"
)

const (
	GeneratedDir  = ".schemas"
	GeneratedFile = "nara.config.schema.json"
)

func GeneratedSchemaRelPath() string {
	return filepath.ToSlash(filepath.Join(GeneratedDir, GeneratedFile))
}

func Format(configPath string) error {
	cfg, err := config.LoadRaw(configPath)
	if err != nil {
		return err
	}

	baseDir := filepath.Dir(configPath)
	if !cfg.UpdateSchemaOnFormat {
		return config.Validate(cfg, configPath, baseDir)
	}

	cfg.SchemaPath = GeneratedSchemaRelPath()
	schemaPath := filepath.Join(baseDir, filepath.FromSlash(cfg.SchemaPath))
	if err := writeFile(schemaPath, Generate()); err != nil {
		return err
	}
	if err := config.Validate(cfg, configPath, baseDir); err != nil {
		return err
	}

	b, err := yaml.Marshal(cfg)
	if err != nil {
		return errors.New(errors.CategoryConfig, configPath, "unable to render config YAML")
	}
	return writeFile(configPath, b)
}

func Generate() []byte {
	doc := map[string]any{
		"$schema": "https://json-schema.org/draft/2020-12/schema",
		"$id":     GeneratedFile,
		"title":   "nara config",
		"type":    "object",
		"required": []string{
			"version",
			"paths",
			"meta",
			"schemas",
			"resolution",
		},
		"additionalProperties": false,
		"properties": map[string]any{
			"version": map[string]any{
				"type":  "integer",
				"const": 1,
			},
			"$schema": map[string]any{
				"type": "string",
			},
			"updateSchemaOnFormat": map[string]any{
				"type": "boolean",
			},
			"paths": map[string]any{
				"type":          "object",
				"minProperties": 1,
				"additionalProperties": map[string]any{
					"type": "string",
				},
			},
			"meta": map[string]any{
				"type":                 "object",
				"additionalProperties": false,
				"required":             []string{"ref", "id", "schema"},
				"properties": map[string]any{
					"ref": map[string]any{"type": "string"},
					"id":  map[string]any{"type": "string"},
					"schema": map[string]any{
						"type": "string",
					},
					"includeKeys": map[string]any{
						"type": "array",
						"items": map[string]any{
							"type": "string",
						},
					},
				},
			},
			"schemas": map[string]any{
				"type":                 "object",
				"additionalProperties": false,
				"required":             []string{"sources"},
				"properties": map[string]any{
					"sources": map[string]any{
						"type":        "array",
						"minItems":    1,
						"items":       map[string]any{"type": "string"},
						"uniqueItems": true,
					},
					"inferFromFilename": map[string]any{
						"type": "boolean",
					},
				},
			},
			"resolution": map[string]any{
				"type":                 "object",
				"additionalProperties": false,
				"required":             []string{"filenamePattern", "extensions"},
				"properties": map[string]any{
					"filenamePattern": map[string]any{
						"type": "string",
					},
					"extensions": map[string]any{
						"type":        "array",
						"minItems":    1,
						"items":       map[string]any{"type": "string"},
						"uniqueItems": true,
					},
				},
			},
		},
	}

	b, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		panic(err)
	}
	return append(b, '\n')
}

func writeFile(path string, b []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	existing, err := os.ReadFile(path)
	if err == nil && bytes.Equal(existing, b) {
		return nil
	}

	return os.WriteFile(path, b, 0o644)
}
