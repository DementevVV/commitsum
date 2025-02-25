# GitHub Commit Summarizer

GitHub Commit Summarizer is a terminal-based user interface (TUI) application for summarizing GitHub commits. It allows users to view, select, and generate formatted summaries of their GitHub commits for the current day. The application is built using the Bubble Tea framework and integrates with the GitHub CLI (gh) for fetching commit information.

## Features

- View today's commits across multiple repositories
- Interactively select repositories of interest
- Generate and copy formatted summaries of selected commits
- Navigate through repositories using keyboard controls
- Supports clipboard integration for summary copying

## Requirements

- **GitHub CLI (gh) must be installed and authenticated**
- Terminal with ANSI color support

## Installation

1. Clone the repository

2. Install dependencies:

   ```sh
   go mod tidy
   ```

3. Build the application:

   ```sh
   go build -o commitsum
   ```

4. Run the application:
   ```sh
   ./commitsum
   ```

## Usage

Use the following keyboard controls to interact with the application:

- `Space`: Toggle repository selection
- `j` or `↓`: Move cursor down
- `k` or `↑`: Move cursor up
- `Enter`: Generate summary
- `c`: Copy summary to clipboard
- `q`: Quit application

## Example

When you run the application, you will see a list of repositories with today's commits. Use the keyboard controls to navigate and select repositories. Once you have selected the repositories of interest, press `Enter` to generate a summary. You can then press `c` to copy the summary to your clipboard.

## Code Overview

### Main Components

- **Model**: Holds the application state, including commits, repository list, cursor position, selection state, and error state.
- **getGitHubUser**: Retrieves the authenticated GitHub username using the GitHub CLI.
- **getCommits**: Fetches commits from the GitHub API and groups them by repository.
- **generateSummary**: Creates a formatted string containing a summary of today's commits from selected repositories.
- **generatePlainSummary**: Creates a plain text summary of selected repository commits.
- **View**: Renders the terminal user interface for the commit summary tool.
- **Update**: Handles UI interactions and keyboard events.

### Styles

The application uses the [lipgloss](http://_vscodecontentref_/0) package for styling the UI elements, including titles, repository names, commit messages, cursor, footer, and checkboxes.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request if you have any improvements or bug fixes.

## License

This project is licensed under the MIT License. See the LICENSE file for details.
