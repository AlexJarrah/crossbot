package crossbot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var CallbackCache = make(map[string]Callback)

// DefaultCacheDirectory creates and returns a temporary directory to store cache
func (c *Config) DefaultCacheDirectory() (string, error) {
	dir := filepath.Join(os.TempDir(), c.ID)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("failed to create cache directory '%s': %w", dir, err)
	}

	return dir, nil
}

func (c *Config) ReadCache(key string, value any) error {
	path := filepath.Join(c.CacheDirectory, key+".json")
	file, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	if err := json.Unmarshal(file, value); err != nil {
		return fmt.Errorf("failed to unmarshal json data: %w", err)
	}

	return nil
}

func (c *Config) WriteCache(key string, value []byte) error {
	path := filepath.Join(c.CacheDirectory, key+".json")
	if _, err := os.Create(path); err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}

	if err := os.WriteFile(path, value, 0o644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func (c *Config) NewCacheID(fields map[string]string) string {
	keys := make([]string, 0, len(fields))
	for k := range fields {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	pairs := make([]string, 0, len(fields))
	for _, k := range keys {
		pairs = append(pairs, fmt.Sprintf("%s:%s", k, fields[k]))
	}

	return strings.Join(pairs, "-")
}
