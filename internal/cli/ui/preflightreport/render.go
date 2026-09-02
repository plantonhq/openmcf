//go:build !codegen
// +build !codegen

// Package preflightreport renders a manifest-set deploy's preflight report:
// every wall check line-itemed, failures grouped under the check that found
// them, and a one-line verdict rendered last so a CI log's tail shows the one
// sentence that matters with the full report in the scrollback above it.
//
// It is a ui subpackage (not part of ui itself) because it imports
// pkg/setdeploy, which reaches internal/manifest — a package that imports ui
// for its load-error display. The subpackage keeps the render layer out of
// that cycle, the same shape ui/moduleverify uses.
package preflightreport

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/plantonhq/planton/pkg/setdeploy"
)

var (
	colorRed    = lipgloss.Color("#FF6B6B")
	colorGreen  = lipgloss.Color("#69DB7C")
	colorYellow = lipgloss.Color("#FFD43B")
	colorGray   = lipgloss.Color("#868E96")

	refusalStyle    = lipgloss.NewStyle().Foreground(colorRed)
	titleStyle      = lipgloss.NewStyle().Foreground(colorGreen).Bold(true)
	warningStyle    = lipgloss.NewStyle().Foreground(colorYellow)
	assumptionStyle = lipgloss.NewStyle().Foreground(colorGray)
)

// Line marks carry the class; the goldens pin them along with the text.
//
//	✔  a verified fact
//	✖  a refusal (blocks the deploy)
//	⚠  a warning (deploy proceeds; the user should know)
//	≈  a stated assumption (unverifiable here; verified at handoff)
const (
	markVerified   = "✔"
	markRefusal    = "✖"
	markWarning    = "⚠"
	markAssumption = "≈"
)

// FormatReportLines renders the whole report as plain lines: a heading per
// check, then its verified facts and entries indented beneath it. Split from
// the printer so the refusal sentences and report shape are pinned by golden
// tests — a wording change is a visible diff, never an accident.
func FormatReportLines(report *setdeploy.Report) []string {
	var lines []string
	for i := range report.Checks {
		check := &report.Checks[i]
		lines = append(lines, checkHeading(check))
		for _, fact := range check.Verified {
			lines = append(lines, fmt.Sprintf("   %s %s", markVerified, fact))
		}
		for _, entry := range check.Entries {
			lines = append(lines, fmt.Sprintf("   %s %s", entryMark(entry.Severity), entryText(entry)))
		}
	}
	return lines
}

// FormatVerdict renders the one-line verdict: pass or refusal, with the
// counts that matter.
func FormatVerdict(report *setdeploy.Report) string {
	refusals := report.RefusalCount()
	if refusals > 0 {
		plural := "s"
		if refusals == 1 {
			plural = ""
		}
		return fmt.Sprintf("%s preflight refused the deploy: %d problem%s named above — nothing was handed to an IaC engine", markRefusal, refusals, plural)
	}
	facts, warnings, assumptions := 0, 0, 0
	for i := range report.Checks {
		facts += len(report.Checks[i].Verified)
		for _, e := range report.Checks[i].Entries {
			switch e.Severity {
			case setdeploy.SeverityWarning:
				warnings++
			case setdeploy.SeverityAssumption:
				assumptions++
			}
		}
	}
	verdict := fmt.Sprintf("%s preflight passed: %d facts verified", markVerified, facts)
	if warnings > 0 {
		verdict += fmt.Sprintf(", %d warnings", warnings)
	}
	if assumptions > 0 {
		verdict += fmt.Sprintf(", %d assumptions stated", assumptions)
	}
	return verdict
}

func checkHeading(check *setdeploy.Check) string {
	mark := markVerified
	for _, e := range check.Entries {
		if e.Severity == setdeploy.SeverityRefusal {
			mark = markRefusal
			break
		}
	}
	return fmt.Sprintf("%s %s", mark, check.Title)
}

func entryMark(severity setdeploy.Severity) string {
	switch severity {
	case setdeploy.SeverityRefusal:
		return markRefusal
	case setdeploy.SeverityWarning:
		return markWarning
	default:
		return markAssumption
	}
}

// entryText renders one entry with its document source prefixed. The entry's
// message already names the field path where one applies (the refusal
// sentence grammar), so only the source rides here.
func entryText(entry setdeploy.Entry) string {
	if entry.Source != "" && !strings.HasPrefix(entry.Message, entry.Source) {
		return entry.Source + ": " + entry.Message
	}
	return entry.Message
}

// Print renders the styled report and verdict to stdout.
func Print(report *setdeploy.Report) {
	fmt.Println()
	fmt.Printf("🧱 %s\n", titleStyle.Render("Preflight"))
	for _, line := range FormatReportLines(report) {
		fmt.Printf("  %s\n", styleLine(line))
	}
	fmt.Println()
	fmt.Println(styleLine(FormatVerdict(report)))
	fmt.Println()
}

// styleLine colors a rendered line by its mark without altering its text —
// the goldens pin the text, the colors ride on top.
func styleLine(line string) string {
	trimmed := strings.TrimLeft(line, " ")
	switch {
	case strings.HasPrefix(trimmed, markRefusal):
		return refusalStyle.Render(line)
	case strings.HasPrefix(trimmed, markWarning):
		return warningStyle.Render(line)
	case strings.HasPrefix(trimmed, markAssumption):
		return assumptionStyle.Render(line)
	default:
		return line
	}
}
