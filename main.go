// Package main provides a TUI application for summarizing GitHub commits.
//
// The application allows users to:
//   - View today's commits across multiple repositories
//   - Interactively select repositories of interest
//   - Generate and copy formatted summaries of selected commits
//   - Navigate through repositories using keyboard controls
//
// Core Features:
//   - Fetches commits using GitHub CLI (gh) authentication
//   - Groups commits by repository
//   - Provides an interactive TUI using Bubble Tea
//   - Supports clipboard integration for summary copying
//
// Usage:
//
//	Use keyboard controls:
//	- Space: Toggle repository selection
//	- j/↓: Move cursor down
//	- k/↑: Move cursor up
//	- Enter: Generate summary
//	- c: Copy summary to clipboard
//	- q: Quit application
//
// Requirements:
//   - GitHub CLI (gh) must be installed and authenticated
//   - Terminal with ANSI color support
//
// The application uses the Bubble Tea framework for the terminal user interface
// and the GitHub API through the gh CLI tool to fetch commit information.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Commit stores repository name and commit message information.
type Commit struct {
	Repo    string
	Message string
}

// APIResponse represents the GitHub API response format.
// APIResponse represents the structure of a GitHub API response for commit searches.
// It contains a list of commit items, where each item includes repository information
// and commit details. This structure is designed to unmarshal JSON responses from
// GitHub's search API endpoint.
//
// The Items field holds an array of commit search results, each containing:
//   - Repository information (including the full repository name)
//   - Commit details (including the commit message)
type APIResponse struct {
	Items []struct {
		Repository struct {
			FullName string `json:"full_name"`
		} `json:"repository"`
		Commit struct {
			Message string `json:"message"`
		} `json:"commit"`
	} `json:"items"`
}

// Model holds the application state.
// model represents the application state for commit summarization.
// It manages repository commits, selection state, and UI display preferences.
//
// Fields:
//   - commits: Repository-grouped commit data where key is repo name and value is commit list
//   - repoList: Maintains order of repository names for consistent display
//   - cursor: Tracks current selection position in repository list
//   - selected: Repository selection state mapping where true indicates selected
//   - showSummary: Controls visibility of commit summary view
//   - err: Stores current error state of the application
//   - copyRequested: Indicates if summary should be copied to system clipboard
type model struct {
	commits       map[string][]Commit
	repoList      []string
	cursor        int
	selected      map[string]bool
	showSummary   bool
	err           error
	copyRequested bool
}

// getGitHubUser retrieves the authenticated GitHub username.
// getGitHubUser retrieves the currently authenticated GitHub username using the GitHub CLI.
// It executes the 'gh api user' command and extracts the login name from the response.
// Returns the username as a string and any error encountered during the process.
// Requires the GitHub CLI (gh) to be installed and authenticated.
func getGitHubUser() (string, error) {
	cmd := exec.Command("gh", "api", "user", "--jq", ".login")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// getCommits fetches commits from GitHub API and groups them by repository.
// getCommits retrieves GitHub commits made by the authenticated user for the current day.
// It uses the GitHub CLI to fetch commits via the GitHub API, and organizes them by repository.
//
// Returns:
//   - map[string][]Commit: A map of repository names to their commits
//   - []string: A sorted list of repository names that have commits
//   - error: Any error encountered during the process
//
// The function performs the following steps:
//  1. Gets the authenticated GitHub username
//  2. Searches for commits made today by the user
//  3. Parses the API response and organizes commits by repository
//  4. Returns both the commit map and a sorted list of repositories
func getCommits() (map[string][]Commit, []string, error) {
	ghUser, err := getGitHubUser()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get GitHub user: %w", err)
	}

	today := time.Now().Format("2006-01-02")
	query := fmt.Sprintf("author:%s+committer-date:%s", ghUser, today)

	cmd := exec.Command("gh", "api", "/search/commits?q="+query, "--header", "Accept: application/vnd.github.cloak-preview", "--paginate")
	out, err := cmd.Output()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch commits: %w", err)
	}

	var resp APIResponse
	if err := json.Unmarshal(out, &resp); err != nil {
		return nil, nil, err
	}

	commitMap := make(map[string][]Commit)
	for _, item := range resp.Items {
		repo := item.Repository.FullName
		message := strings.Split(item.Commit.Message, "\n")[0]
		commitMap[repo] = append(commitMap[repo], Commit{Repo: repo, Message: message})
	}

	var repoList []string
	for repo := range commitMap {
		repoList = append(repoList, repo)
	}
	sort.Strings(repoList)

	return commitMap, repoList, nil
}

// initModel initializes the application state.
// initModel initializes and returns a new model instance with default values.
// It retrieves commits and repository list data, and initializes an empty selection map.
// The model is used to track the application state including commit history,
// available repositories, selected items, and UI display settings.
func initModel() model {
	commits, repoList, err := getCommits()
	return model{
		commits:     commits,
		repoList:    repoList,
		err:         err,
		selected:    make(map[string]bool),
		showSummary: false,
	}
}

// Init initializes Bubble Tea.
// Init implements tea.Model interface and initializes the model.
// Returns a tea.Cmd which in this case is nil as no initial commands are needed.
func (m model) Init() tea.Cmd {
	return nil
}

