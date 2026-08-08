package main

import (
	"github.com/pkg/errors"
	"github.com/plantonhq/planton/catalog/azure/azurevirtualwan/iac/pulumi/module"
	azurevirtualwanv1alpha1 "github.com/plantonhq/planton/catalog/azure/azurevirtualwan/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/stackinput"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stackInput := &azurevirtualwanv1alpha1.AzureVirtualWanStackInput{}

		if err := stackinput.LoadStackInput(ctx, stackInput); err != nil {
			return errors.Wrap(err, "failed to load stack-input")
		}

		return module.Resources(ctx, stackInput)
	})
}
