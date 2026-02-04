// Package ui provides the terminal user interface.
package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// Color palette - modern soft gradients, muted tones.
var (
	// Primary colors - soft purple/indigo gradient.
	colorPrimary      = lipgloss.Color("#8B5CF6") // Vibrant violet.
	colorPrimaryLight = lipgloss.Color("#A78BFA") // Light violet.

	// Accent colors - cyan/teal for highlights.
	colorAccent      = lipgloss.Color("#06B6D4") // Cyan.
	colorAccentLight = lipgloss.Color("#22D3EE") // Light cyan.

	// Status colors - modern, softer tones.
	colorSuccess   = lipgloss.Color("#10B981") // Emerald.
	colorSuccessBg = lipgloss.Color("#064E3B") // Dark emerald background.
	colorWarning   = lipgloss.Color("#F59E0B") // Amber.
	colorWarningBg = lipgloss.Color("#78350F") // Dark amber background.
	colorError     = lipgloss.Color("#EF4444") // Red.
	colorErrorBg   = lipgloss.Color("#7F1D1D") // Dark red background.

	// Progress bar gradient colors.
	colorBarStart  = lipgloss.Color("#06B6D4") // Cyan.
	colorBarMiddle = lipgloss.Color("#8B5CF6") // Violet.
	colorBarEnd    = lipgloss.Color("#EC4899") // Pink.

	// Neutral colors - refined grayscale.
	colorText       = lipgloss.Color("#F9FAFB") // Near white.
	colorTextDim    = lipgloss.Color("#9CA3AF") // Dim text.
	colorTextMuted  = lipgloss.Color("#6B7280") // Muted text.
	colorTextSubtle = lipgloss.Color("#4B5563") // Subtle text.
)

// UI symbols.
const (
	iconArrowRight = "› "
	iconCheckBox   = "◉"
	iconUncheckBox = "○"
	iconCommit     = "•"
	iconSuccess    = "✓"
	iconWarning    = "⚠"
	iconError      = "✗"
	iconInfo       = "ℹ"
	iconBreadcrumb = " › "
	iconDivider    = "─"

	// Progress bar characters (using Unicode block elements).
	barFull      = "█"
	barThreeQtr  = "▓"
	barHalf      = "▒"
	barQuarter   = "░"
	barEmpty     = "░"
	barCapLeft   = "▐"
	barCapRight  = "▌"
)

// Styles - modern, cohesive design system.
var (
	// Repository styling.
	styleRepo = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorAccent)

	// Commit message styling.
	styleCommit = lipgloss.NewStyle().
			Foreground(colorText)

	// Cursor and selection.
	styleCursor = lipgloss.NewStyle().
			Foreground(colorPrimary).
			Bold(true)

	// Checkbox styling.
	styleCheckbox = lipgloss.NewStyle().
			Foreground(colorSuccess)

	styleCheckboxUnchecked = lipgloss.NewStyle().
				Foreground(colorTextMuted)

	// Highlight for selected items.
	styleHighlight = lipgloss.NewStyle().
			Foreground(colorAccentLight)

	// Footer and help text.
	styleFooter = lipgloss.NewStyle().
			Foreground(colorTextMuted).
			MarginTop(1)

	styleHelpKey = lipgloss.NewStyle().
			Foreground(colorPrimaryLight).
			Bold(true)

	styleHelpText = lipgloss.NewStyle().
			Foreground(colorTextMuted)

	styleHelpDivider = lipgloss.NewStyle().
				Foreground(colorTextSubtle).
				SetString(" │ ")

	// Date input styling.
	styleDateLabel = lipgloss.NewStyle().
			Foreground(colorTextDim)

	// Banner-style messages with icons and backgrounds.
	styleSuccessBanner = lipgloss.NewStyle().
				Foreground(colorSuccess).
				Background(colorSuccessBg).
				Bold(true).
				Padding(0, 1)

	styleWarningBanner = lipgloss.NewStyle().
				Foreground(colorWarning).
				Background(colorWarningBg).
				Bold(true).
				Padding(0, 1)

	styleErrorBanner = lipgloss.NewStyle().
				Foreground(colorError).
				Background(colorErrorBg).
				Bold(true).
				Padding(0, 1)

	// Stats styling.
	styleStatsValue = lipgloss.NewStyle().
			Foreground(colorAccent).
			Bold(true)

	styleStatsLabel = lipgloss.NewStyle().
			Foreground(colorTextDim)

	// Progress bar style for empty portion.
	styleBarEmpty = lipgloss.NewStyle().
			Foreground(colorTextSubtle)

	// Box style for screens.
	styleBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorPrimary).
			Padding(1, 2)

	// Input box style.
	styleInputBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorPrimaryLight).
			Padding(0, 1)

	// Header with breadcrumb navigation.
	styleHeader = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorPrimary)

	styleHeaderApp = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorTextMuted)

	styleHeaderDivider = lipgloss.NewStyle().
				Foreground(colorTextSubtle)

	// Help bar container style.
	styleHelpBar = lipgloss.NewStyle().
			Foreground(colorTextMuted).
			BorderTop(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(colorTextSubtle).
			PaddingTop(1).
			MarginTop(1)

	// Section divider.
	styleDivider = lipgloss.NewStyle().
			Foreground(colorTextSubtle)

	// List header.
	styleListHeader = lipgloss.NewStyle().
			Foreground(colorTextDim).
			Italic(true)
)
