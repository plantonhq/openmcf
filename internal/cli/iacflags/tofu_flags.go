package iacflags

import (
	"github.com/plantonhq/project-planton/apis/org/project_planton/shared/iac/terraform"
	"github.com/plantonhq/project-planton/internal/cli/flag"
	"github.com/spf13/cobra"
)

// AddTofuApplyFlags adds Tofu/Terraform flags for apply and destroy commands.
func AddTofuApplyFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().Bool(string(flag.AutoApprove), false,
		"Skip interactive approval of plan before applying (Tofu/Terraform)")
}

// AddTofuPlanFlags adds Tofu/Terraform flags for the plan command.
func AddTofuPlanFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().Bool(string(flag.Destroy), false,
		"Create a destroy plan instead of apply plan (Tofu/Terraform)")
}

// AddTofuInitFlags adds Tofu/Terraform flags specific to the init command.
// These flags configure state backend settings during initialization.
func AddTofuInitFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().String(string(flag.BackendType),
		terraform.TerraformBackendType_local.String(),
		"Specifies the backend type (Tofu/Terraform) - 'local', 's3', 'gcs', 'azurerm', etc.")

	cmd.PersistentFlags().StringArray(string(flag.BackendConfig), []string{},
		"Backend configuration key=value pairs (Tofu/Terraform)")

	cmd.PersistentFlags().Bool(string(flag.Reconfigure), false,
		"Reconfigure backend, ignoring any saved configuration (Tofu/Terraform)")
}
