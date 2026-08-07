package main

import (
	"github.com/pkg/errors"
	"github.com/plantonhq/planton/catalog/oci/ociblockvolume/iac/pulumi/module"
	ociblockvolumev1alpha1 "github.com/plantonhq/planton/catalog/oci/ociblockvolume/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/stackinput"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stackInput := &ociblockvolumev1alpha1.OciBlockVolumeStackInput{}

		if err := stackinput.LoadStackInput(ctx, stackInput); err != nil {
			return errors.Wrap(err, "failed to load stack-input")
		}

		return module.Resources(ctx, stackInput)
	})
}
