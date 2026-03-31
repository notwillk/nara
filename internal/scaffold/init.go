package scaffold

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/notwillk/nara/internal/configschema"
)

const (
	defaultConfig = `version: 1
$schema: .schemas/nara.config.schema.json
updateSchemaOnFormat: true
paths:
  "/": "."
meta:
  ref: "$ref"
  id: "$id"
  schema: "$schema"
  includeKeys: []
schemas:
  sources:
    - schemas/*.cue
  inferFromFilename: true
resolution:
  filenamePattern: <id>.<schema>
  extensions:
    - .yaml
    - .yml
    - .json
`
	defaultSchema = `$id: string
$schema: "note"
title: string
body?: string
`
	defaultEntry = `$id: hello
$schema: note
title: Hello
body: Example entity
`
)

func Init(configPath string) error {
	baseDir := filepath.Dir(configPath)
	if err := writeNewFile(configPath, []byte(defaultConfig)); err != nil {
		return err
	}
	if err := writeNewFile(filepath.Join(baseDir, "schemas", "note.cue"), []byte(defaultSchema)); err != nil {
		return err
	}
	if err := writeNewFile(filepath.Join(baseDir, "examples", "hello.note.yaml"), []byte(defaultEntry)); err != nil {
		return err
	}
	return configschema.Format(configPath)
}

func writeNewFile(path string, b []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		if os.IsExist(err) {
			return fmt.Errorf("%s already exists", path)
		}
		return err
	}
	defer f.Close()

	_, err = f.Write(b)
	return err
}
