package schema

import (
	"path/filepath"
	"sort"
	"strings"

	"github.com/notwillk/nara/internal/config"
	"github.com/notwillk/nara/internal/errors"
)

type Schema struct {
	Name string
	Path string
}

// Discover loads schema files from configured sources.
//
// Rules (PRD v1):
// - filename stem = schema name
// - one schema per file (no sub-schema addressing)
// - fail on duplicate schema names
//
// Note: this ticket does "discovery" only (filenames + uniqueness).
// Validation of CUE semantics happens in later tickets.
func Discover(cfg *config.Config, baseDir string) ([]Schema, error) {
	if cfg == nil {
		return nil, errors.New(errors.CategoryConfig, "", "config is nil")
	}
	if len(cfg.Schemas.Sources) == 0 {
		return nil, errors.New(errors.CategoryConfig, "", "schemas.sources is required")
	}

	found := map[string]string{} // schemaName -> path
	added := map[string]bool{}   // file path de-dupe across overlapping globs
	out := make([]Schema, 0)

	for _, pattern := range cfg.Schemas.Sources {
		// Patterns are relative to config directory (PRD default in examples).
		globPattern := pattern
		if !filepath.IsAbs(globPattern) {
			globPattern = filepath.Join(baseDir, globPattern)
		}

		matches, err := filepath.Glob(globPattern)
		if err != nil {
			return nil, errors.New(errors.CategorySchema, "", "invalid schema glob pattern")
		}

		if len(matches) == 0 {
			return nil, errors.New(errors.CategorySchema, "", "no schema files matched: "+pattern)
		}

		for _, p := range matches {
			// Only .cue files are supported.
			if strings.ToLower(filepath.Ext(p)) != ".cue" {
				continue
			}
			abs, _ := filepath.Abs(p)
			if added[abs] {
				continue
			}
			added[abs] = true

			stem := strings.TrimSuffix(filepath.Base(p), filepath.Ext(p))
			if prev, ok := found[stem]; ok {
				// Fail on duplicates per PRD.
				return nil, errors.Wrap(errors.CategorySchema, prev, 0, "", "duplicate schema name: "+stem)
			}
			found[stem] = p
			out = append(out, Schema{Name: stem, Path: p})
		}
	}

	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out, nil
}

