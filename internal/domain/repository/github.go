// Package repository defines the interfaces for external data access.
package repository

import "github.com/DementevVV/commitsum/internal/domain/entity"

// GitHubRepository defines the interface for GitHub data access.
type GitHubRepository interface {
	// GetUser returns the currently authenticated GitHub username.
	GetUser() (string, error)

	// FetchCommitsByAuthorAndDate fetches commits for a given author and date range.
	FetchCommitsByAuthorAndDate(author, dateRange string) (*entity.CommitData, error)
}
