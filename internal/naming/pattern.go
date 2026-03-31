package naming

import (
	"fmt"
	"path/filepath"
	"strings"
)

const (
	TokenID     = "<id>"
	TokenSchema = "<schema>"
)

func ValidateFilenamePattern(pattern string) error {
	if pattern == "" {
		return fmt.Errorf("resolution.filenamePattern is required")
	}
	if strings.Count(pattern, TokenID) != 1 || strings.Count(pattern, TokenSchema) != 1 {
		return fmt.Errorf("resolution.filenamePattern must contain %s and %s exactly once", TokenID, TokenSchema)
	}
	if strings.Contains(pattern, "/") || strings.Contains(pattern, `\`) {
		return fmt.Errorf("resolution.filenamePattern must not contain path separators")
	}

	idIdx := strings.Index(pattern, TokenID)
	schemaIdx := strings.Index(pattern, TokenSchema)
	if idIdx == -1 || schemaIdx == -1 {
		return fmt.Errorf("resolution.filenamePattern must contain %s and %s", TokenID, TokenSchema)
	}

	var middle string
	if idIdx < schemaIdx {
		middle = pattern[idIdx+len(TokenID) : schemaIdx]
	} else {
		middle = pattern[schemaIdx+len(TokenSchema) : idIdx]
	}
	if middle == "" {
		return fmt.Errorf("resolution.filenamePattern must include a literal separator between %s and %s", TokenID, TokenSchema)
	}

	return nil
}

func InferFromPath(path string, pattern string) (id string, schema string, ok bool) {
	stem := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	if stem == "" {
		return "", "", false
	}

	if pattern == "" {
		pattern = TokenID + "." + TokenSchema
	}
	if err := ValidateFilenamePattern(pattern); err != nil {
		return "", "", false
	}

	idIdx := strings.Index(pattern, TokenID)
	schemaIdx := strings.Index(pattern, TokenSchema)
	if idIdx < schemaIdx {
		prefix := pattern[:idIdx]
		middle := pattern[idIdx+len(TokenID) : schemaIdx]
		suffix := pattern[schemaIdx+len(TokenSchema):]
		first, second, matched := splitPattern(stem, prefix, middle, suffix, true)
		return first, second, matched
	}

	prefix := pattern[:schemaIdx]
	middle := pattern[schemaIdx+len(TokenSchema) : idIdx]
	suffix := pattern[idIdx+len(TokenID):]
	first, second, matched := splitPattern(stem, prefix, middle, suffix, false)
	return second, first, matched
}

func splitPattern(stem string, prefix string, middle string, suffix string, firstIsID bool) (first string, second string, ok bool) {
	if !strings.HasPrefix(stem, prefix) || !strings.HasSuffix(stem, suffix) {
		return "", "", false
	}

	body := strings.TrimPrefix(stem, prefix)
	body = strings.TrimSuffix(body, suffix)
	if body == "" {
		return "", "", false
	}

	var splitAt int
	if firstIsID {
		splitAt = strings.LastIndex(body, middle)
	} else {
		splitAt = strings.Index(body, middle)
	}
	if splitAt <= 0 || splitAt+len(middle) >= len(body) {
		return "", "", false
	}

	first = body[:splitAt]
	second = body[splitAt+len(middle):]
	if first == "" || second == "" {
		return "", "", false
	}
	return first, second, true
}
