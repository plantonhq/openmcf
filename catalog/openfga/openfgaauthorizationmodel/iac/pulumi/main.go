package main

import (
	"github.com/pkg/errors"
	"github.com/plantonhq/planton/catalog/openfga/openfgaauthorizationmodel/iac/pulumi/module"
	openfgaauthorizationmodelv1alpha1 "github.com/plantonhq/planton/catalog/openfga/openfgaauthorizationmodel/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/stackinput"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// main is the entry point for the OpenFGA Authorization Model Pulumi module.
//
// IMPORTANT: OpenFGA does not have a Pulumi provider. This module is a
// pass-through placeholder that does not create any resources.
//
// To deploy OpenFGA resources, use Terraform/Tofu as the provisioner:
//
//	planton apply --manifest openfga-authorization-model.yaml --provisioner tofu
func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stackInput := &openfgaauthorizationmodelv1alpha1.OpenFgaAuthorizationModelStackInput{}

		if err := stackinput.LoadStackInput(ctx, stackInput); err != nil {
			return errors.Wrap(err, "failed to load stack-input")
		}

		return module.Resources(ctx, stackInput)
	})
}
