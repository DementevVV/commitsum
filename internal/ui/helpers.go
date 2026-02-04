package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// renderHeader renders a header with app name and screen title.
func renderHeader(screenTitle string) string {
	app := styleHeaderApp.Render("commitsum")
	divider := styleHeaderDivider.Render(iconBreadcrumb)
	title := styleHeader.Render(screenTitle)
	return app + divider + title + "\n\n"
}

// renderHelpBar renders a help bar with key-value pairs and top border.
func renderHelpBar(items [][]string) string {
	var parts []string
	for _, item := range items {
		key := styleHelpKey.Render(item[0])
		text := styleHelpText.Render(item[1])
		parts = append(parts, key+" "+text)
	}
	content := strings.Join(parts, styleHelpDivider.String())
	return styleHelpBar.Render(content)
}

// renderDivider renders a horizontal divider line.
func renderDivider(width int) string {
	if width <= 0 {
		width = 40
	}
	return styleDivider.Render(strings.Repeat(iconDivider, width))
}

// renderSuccessBanner renders a success message with icon.
func renderSuccessBanner(msg string) string {
	return styleSuccessBanner.Render(iconSuccess + " " + msg)
}

// renderWarningBanner renders a warning message with icon.
func renderWarningBanner(msg string) string {
	return styleWarningBanner.Render(iconWarning + " " + msg)
}

// renderErrorBanner renders an error message with icon.
func renderErrorBanner(msg string) string {
	return styleErrorBanner.Render(iconError + " " + msg)
}

// renderListHeader renders a list header with count.
func renderListHeader(label string, count int) string {
	return styleListHeader.Render(fmt.Sprintf("%s (%d)", label, count))
}

// renderProgressBar creates a beautiful gradient progress bar.
func renderProgressBar(value, maxValue, width int) string {
	if maxValue == 0 {
		maxValue = 1
	}
	if width <= 0 {
		width = 20
	}

	// Calculate filled width proportionally.
	filledWidth := (value * width) / maxValue
	if filledWidth > width {
		filledWidth = width
	}
	if value > 0 && filledWidth == 0 {
		filledWidth = 1 // At least 1 block if there's any value.
	}

	emptyWidth := width - filledWidth

	// Create gradient effect across the filled portion.
	var bar strings.Builder

	for i := range filledWidth {
		// Calculate color based on position for gradient effect.
		color := getGradientColor(i, filledWidth)
		bar.WriteString(lipgloss.NewStyle().Foreground(color).Render(barFull))
	}

	// Add empty portion.
	if emptyWidth > 0 {
		bar.WriteString(styleBarEmpty.Render(strings.Repeat(barEmpty, emptyWidth)))
	}

	return bar.String()
}

// getGradientColor returns a color for gradient effect based on position.
func getGradientColor(pos, total int) lipgloss.Color {
	if total <= 1 {
		return colorBarStart
	}

	// Define gradient stops: cyan -> violet -> pink.
	gradient := []lipgloss.Color{
		colorBarStart,  // Cyan.
		colorBarMiddle, // Violet.
		colorBarEnd,    // Pink.
	}

	// Calculate which segment we're in.
	ratio := float64(pos) / float64(total-1)
	if ratio >= 1.0 {
		return gradient[len(gradient)-1]
	}

	// Map to gradient index.
	idx := int(ratio * float64(len(gradient)-1))
	if idx >= len(gradient)-1 {
		idx = len(gradient) - 2
	}

	return gradient[idx]
}

// findMaxCommits returns the maximum commit count from the stats.
func findMaxCommits(commitsPerRepo map[string]int) int {
	maxVal := 0
	for _, count := range commitsPerRepo {
		if count > maxVal {
			maxVal = count
		}
	}
	return maxVal
}
