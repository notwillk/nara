package config

// Config is the top-level config schema for nara.
//
// See PRD §5 for fields and semantics.
type Config struct {
	Version              int                 `yaml:"version" json:"version"`
	SchemaPath           string             `yaml:"$schema" json:"$schema"`
	UpdateSchemaOnFormat bool              `yaml:"updateSchemaOnFormat" json:"updateSchemaOnFormat"`

	Paths       map[string]string `yaml:"paths" json:"paths"`
	Meta        Meta               `yaml:"meta" json:"meta"`
	Schemas     Schemas           `yaml:"schemas" json:"schemas"`
	Resolution  Resolution        `yaml:"resolution" json:"resolution"`
}

type Meta struct {
	Ref         string   `yaml:"ref" json:"ref"`
	ID          string   `yaml:"id" json:"id"`
	Schema      string   `yaml:"schema" json:"schema"`
	IncludeKeys []string `yaml:"includeKeys" json:"includeKeys"`
}

type Schemas struct {
	Sources            []string `yaml:"sources" json:"sources"`
	InferFromFilename bool     `yaml:"inferFromFilename" json:"inferFromFilename"`
}

type Resolution struct {
	FilenamePattern string   `yaml:"filenamePattern" json:"filenamePattern"`
	Extensions       []string `yaml:"extensions" json:"extensions"`
}

