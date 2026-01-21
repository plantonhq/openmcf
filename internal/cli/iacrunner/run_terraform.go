package iacrunner

import (
	"github.com/plantonhq/project-planton/apis/org/project_planton/shared/iac/terraform"
	"github.com/plantonhq/project-planton/pkg/iac/provisioner"
	"github.com/spf13/cobra"
)

// RunTerraform executes a Terraform operation using the resolved context.
func RunTerraform(ctx *Context, cmd *cobra.Command, operation terraform.TerraformOperationType) error {
	return runHcl(ctx, cmd, operation, provisioner.HclBinaryTerraform)
}
