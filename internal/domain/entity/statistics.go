package entity

// Statistics holds commit statistics.
type Statistics struct {
	TotalCommits      int            `json:"total_commits"`
	TotalRepositories int            `json:"total_repositories"`
	CommitsPerRepo    map[string]int `json:"commits_per_repo"`
	MostActiveRepo    string         `json:"most_active_repo"`
	MaxCommits        int            `json:"max_commits"`
}
