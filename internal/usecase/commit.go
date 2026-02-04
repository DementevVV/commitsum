// Package usecase contains the application business logic.
package usecase

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/DementevVV/commitsum/internal/domain/entity"
	"github.com/DementevVV/commitsum/internal/domain/repository"
)

// CommitUseCase handles commit-related business logic.
type CommitUseCase struct {
	github repository.GitHubRepository
	cache  repository.CacheRepository
}

// NewCommitUseCase creates a new CommitUseCase.
func NewCommitUseCase(github repository.GitHubRepository, cache repository.CacheRepository) *CommitUseCase {
	return &CommitUseCase{
		github: github,
		cache:  cache,
	}
}

// GetCommitsForRange fetches commits for a date range.
func (uc *CommitUseCase) GetCommitsForRange(startDate, endDate string) (*entity.CommitData, error) {
	// Validate date range.
	if err := uc.validateDateRange(startDate, endDate); err != nil {
		return nil, err
	}

	// Get GitHub user.
	ghUser, err := uc.github.GetUser()
	if err != nil {
		return nil, fmt.Errorf("failed to get GitHub user: %w", err)
	}

	// Build date range query.
	dateRange := startDate
	if startDate != endDate {
		dateRange = fmt.Sprintf("%s..%s", startDate, endDate)
	}

	// Try cache first.
	if uc.cache != nil {
		if data, found, err := uc.cache.GetCommits(ghUser, dateRange); err == nil && found {
			return data, nil
		}
	}

	// Fetch from GitHub.
	data, err := uc.github.FetchCommitsByAuthorAndDate(ghUser, dateRange)
	if err != nil {
		return nil, err
	}

	// Store in cache.
	if uc.cache != nil {
		_ = uc.cache.SetCommits(ghUser, dateRange, data)
	}

	return data, nil
}

func (uc *CommitUseCase) validateDateRange(startDate, endDate string) error {
	startTime, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return fmt.Errorf("invalid start date format: %w", err)
	}

	endTime, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return fmt.Errorf("invalid end date format: %w", err)
	}

	if startTime.After(endTime) {
		return fmt.Errorf("start date cannot be after end date")
	}

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())
	if endTime.After(today) {
		return fmt.Errorf("end date cannot be in the future")
	}

	return nil
}

// FilterReposByPattern filters repositories by glob pattern.
func (uc *CommitUseCase) FilterReposByPattern(repos []string, pattern string) []string {
	if pattern == "" {
		return repos
	}

	var filtered []string
	for _, repo := range repos {
		if matched, _ := matchPattern(pattern, repo); matched {
			filtered = append(filtered, repo)
		}
	}
	return filtered
}

// matchPattern matches a repository name against a pattern.
func matchPattern(pattern, name string) (bool, error) {
	pattern = strings.ToLower(pattern)
	name = strings.ToLower(name)

	// Simple contains check for non-glob patterns.
	if !strings.ContainsAny(pattern, "*?[]") {
		return strings.Contains(name, pattern), nil
	}

	// Convert basic glob to regex; allow * to match across '/'.
	rePattern := regexp.QuoteMeta(pattern)
	rePattern = strings.ReplaceAll(rePattern, `\*`, ".*")
	rePattern = strings.ReplaceAll(rePattern, `\?`, ".")
	rePattern = "^" + rePattern + "$"

	re, err := regexp.Compile(rePattern)
	if err != nil {
		// Fallback to contains for invalid patterns.
		clean := strings.ReplaceAll(strings.ReplaceAll(pattern, "*", ""), "?", "")
		return strings.Contains(name, clean), nil
	}
	return re.MatchString(name), nil
}

// CalculateStatistics calculates statistics for selected commits.
func (uc *CommitUseCase) CalculateStatistics(commits map[string][]entity.Commit, selected map[string]bool) *entity.Statistics {
	stats := &entity.Statistics{
		CommitsPerRepo: make(map[string]int),
	}

	for repo, repoCommits := range commits {
		if !selected[repo] {
			continue
		}
		count := len(repoCommits)
		stats.CommitsPerRepo[repo] = count
		stats.TotalCommits += count
		stats.TotalRepositories++

		if count > stats.MaxCommits {
			stats.MaxCommits = count
			stats.MostActiveRepo = repo
		}
	}

	return stats
}

// GetSelectedReposSorted returns a sorted slice of selected repository names.
func (uc *CommitUseCase) GetSelectedReposSorted(commits map[string][]entity.Commit, selected map[string]bool) []string {
	var repos []string
	for repo := range commits {
		if selected[repo] {
			repos = append(repos, repo)
		}
	}
	sort.Strings(repos)
	return repos
}
