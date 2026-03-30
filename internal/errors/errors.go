package errors

import "fmt"

// Category classifies errors for user-facing reporting.
type Category string

const (
	CategoryConfig     Category = "config error"
	CategorySchema     Category = "schema error"
	CategoryResolution Category = "resolution error"
	CategoryValidation Category = "validation error"
)

// NaraError is the standard error type emitted by nara.
//
// Line/Field are best-effort; callers should populate when available.
type NaraError struct {
	Category Category
	File     string // usually a config/entity/schema file path
	Line     int
	Field    string
	Message  string
}

func (e *NaraError) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Line > 0 && e.Field != "" {
		return fmt.Sprintf("%s: %s:%d:%s: %s", e.Category, e.File, e.Line, e.Field, e.Message)
	}
	if e.Line > 0 {
		return fmt.Sprintf("%s: %s:%d: %s", e.Category, e.File, e.Line, e.Message)
	}
	if e.Field != "" {
		return fmt.Sprintf("%s: %s:%s: %s", e.Category, e.File, e.Field, e.Message)
	}
	return fmt.Sprintf("%s: %s: %s", e.Category, e.File, e.Message)
}

func New(category Category, file string, message string) *NaraError {
	return &NaraError{
		Category: category,
		File:     file,
		Message:  message,
	}
}

func Wrap(category Category, file string, line int, field string, message string) *NaraError {
	return &NaraError{
		Category: category,
		File:     file,
		Line:     line,
		Field:    field,
		Message:  message,
	}
}

