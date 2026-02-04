// Package cache provides caching implementations.
package cache

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/DementevVV/commitsum/internal/infrastructure/logger"
)

// Entry represents a cache entry.
type Entry struct {
	Data      interface{}   `json:"data"`
	Timestamp time.Time     `json:"timestamp"`
	TTL       time.Duration `json:"ttl"`
}

// IsExpired checks if the entry has expired.
func (e *Entry) IsExpired() bool {
	return time.Since(e.Timestamp) > e.TTL
}

// FileCache represents a file-based cache.
type FileCache struct {
	dir string
}

// NewFileCache creates a new file cache.
func NewFileCache() (*FileCache, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	cacheDir := filepath.Join(homeDir, ".config", "commitsum", "cache")
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create cache directory: %w", err)
	}

	return &FileCache{dir: cacheDir}, nil
}

// GetCacheKey generates a cache key based on parameters.
func (c *FileCache) GetCacheKey(prefix string, params ...string) string {
	combined := prefix
	for _, param := range params {
		combined += "-" + param
	}

	hash := md5.Sum([]byte(combined))
	return fmt.Sprintf("%x.json", hash)
}

// getCacheFilePath returns the path to a cache file.
func (c *FileCache) getCacheFilePath(key string) string {
	return filepath.Join(c.dir, key)
}

// Set stores data in the cache.
func (c *FileCache) Set(key string, data interface{}, ttl time.Duration) error {
	entry := &Entry{
		Data:      data,
		Timestamp: time.Now(),
		TTL:       ttl,
	}

	jsonData, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal cache entry: %w", err)
	}

	filePath := c.getCacheFilePath(key)
	if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write cache file: %w", err)
	}

	logger.Debug("Cache entry saved", "key", key, "ttl_minutes", ttl.Minutes())
	return nil
}

// Get retrieves data from the cache.
func (c *FileCache) Get(key string, target interface{}) (bool, error) {
	filePath := c.getCacheFilePath(key)

	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil // Cache not found.
		}
		return false, fmt.Errorf("failed to read cache file: %w", err)
	}

	var entry Entry
	if err := json.Unmarshal(data, &entry); err != nil {
		// Remove corrupted cache file.
		_ = os.Remove(filePath)
		return false, fmt.Errorf("failed to unmarshal cache entry: %w", err)
	}

	// Check expiration.
	if entry.IsExpired() {
		_ = os.Remove(filePath)
		logger.Debug("Cache entry expired and removed", "key", key)
		return false, nil
	}

	// Deserialize data into target structure.
	entryData, err := json.Marshal(entry.Data)
	if err != nil {
		return false, fmt.Errorf("failed to marshal entry data: %w", err)
	}

	if err := json.Unmarshal(entryData, target); err != nil {
		return false, fmt.Errorf("failed to unmarshal target data: %w", err)
	}

	logger.Debug("Cache hit", "key", key)
	return true, nil
}

// Delete removes an entry from the cache.
func (c *FileCache) Delete(key string) error {
	filePath := c.getCacheFilePath(key)
	err := os.Remove(filePath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete cache file: %w", err)
	}
	return nil
}

// Clear removes all cache entries.
func (c *FileCache) Clear() error {
	files, err := filepath.Glob(filepath.Join(c.dir, "*.json"))
	if err != nil {
		return fmt.Errorf("failed to list cache files: %w", err)
	}

	for _, file := range files {
		_ = os.Remove(file)
	}

	logger.Info("Cache cleared", "files_removed", len(files))
	return nil
}

// CleanExpired removes expired entries.
func (c *FileCache) CleanExpired() error {
	files, err := filepath.Glob(filepath.Join(c.dir, "*.json"))
	if err != nil {
		return fmt.Errorf("failed to list cache files: %w", err)
	}

	removedCount := 0
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			continue
		}

		var entry Entry
		if err := json.Unmarshal(data, &entry); err != nil {
			// Remove corrupted files.
			_ = os.Remove(file)
			removedCount++
			continue
		}

		if entry.IsExpired() {
			_ = os.Remove(file)
			removedCount++
		}
	}

	if removedCount > 0 {
		logger.Info("Expired cache entries cleaned", "removed_count", removedCount)
	}

	return nil
}

// GetStats returns cache statistics.
func (c *FileCache) GetStats() (map[string]interface{}, error) {
	files, err := filepath.Glob(filepath.Join(c.dir, "*.json"))
	if err != nil {
		return nil, err
	}

	stats := map[string]interface{}{
		"total_files":      len(files),
		"expired_files":    0,
		"total_size_bytes": int64(0),
	}

	expiredCount := 0
	var totalSize int64

	for _, file := range files {
		fileInfo, err := os.Stat(file)
		if err != nil {
			continue
		}

		totalSize += fileInfo.Size()

		data, err := os.ReadFile(file)
		if err != nil {
			continue
		}

		var entry Entry
		if err := json.Unmarshal(data, &entry); err != nil {
			expiredCount++
			continue
		}

		if entry.IsExpired() {
			expiredCount++
		}
	}

	stats["expired_files"] = expiredCount
	stats["total_size_bytes"] = totalSize

	return stats, nil
}

// Dir returns the cache directory.
func (c *FileCache) Dir() string {
	return c.dir
}
