//go:build !codegen
// +build !codegen

// Package moduleverify renders `planton module verify` results: findings
// grouped by severity, each a plain-language sentence naming the file, the
// cause, and the way out.
package moduleverify

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/plantonhq/planton/pkg/iac/moduleverify"
)

var (
	colorRed    = lipgloss.Color("#FF6B6B")
	colorGreen  = lipgloss.Color("#69DB7C")
	colorYellow = lipgloss.Color("#FFD43B")
	colorBlue   = lipgloss.Color("#74C0FC")
	colorGray   = lipgloss.Color("#868E96")
	colorWhite  = lipgloss.Color("#DEE2E6")

	errorStyle   = lipgloss.NewStyle().Foreground(colorRed).Bold(true)
	successStyle = lipgloss.NewStyle().Foreground(colorGreen).Bold(true)
	warningStyle = lipgloss.NewStyle().Foreground(colorYellow).Bold(true)
	fileStyle    = lipgloss.NewStyle().Foreground(colorBlue)
	textStyle    = lipgloss.NewStyle().Foreground(colorWhite)
	dimStyle     = lipgloss.NewStyle().Foreground(colorGray)
)

const (
	iconError   = "✗"
	iconSuccess = "✓"
	iconWarning = "!"
)

// RenderResult prints the verification report: errors first (they fail the
// run), then warnings, then the notices for anything that was skipped.
func RenderResult(result *moduleverify.Result) {
	fmt.Fprintln(os.Stderr)

	var verifyErrors, verifyWarnings []moduleverify.Violation
	for _, v := range result.Violations {
		if v.Severity == moduleverify.SeverityError {
			verifyErrors = append(verifyErrors, v)
		} else {
			verifyWarnings = append(verifyWarnings, v)
		}
	}

	if len(verifyErrors) == 0 {
		fmt.Fprintf(os.Stderr, "%s  %s\n",
			successStyle.Render(iconSuccess),
			successStyle.Render(fmt.Sprintf("The module conforms to the %s contract", result.KindName)))
	} else {
		fmt.Fprintf(os.Stderr, "%s  %s\n",
			errorStyle.Render(iconError),
			errorStyle.Render(fmt.Sprintf("The module breaks the %s contract", result.KindName)))
	}
	fmt.Fprintf(os.Stderr, "   %s %s   %s %s\n",
		dimStyle.Render("module:"), fileStyle.Render(result.ModuleDir),
		dimStyle.Render("engine:"), textStyle.Render(result.Provisioner.String()))

	renderViolations("Errors (these fail deployments)", errorStyle, iconError, verifyErrors)
	renderViolations("Warnings (worth a look, not fatal)", warningStyle, iconWarning, verifyWarnings)

	if len(result.Notices) > 0 {
		fmt.Fprintln(os.Stderr)
		for _, notice := range result.Notices {
			fmt.Fprintf(os.Stderr, "   %s %s\n", dimStyle.Render("note:"), dimStyle.Render(notice))
		}
	}

	fmt.Fprintln(os.Stderr)
}

func renderViolations(title string, style lipgloss.Style, icon string, violations []moduleverify.Violation) {
	if len(violations) == 0 {
		return
	}
	fmt.Fprintln(os.Stderr)
	fmt.Fprintf(os.Stderr, "   %s\n", style.Render(title))
	for _, v := range violations {
		location := ""
		if v.File != "" {
			location = fileStyle.Render(v.File) + dimStyle.Render(" — ")
		}
		fmt.Fprintf(os.Stderr, "   %s %s%s\n", style.Render(icon), location, textStyle.Render(v.Summary))
	}
}
