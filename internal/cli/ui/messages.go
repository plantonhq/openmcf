package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// separator returns a styled separator line
func separator(style lipgloss.Style) string {
	return style.Render(strings.Repeat(separatorChar, separatorLength))
}

// Error prints a styled error message and exits with code 1
func Error(title, message string, hints ...string) {
	printError(title, message, hints...)
	os.Exit(1)
}

// ErrorWithoutExit prints a styled error message without exiting
func ErrorWithoutExit(title, message string, hints ...string) {
	printError(title, message, hints...)
}

func printError(title, message string, hints ...string) {
	fmt.Fprintln(os.Stderr)

	// Icon and title
	fmt.Fprintf(os.Stderr, "%s  %s\n",
		errorIcon.Render(iconError),
		errorTitle.Render(title))

	// Message
	if message != "" {
		fmt.Fprintln(os.Stderr)
		for _, line := range strings.Split(message, "\n") {
			fmt.Fprintf(os.Stderr, "   %s\n", errorMessage.Render(line))
		}
	}

	// Hints
	if len(hints) > 0 {
		fmt.Fprintln(os.Stderr)
		for _, hint := range hints {
			if hint != "" {
				fmt.Fprintf(os.Stderr, "   %s %s\n",
					dimStyle.Render("Hint:"),
					hintStyle.Render(hint))
			}
		}
	}

	fmt.Fprintln(os.Stderr)
}

// Success prints a styled success message
func Success(title string, details ...string) {
	fmt.Println()

	// Icon and title
	fmt.Printf("%s  %s\n",
		successIcon.Render(iconSuccess),
		successTitle.Render(title))

	// Details
	if len(details) > 0 {
		fmt.Println()
		for _, detail := range details {
			for _, line := range strings.Split(detail, "\n") {
				fmt.Printf("   %s\n", successMessage.Render(line))
			}
		}
	}

	fmt.Println()
}

// Warning prints a styled warning message
func Warning(title, message string) {
	fmt.Println()

	// Icon and title
	fmt.Printf("%s  %s\n",
		warningIcon.Render(iconWarning),
		warningTitle.Render(title))

	// Message
	if message != "" {
		fmt.Println()
		for _, line := range strings.Split(message, "\n") {
			fmt.Printf("   %s\n", dimStyle.Render(line))
		}
	}

	fmt.Println()
}

// Info prints a styled info message
func Info(message string) {
	fmt.Printf("%s  %s\n",
		infoIcon.Render(iconInfo),
		infoMessage.Render(message))
}

// ErrorBanner prints a styled error banner with title, message, and optional command suggestion
func ErrorBanner(title, message, suggestedCmd, tip string) {
	sep := separator(errorIcon)

	fmt.Println()
	fmt.Println(sep)
	fmt.Printf("%s  %s\n", errorIcon.Render(iconError), errorTitle.Render(title))
	fmt.Println(sep)

	if message != "" {
		for _, line := range strings.Split(message, "\n") {
			fmt.Println(errorMessage.Render(line))
		}
		fmt.Println()
	}

	if suggestedCmd != "" {
		fmt.Println(errorMessage.Render("To fix this, run:"))
		fmt.Println()
		fmt.Printf("    %s\n", cmdStyle.Render(suggestedCmd))
		fmt.Println()
	}

	if tip != "" {
		fmt.Printf("%s %s\n", infoIcon.Render(iconTip), infoMessage.Render("Tip: "+tip))
	}

	fmt.Println(sep)
}

// InfoBanner prints a styled info banner with title, message, and optional command suggestion
func InfoBanner(title, message, suggestedCmd, tip string) {
	sep := separator(infoIcon)

	fmt.Println()
	fmt.Println(sep)
	fmt.Printf("%s  %s\n", infoIcon.Render("ℹ️"), infoTitle.Render(title))
	fmt.Println(sep)

	if message != "" {
		for _, line := range strings.Split(message, "\n") {
			fmt.Println(infoMessage.Render(line))
		}
		fmt.Println()
	}

	if suggestedCmd != "" {
		fmt.Println(infoMessage.Render("Example:"))
		fmt.Println()
		fmt.Printf("    %s\n", cmdStyle.Render(suggestedCmd))
		fmt.Println()
	}

	if tip != "" {
		fmt.Printf("%s %s\n", infoIcon.Render(iconTip), infoMessage.Render("Tip: "+tip))
	}

	fmt.Println(sep)
}

// WarningBanner prints a styled warning banner with title, message, and optional command suggestion
func WarningBanner(title, message, suggestedCmd, tip string) {
	sep := separator(warningIcon)

	fmt.Println()
	fmt.Println(sep)
	fmt.Printf("%s  %s\n", warningIcon.Render("⚠️"), warningTitle.Render(title))
	fmt.Println(sep)

	if message != "" {
		for _, line := range strings.Split(message, "\n") {
			fmt.Println(warningMessage.Render(line))
		}
		fmt.Println()
	}

	if suggestedCmd != "" {
		fmt.Println(warningMessage.Render("To proceed, run:"))
		fmt.Println()
		fmt.Printf("    %s\n", cmdStyle.Render(suggestedCmd))
		fmt.Println()
	}

	if tip != "" {
		fmt.Printf("%s %s\n", infoIcon.Render(iconTip), infoMessage.Render("Tip: "+tip))
	}

	fmt.Println(sep)
}

// Path formats a path with styling (for use in messages)
func Path(p string) string {
	return pathStyle.Render(p)
}

// Cmd formats a command with styling (for use in messages)
func Cmd(c string) string {
	return cmdStyle.Render(c)
}

// Dim formats text as dimmed (for use in messages)
func Dim(s string) string {
	return dimStyle.Render(s)
}
