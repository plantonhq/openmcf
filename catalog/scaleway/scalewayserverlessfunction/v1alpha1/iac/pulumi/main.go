package main

import (
	"github.com/pkg/errors"
	scalewayserverlessfunctionv1alpha1 "github.com/plantonhq/planton/catalog/scaleway/scalewayserverlessfunction/v1alpha1"
	"github.com/plantonhq/planton/catalog/scaleway/scalewayserverlessfunction/v1alpha1/iac/pulumi/module"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/stackinput"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stackInput := &scalewayserverlessfunctionv1alpha1.ScalewayServerlessFunctionStackInput{}

		if err := stackinput.LoadStackInput(ctx, stackInput); err != nil {
			return errors.Wrap(err, "failed to load stack-input")
		}

		return module.Resources(ctx, stackInput)
	})
}
