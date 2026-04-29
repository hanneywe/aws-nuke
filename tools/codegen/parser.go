package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Parse reads a YAML DSL file at the given path and unmarshals it into a ResourceDef.
func Parse(path string) (*ResourceDef, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading DSL file: %w", err)
	}

	var def ResourceDef
	if err := yaml.Unmarshal(data, &def); err != nil {
		return nil, fmt.Errorf("parsing YAML: %w", err)
	}

	return &def, nil
}

// PrettyPrint serializes a ResourceDef back to YAML.
func PrettyPrint(def *ResourceDef) ([]byte, error) {
	data, err := yaml.Marshal(def)
	if err != nil {
		return nil, fmt.Errorf("marshaling to YAML: %w", err)
	}

	return data, nil
}
