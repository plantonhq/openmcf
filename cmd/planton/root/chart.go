//go:build !codegen
// +build !codegen

package root

import (
	"errors"
	"fmt"

	"github.com/plantonhq/planton/internal/cli/cliprint"
	"github.com/plantonhq/planton/pkg/infrachart"
	"github.com/spf13/cobra"
)

// Chart is the standalone binary's offline infra-chart tooling. It is deliberately NOT
// part of the engine command set (RegisterCommands): the Planton Platform CLI mounts its
// own `chart` command tree (build/publish against the control plane) and adds the shared
// offline subcommands to it via their exported constructors -- NewChartValidateCommand
// below -- so both binaries carry the exact same validate command with no drift.
var Chart = &cobra.Command{
	Use:   "chart",
	Short: "work with infra-charts offline",
}

func init() {
	Chart.AddCommand(NewChartValidateCommand())
}

// NewChartValidateCommand builds the `chart validate` command. It is a constructor
// rather than a package variable because cobra commands are stateful (flag values live
// on the instance) and this command mounts in two binaries: the standalone planton
// binary's `chart` tree and the Planton Platform CLI's `chart` tree. Each host mounts
// its own instance.
//
// The command is fully offline -- it dials nothing and needs no configuration -- so
// hosts that guard backend-requiring commands must exempt it.
func NewChartValidateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate [chart-dir ...]",
		Short: "render and validate infra-charts against the compiled-in kind registry",
		Long: `Render each chart's templates with the defaults declared in values.yaml and validate
every produced manifest offline: the kind must exist, every field must exist on the
kind's spec (unknown or renamed fields fail), the spec must pass its validation rules,
and every valueFrom reference must resolve to a real field on the referenced kind.

Each bool param is additionally flipped once so conditional manifests are exercised in
both branches.

This requires no backend and no network: everything validates against the schemas
compiled into this binary. The control plane performs the same pipeline authoritatively
when a chart is published.`,
		Example: `
	# Validate one chart
	planton chart validate charts/aws/eks-environment

	# Validate several charts
	planton chart validate charts/aws/eks-environment charts/aws/ecs-environment

	# Validate every chart under a directory tree
	planton chart validate --all charts/`,
		RunE: chartValidateHandler,
		// The handler prints the full per-chart report itself; the returned error is a
		// one-line summary for the host root's error path, so neither cobra's usage
		// dump nor its error echo should repeat it.
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	cmd.Flags().Bool("all", false, "treat the given directories as roots and validate every chart found under them")
	cmd.Flags().Bool("verbose", false, "also list charts and variants that validated cleanly")
	return cmd
}

func chartValidateHandler(cmd *cobra.Command, args []string) error {
	all, _ := cmd.Flags().GetBool("all")
	verbose, _ := cmd.Flags().GetBool("verbose")

	if len(args) == 0 {
		return errors.New("provide at least one chart directory (or a root directory with --all)")
	}

	dirs, err := infrachart.DiscoverCharts(args, all)
	if err != nil {
		return err
	}
	if len(dirs) == 0 {
		return errors.New("no charts found (a chart directory contains a Chart.yaml)")
	}

	passed, failed := 0, 0
	for _, dir := range dirs {
		report, err := infrachart.ValidateChart(dir)
		if err != nil {
			failed++
			cliprint.PrintError(fmt.Sprintf("✘ %s: %v", dir, err))
			continue
		}
		if report.Valid() {
			passed++
			if verbose {
				cliprint.PrintSuccess(fmt.Sprintf("✔ %s (%d variants)", dir, len(report.Variants)))
			}
			continue
		}
		failed++
		cliprint.PrintError(fmt.Sprintf("✘ %s", dir))
		for _, variant := range report.Variants {
			for _, docErr := range variant.Errors {
				label := docErr.Template
				if docErr.Kind != "" {
					label = fmt.Sprintf("%s (%s/%s)", docErr.Template, docErr.Kind, docErr.Name)
				}
				fmt.Printf("  [%s] %s:\n%v\n", variant.Name, label, docErr.Err)
			}
		}
	}

	fmt.Printf("\nchart validation: %d passed, %d failed out of %d\n", passed, failed, passed+failed)
	if failed > 0 {
		return fmt.Errorf("%d of %d charts failed validation", failed, passed+failed)
	}
	return nil
}
