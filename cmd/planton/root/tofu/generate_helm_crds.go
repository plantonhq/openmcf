package tofu

import (
	"fmt"
	"github.com/plantonhq/planton/internal/cli/ui"
	"os"

	"github.com/plantonhq/planton/internal/cli/flag"
	"github.com/plantonhq/planton/pkg/iac/tofu/generators"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var GenerateHelmCrds = &cobra.Command{
	Use:   "generate-helm-crds",
	Short: "Write the canonical helm_crds.tf a Terraform module carries to derive its CRDs from the pinned chart",
	Long: `The "generate-helm-crds" command emits the canonical helm_crds.tf: the
Terraform half of the catalog's primitive for Helm charts that carry CRDs.

The file carries no per-kind content. It renders the pinned chart with its
CRD switch on (or fetches the pinned upstream bundle), keeps only the
CustomResourceDefinition documents, stamps each with its source chart and
version, and applies each one as a module-owned kept resource ahead of the
release. A module supplies its identity through the locals the file's header
names, and every module's copy is held byte-identical to the generator by the
drift test in pkg/iac/tofu/generators.`,
	Example: `
  # Write the file into a module (the only supported location)
  planton tofu generate-helm-crds --output-file catalog/kubernetes/<kind>/iac/tf/helm_crds.tf

  # Print it
  planton tofu generate-helm-crds
`,
	Args: cobra.NoArgs,
	Run:  generateHelmCrdsHandler,
}

func init() {
	GenerateHelmCrds.Flags().String(string(flag.OutputFile), "", "output file (the module's iac/tf/helm_crds.tf)")
}

func generateHelmCrdsHandler(cmd *cobra.Command, args []string) {
	outputFile, err := cmd.Flags().GetString(string(flag.OutputFile))
	flag.HandleFlagErr(err, flag.OutputFile)

	content := generators.HelmCRDsTF()
	if outputFile != "" {
		if err := os.WriteFile(outputFile, []byte(content), 0644); err != nil {
			ui.Failure(
				fmt.Sprintf("the generated file could not be written to %s: %v", outputFile, err),
				"the output path is not writable, or its parent directory does not exist",
				"create the parent directory or point --output-file at a writable path",
			)
		}
		log.Infof("canonical %s written to %s", generators.HelmCRDsTFFileName, outputFile)
		return
	}
	// fmt.Println writes to stdout; the builtin println writes to stderr,
	// which silently breaks `> helm_crds.tf` capture in scripts.
	fmt.Println(content)
}
