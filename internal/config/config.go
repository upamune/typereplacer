package config

import (
	"os"

	"github.com/goccy/go-yaml"
)

// Field represents a single field within a struct:
// e.g., name="SomeField", type="string".
type Field struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
}

// Struct represents a single struct with a name and multiple fields.
type Struct struct {
	Name   string  `yaml:"name"`
	Fields []Field `yaml:"fields"`
}

// Config is the top-level config, listing imports (packages) to analyze
// and a slice of structs (each having fields).
type Config struct {
	Imports []string `yaml:"imports"`
	Structs []Struct `yaml:"structs"`
}

// LoadConfig reads a YAML file from disk and unmarshals it into Config.
// Returns an error if the file is invalid or not found.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
