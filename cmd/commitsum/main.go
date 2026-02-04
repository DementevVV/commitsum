// Package main is the entry point for the CommitSum CLI application.
//
// CommitSum provides a terminal user interface for viewing, selecting,
// and summarizing GitHub commits from a specified date range.
package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/DementevVV/commitsum/internal/infrastructure/cache"
	"github.com/DementevVV/commitsum/internal/infrastructure/clipboard"
	"github.com/DementevVV/commitsum/internal/infrastructure/config"
	"github.com/DementevVV/commitsum/internal/infrastructure/github"
	"github.com/DementevVV/commitsum/internal/infrastructure/logger"
	"github.com/DementevVV/commitsum/internal/ui"
	"github.com/DementevVV/commitsum/internal/usecase"
)

// Version information (set during build).
var (
	Version   = "dev"
	BuildTime = "unknown"
)

func main() {
	// Initialize logging.
	logLevel := logger.LevelInfo
	if os.Getenv("DEBUG") != "" {
		logLevel = logger.LevelDebug
	}

	if err := logger.Init(logLevel, Version, BuildTime); err != nil {
		fmt.Printf("Warning: Failed to initialize logger: %v\n", err)
	}

	defer func() {
		if err := logger.Close(); err != nil {
			fmt.Printf("Warning: Failed to close logger: %v\n", err)
		}
	}()

	// Load configuration.
	cfg := config.Load()

	// Initialize infrastructure dependencies.
	githubClient := github.NewClient()
	commitsCache, err := cache.NewCommitsCache()
	if err != nil {
		logger.Warn("Failed to initialize cache", "error", err.Error())
	}
	clipboardService := clipboard.New()

	// Initialize use cases.
	commitUC := usecase.NewCommitUseCase(githubClient, commitsCache)
	exportUC := usecase.NewExportUseCase()

	// Initialize TUI model.
	model := ui.NewModel(cfg, commitUC, exportUC, clipboardService)

	// Run the application.
	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		logger.Error("Application error", "error", err.Error())
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	logger.Info("Application terminated successfully")
}