// Update handles UI interactions and keyboard events.
// It processes user input and updates the application state accordingly.
//
// Key bindings:
//   - q: Quit the application
//   - enter: Switch to summary view
//   - space: Toggle selection of current repository
//   - j/down: Move cursor down through repository list
//   - k/up: Move cursor up through repository list
//   - c: Copy summary to clipboard (when in summary view)
//
// Returns:
//   - tea.Model: Updated application model
//   - tea.Cmd: Command to be executed by the Bubble Tea framework
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "enter":
			m.showSummary = true
		case " ":
			currentRepo := m.repoList[m.cursor]
			m.selected[currentRepo] = !m.selected[currentRepo]
		case "j", "down":
			if m.cursor < len(m.repoList)-1 {
				m.cursor++
			}
		case "k", "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "c":
			if m.showSummary {
				cmd := exec.Command("pbcopy")
				cmd.Stdin = strings.NewReader(m.generatePlainSummary())
				return m, tea.ExecProcess(cmd, nil)
			}
		}
	}
	return m, nil
}

// generateSummary creates the summary view inside Bubble Tea.
// generateSummary creates a formatted string containing a summary of today's commits
// from selected repositories. The summary includes:
//   - A title header with emoji
//   - For each selected repository:
//   - Repository name
//   - Bullet points listing commit messages
//   - A footer with available actions
//
// If no repositories are selected, it shows a notification message instead.
//
// Returns a formatted string containing the complete summary.
func (m model) generateSummary() string {
	output := titleStyle.Render("📌 Today's Commits Summary") + "\n\n"

	for repo, repoCommits := range m.commits {
		if m.selected[repo] {
			output += repoStyle.Render(fmt.Sprintf("[%s]:", repo)) + "\n"

			for _, commit := range repoCommits {
				output += commitStyle.Render(fmt.Sprintf("  • %s", commit.Message)) + "\n"
			}
			output += "\n"
		}
	}

	if output == titleStyle.Render("📌 Today's Commits Summary")+"\n\n" {
		output += footerStyle.Render("No repositories selected.\nPress [q] to quit.") + "\n"
	} else {
		output += footerStyle.Render("\nPress [q] to quit, [c] to copy.") + "\n"
	}

	return output
}

// generatePlainSummary creates a plain text summary of selected repository commits.
// It formats the output with repository names as headers followed by commit messages
// as bullet points. Only repositories that are marked as selected in the model
// will be included in the summary.
//
// Returns a string containing the formatted commit summary.
func (m model) generatePlainSummary() string {
	var output strings.Builder
	output.WriteString("Today's Commits Summary\n\n")

	for repo, repoCommits := range m.commits {
		if m.selected[repo] {
			output.WriteString(fmt.Sprintf("[%s]:\n", repo))
			for _, commit := range repoCommits {
				output.WriteString(fmt.Sprintf("  • %s\n", commit.Message))
			}
			output.WriteString("\n")
		}
	}
	return output.String()
}

// Styles
// Style definitions for the commit summary UI:
//
//	titleStyle    - Orange bold text for titles
//	repoStyle     - Cyan bold text for repository names
//	commitStyle   - White text for commit messages
//	cursorStyle   - Red-orange bold text for cursor/selection
//	footerStyle   - Gray text for footer information
//	highlight     - Green text for highlighting
//	checkboxStyle - Light green text for checkboxes
//	uncheckedMark - Circle symbol (○) for unchecked items
//	checkedMark   - Filled circle symbol (●) for checked items
var (
	titleStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFA500"))
	repoStyle     = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00FFFF"))
	commitStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))
	cursorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF4500")).Bold(true)
	footerStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#777777"))
	highlight     = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00"))
	checkboxStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#98c379"))
	uncheckedMark = "○"
	checkedMark   = "●"
)

// View renders the UI.
// View renders the terminal user interface for the commit summary tool.
// It handles multiple display states:
//   - Error state: Shows error message if any error occurred
//   - Summary state: Displays the generated commit summary when showSummary is true
//   - Empty state: Shows message when no commits are found
//   - Main view: Displays an interactive list of repositories with:
//   - Checkboxes to select/unselect repositories
//   - Cursor indicating current selection
//   - Nested commit messages for selected repositories
//   - Footer with navigation instructions
//
// Returns a styled string representing the complete UI view.
func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("%s\n%s", titleStyle.Render("Error"), m.err.Error())
	}

	if m.showSummary {
		return m.generateSummary()
	}

	if len(m.repoList) == 0 {
		return titleStyle.Render("No commits found for today.\nPress q to quit.")
	}

	s := titleStyle.Render("📌 Today's GitHub Commits") + "\n\n"

	for i, repo := range m.repoList {
		checkbox := checkboxStyle.Render(uncheckedMark)
		if m.selected[repo] {
			checkbox = checkboxStyle.Render(checkedMark)
		}

		if i == m.cursor {
			s += cursorStyle.Render("➡ ") + checkbox + " " + repoStyle.Render(repo) + "\n"
		} else {
			s += "   " + checkbox + " " + repoStyle.Render(repo) + "\n"
		}

		if m.selected[repo] {
			for _, commit := range m.commits[repo] {
				s += "    " + highlight.Render("• ") + commitStyle.Render(commit.Message) + "\n"
			}
		}
	}

	s += footerStyle.Render("\n[ space ] Select/Unselect  [ j / ↓ ] Down  [ k / ↑ ] Up  [ enter ] Show Summary  [ q ] Quit\n")
	return s
}

// main runs the Bubble Tea program.
// main initializes and starts a new Bubble Tea program using the initModel.
// If the program encounters an error during execution, it prints the error
// message to standard output and exits with status code 1.
func main() {
	p := tea.NewProgram(initModel())
	if err := p.Start(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
