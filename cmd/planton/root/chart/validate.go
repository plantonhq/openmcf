// Package chart implements the `planton chart` subcommands — offline
// authoring tools for the InfraCharts under charts/.
package chart

import (
	"fmt"
	"os"
	"sort"

	"github.com/fatih/color"
	"github.com/plantonhq/planton/internal/cli/cliprint"
	"github.com/plantonhq/planton/pkg/infrachart"
	"github.com/spf13/cobra"
)

var (
	validateSetFlags []string
	validateOrg      string
	validateEnv      string
)

var Validate = &cobra.Command{
	Use:   "validate <chart-dir>",
	Short: "render an InfraChart with its default values and validate every manifest and reference offline",
	Long: `Renders the chart's templates with the default values from values.yaml
(plus any --set overrides), then validates the result without a control plane:

  1. Templates must render inside the platform's sandboxed template subset.
  2. Every rendered document must parse, name a registered kind, carry a
     metadata.name, unmarshal strictly, and pass its spec's validation rules.
  3. Every valueFrom reference must resolve: the field path must exist on the
     target kind, and references to a field's default kind must use the
     annotated composition key.
  4. Chart-internal references must form a dependency graph without cycles.

Validation runs against the protos compiled into this binary. The platform's
'chart build' remains the authoritative server-side gate before publishing.`,
	Example: `
	# Validate a chart with its defaults
	planton chart validate charts/gcp/terraform-state-backend

	# Exercise a feature toggle in its non-default position
	planton chart validate charts/gcp/cloud-run-service --set dnsEnabled=false
	`,
	Args: cobra.ExactArgs(1),
	Run:  validateHandler,
}

func init() {
	Validate.Flags().StringArrayVar(&validateSetFlags, "set", nil,
		"override a param value (key=value, repeatable); org and env may also be overridden")
	Validate.Flags().StringVar(&validateOrg, "org", "acme",
		"value bound to the reserved org template variable")
	Validate.Flags().StringVar(&validateEnv, "env", "dev",
		"value bound to the reserved env template variable")
}

func validateHandler(cmd *cobra.Command, args []string) {
	setOverrides, err := infrachart.ParseSetFlags(validateSetFlags)
	if err != nil {
		cliprint.PrintError(err.Error())
		os.Exit(1)
	}

	report, err := infrachart.Validate(args[0], infrachart.Options{
		Org: validateOrg,
		Env: validateEnv,
		Set: setOverrides,
	})
	if err != nil {
		cliprint.PrintError(err.Error())
		os.Exit(1)
	}

	printReport(args[0], report)
	if report.HasErrors() {
		os.Exit(1)
	}
}

func printReport(chartDir string, report *infrachart.Report) {
	bold := color.New(color.Bold).SprintFunc()
	red := color.New(color.FgRed, color.Bold).SprintFunc()
	yellow := color.New(color.FgYellow, color.Bold).SprintFunc()
	green := color.New(color.FgGreen, color.Bold).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()

	fmt.Printf("%s %s (%s)\n", bold("Chart:"), report.ChartName, chartDir)
	fmt.Printf("%s %d resource(s) rendered\n", bold("Rendered:"), len(report.Docs))
	for _, doc := range report.Docs {
		fmt.Printf("  • %s %s  %s\n", doc.Kind, cyan(doc.Name), color.New(color.Faint).Sprintf("(%s)", doc.File))
	}

	if len(report.Issues) == 0 {
		fmt.Println()
		fmt.Println(green("✔ chart is valid — all manifests and references check out"))
		return
	}

	// Group issues by file for readable attribution.
	byFile := map[string][]infrachart.Issue{}
	var files []string
	for _, issue := range report.Issues {
		if _, seen := byFile[issue.File]; !seen {
			files = append(files, issue.File)
		}
		byFile[issue.File] = append(byFile[issue.File], issue)
	}
	sort.Strings(files)

	errorCount, warningCount := 0, 0
	fmt.Println()
	for _, file := range files {
		fmt.Println(bold("── " + file))
		for _, issue := range byFile[file] {
			prefix := yellow("warning")
			if issue.Severity == infrachart.SeverityError {
				prefix = red("error")
				errorCount++
			} else {
				warningCount++
			}
			subject := ""
			if issue.ResourceKind != "" {
				subject = fmt.Sprintf(" [%s %s]", issue.ResourceKind, issue.ResourceName)
			}
			fmt.Printf("  %s%s: %s\n", prefix, subject, issue.Message)
		}
	}

	fmt.Println()
	if errorCount > 0 {
		fmt.Println(red(fmt.Sprintf("✘ %d error(s), %d warning(s)", errorCount, warningCount)))
	} else {
		fmt.Println(green("✔ chart is valid") + yellow(fmt.Sprintf(" — %d warning(s)", warningCount)))
	}
}
