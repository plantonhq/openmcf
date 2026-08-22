package main

import (
	"github.com/pkg/errors"
	"github.com/plantonhq/planton/catalog/digitalocean/digitaloceancdn/iac/pulumi/module"
	digitaloceancdnv1alpha1 "github.com/plantonhq/planton/catalog/digitalocean/digitaloceancdn/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/stackinput"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stackInput := &digitaloceancdnv1alpha1.DigitalOceanCdnStackInput{}

		if err := stackinput.LoadStackInput(ctx, stackInput); err != nil {
			return errors.Wrap(err, "failed to load stack-input")
		}

		return module.Resources(ctx, stackInput)
	})
}
