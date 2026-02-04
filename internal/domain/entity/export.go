package entity

import "time"

// ExportFormat represents the output format type.
type ExportFormat string

const (
	FormatText     ExportFormat = "text"
	FormatMarkdown ExportFormat = "markdown"
	FormatJSON     ExportFormat = "json"
)

// CommitExport represents a commit for export.
type CommitExport struct {
	Repository string `json:"repository"`
	Message    string `json:"message"`
}

// SummaryExport represents the full summary for export.
type SummaryExport struct {
	Date         string                    `json:"date"`
	DateRange    string                    `json:"date_range,omitempty"`
	TotalRepos   int                       `json:"total_repos"`
	TotalCommits int                       `json:"total_commits"`
	Commits      map[string][]CommitExport `json:"commits"`
	Stats        *Statistics               `json:"stats,omitempty"`
	GeneratedAt  string                    `json:"generated_at"`
}

// NewSummaryExport creates a new SummaryExport instance.
func NewSummaryExport(dateStr string) *SummaryExport {
	return &SummaryExport{
		Date:        dateStr,
		Commits:     make(map[string][]CommitExport),
		GeneratedAt: time.Now().Format(time.RFC3339),
	}
}
