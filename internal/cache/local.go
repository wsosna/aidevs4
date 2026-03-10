package cache

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
)

const defaultDir = ".tmp/cache"

type LocalCache struct {
	dir string
}

var DefaultLocalCache = NewLocalCache(defaultDir)

func NewLocalCache(dir string) *LocalCache {
	if err := os.MkdirAll(dir, 0755); err != nil {
		panic(fmt.Errorf("failed to create cache dir %s: %w", dir, err))
	}
	return &LocalCache{dir: dir}
}

func (c *LocalCache) Get(key string) (string, bool) {
	path := c.path(key)
	data, err := os.ReadFile(path)
	if err != nil {
		return "", false
	}
	return string(data), true
}

func (c *LocalCache) Set(key, value string) error {
	return os.WriteFile(c.path(key), []byte(value), 0644)
}

func (c *LocalCache) path(key string) string {
	hash := sha256.Sum256([]byte(key))
	return filepath.Join(c.dir, fmt.Sprintf("%x", hash))
}
