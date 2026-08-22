package main

import (
	"github.com/pkg/errors"
	"github.com/plantonhq/planton/catalog/digitalocean/digitaloceandatabasedb/iac/pulumi/module"
	digitaloceandatabasedbv1alpha1 "github.com/plantonhq/planton/catalog/digitalocean/digitaloceandatabasedb/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/stackinput"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stackInput := &digitaloceandatabasedbv1alpha1.DigitalOceanDatabaseDbStackInput{}

		if err := stackinput.LoadStackInput(ctx, stackInput); err != nil {
			return errors.Wrap(err, "failed to load stack-input")
		}

		return module.Resources(ctx, stackInput)
	})
}
