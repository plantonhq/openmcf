package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// Color palette - matching Planton CLI and gitr
var (
	colorRed     = lipgloss.Color("#FF6B6B")
	colorGreen   = lipgloss.Color("#69DB7C")
	colorYellow  = lipgloss.Color("#FFD43B")
	colorOrange  = lipgloss.Color("#FFA94D")
	colorBlue    = lipgloss.Color("#74C0FC")
	colorGray    = lipgloss.Color("#868E96")
	colorDimGray = lipgloss.Color("#495057")
	colorWhite   = lipgloss.Color("#DEE2E6")
)

// Icons
const (
	iconError   = "✗"
	iconSuccess = "✓"
	iconWarning = "!"
	iconInfo    = "→"
	iconTip     = "💡"
)

// Text styles
var (
	// Error styles
	errorIcon = lipgloss.NewStyle().
			Foreground(colorRed).
			Bold(true)

	errorTitle = lipgloss.NewStyle().
			Foreground(colorRed).
			Bold(true)

	errorMessage = lipgloss.NewStyle().
			Foreground(colorWhite)

	// Success styles
	successIcon = lipgloss.NewStyle().
			Foreground(colorGreen).
			Bold(true)

	successTitle = lipgloss.NewStyle().
			Foreground(colorGreen).
			Bold(true)

	successMessage = lipgloss.NewStyle().
			Foreground(colorWhite)

	// Warning styles
	warningIcon = lipgloss.NewStyle().
			Foreground(colorYellow).
			Bold(true)

	warningTitle = lipgloss.NewStyle().
			Foreground(colorYellow).
			Bold(true)

	warningMessage = lipgloss.NewStyle().
			Foreground(colorWhite)

	// Info styles
	infoIcon = lipgloss.NewStyle().
			Foreground(colorBlue).
			Bold(true)

	infoTitle = lipgloss.NewStyle().
			Foreground(colorBlue).
			Bold(true)

	infoMessage = lipgloss.NewStyle().
			Foreground(colorWhite)

	// Path style - for file/directory paths
	pathStyle = lipgloss.NewStyle().
			Foreground(colorBlue).
			Bold(true)

	// Command style - for CLI commands
	cmdStyle = lipgloss.NewStyle().
			Foreground(colorYellow)

	// Dim text
	dimStyle = lipgloss.NewStyle().
			Foreground(colorGray)

	// Hint style
	hintStyle = lipgloss.NewStyle().
			Foreground(colorDimGray).
			Italic(true)

	// Attention styles (for updates, important notices)
	attentionIcon = lipgloss.NewStyle().
			Foreground(colorOrange).
			Bold(true)

	attentionMessage = lipgloss.NewStyle().
				Foreground(colorWhite)
)

// Separator line for banners
const separatorChar = "═"
const separatorLength = 80
