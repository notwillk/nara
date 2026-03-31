package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/notwillk/nara/internal/errors"
	"github.com/notwillk/nara/internal/naming"
	"gopkg.in/yaml.v3"
)

// Load reads nara.yaml and returns the parsed config.
func Load(configPath string) (*Config, error) {
	cfg, err := LoadRaw(configPath)
	if err != nil {
		return nil, err
	}

	baseDir := filepath.Dir(configPath)
	if err := Validate(cfg, configPath, baseDir); err != nil {
		return nil, err
	}

	return cfg, nil
}

func LoadRaw(configPath string) (*Config, error) {
	b, err := os.ReadFile(configPath)
	if err != nil {
		return nil, errors.New(errors.CategoryConfig, configPath, "unable to read config file")
	}

	var cfg Config
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return nil, errors.New(errors.CategoryConfig, configPath, "config YAML parse error")
	}
	return &cfg, nil
}

func Validate(cfg *Config, configPath string, baseDir string) error {
	if cfg == nil {
		return errors.New(errors.CategoryConfig, configPath, "config is empty")
	}

	if cfg.Version != 1 {
		return errors.New(errors.CategoryConfig, configPath, "version must be 1")
	}

	if cfg.SchemaPath != "" {
		schemaAbs := cfg.SchemaPath
		if !filepath.IsAbs(schemaAbs) {
			schemaAbs = filepath.Join(baseDir, schemaAbs)
		}

		sb, err := os.ReadFile(schemaAbs)
		if err != nil {
			return errors.New(errors.CategoryConfig, configPath, "$schema file must exist")
		}

		var v any
		if err := json.Unmarshal(sb, &v); err != nil {
			return errors.New(errors.CategoryConfig, configPath, "$schema must be valid JSON")
		}
	}

	if len(cfg.Paths) == 0 {
		return errors.New(errors.CategoryConfig, configPath, "paths is required")
	}
	if cfg.Meta.Ref == "" || cfg.Meta.ID == "" || cfg.Meta.Schema == "" {
		return errors.New(errors.CategoryConfig, configPath, "meta.ref/meta.id/meta.schema are required")
	}
	if len(cfg.Schemas.Sources) == 0 {
		return errors.New(errors.CategoryConfig, configPath, "schemas.sources is required")
	}
	if cfg.Resolution.FilenamePattern == "" {
		return errors.New(errors.CategoryConfig, configPath, "resolution.filenamePattern is required")
	}
	if len(cfg.Resolution.Extensions) == 0 {
		return errors.New(errors.CategoryConfig, configPath, "resolution.extensions is required")
	}
	if err := naming.ValidateFilenamePattern(cfg.Resolution.FilenamePattern); err != nil {
		return errors.New(errors.CategoryConfig, configPath, err.Error())
	}

	for _, src := range cfg.Schemas.Sources {
		if _, err := filepath.Glob(filepath.Join(baseDir, src)); err != nil {
			return errors.New(errors.CategoryConfig, configPath, "schemas.sources contains invalid glob pattern")
		}
	}

	for _, p := range cfg.Paths {
		if p == "" {
			return errors.New(errors.CategoryConfig, configPath, "paths entries must be non-empty")
		}
		if filepath.IsAbs(p) {
			return errors.New(errors.CategoryConfig, configPath, "paths entries must be relative")
		}
	}

	for _, ext := range cfg.Resolution.Extensions {
		if ext == "" || !strings.HasPrefix(ext, ".") {
			return errors.New(errors.CategoryConfig, configPath, "resolution.extensions entries must start with .")
		}
	}

	return nil
}
