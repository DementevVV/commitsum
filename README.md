# 📊 GitHub Commit Summarizer

![SSH Notification Mascot](docs/images/mascot.png)

![Go Version](https://img.shields.io/badge/Go-00ADD8?logo=Go&logoColor=white&style=for-the-badge)
![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)
![Platform](https://img.shields.io/badge/Platform-Linux%20%7C%20macOS%20%7C%20Windows-blue?style=for-the-badge)

Commitsum turns your GitHub commits into a clean, shareable summary in seconds. It’s a beautiful Go-powered CLI with a modern Bubble Tea TUI, plus local caching and detailed logs so you can move fast and troubleshoot quickly.

## 🎬 Demo

See the complete keyboard-driven flow: date selection → repo selection/filtering → summary → export.

![Commitsum Demo](docs/images/commitsum-demo.gif)

---

## ✨ Features

- 📅 **Flexible date selection** — Today, yesterday, last week, last month, or custom date
- 🔍 **Repository filtering** — Filter repos by pattern (e.g., `*project*` or `org/*`)
- 🎯 **Multi-repository support** — See all your commits across different repositories
- ✅ **Smart selection** — Select all, none, or individual repositories
- 📋 **One-click copy** — Cross-platform clipboard support (macOS, Linux, Windows)
- 📤 **Multiple export formats** — Export to Text, Markdown, or JSON
- 📊 **Commit statistics** — Visualize commits per repository with charts
- 🗂️ **Local caching** — Speeds up repeated queries with a short-lived cache
- 🧾 **Logs for debugging** — Daily log files stored locally
- ⚙️ **Configuration file** — Optional. You can create `~/.config/commitsum/config.json` manually to set defaults
- 🎨 **Modern terminal UI** — Beautiful interface with soft purple gradient theme
- ⌨️ **Keyboard navigation** — Efficient keyboard-driven workflow

## 🚀 Quick Start

### Prerequisites

- Go
- GitHub CLI (gh) must be installed and authenticated
- Terminal with ANSI color support

### Installation

#### Option 1: Download a prebuilt binary

Download the appropriate archive for your OS/CPU from the GitHub Releases page and extract it:

https://github.com/DementevVV/commitsum/releases

#### Option 1b: One-line install (macOS/Linux)

```bash
curl -fsSL https://raw.githubusercontent.com/DementevVV/commitsum/master/install.sh | sh
```

#### Option 1c: One-line install (Windows PowerShell)

```powershell
irm https://raw.githubusercontent.com/DementevVV/commitsum/master/install.ps1 | iex
```

#### Option 2: Build from source

```bash
# Clone the repository
git clone https://github.com/DementevVV/commitsum.git
cd commitsum

# Install dependencies
go mod tidy

# Build the binary
go build -o commitsum ./cmd/commitsum

# Run the application
./commitsum
```

### Alternative: Direct build and run

```bash
go run ./cmd/commitsum
```

### First Run

1. **Select time range** — Choose from presets or enter custom date
   - Today, Yesterday, Last 7 days, Last 30 days
   - Or enter a custom date (YYYY-MM-DD format)
2. **Review commits** — Browse your commits across all repositories
3. **Filter repositories** — Press `f` to filter by pattern (optional)
4. **Select repositories** — Use `space` to toggle, `a` for all, `n` for none
5. **Generate summary** — Press `Enter` to view the formatted summary
6. **Export or copy** — Press `c` to copy, `e` to export to file

That's it! You now have a beautiful summary of your day's work.

## 📖 Usage

### Date Range Selection

| Key        | Action            |
| ---------- | ----------------- |
| `j` or `↓` | Move cursor down  |
| `k` or `↑` | Move cursor up    |
| `enter`    | Select date range |
| `esc`      | Quit application  |
| `q`        | Quit application  |

### Repository Selection

| Key        | Action                     |
| ---------- | -------------------------- |
| `space`    | Select/unselect repository |
| `a`        | Select all repositories    |
| `n`        | Deselect all               |
| `f` or `/` | Filter by pattern          |
| `s`        | Show statistics            |
| `r`        | Change date range          |
| `j` or `↓` | Move cursor down           |
| `k` or `↑` | Move cursor up             |
| `enter`    | Show summary               |
| `q`        | Quit application           |

### Summary Screen

| Key   | Action            |
| ----- | ----------------- |
| `c`   | Copy to clipboard |
| `e`   | Export to file    |
| `s`   | Show statistics   |
| `b`   | Back to selection |
| `esc` | Back to selection |
| `q`   | Quit application  |

### Export Screen

| Key     | Action                  |
| ------- | ----------------------- |
| `enter` | Save to file            |
| `c`     | Copy in selected format |
| `b`     | Back to summary         |
| `esc`   | Back to summary         |
| `q`     | Quit application        |

## 📋 Export Formats

### Text Format (.txt)

```text
Commit Summary - 2026-02-02

[username/project-one]
  - Add new feature for user authentication
  - Fix bug in login flow

---
Statistics: 5 commits across 2 repositories
Most active: username/project-one (3 commits)
```

### Markdown Format (.md)

```markdown
# Commit Summary

**Date:** 2026-02-02

## Statistics

- **Total Commits:** 5
- **Repositories:** 2
- **Most Active:** username/project-one (3 commits)

## Commits

### username/project-one

- Add new feature for user authentication
- Fix bug in login flow

---

_Generated by commitsum on 2026-02-02 09:41:12_
```

### JSON Format (.json)

```json
{
  "date": "2026-02-02",
  "total_repos": 2,
  "total_commits": 5,
  "commits": {
    "username/project-one": [
      { "repository": "username/project-one", "message": "Add new feature" }
    ]
  },
  "stats": {
    "total_commits": 5,
    "total_repositories": 2,
    "most_active_repo": "username/project-one",
    "max_commits": 3,
    "commits_per_repo": { "username/project-one": 3 }
  },
  "generated_at": "2026-02-02T09:41:12Z"
}
```

## ⚙️ Configuration

Configuration is optional and is read from `~/.config/commitsum/config.json` if the file exists. You can create it manually:

```json
{
  "default_date_range": "today",
  "repo_filter": "",
  "output_format": "text",
  "custom_template": "",
  "auto_copy": false,
  "show_stats": true
}
```

| Option               | Description                                                                  |
| -------------------- | ---------------------------------------------------------------------------- |
| `default_date_range` | Default preset: `today`, `yesterday`, `week`, `month` _(reserved for UI)_    |
| `repo_filter`        | Default repository filter pattern (pre-fills the filter input)               |
| `output_format`      | Default export format: `text`, `markdown`, `json` _(reserved for export UI)_ |
| `custom_template`    | Custom template for exports _(use case available, UI pending)_               |
| `auto_copy`          | Automatically copy summary to clipboard _(reserved for UI)_                  |
| `show_stats`         | Show statistics in summaries _(reserved for UI)_                             |

## 🔧 Development

### Building from Source

```bash
# Build for current platform
go build -o commitsum ./cmd/commitsum

# Build for all platforms
GOOS=linux GOARCH=amd64 go build -o commitsum-linux-amd64 ./cmd/commitsum
GOOS=darwin GOARCH=amd64 go build -o commitsum-darwin-amd64 ./cmd/commitsum
GOOS=darwin GOARCH=arm64 go build -o commitsum-darwin-arm64 ./cmd/commitsum
```

### Project Structure

```text
commitsum/
├── cmd/
│   └── commitsum/         # Application entry point
├── internal/
│   ├── domain/            # Entities and domain contracts
│   ├── infrastructure/    # GitHub client, config, cache, clipboard, logger
│   ├── ui/                # Bubble Tea UI state, views, and styles
│   └── usecase/           # Business logic (commits + export)
├── docs/                  # Images and docs assets
├── go.mod                 # Go dependencies
├── go.sum                 # Dependency checksums
├── Makefile               # Build automation
├── README.md              # This file
└── LICENSE                # MIT License
```

## 🔍 How It Works

1. **GitHub CLI Integration** — Uses `gh` CLI to authenticate and fetch commit data
2. **GitHub Search API** — Queries commits by author and date using GitHub's search API (up to 1000 results)
3. **Local Cache** — Stores short-lived results in `~/.config/commitsum/cache` for faster repeat runs
4. **Interactive UI** — Bubble Tea framework provides the terminal user interface
5. **Lipgloss Styling** — Modern terminal styling with soft purple/violet gradient theme
6. **Logs** — Writes daily logs to `~/.config/commitsum/logs`

## 🛠️ Troubleshooting

### GitHub CLI not authenticated

```bash
# Login to GitHub CLI
gh auth login
```

### No commits found

- Ensure you have commits on the selected date
- Verify GitHub CLI has access to your repositories
- Check that your commits are authored with the correct GitHub email

### Date format errors

- Use the format YYYY-MM-DD (e.g., 2026-02-02)
- Year must be 4 digits, month and day must be 2 digits

### Clipboard not working (Linux)

- Install one of: `xclip`, `xsel`, or `wl-copy`

### Need more details

- Logs are written to `~/.config/commitsum/logs`
- Set `DEBUG=1` to also print logs to stderr

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

### Development Setup

1. Fork the repository
2. Clone your fork: `git clone https://github.com/YOUR_USERNAME/commitsum.git`
3. Create a feature branch: `git checkout -b feature/amazing-feature`
4. Make your changes and test thoroughly
5. Commit your changes: `git commit -m 'Add amazing feature'`
6. Push to the branch: `git push origin feature/amazing-feature`
7. Open a Pull Request

Please ensure your code:

- Follows Go best practices and conventions
- Includes comments for exported functions
- Is properly formatted with `gofmt`

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) — The TUI framework powering the interface
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) — Terminal styling and layout
- [Bubbles](https://github.com/charmbracelet/bubbles) — Reusable TUI components (textinput)
- [GitHub CLI](https://cli.github.com/) — GitHub integration and authentication

## 📞 Support

If you have any questions or issues, please:

1. Check existing [Issues](https://github.com/DementevVV/commitsum/issues)
2. Create a new issue with details (OS, Go version, error messages)
3. Provide steps to reproduce any problems

## 🚀 Future Enhancements

- [x] ~~Cross-platform clipboard support (Linux, Windows)~~
- [x] ~~Export summaries to file (Markdown, JSON)~~
- [x] ~~Date range selection (e.g., last week, last month)~~
- [x] ~~Filter commits by repository pattern~~
- [x] ~~Commit statistics and visualization~~
- [x] ~~Configuration file support~~
- [ ] Template-based exports (use case available, UI pending)
- [ ] Git integration (local repository commits)
- [ ] Multiple GitHub accounts support
- [ ] Interactive commit message editing
- [ ] Slack/Discord integration
- [ ] Daily/weekly digest scheduling

---

Made with ❤️ by [DementevVV](https://github.com/DementevVV)
