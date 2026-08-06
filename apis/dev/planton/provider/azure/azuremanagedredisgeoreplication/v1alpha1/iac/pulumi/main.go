package main

import (
	"github.com/pkg/errors"
	azuremanagedredisgeoreplicationv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/azure/azuremanagedredisgeoreplication/v1alpha1"
	"github.com/plantonhq/planton/apis/dev/planton/provider/azure/azuremanagedredisgeoreplication/v1alpha1/iac/pulumi/module"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/stackinput"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stackInput := &azuremanagedredisgeoreplicationv1alpha1.AzureManagedRedisGeoReplicationStackInput{}

		if err := stackinput.LoadStackInput(ctx, stackInput); err != nil {
			return errors.Wrap(err, "failed to load stack-input")
		}

		return module.Resources(ctx, stackInput)
	})
}
