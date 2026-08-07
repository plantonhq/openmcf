package main

import (
	"github.com/pkg/errors"
	"github.com/plantonhq/planton/catalog/scaleway/scalewayrdbinstance/iac/pulumi/module"
	scalewayrdbinstancev1alpha1 "github.com/plantonhq/planton/catalog/scaleway/scalewayrdbinstance/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/stackinput"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stackInput := &scalewayrdbinstancev1alpha1.ScalewayRdbInstanceStackInput{}

		if err := stackinput.LoadStackInput(ctx, stackInput); err != nil {
			return errors.Wrap(err, "failed to load stack-input")
		}

		return module.Resources(ctx, stackInput)
	})
}
