//go:build !codegen
// +build !codegen

package root

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/plantonhq/planton/internal/cli/cliprint"
	"github.com/plantonhq/planton/pkg/infrachart"
	"github.com/spf13/cobra"
)

// Chart is the standalone binary's offline infra-chart tooling. It is deliberately NOT
// part of the engine command set (RegisterCommands): the Planton Platform CLI embeds that
// set and mounts its own `chart` command tree (build/publish against the control plane),
// so the offline subcommands join the shared seam only when CLI convergence unifies the
// two trees.
var Chart = &cobra.Command{
	Use:   "chart",
	Short: "work with infra-charts offline",
}

var chartValidate = &cobra.Command{
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
	Run: chartValidateHandler,
}

func init() {
	chartValidate.Flags().Bool("all", false, "treat the given directories as roots and validate every chart found under them")
	chartValidate.Flags().Bool("verbose", false, "also list charts and variants that validated cleanly")
	Chart.AddCommand(chartValidate)
}

func chartValidateHandler(cmd *cobra.Command, args []string) {
	all, _ := cmd.Flags().GetBool("all")
	verbose, _ := cmd.Flags().GetBool("verbose")

	if len(args) == 0 {
		cliprint.PrintError("provide at least one chart directory (or a root directory with --all)")
		os.Exit(1)
	}

	dirs, err := resolveChartDirs(args, all)
	if err != nil {
		cliprint.PrintError(err.Error())
		os.Exit(1)
	}
	if len(dirs) == 0 {
		cliprint.PrintError("no charts found (a chart directory contains a Chart.yaml)")
		os.Exit(1)
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
		os.Exit(1)
	}
}

// resolveChartDirs expands the command arguments into chart directories: each argument is
// either a chart directory itself or, with --all, a root walked for every chart under it.
func resolveChartDirs(args []string, all bool) ([]string, error) {
	if !all {
		return args, nil
	}
	var dirs []string
	for _, rootDir := range args {
		err := filepath.WalkDir(rootDir, func(path string, d os.DirEntry, err error) error {
			if err != nil || !d.IsDir() {
				return err
			}
			if infrachart.IsChartDir(path) {
				dirs = append(dirs, path)
				return filepath.SkipDir
			}
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("walking %s: %w", rootDir, err)
		}
	}
	sort.Strings(dirs)
	return dirs, nil
}
