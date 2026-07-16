package ui

import (
	"fmt"

	"github.com/plantonhq/planton/pkg/yamldiag"
)

// ManifestDiagnosis displays a YAML-aware manifest diagnosis: each mismatch
// with its real line number, field path, expected shape, suggestion, and the
// schema-reference command that answers "what CAN I write here?". This is
// the styled sibling of the diagnosis text embedded in the error itself --
// same content, terminal layout on top.
func ManifestDiagnosis(manifestPath, kindName string, mismatches []yamldiag.Mismatch) {
	sep := separator(errorIcon)

	fmt.Println()
	fmt.Println(sep)
	fmt.Printf("%s  %s\n", errorIcon.Render(iconError), errorTitle.Render("Manifest Does Not Fit Its Schema"))
	fmt.Println(sep)

	fmt.Printf("Manifest: %s (kind: %s)\n", Path(manifestPath), errorTitle.Render(kindName))
	fmt.Println()

	for _, m := range mismatches {
		fmt.Printf("%s %s %s\n",
			errorIcon.Render(iconError),
			errorTitle.Render(m.Path),
			Dim(fmt.Sprintf("line %d", m.Line)))
		fmt.Printf("   %s\n", errorMessage.Render(m.Problem))
		if m.Suggestion != "" {
			fmt.Printf("   Did you mean: %s\n", Cmd(m.Suggestion))
		}
		fmt.Println()
	}

	fmt.Printf("%s %s\n", infoIcon.Render(iconTip),
		infoMessage.Render(fmt.Sprintf("Tip: `planton explain %s.<field.path>` shows what any field accepts", kindName)))
	fmt.Println(sep)
}
