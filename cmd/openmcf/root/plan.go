package root

import (
	"os"

	"github.com/plantonhq/openmcf/apis/org/openmcf/shared/iac/pulumi"
	"github.com/plantonhq/openmcf/apis/org/openmcf/shared/iac/terraform"
	"github.com/plantonhq/openmcf/internal/cli/cliprint"
	"github.com/plantonhq/openmcf/internal/cli/iacflags"
	"github.com/plantonhq/openmcf/internal/cli/iacrunner"
	climanifest "github.com/plantonhq/openmcf/internal/cli/manifest"
	"github.com/plantonhq/openmcf/internal/manifest"
	"github.com/plantonhq/openmcf/pkg/iac/provisioner"
	"github.com/spf13/cobra"
)

var Plan = &cobra.Command{
	Use:     "plan",
	Aliases: []string{"preview"},
	Short:   "preview infrastructure changes using the provisioner specified in manifest",
	Long: `Preview infrastructure changes by automatically routing to the appropriate provisioner
(Pulumi, Tofu, or Terraform) based on the manifest label 'openmcf.org/provisioner'.

If the provisioner label is not present, you will be prompted to select one interactively.

This command has 'preview' as an alias for Pulumi-style experience.`,
	Example: `
	# Preview changes with manifest file
	openmcf plan -f manifest.yaml
	openmcf preview -f manifest.yaml
	openmcf plan --manifest manifest.yaml

	# Preview with stack input file (extracts manifest from target field)
	openmcf plan -i stack-input.yaml

	# Preview with kustomize
	openmcf plan --kustomize-dir _kustomize --overlay prod

	# Preview with field overrides
	openmcf plan -f manifest.yaml --set spec.version=v1.2.3

	# Preview destroy plan (Tofu)
	openmcf plan -f manifest.yaml --destroy
	`,
	Run: planHandler,
}

func init() {
	iacflags.AddManifestSourceFlags(Plan)
	iacflags.AddProviderConfigFlags(Plan)
	iacflags.AddExecutionFlags(Plan)
	iacflags.AddPulumiFlags(Plan)
	iacflags.AddTofuPlanFlags(Plan)
	iacflags.AddTofuInitFlags(Plan)
}

func planHandler(cmd *cobra.Command, args []string) {
	ctx, err := iacrunner.ResolveContext(cmd)
	if err != nil {
		// Only print error if it wasn't already handled (clipboard/manifest load errors are pre-handled)
		if !climanifest.IsClipboardError(err) && !manifest.IsManifestLoadError(err) {
			cliprint.PrintError(err.Error())
		}
		os.Exit(1)
	}
	defer ctx.Cleanup()

	switch ctx.ProvisionerType {
	case provisioner.ProvisionerTypePulumi:
		// For preview, we use update operation with isPreview=true
		if err := iacrunner.RunPulumi(ctx, cmd, pulumi.PulumiOperationType_update, true); err != nil {
			os.Exit(1)
		}
	case provisioner.ProvisionerTypeTofu:
		if err := iacrunner.RunTofu(ctx, cmd, terraform.TerraformOperationType_plan); err != nil {
			os.Exit(1)
		}
	case provisioner.ProvisionerTypeTerraform:
		if err := iacrunner.RunTerraform(ctx, cmd, terraform.TerraformOperationType_plan); err != nil {
			os.Exit(1)
		}
	default:
		cliprint.PrintError("Unknown provisioner type")
		os.Exit(1)
	}
}
