package main

import (
	"github.com/pkg/errors"
	"github.com/plantonhq/planton/catalog/azure/azurefabriccapacity/iac/pulumi/module"
	azurefabriccapacityv1alpha1 "github.com/plantonhq/planton/catalog/azure/azurefabriccapacity/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/stackinput"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stackInput := &azurefabriccapacityv1alpha1.AzureFabricCapacityStackInput{}

		if err := stackinput.LoadStackInput(ctx, stackInput); err != nil {
			return errors.Wrap(err, "failed to load stack-input")
		}

		return module.Resources(ctx, stackInput)
	})
}
