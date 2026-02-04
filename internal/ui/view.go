package ui

import (
	"fmt"

	"github.com/DementevVV/commitsum/internal/domain/entity"
)

// View renders the current state of the application model.
func (m *Model) View() string {
	switch m.screen {
	case screenDateRange:
		return m.viewDateRange()
	case screenDateSelect:
		return m.viewDateSelect()
	case screenRepoFilter:
		return m.viewRepoFilter()
	case screenRepoList:
		return m.viewRepoList()
	case screenSummary:
		return m.viewSummary()
	case screenExport:
		return m.viewExport()
	case screenStats:
		return m.viewStats()
	case screenLoading:
		return m.viewLoading()
	}

	return ""
}

func (m *Model) viewDateRange() string {
	s := renderHeader("Select Time Range")
	s += styleDateLabel.Render("Choose a preset or custom date range:") + "\n\n"

	for i, preset := range entity.DateRangePresets {
		cursor := "  "
		if i == m.dateRangeIdx {
			cursor = styleCursor.Render(iconArrowRight)
		}

		label := preset.Label
		if preset.Key != "custom" {
			dr := entity.GetDateRange(preset.Key)
			label += " " + styleFooter.Render("("+entity.FormatDateDisplay(dr.StartDate, dr.EndDate)+")")
		}

		s += cursor + styleRepo.Render(label) + "\n"
	}

	s += renderHelpBar([][]string{
		{"j/k", "navigate"},
		{"enter", "select"},
		{"q", "quit"},
	})

	return "\n" + styleBox.Render(s) + "\n"
}

func (m *Model) viewDateSelect() string {
	s := renderHeader("Custom Date")

	if m.err != nil {
		s += renderErrorBanner(m.err.Error()) + "\n\n"
	}

	s += styleDateLabel.Render("Enter custom date:") + "\n\n"

	inputBox := styleInputBox.Render(m.dateInput.View())

	s += inputBox + "\n\n"
	s += styleFooter.Render("Format: YYYY-MM-DD (e.g., 2026-02-02)") + "\n"
	s += renderHelpBar([][]string{
		{"enter", "confirm"},
		{"esc", "back"},
	})

	return "\n" + styleBox.Render(s) + "\n"
}

func (m *Model) viewRepoFilter() string {
	s := renderHeader("Filter Repositories")
	s += styleDateLabel.Render("Enter filter pattern:") + "\n\n"

	inputBox := styleInputBox.Render(m.filterInput.View())

	s += inputBox + "\n\n"
	s += styleFooter.Render("Use * as wildcard (e.g., *project* or org/*)") + "\n"
	s += renderHelpBar([][]string{
		{"enter", "apply"},
		{"esc", "cancel"},
	})

	return "\n" + styleBox.Render(s) + "\n"
}

func (m *Model) viewRepoList() string {
	repos := m.getDisplayRepos()

	if m.err != nil {
		s := renderHeader("Error")
		s += renderErrorBanner(m.err.Error()) + "\n"
		s += renderHelpBar([][]string{{"r", "retry"}, {"q", "quit"}})
		return "\n" + styleBox.Render(s) + "\n"
	}

	if len(repos) == 0 {
		dateStr := entity.FormatDateDisplay(m.startDate, m.endDate)
		s := renderHeader("No Commits Found")
		s += styleFooter.Render("No commits found for "+dateStr) + "\n"
		s += renderHelpBar([][]string{{"r", "change date"}, {"q", "quit"}})
		return "\n" + styleBox.Render(s) + "\n"
	}

	dateDisplay := entity.FormatDateDisplay(m.startDate, m.endDate)
	s := renderHeader("Commits for " + dateDisplay)

	// List header with count.
	totalCommits := 0
	for _, repo := range repos {
		totalCommits += len(m.commits[repo])
	}
	s += renderListHeader("Repositories", len(repos)) + "  " +
		styleFooter.Render(fmt.Sprintf("• %d commits total", totalCommits)) + "\n"
	s += renderDivider(50) + "\n\n"

	// Show filter if active.
	if m.filterActive && m.filterInput.Value() != "" {
		s += styleFooter.Render("Filter: "+m.filterInput.Value()) + "\n\n"
	}
	if m.warning != "" {
		s += renderWarningBanner(m.warning) + "\n\n"
	}

	for i, repo := range repos {
		checkbox := styleCheckboxUnchecked.Render(iconUncheckBox)
		if m.selected[repo] {
			checkbox = styleCheckbox.Render(iconCheckBox)
		}

		commitCount := styleFooter.Render(fmt.Sprintf(" (%d)", len(m.commits[repo])))

		if i == m.cursor {
			s += styleCursor.Render(iconArrowRight) + checkbox + " " + styleRepo.Render(repo) + commitCount + "\n"
		} else {
			s += "  " + checkbox + " " + styleRepo.Render(repo) + commitCount + "\n"
		}

		if m.selected[repo] {
			for _, commit := range m.commits[repo] {
				s += "     " + styleHighlight.Render(iconCommit) + " " + styleCommit.Render(commit.Message) + "\n"
			}
		}
	}

	s += renderHelpBar([][]string{
		{"space", "select"},
		{"a/n", "all/none"},
		{"f", "filter"},
		{"enter", "summary"},
		{"q", "quit"},
	})
	return "\n" + styleBox.Render(s) + "\n"
}

