package cuevalidate

import (
	"os"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"github.com/notwillk/nara/internal/errors"
	"github.com/notwillk/nara/internal/schema"
)

type Validator struct {
	ctx     *cue.Context
	schemas map[string]cue.Value
}

// New creates a validator from discovered schema files.
func New(discovered []schema.Schema) (*Validator, error) {
	ctx := cuecontext.New()
	v := &Validator{
		ctx:     ctx,
		schemas: map[string]cue.Value{},
	}

	for _, s := range discovered {
		b, err := os.ReadFile(s.Path)
		if err != nil {
			return nil, errors.New(errors.CategorySchema, s.Path, "unable to read cue schema")
		}
		value := ctx.CompileString(string(b), cue.Filename(s.Path))
		if value.Err() != nil {
			return nil, errors.New(errors.CategorySchema, s.Path, value.Err().Error())
		}
		v.schemas[s.Name] = value
	}
	return v, nil
}

// ValidateEntity validates one expanded entity map against named schema.
func (v *Validator) ValidateEntity(schemaName string, payload map[string]any, filePath string) error {
	sv, ok := v.schemas[schemaName]
	if !ok {
		return errors.New(errors.CategoryValidation, filePath, "unknown schema: "+schemaName)
	}

	enc := v.ctx.Encode(payload)
	u := sv.Unify(enc)
	if err := u.Validate(); err != nil {
		return errors.New(errors.CategoryValidation, filePath, err.Error())
	}
	return nil
}

