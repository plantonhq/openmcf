package main

import (
	"github.com/pkg/errors"
	ocisubnetv1alpha1 "github.com/plantonhq/planton/catalog/oci/ocisubnet/v1alpha1"
	"github.com/plantonhq/planton/catalog/oci/ocisubnet/v1alpha1/iac/pulumi/module"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/stackinput"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stackInput := &ocisubnetv1alpha1.OciSubnetStackInput{}

		if err := stackinput.LoadStackInput(ctx, stackInput); err != nil {
			return errors.Wrap(err, "failed to load stack-input")
		}

		return module.Resources(ctx, stackInput)
	})
}
