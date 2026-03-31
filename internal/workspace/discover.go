package workspace

import (
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/notwillk/nara/internal/config"
	"github.com/notwillk/nara/internal/errors"
	"github.com/notwillk/nara/internal/graph"
	"github.com/notwillk/nara/internal/naming"
)

type EntryFile struct {
	ID     string
	Schema string
	Path   string
}

func DiscoverEntries(cfg *config.Config, baseDir string, allowedSchemas map[string]bool) ([]EntryFile, error) {
	if cfg == nil {
		return nil, errors.New(errors.CategoryConfig, "", "config is nil")
	}

	roots := uniqueRoots(cfg, baseDir)
	seen := map[string]bool{}
	out := make([]EntryFile, 0)
	for _, root := range roots {
		info, err := os.Stat(root)
		if err != nil {
			return nil, errors.New(errors.CategoryResolution, root, "configured path does not exist")
		}

		if !info.IsDir() {
			entry, ok := entryFromPath(root, cfg, allowedSchemas)
			if ok && !seen[entry.Path] {
				seen[entry.Path] = true
				out = append(out, entry)
			}
			continue
		}

		err = filepath.WalkDir(root, func(path string, d os.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if d.IsDir() {
				return nil
			}

			entry, ok := entryFromPath(path, cfg, allowedSchemas)
			if !ok || seen[entry.Path] {
				return nil
			}
			seen[entry.Path] = true
			out = append(out, entry)
			return nil
		})
		if err != nil {
			return nil, errors.New(errors.CategoryResolution, root, err.Error())
		}
	}

	slices.SortFunc(out, func(a EntryFile, b EntryFile) int {
		return strings.Compare(a.Path, b.Path)
	})
	return out, nil
}

func EntriesForGlobs(globs []string, cfg *config.Config, allowedSchemas map[string]bool) ([]EntryFile, error) {
	paths, err := graph.NormalizeEntryPaths(globs)
	if err != nil {
		return nil, err
	}

	out := make([]EntryFile, 0, len(paths))
	for _, path := range paths {
		entry, ok := entryFromPath(path, cfg, allowedSchemas)
		if !ok {
			return nil, errors.New(errors.CategoryValidation, path, "filename does not match resolution.filenamePattern")
		}
		out = append(out, entry)
	}

	slices.SortFunc(out, func(a EntryFile, b EntryFile) int {
		return strings.Compare(a.Path, b.Path)
	})
	return out, nil
}

func uniqueRoots(cfg *config.Config, baseDir string) []string {
	seen := map[string]bool{}
	roots := make([]string, 0, len(cfg.Paths))
	for _, p := range cfg.Paths {
		root := p
		if !filepath.IsAbs(root) {
			root = filepath.Join(baseDir, root)
		}
		root = filepath.Clean(root)
		if seen[root] {
			continue
		}
		seen[root] = true
		roots = append(roots, root)
	}
	slices.Sort(roots)
	return roots
}

func entryFromPath(path string, cfg *config.Config, allowedSchemas map[string]bool) (EntryFile, bool) {
	if cfg == nil || !hasSupportedExt(path, cfg.Resolution.Extensions) {
		return EntryFile{}, false
	}

	id, schemaName, ok := naming.InferFromPath(path, cfg.Resolution.FilenamePattern)
	if !ok {
		return EntryFile{}, false
	}
	if allowedSchemas != nil && !allowedSchemas[schemaName] {
		return EntryFile{}, false
	}

	abs, err := filepath.Abs(path)
	if err != nil {
		abs = filepath.Clean(path)
	}

	return EntryFile{
		ID:     id,
		Schema: schemaName,
		Path:   abs,
	}, true
}

func hasSupportedExt(path string, exts []string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	for _, candidate := range exts {
		if strings.ToLower(candidate) == ext {
			return true
		}
	}
	return false
}
