package main

import (
	"github.com/pkg/errors"
	"github.com/plantonhq/planton/catalog/atlas/atlasmongodb/iac/pulumi/module"
	atlasmongodbv1alpha1 "github.com/plantonhq/planton/catalog/atlas/atlasmongodb/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/stackinput"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stackInput := &atlasmongodbv1alpha1.AtlasMongodbStackInput{}

		if err := stackinput.LoadStackInput(ctx, stackInput); err != nil {
			return errors.Wrap(err, "failed to load stack-input")
		}

		return module.Resources(ctx, stackInput)
	})
}
