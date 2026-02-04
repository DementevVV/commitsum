package repository

import "github.com/DementevVV/commitsum/internal/domain/entity"

// CacheRepository defines the interface for caching commits.
type CacheRepository interface {
	// GetCommits retrieves cached commits for a given author and date range.
	GetCommits(author, dateRange string) (*entity.CommitData, bool, error)

	// SetCommits stores commits in the cache.
	SetCommits(author, dateRange string, data *entity.CommitData) error

	// Invalidate removes cached data for a user.
	Invalidate(author string) error

	// Clear removes all cached data.
	Clear() error
}