func (m *Model) viewExport() string {
	s := renderHeader("Export Summary")
	s += styleDateLabel.Render("Select export format:") + "\n\n"

	formats := []struct {
		name string
		desc string
	}{
		{"Text", "Plain text format (.txt)"},
		{"Markdown", "Markdown format (.md)"},
		{"JSON", "JSON format (.json)"},
	}

	for i, f := range formats {
		cursor := "  "
		if i == m.exportFormat {
			cursor = styleCursor.Render(iconArrowRight)
		}
		s += cursor + styleRepo.Render(f.name) + " " + styleFooter.Render(f.desc) + "\n"
	}

	if m.message != "" {
		s += "\n" + renderSuccessBanner(m.message) + "\n"
	}

	s += renderHelpBar([][]string{
		{"enter", "save file"},
		{"c", "copy"},
		{"b", "back"},
	})

	return "\n" + styleBox.Render(s) + "\n"
}

func (m *Model) viewLoading() string {
	dateStr := entity.FormatDateDisplay(m.startDate, m.endDate)

	s := renderHeader("Loading")
	s += m.spinner.View() + " " + styleDateLabel.Render("Fetching commits for "+dateStr+"...") + "\n\n"
	s += styleFooter.Render("Connecting to GitHub API") + "\n"
	s += renderHelpBar([][]string{
		{"esc", "cancel"},
		{"q", "quit"},
	})

	return "\n" + styleBox.Render(s) + "\n"
}

func (m *Model) viewStats() string {
	s := renderHeader("Statistics")

	if m.stats == nil {
		s += styleFooter.Render("No statistics available") + "\n"
		s += renderHelpBar([][]string{
			{"b", "back"},
			{"q", "quit"},
		})
		return "\n" + styleBox.Render(s) + "\n"
	}

	stats := m.stats

	// Summary stats with nice formatting.
	s += styleStatsLabel.Render("Total Commits:      ") + styleStatsValue.Render(fmt.Sprintf("%d", stats.TotalCommits)) + "\n"
	s += styleStatsLabel.Render("Total Repositories: ") + styleStatsValue.Render(fmt.Sprintf("%d", stats.TotalRepositories)) + "\n"

	if stats.MostActiveRepo != "" {
		s += styleStatsLabel.Render("Most Active:        ") + styleStatsValue.Render(stats.MostActiveRepo) +
			styleFooter.Render(fmt.Sprintf(" (%d commits)", stats.MaxCommits)) + "\n"
	}

	s += "\n" + renderDivider(50) + "\n\n"
	s += styleDateLabel.Render("Commits per Repository:") + "\n\n"

	// Find max value for proportional bars.
	maxCommits := findMaxCommits(stats.CommitsPerRepo)
	barWidth := 25 // Width of the progress bar.

	// Find the longest repo name for alignment.
	maxRepoLen := 0
	for repo := range stats.CommitsPerRepo {
		if len(repo) > maxRepoLen {
			maxRepoLen = len(repo)
		}
	}

	for repo, count := range stats.CommitsPerRepo {
		// Pad repo name for alignment.
		paddedRepo := repo
		for len(paddedRepo) < maxRepoLen {
			paddedRepo += " "
		}

		// Calculate percentage.
		pct := 0
		if stats.TotalCommits > 0 {
			pct = (count * 100) / stats.TotalCommits
		}

		bar := renderProgressBar(count, maxCommits, barWidth)
		s += "  " + styleRepo.Render(paddedRepo) + " " + bar + " " +
			styleStatsValue.Render(fmt.Sprintf("%2d", count)) + " " +
			styleFooter.Render(fmt.Sprintf("(%2d%%)", pct)) + "\n"
	}

	s += renderHelpBar([][]string{
		{"b", "back"},
		{"q", "quit"},
	})

	return "\n" + styleBox.Render(s) + "\n"
}

func (m *Model) viewSummary() string {
	dateStr := entity.FormatDateDisplay(m.startDate, m.endDate)
	s := renderHeader("Summary for " + dateStr)

	hasSelection := false
	repos := m.commitUC.GetSelectedReposSorted(m.commits, m.selected)

	if len(repos) > 0 {
		// Count total commits.
		totalCommits := 0
		for _, repo := range repos {
			totalCommits += len(m.commits[repo])
		}
		s += renderListHeader("Selected repositories", len(repos)) + "  " +
			styleFooter.Render(fmt.Sprintf("• %d commits", totalCommits)) + "\n"
		s += renderDivider(50) + "\n\n"
	}

	for _, repo := range repos {
		repoCommits := m.commits[repo]
		hasSelection = true
		s += styleRepo.Render("▸ "+repo) + "\n"

		for _, commit := range repoCommits {
			s += "  " + styleHighlight.Render(iconCommit) + " " + styleCommit.Render(commit.Message) + "\n"
		}
		s += "\n"
	}

	if !hasSelection {
		s += styleFooter.Render("No repositories selected.") + "\n\n"
	}

	if m.message != "" {
		s += renderSuccessBanner(m.message) + "\n"
	}

	s += renderHelpBar([][]string{
		{"c", "copy"},
		{"e", "export"},
		{"s", "stats"},
		{"b", "back"},
		{"q", "quit"},
	})

	return "\n" + styleBox.Render(s) + "\n"
}
