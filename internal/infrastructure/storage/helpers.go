package storage

import (
	"encoding/json"
	"os"
)

// readJSON reads and unmarshals a JSON file into the specified type.
func readJSON[T any](path string) (T, error) {
	var result T

	data, err := os.ReadFile(path)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(data, &result)
	return result, err
}

// writeJSON marshals data and writes it to a JSON file with pretty printing.
func writeJSON(path string, data any) error {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, bytes, 0644)
}
