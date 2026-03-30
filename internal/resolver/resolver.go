package resolver

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/notwillk/nara/internal/config"
	"github.com/notwillk/nara/internal/errors"
)

// Resolve resolves a $ref string to a concrete entity file path.
//
// Order:
// 1) ~alias
// 2) /root
// 3) relative ./ ../
// 4) bare
func Resolve(currentPath string, currentSchema string, ref string, cfg *config.Config, baseDir string) (string, error) {
	ref = strings.TrimSpace(ref)
	if ref == "" {
		return "", errors.New(errors.CategoryResolution, currentPath, "empty ref")
	}

	switch {
	case strings.HasPrefix(ref, "~"):
		return resolveAlias(currentPath, currentSchema, ref, cfg, baseDir)
	case strings.HasPrefix(ref, "/"):
		rootBase := resolveBasePath(cfg.Paths["/"], baseDir)
		return resolveFromBase(rootBase, strings.TrimPrefix(ref, "/"), currentSchema, cfg, currentPath)
	case strings.HasPrefix(ref, "./") || strings.HasPrefix(ref, "../"):
		dir := filepath.Dir(currentPath)
		return resolveRelative(dir, ref, currentSchema, cfg, currentPath)
	default:
		// bare sibling id
		dir := filepath.Dir(currentPath)
		return findByIDAndSchema(dir, ref, currentSchema, cfg, currentPath)
	}
}

func resolveAlias(currentPath string, currentSchema string, ref string, cfg *config.Config, baseDir string) (string, error) {
	parts := strings.SplitN(strings.TrimPrefix(ref, "~"), "/", 2)
	aliasKey := "~" + parts[0]
	aliasBase := cfg.Paths[aliasKey]
	if aliasBase == "" && parts[0] == "" {
		aliasBase = cfg.Paths["~"]
	}
	if aliasBase == "" {
		return "", errors.New(errors.CategoryResolution, currentPath, "unknown alias in ref: "+ref)
	}
	base := resolveBasePath(aliasBase, baseDir)
	target := ""
	if len(parts) == 2 {
		target = parts[1]
	}
	return resolveFromBase(base, target, currentSchema, cfg, currentPath)
}

func resolveRelative(dir string, ref string, currentSchema string, cfg *config.Config, currentPath string) (string, error) {
	joined := filepath.Clean(filepath.Join(dir, ref))
	if hasSupportedExt(joined, cfg.Resolution.Extensions) {
		if fileExists(joined) {
			return joined, nil
		}
		return "", errors.New(errors.CategoryResolution, currentPath, "missing file for ref: "+ref)
	}
	// No extension in ref; treat basename as id in destination dir.
	return findByIDAndSchema(filepath.Dir(joined), filepath.Base(joined), currentSchema, cfg, currentPath)
}

func resolveFromBase(base string, target string, schemaName string, cfg *config.Config, currentPath string) (string, error) {
	candidate := filepath.Join(base, target)
	if hasSupportedExt(candidate, cfg.Resolution.Extensions) {
		if fileExists(candidate) {
			return candidate, nil
		}
		return "", errors.New(errors.CategoryResolution, currentPath, "missing file for ref target")
	}

	// target may be a path ending with id; infer id from leaf.
	id := filepath.Base(strings.Trim(target, "/"))
	searchDir := filepath.Dir(candidate)
	if target == "" {
		id = ""
		searchDir = base
	}
	if id == "" {
		return "", errors.New(errors.CategoryResolution, currentPath, "invalid ref target")
	}
	return findByIDAndSchema(searchDir, id, schemaName, cfg, currentPath)
}

func findByIDAndSchema(dir string, id string, schemaName string, cfg *config.Config, currentPath string) (string, error) {
	pattern := cfg.Resolution.FilenamePattern
	if pattern == "" {
		pattern = "<id>.<schema>"
	}
	stem := strings.ReplaceAll(strings.ReplaceAll(pattern, "<id>", id), "<schema>", schemaName)
	for _, ext := range cfg.Resolution.Extensions {
		c := filepath.Join(dir, stem+ext)
		if fileExists(c) {
			return c, nil
		}
	}
	return "", errors.New(errors.CategoryResolution, currentPath, "unable to resolve ref: "+id)
}

func resolveBasePath(p string, baseDir string) string {
	if filepath.IsAbs(p) {
		return p
	}
	return filepath.Join(baseDir, p)
}

func hasSupportedExt(path string, exts []string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	for _, e := range exts {
		if strings.ToLower(e) == ext {
			return true
		}
	}
	return false
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

