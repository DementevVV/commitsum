// Package github provides GitHub CLI client implementation.
package github

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"sort"
	"strings"
	"time"

	"github.com/DementevVV/commitsum/internal/domain/entity"
	"github.com/DementevVV/commitsum/internal/domain/repository"
)

// commitSearchItem represents a single commit search result from the GitHub CLI.
type commitSearchItem struct {
	Repository struct {
		FullName      string `json:"full_name"`
		NameWithOwner string `json:"nameWithOwner"`
		Name          string `json:"name"`
	} `json:"repository"`
	Commit struct {
		Message         string `json:"message"`
		MessageHeadline string `json:"messageHeadline"`
	} `json:"commit"`
}

// Client encapsulates GitHub API operations via the gh CLI.
type Client struct {
	timeout time.Duration
	limit   int
}

// Ensure Client implements GitHubRepository.
var _ repository.GitHubRepository = (*Client)(nil)

// NewClient creates a new GitHub client with default settings.
func NewClient() *Client {
	return &Client{
		timeout: 20 * time.Second,
		limit:   1000,
	}
}

// GetUser retrieves the currently authenticated GitHub username using the GitHub CLI.
func (c *Client) GetUser() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "gh", "api", "user", "--jq", ".login")
	out, err := cmd.Output()
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return "", fmt.Errorf("gh api user timed out after %s", c.timeout)
		}
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}

// FetchCommitsByAuthorAndDate fetches commits for a given author and date range.
func (c *Client) FetchCommitsByAuthorAndDate(author, dateRange string) (*entity.CommitData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	cmd := exec.CommandContext(
		ctx,
		"gh",
		"search",
		"commits",
		"--author", author,
		"--committer-date", dateRange,
		"--json", "repository,commit",
		"--limit", fmt.Sprintf("%d", c.limit),
	)

	out, err := cmd.CombinedOutput()
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return nil, fmt.Errorf("gh search commits timed out after %s", c.timeout)
		}
		return nil, fmt.Errorf("failed to fetch commits: %w\n%s", err, strings.TrimSpace(string(out)))
	}

	items, err := c.parseCommitSearchItems(out)
	if err != nil {
		return nil, err
	}

	var warning string
	if len(items) >= c.limit {
		warning = fmt.Sprintf("Results capped at %d commits by GitHub; summary may be incomplete.", c.limit)
	}

	commitMap := make(map[string][]entity.Commit)
	for _, item := range items {
		repo := item.Repository.NameWithOwner
		if repo == "" {
			repo = item.Repository.FullName
		}
		if repo == "" {
			repo = item.Repository.Name
		}

		message := item.Commit.MessageHeadline
		if message == "" {
			message = strings.Split(item.Commit.Message, "\n")[0]
		}

		if repo == "" || message == "" {
			continue
		}

		commitMap[repo] = append(commitMap[repo], entity.Commit{Repository: repo, Message: message})
	}

	var repoList []string
	for repo := range commitMap {
		repoList = append(repoList, repo)
	}
	sort.Strings(repoList)

	return &entity.CommitData{
		Commits:  commitMap,
		RepoList: repoList,
		Warning:  warning,
	}, nil
}

func (c *Client) parseCommitSearchItems(data []byte) ([]commitSearchItem, error) {
	if len(bytes.TrimSpace(data)) == 0 {
		return nil, nil
	}

	trimmed := bytes.TrimSpace(data)
	var items []commitSearchItem

	if len(trimmed) > 0 && trimmed[0] == '[' {
		if err := json.Unmarshal(trimmed, &items); err != nil {
			return nil, err
		}
		return items, nil
	}

	dec := json.NewDecoder(bytes.NewReader(trimmed))
	for {
		var item commitSearchItem
		if err := dec.Decode(&item); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}
