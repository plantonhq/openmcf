package module

import (
	openfgaauthorizationmodelv1 "github.com/plantonhq/openmcf/apis/org/openmcf/provider/openfga/openfgaauthorizationmodel/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Resources is a pass-through placeholder for OpenFGA Authorization Model.
//
// IMPORTANT: OpenFGA does not have a Pulumi provider. This module exists only
// to maintain consistency with the OpenMCF deployment component structure.
// It does not create any resources.
//
// To deploy OpenFGA resources, use Terraform/Tofu as the provisioner:
//
//	openmcf apply --manifest openfga-authorization-model.yaml --provisioner tofu
//
// Reference: https://github.com/openfga/terraform-provider-openfga
func Resources(ctx *pulumi.Context, stackInput *openfgaauthorizationmodelv1.OpenFgaAuthorizationModelStackInput) error {
	// Log that this is a pass-through module
	ctx.Log.Warn("OpenFGA does not have a Pulumi provider. This module is a pass-through placeholder.", nil)
	ctx.Log.Warn("Use Terraform/Tofu as the provisioner to deploy OpenFGA resources.", nil)
	ctx.Log.Info("To deploy: openmcf apply --manifest <manifest.yaml> --provisioner tofu", nil)

	// Export empty outputs to indicate no resources were created
	ctx.Export("notice", pulumi.String("OpenFGA Authorization Model was not created. No Pulumi provider available. Use Terraform/Tofu provisioner."))

	return nil
}
