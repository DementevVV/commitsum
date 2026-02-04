package ui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/DementevVV/commitsum/internal/domain/entity"
)

// Update handles all user interactions and state changes.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Global quit.
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

		// Clear message on any key.
		m.message = ""
	}

	switch m.screen {
	case screenDateRange:
		return m.updateDateRange(msg)
	case screenDateSelect:
		return m.updateDateSelect(msg)
	case screenRepoFilter:
		return m.updateRepoFilter(msg)
	case screenRepoList:
		return m.updateRepoList(msg)
	case screenSummary:
		return m.updateSummary(msg)
	case screenExport:
		return m.updateExport(msg)
	case screenStats:
		return m.updateStats(msg)
	case screenLoading:
		return m.updateLoading(msg)
	}

	return m, nil
}

func (m *Model) updateDateRange(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			return m, tea.Quit
		case "j", "down":
			if m.dateRangeIdx < len(entity.DateRangePresets)-1 {
				m.dateRangeIdx++
			}
		case "k", "up":
			if m.dateRangeIdx > 0 {
				m.dateRangeIdx--
			}
		case "enter":
			preset := entity.DateRangePresets[m.dateRangeIdx].Key
			if preset == "custom" {
				m.err = nil
				m.screen = screenDateSelect
				m.dateInput.Focus()
				return m, textinput.Blink
			}
			dr := entity.GetDateRange(preset)
			m.startDate = dr.StartDate
			m.endDate = dr.EndDate
			return m.loadCommits()
		}
	}
	return m, nil
}

func (m *Model) updateDateSelect(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			dateValue := m.dateInput.Value()
			parsedDate, err := time.Parse("2006-01-02", dateValue)
			if err != nil {
				m.err = fmt.Errorf("invalid date format, please use YYYY-MM-DD")
				return m, nil
			}

			// Check that the date is not in the future.
			now := time.Now()
			today := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())
			if parsedDate.After(today) {
				m.err = fmt.Errorf("date cannot be in the future")
				return m, nil
			}

			m.startDate = dateValue
			m.endDate = dateValue
			m.err = nil
			return m.loadCommits()
		case tea.KeyEsc:
			m.err = nil
			m.screen = screenDateRange
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.dateInput, cmd = m.dateInput.Update(msg)
	return m, cmd
}

func (m *Model) updateRepoFilter(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			pattern := m.filterInput.Value()
			if pattern == "" {
				m.filterActive = false
				m.filteredRepos = m.repoList
			} else {
				m.filterActive = true
				m.filteredRepos = m.commitUC.FilterReposByPattern(m.repoList, pattern)
			}
			m.cursor = 0
			m.screen = screenRepoList
			return m, nil
		case tea.KeyEsc:
			m.filterActive = false
			m.filterInput.SetValue("")
			m.filteredRepos = m.repoList
			m.screen = screenRepoList
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.filterInput, cmd = m.filterInput.Update(msg)
	return m, cmd
}

func (m *Model) updateRepoList(msg tea.Msg) (tea.Model, tea.Cmd) {
	repos := m.getDisplayRepos()

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "enter":
			m.screen = screenSummary
			m.stats = m.commitUC.CalculateStatistics(m.commits, m.selected)
		case " ":
			if len(repos) > 0 {
				currentRepo := repos[m.cursor]
				m.selected[currentRepo] = !m.selected[currentRepo]
			}
		case "j", "down":
			if m.cursor < len(repos)-1 {
				m.cursor++
			}
		case "k", "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "a":
			// Select all.
			for _, repo := range repos {
				m.selected[repo] = true
			}
		case "n":
			// Select none.
			for _, repo := range repos {
				m.selected[repo] = false
			}
		case "f", "/":
			m.screen = screenRepoFilter
			m.filterInput.Focus()
			return m, textinput.Blink
		case "s":
			// Stats.
			m.stats = m.commitUC.CalculateStatistics(m.commits, m.selected)
			m.screen = screenStats
		case "r":
			// Refresh - go back to date selection.
			m.err = nil
			m.screen = screenDateRange
			m.cursor = 0
		}
	}
	return m, nil
}

func (m *Model) updateSummary(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "esc", "b":
			m.screen = screenRepoList
		case "c":
			content, err := m.generateExportContent(entity.FormatText)
			if err != nil {
				m.message = "Failed to generate content: " + err.Error()
			} else if err := m.clipboard.Copy(content); err != nil {
				m.message = "Failed to copy: " + err.Error()
			} else {
				m.message = "Copied to clipboard!"
			}
		case "e":
			m.screen = screenExport
			m.exportFormat = 0
		case "s":
			m.stats = m.commitUC.CalculateStatistics(m.commits, m.selected)
			m.screen = screenStats
		}
	}
	return m, nil
}

func (m *Model) updateExport(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "esc", "b":
			m.screen = screenSummary
		case "j", "down":
			if m.exportFormat < len(m.exportFormats)-1 {
				m.exportFormat++
			}
		case "k", "up":
			if m.exportFormat > 0 {
				m.exportFormat--
			}
		case "enter":
			format := entity.ExportFormat(m.exportFormats[m.exportFormat])
			content, err := m.generateExportContent(format)
			if err != nil {
				m.message = "Failed to generate content: " + err.Error()
				m.screen = screenSummary
				return m, nil
			}

			filename := m.exportUC.GenerateFilename(m.startDate, format)

			if err := m.exportUC.SaveToFile(content, filename); err != nil {
				m.message = "Failed to save: " + err.Error()
			} else {
				m.message = "Saved to " + filename
			}
			m.screen = screenSummary
		case "c":
			format := entity.ExportFormat(m.exportFormats[m.exportFormat])
			content, err := m.generateExportContent(format)
			if err != nil {
				m.message = "Failed to generate content: " + err.Error()
			} else if err := m.clipboard.Copy(content); err != nil {
				m.message = "Failed to copy: " + err.Error()
			} else {
				m.message = "Copied to clipboard!"
			}
		}
	}
	return m, nil
}

func (m *Model) updateStats(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "esc", "b":
			m.screen = screenRepoList
		}
	}
	return m, nil
}

func (m *Model) loadCommits() (*Model, tea.Cmd) {
	m.loading = true
	m.screen = screenLoading
	m.err = nil

	return m, tea.Batch(
		m.spinner.Tick,
		func() tea.Msg {
			data, err := m.commitUC.GetCommitsForRange(m.startDate, m.endDate)
			if err != nil {
				return commitsLoadedMsg{err: err}
			}
			return commitsLoadedMsg{
				commits:  data.Commits,
				repoList: data.RepoList,
				warning:  data.Warning,
				err:      nil,
			}
		},
	)
}

func (m *Model) updateLoading(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case commitsLoadedMsg:
		m.loading = false
		m.commits = msg.commits
		m.repoList = msg.repoList
		m.warning = msg.warning
		if m.filterInput.Value() != "" {
			m.filterActive = true
			m.filteredRepos = m.commitUC.FilterReposByPattern(m.repoList, m.filterInput.Value())
		} else {
			m.filterActive = false
			m.filteredRepos = msg.repoList
		}
		m.err = msg.err
		m.screen = screenRepoList
		m.cursor = 0
		return m, nil
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "esc":
			// Cancel loading, return to date range selection.
			m.loading = false
			m.err = nil
			m.screen = screenDateRange
			return m, nil
		}
	}
	return m, nil
}
