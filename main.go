// Package main implements a GitHub commit summarization tool.
//
// CommitSum provides a terminal user interface for viewing, selecting,
// and summarizing GitHub commits from a specified date. It leverages
// the GitHub CLI for data retrieval and presents an interactive interface
// for repository selection and summary generation.
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

// Commit represents a repository commit with its message.
// It contains the repository identifier and the commit message content.
type Commit struct {
	Repo    string
	Message string
}

// APIResponse represents the structure of GitHub's search API response for commits.
// It contains a list of commit items, each with associated repository and commit message data.
//
// This struct is used for unmarshaling the JSON response from the GitHub API when
// searching for commits. The JSON tags match the API's response format.
//
// Fields:
//   - Items: An array of commit search results, each containing repository and commit data
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

// Model represents the application state for commit summarization.
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
//   - dateInput: Stores user input for date confirmation
//   - dateConfirmed: Tracks user confirmation of date input
type model struct {
	commits       map[string][]Commit
	repoList      []string
	cursor        int
	selected      map[string]bool
	showSummary   bool
	err           error
	copyRequested bool
	dateInput     string
	dateConfirmed bool
}

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

// getCommits fetches GitHub commits for a specific date from the authenticated user.
//
// The function queries GitHub's API for commits authored by the currently authenticated
// user on the specified date. It organizes these commits by repository and returns them
// in a structured format for display and interaction.
//
// Parameters:
//   - date: A string representing the date in "YYYY-MM-DD" format to search commits for
//
// Returns:
//   - map[string][]Commit: A map where keys are repository names and values are slices of commits
//   - []string: A sorted list of repository names for consistent display order
//   - error: Any error encountered during the API request or processing
//
// The function depends on the GitHub CLI (gh) being installed and properly authenticated.
// It uses the GitHub search API with the 'cloak-preview' header to access commit information.
func getCommits(date string) (map[string][]Commit, []string, error) {
	ghUser, err := getGitHubUser()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get GitHub user: %w", err)
	}

	query := fmt.Sprintf("author:%s+committer-date:%s", ghUser, date)

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

// initModel creates and initializes a new application model with default values.
//
// This function sets up the initial state for the commit summarization tool,
// including today's date as the default date input and empty data structures
// for repository selection and visualization.
//
// Returns:
//   - model: A fully initialized model instance with the following defaults:
//   - Today's date as the dateInput in YYYY-MM-DD format
//   - dateConfirmed set to false to show the date selection screen first
//   - An empty selection map for repositories
//   - showSummary set to false to start in repository selection view
//
// The returned model is ready to be used with the Bubble Tea framework.
func initModel() model {
	today := time.Now().Format("2006-01-02")
	return model{
		dateInput:     today,
		dateConfirmed: false,
		selected:      make(map[string]bool),
		showSummary:   false,
	}
}

// Init implements the Bubble Tea model interface initialization hook.
// This method is called once when the program starts and provides an
// opportunity to return initial commands for the Bubble Tea runtime.
//
// In this application, no initial commands are needed since the model
// is fully initialized in the initModel function and starts in the
// date selection state.
//
// Returns:
//   - tea.Cmd: nil, as no initial commands are required
func (m model) Init() tea.Cmd {
	return nil
}

