package cache

import (
	"os"
	"path/filepath"
	"time"

	"github.com/DementevVV/commitsum/internal/domain/entity"
	"github.com/DementevVV/commitsum/internal/domain/repository"
	"github.com/DementevVV/commitsum/internal/infrastructure/logger"
)

// cachedCommitData represents cached commit data.
type cachedCommitData struct {
	Commits  map[string][]entity.Commit `json:"commits"`
	RepoList []string                   `json:"repo_list"`
	Warning  string                     `json:"warning"`
}

// CommitsCache represents a specialized cache for commits.
type CommitsCache struct {
	cache *FileCache
}

// Ensure CommitsCache implements CacheRepository.
var _ repository.CacheRepository = (*CommitsCache)(nil)

// NewCommitsCache creates a new commits cache.
func NewCommitsCache() (*CommitsCache, error) {
	cache, err := NewFileCache()
	if err != nil {
		return nil, err
	}
	return &CommitsCache{cache: cache}, nil
}

// GetCommits retrieves cached commits.
func (cc *CommitsCache) GetCommits(author, dateRange string) (*entity.CommitData, bool, error) {
	key := cc.cache.GetCacheKey("commits", author, dateRange)

	var data cachedCommitData
	found, err := cc.cache.Get(key, &data)
	if err != nil {
		return nil, false, err
	}

	if found {
		logger.Debug("Commits cache hit", "author", author, "date_range", dateRange)
		return &entity.CommitData{
			Commits:  data.Commits,
			RepoList: data.RepoList,
			Warning:  data.Warning,
		}, true, nil
	}

	return nil, false, nil
}

// SetCommits stores commits in the cache.
func (cc *CommitsCache) SetCommits(author, dateRange string, commitData *entity.CommitData) error {
	key := cc.cache.GetCacheKey("commits", author, dateRange)

	data := &cachedCommitData{
		Commits:  commitData.Commits,
		RepoList: commitData.RepoList,
		Warning:  commitData.Warning,
	}

	// Cache for 5 minutes for today, 1 hour for older dates.
	ttl := 5 * time.Minute
	if !isToday(dateRange) {
		ttl = time.Hour
	}

	err := cc.cache.Set(key, data, ttl)
	if err == nil {
		logger.Debug("Commits cached", "author", author, "date_range", dateRange, "ttl_minutes", ttl.Minutes())
	}

	return err
}

// Invalidate removes all cached data for a user.
func (cc *CommitsCache) Invalidate(author string) error {
	files, err := filepath.Glob(filepath.Join(cc.cache.Dir(), "*.json"))
	if err != nil {
		return err
	}

	removedCount := 0
	for _, file := range files {
		if err := os.Remove(file); err == nil {
			removedCount++
		}
	}

	logger.Info("User cache invalidated", "author", author, "files_removed", removedCount)
	return nil
}

// Clear removes all cached data.
func (cc *CommitsCache) Clear() error {
	return cc.cache.Clear()
}

// isToday checks if the date is today.
func isToday(dateRange string) bool {
	today := time.Now().Format("2006-01-02")
	return dateRange == today
}
