package iacrunner

import (
	"os"

	"github.com/pkg/errors"
	"github.com/plantonhq/project-planton/apis/org/project_planton/shared/iac/terraform"
	"github.com/plantonhq/project-planton/internal/cli/cliprint"
	"github.com/plantonhq/project-planton/internal/cli/flag"
	"github.com/plantonhq/project-planton/pkg/iac/provisioner"
	"github.com/plantonhq/project-planton/pkg/iac/tofu/tofumodule"
	"github.com/spf13/cobra"
)

// RunTofu executes an OpenTofu operation using the resolved context.
func RunTofu(ctx *Context, cmd *cobra.Command, operation terraform.TerraformOperationType) error {
	return runHcl(ctx, cmd, operation, provisioner.HclBinaryTofu)
}

// runHcl executes an HCL-based IaC operation (tofu or terraform) using the resolved context.
func runHcl(ctx *Context, cmd *cobra.Command, operation terraform.TerraformOperationType, binary provisioner.HclBinary) error {
	isAutoApprove, err := cmd.Flags().GetBool(string(flag.AutoApprove))
	if err != nil {
		return errors.Wrap(err, "failed to get auto-approve flag")
	}

	// For plan operation, check if it's a destroy plan
	isDestroyPlan := false
	if operation == terraform.TerraformOperationType_plan {
		isDestroyPlan, _ = cmd.Flags().GetBool(string(flag.Destroy))
		// Plan is always auto-approve (non-interactive)
		isAutoApprove = true
	}

	// Check if binary is available before proceeding
	if err := binary.CheckAvailable(); err != nil {
		cliprint.PrintError(err.Error())
		os.Exit(1)
	}

	cliprint.PrintHandoff(binary.DisplayName())

	err = tofumodule.RunCommand(
		string(binary),
		ctx.ModuleDir,
		ctx.ManifestPath,
		operation,
		ctx.ValueOverrides,
		isAutoApprove,
		isDestroyPlan,
		ctx.ModuleVersion,
		ctx.NoCleanup,
		ctx.KubeContext,
		ctx.ProviderConfigOpts...,
	)
	if err != nil {
		printHclFailure(binary)
		os.Exit(1)
	}

	printHclSuccess(binary)
	return nil
}

// printHclSuccess prints a success message for the appropriate binary.
func printHclSuccess(binary provisioner.HclBinary) {
	switch binary {
	case provisioner.HclBinaryTofu:
		cliprint.PrintTofuSuccess()
	case provisioner.HclBinaryTerraform:
		cliprint.PrintTerraformSuccess()
	}
}

// printHclFailure prints a failure message for the appropriate binary.
func printHclFailure(binary provisioner.HclBinary) {
	switch binary {
	case provisioner.HclBinaryTofu:
		cliprint.PrintTofuFailure()
	case provisioner.HclBinaryTerraform:
		cliprint.PrintTerraformFailure()
	}
}
