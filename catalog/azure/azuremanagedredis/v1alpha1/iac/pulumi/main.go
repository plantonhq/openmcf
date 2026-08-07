package main

import (
	"github.com/pkg/errors"
	azuremanagedredisv1alpha1 "github.com/plantonhq/planton/catalog/azure/azuremanagedredis/v1alpha1"
	"github.com/plantonhq/planton/catalog/azure/azuremanagedredis/v1alpha1/iac/pulumi/module"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/stackinput"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stackInput := &azuremanagedredisv1alpha1.AzureManagedRedisStackInput{}

		if err := stackinput.LoadStackInput(ctx, stackInput); err != nil {
			return errors.Wrap(err, "failed to load stack-input")
		}

		return module.Resources(ctx, stackInput)
	})
}
