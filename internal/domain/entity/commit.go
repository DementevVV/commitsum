// Package entity contains the core domain entities.
package entity

// Commit represents a repository commit with its message.
type Commit struct {
	Repository string
	Message    string
}

// CommitData represents commits grouped by repository.
type CommitData struct {
	Commits  map[string][]Commit
	RepoList []string
	Warning  string
}
