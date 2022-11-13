// Package cmdutil contains utilities for writing CLIs.
package cmdutil

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// ReadJSONFile parses the JSON contents of a file into the destination argument
// and returns dst or an error.
func ReadJSONFile[T any](filename string, dst T) (T, error) {
	jsonBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return dst, fmt.Errorf("failed to read input file %q: %w", filename, err)
	}
	if err := json.Unmarshal(jsonBytes, dst); err != nil {
		return dst, fmt.Errorf("failed to unmarshal JSON from file %q: %w", filename, err)
	}
	return dst, nil
}

// WriteJSONFile writes a value as JSON to a file.
func WriteJSONFile(filename string, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON data for file %q: %w", filename, err)
	}
	data = append(data, '\n')
	if err := ioutil.WriteFile(filename, data, 0664); err != nil {
		return fmt.Errorf("failed to write JSON file %q: %w", filename, err)
	}
	return nil
}
