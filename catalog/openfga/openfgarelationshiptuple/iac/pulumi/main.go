package main

import (
	"github.com/pkg/errors"
	"github.com/plantonhq/planton/catalog/openfga/openfgarelationshiptuple/iac/pulumi/module"
	openfgarelationshiptuplev1alpha1 "github.com/plantonhq/planton/catalog/openfga/openfgarelationshiptuple/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/stackinput"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// main is the entry point for the OpenFGA Relationship Tuple Pulumi module.
//
// IMPORTANT: OpenFGA does not have a Pulumi provider. This module is a
// pass-through placeholder that does not create any resources.
//
// To deploy OpenFGA resources, use Terraform/Tofu as the provisioner:
//
//	planton apply --manifest openfga-relationship-tuple.yaml --provisioner tofu
func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stackInput := &openfgarelationshiptuplev1alpha1.OpenFgaRelationshipTupleStackInput{}

		if err := stackinput.LoadStackInput(ctx, stackInput); err != nil {
			return errors.Wrap(err, "failed to load stack-input")
		}

		return module.Resources(ctx, stackInput)
	})
}
