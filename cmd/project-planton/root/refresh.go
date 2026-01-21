package root

import (
	"os"

	"github.com/plantonhq/project-planton/apis/org/project_planton/shared/iac/pulumi"
	"github.com/plantonhq/project-planton/apis/org/project_planton/shared/iac/terraform"
	"github.com/plantonhq/project-planton/internal/cli/cliprint"
	"github.com/plantonhq/project-planton/internal/cli/iacflags"
	"github.com/plantonhq/project-planton/internal/cli/iacrunner"
	climanifest "github.com/plantonhq/project-planton/internal/cli/manifest"
	"github.com/plantonhq/project-planton/internal/manifest"
	"github.com/plantonhq/project-planton/pkg/iac/provisioner"
	"github.com/spf13/cobra"
)

var Refresh = &cobra.Command{
	Use:   "refresh",
	Short: "sync state with cloud reality using the provisioner specified in manifest",
	Long: `Refresh infrastructure state by automatically routing to the appropriate provisioner
(Pulumi, Tofu, or Terraform) based on the manifest label 'project-planton.org/provisioner'.

This command queries your cloud provider for the current state of managed resources and
updates the state file to reflect reality. It does NOT modify any cloud resources.

If the provisioner label is not present, you will be prompted to select one interactively.`,
	Example: `
	# Refresh state with manifest file
	project-planton refresh -f manifest.yaml
	project-planton refresh --manifest manifest.yaml

	# Refresh with stack input file (extracts manifest from target field)
	project-planton refresh -i stack-input.yaml

	# Refresh with kustomize
	project-planton refresh --kustomize-dir _kustomize --overlay prod

	# Refresh with field overrides
	project-planton refresh -f manifest.yaml --set spec.version=v1.2.3
	`,
	Run: refreshHandler,
}

func init() {
	iacflags.AddManifestSourceFlags(Refresh)
	iacflags.AddProviderConfigFlags(Refresh)
	iacflags.AddExecutionFlags(Refresh)
	iacflags.AddPulumiFlags(Refresh)
	iacflags.AddTofuInitFlags(Refresh)
}

func refreshHandler(cmd *cobra.Command, args []string) {
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
		if err := iacrunner.RunPulumi(ctx, cmd, pulumi.PulumiOperationType_refresh, false); err != nil {
			os.Exit(1)
		}
	case provisioner.ProvisionerTypeTofu:
		if err := iacrunner.RunTofu(ctx, cmd, terraform.TerraformOperationType_refresh); err != nil {
			os.Exit(1)
		}
	case provisioner.ProvisionerTypeTerraform:
		if err := iacrunner.RunTerraform(ctx, cmd, terraform.TerraformOperationType_refresh); err != nil {
			os.Exit(1)
		}
	default:
		cliprint.PrintError("Unknown provisioner type")
		os.Exit(1)
	}
}
