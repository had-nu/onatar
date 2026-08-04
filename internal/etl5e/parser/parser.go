package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ParseFile reads and parses a JSON file from the 5etools data directory.
// Returns the raw JSON as a map for flexible parsing.
func ParseFile(root, path string) (map[string]any, error) {
	fullPath := filepath.Join(root, path)
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", fullPath, err)
	}

	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("parse %s: %w", fullPath, err)
	}
	return result, nil
}

// ParseFileFullPath reads and parses a JSON file given its full path.
func ParseFileFullPath(fullPath string) (map[string]any, error) {
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", fullPath, err)
	}

	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("parse %s: %w", fullPath, err)
	}
	return result, nil
}

// ParseFilesInDir parses all .json files in a directory matching a pattern.
func ParseFilesInDir(root, dir, pattern string) ([]map[string]any, error) {
	fullDir := filepath.Join(root, dir)
	entries, err := os.ReadDir(fullDir)
	if err != nil {
		return nil, err
	}

	var results []map[string]any
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		if pattern != "" && !strings.Contains(entry.Name(), pattern) {
			continue
		}
		parsed, err := ParseFile(root, filepath.Join(dir, entry.Name()))
		if err != nil {
			return nil, err
		}
		results = append(results, parsed)
	}
	return results, nil
}

// UnwrapArray extracts an array from a map by key, handling the 5etools wrapper format.
// 5etools wraps arrays in objects with a single key, e.g. {"class": [...], "subclass": [...]}.
func UnwrapArray(data map[string]any, key string) ([]any, bool) {
	if data == nil {
		return nil, false
	}
	val, ok := data[key]
	if !ok {
		return nil, false
	}
	arr, ok := val.([]any)
	if !ok {
		// Single object wrapped in array? Try to wrap.
		if m, ok := val.(map[string]any); ok {
			return []any{m}, true
		}
		return nil, false
	}
	return arr, true
}

// GetString safely extracts a string from a map.
func GetString(m map[string]any, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// GetInt safely extracts an int from a map.
func GetInt(m map[string]any, key string) int {
	if v, ok := m[key]; ok {
		switch v := v.(type) {
		case float64:
			return int(v)
		case int:
			return v
		case int64:
			return int(v)
		}
	}
	return 0
}

// GetBool safely extracts a bool from a map.
func GetBool(m map[string]any, key string) bool {
	if v, ok := m[key]; ok {
		if b, ok := v.(bool); ok {
			return b
		}
	}
	return false
}

// GetStringSlice safely extracts a []string from a map.
func GetStringSlice(m map[string]any, key string) []string {
	if v, ok := m[key]; ok {
		if arr, ok := v.([]any); ok {
			var out []string
			for _, item := range arr {
				if s, ok := item.(string); ok {
					out = append(out, s)
				}
			}
			return out
		}
	}
	return nil
}