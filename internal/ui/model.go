package ui

import (
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/DementevVV/commitsum/internal/domain/entity"
	"github.com/DementevVV/commitsum/internal/domain/repository"
	"github.com/DementevVV/commitsum/internal/infrastructure/config"
	"github.com/DementevVV/commitsum/internal/usecase"
)

// Screen states.
type screenState int

const (
	screenDateRange screenState = iota
	screenDateSelect
	screenRepoFilter
	screenRepoList
	screenSummary
	screenExport
	screenStats
	screenLoading
)

// Model represents the application state for the TUI.
type Model struct {
	// Data.
	commits       map[string][]entity.Commit
	repoList      []string
	filteredRepos []string

	// Selection state.
	cursor   int
	selected map[string]bool

	// Screen state.
	screen screenState

	// Inputs.
	dateInput    textinput.Model
	filterInput  textinput.Model
	spinner      spinner.Model
	filterActive bool

	// Date range.
	dateRangeIdx int
	startDate    string
	endDate      string

	// Export.
	exportFormat  int
	exportFormats []string

	// Config & Stats.
	config config.Config
	stats  *entity.Statistics

	// Use cases.
	commitUC  *usecase.CommitUseCase
	exportUC  *usecase.ExportUseCase
	clipboard repository.ClipboardRepository

	// Status.
	err     error
	message string
	warning string
	loading bool
}

// commitsLoadedMsg is sent when commits finish loading.
type commitsLoadedMsg struct {
	commits  map[string][]entity.Commit
	repoList []string
	warning  string
	err      error
}

// NewModel creates and initializes a new UI model.
func NewModel(cfg config.Config, commitUC *usecase.CommitUseCase, exportUC *usecase.ExportUseCase, clipboard repository.ClipboardRepository) *Model {
	today := time.Now().Format("2006-01-02")

	// Initialize date text input.
	ti := textinput.New()
	ti.Placeholder = "YYYY-MM-DD"
	ti.Focus()
	ti.CharLimit = 10
	ti.Width = 20
	ti.SetValue(today)
	ti.Prompt = ""
	ti.PromptStyle = lipgloss.NewStyle().Foreground(colorPrimaryLight)
	ti.TextStyle = lipgloss.NewStyle().Foreground(colorPrimary).Bold(true)
	ti.PlaceholderStyle = lipgloss.NewStyle().Foreground(colorTextMuted)
	ti.Cursor.Style = lipgloss.NewStyle().Foreground(colorAccent)

	// Initialize filter text input.
	fi := textinput.New()
	fi.Placeholder = "e.g., *project* or org/*"
	fi.CharLimit = 50
	fi.Width = 30
	fi.Prompt = ""
	fi.PromptStyle = lipgloss.NewStyle().Foreground(colorPrimaryLight)
	fi.TextStyle = lipgloss.NewStyle().Foreground(colorPrimary)
	fi.PlaceholderStyle = lipgloss.NewStyle().Foreground(colorTextMuted)
	fi.Cursor.Style = lipgloss.NewStyle().Foreground(colorAccent)
	if cfg.RepoFilter != "" {
		fi.SetValue(cfg.RepoFilter)
	}

	// Initialize spinner.
	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(colorPrimary)

	return &Model{
		dateInput:     ti,
		filterInput:   fi,
		spinner:       sp,
		screen:        screenDateRange,
		selected:      make(map[string]bool),
		config:        cfg,
		exportFormats: []string{"text", "markdown", "json"},
		startDate:     today,
		endDate:       today,
		commitUC:      commitUC,
		exportUC:      exportUC,
		clipboard:     clipboard,
	}
}

// Init implements the Bubble Tea model interface.
func (m *Model) Init() tea.Cmd {
	return textinput.Blink
}

// getDisplayRepos returns the repos to display based on filter state.
func (m *Model) getDisplayRepos() []string {
	if m.filterActive {
		return m.filteredRepos
	}
	return m.repoList
}

// generateExportContent generates content for export.
func (m *Model) generateExportContent(format entity.ExportFormat) (string, error) {
	dateStr := entity.FormatDateDisplay(m.startDate, m.endDate)
	stats := m.commitUC.CalculateStatistics(m.commits, m.selected)

	switch format {
	case entity.FormatMarkdown:
		return m.exportUC.ExportToMarkdown(m.commits, m.selected, dateStr, stats), nil
	case entity.FormatJSON:
		return m.exportUC.ExportToJSON(m.commits, m.selected, dateStr, stats)
	default:
		return m.exportUC.ExportToText(m.commits, m.selected, dateStr, stats), nil
	}
}