// Update handles all user interactions and state changes in the application.
// This method processes keyboard input and updates the model state accordingly.
//
// The method handles two main application states:
//   - Date selection mode: When dateConfirmed is false, user enters the date
//     to fetch commits for. Validates input format and fetches commits on confirmation.
//   - Repository selection mode: When dateConfirmed is true, user can navigate the
//     repository list, select/deselect repos, view summaries, and copy content.
//
// Key bindings in date selection mode:
//   - Enter: Confirms date input, validates format, and fetches commits
//   - Backspace/Delete: Removes characters from date input
//   - Runes: Adds typed characters to date input
//
// Key bindings in repository selection mode:
//   - q: Quits the application
//   - Enter: Shows commit summary
//   - Space: Toggles selection of current repository
//   - j/down: Moves cursor down
//   - k/up: Moves cursor up
//   - c: Copies summary to clipboard (when summary is shown)
//
// Parameters:
//   - msg: The message to process, typically a keyboard input event
//
// Returns:
//   - tea.Model: The updated model after processing the message
//   - tea.Cmd: Any command to be executed by the Bubble Tea framework
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Date selection mode
	if !m.dateConfirmed {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyEnter:
				// Validate date format
				_, err := time.Parse("2006-01-02", m.dateInput)
				if err != nil {
					m.err = fmt.Errorf("invalid date format, please use YYYY-MM-DD")
					return m, nil
				}

				// Load commits with selected date
				m.dateConfirmed = true
				commits, repoList, err := getCommits(m.dateInput)
				m.commits = commits
				m.repoList = repoList
				m.err = err
				return m, nil

			case tea.KeyBackspace, tea.KeyDelete:
				if len(m.dateInput) > 0 {
					m.dateInput = m.dateInput[:len(m.dateInput)-1]
				}

			case tea.KeyRunes:
				m.dateInput += string(msg.Runes)
			}
			return m, nil
		}
	}

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

// generateSummary creates a formatted display string of selected repository commits.
//
// This method generates a stylized text representation of all commits from repositories
// that have been selected by the user. The output includes:
//   - A title header indicating these are today's commit summaries
//   - Repository names styled in cyan
//   - Commit messages with bullet points and proper indentation
//   - Footer text with available actions
//
// The method handles the empty selection case by showing a message that no repositories
// are selected, and includes appropriate navigation instructions in the footer.
//
// Returns:
//   - string: A fully formatted and styled summary of selected repository commits
//     ready to be displayed in the terminal UI
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

// generatePlainSummary creates an unstyled text summary of selected repository commits.
//
// This method produces a plain text representation of commits from selected repositories,
// without any terminal styling or color codes. It's primarily used for clipboard operations
// where ANSI escape sequences would be undesirable.
//
// The output includes:
//   - A simple title indicating these are today's commit summaries
//   - Repository names in square brackets
//   - Commit messages with bullet points and consistent indentation
//   - Proper spacing between repository sections
//
// Returns:
//   - string: An unformatted plain text summary suitable for clipboard use or
//     pasting into other applications
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

// View renders the current state of the application model as a styled string.
// This method produces the terminal user interface for the commit summary tool.
//
// The View method handles different application states:
//   - Error state: Displays error messages when operations fail
//   - Date selection state: Shows date input field with current value when dateConfirmed is false
//   - Summary state: Shows formatted commit summary when showSummary is true
//   - Selection state: Shows repository list with checkboxes and commit details when selected
//
// The repository selection view includes:
//   - A title header
//   - Repository names with selection checkboxes
//   - Commit messages for selected repositories
//   - Navigation and action instructions in the footer
//
// When no repositories are found, an appropriate message is displayed.
//
// Returns:
//   - string: A formatted string representing the current UI state,
//     including all styling and layout elements for terminal display
func (m model) View() string {
	if m.err != nil {
		return titleStyle.Render("Error") + "\n" + m.err.Error()
	}

	if !m.dateConfirmed {
		s := titleStyle.Render("GitHub Commit Summary") + "\n\n"
		s += "Enter date for commit summary (YYYY-MM-DD):\n"
		s += repoStyle.Render(m.dateInput) + "\n\n"
		s += footerStyle.Render("Press Enter to confirm, or edit date")
		return s
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

// main initializes and runs the commit summarization application.
//
// This function creates a new Bubble Tea program with the initial application model,
// starts the event loop, and handles any initialization errors that may occur.
//
// The function follows these steps:
//  1. Creates a new Bubble Tea program with the default model
//  2. Starts the program's event loop
//  3. Handles any errors by displaying them and exiting with a non-zero status code
//
// No parameters or return values as this is the program entry point.
func main() {
	p := tea.NewProgram(initModel())
	if err := p.Start(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
