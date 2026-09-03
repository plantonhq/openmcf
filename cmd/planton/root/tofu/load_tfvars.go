package tofu

import (
	"fmt"
	"github.com/plantonhq/planton/internal/cli/ui"

	"github.com/plantonhq/planton/internal/manifest"
	"github.com/plantonhq/planton/pkg/iac/tofu/generators"
	"github.com/spf13/cobra"
)

var LoadTfVars = &cobra.Command{
	Use:   "load-tfvars",
	Short: "load a planton manifest into tfvars format",
	Example: `
	planton tofu load-tfvars --manifest manifest.yaml
	`,
	Args: cobra.ExactArgs(1), //path of the manifest to load
	Run:  loadTfVarsHandler,
}

func loadTfVarsHandler(cmd *cobra.Command, args []string) {
	manifestPath := args[0]
	updatedManifest, err := manifest.LoadWithOverrides(manifestPath, map[string]string{})
	if err != nil {
		ui.Failure(
			fmt.Sprintf("the manifest at %s could not be loaded: %v", manifestPath, err),
			"the file is missing, unreadable, or does not load into the kind it declares",
			fmt.Sprintf("check the path, then run `planton validate-manifest -f %s` for the field-level report", manifestPath),
		)
	}
	tfvarsString, err := generators.RenderTFVars(updatedManifest)
	if err != nil {
		ui.Failure(
			fmt.Sprintf("the manifest at %s could not be rendered as tfvars: %v", manifestPath, err),
			"the kind's spec loaded but one of its fields has no HCL representation",
			"report it at https://github.com/plantonhq/planton/issues naming the kind and the field",
		)
	}
	// stdout, not the builtin println (which writes to stderr) -- the whole
	// point of this command is piping the tfvars into a file or another tool.
	fmt.Println(tfvarsString)
}
