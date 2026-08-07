package main

import (
	"github.com/pkg/errors"
	hetznercloudplacementgroupv1alpha1 "github.com/plantonhq/planton/catalog/hetznercloud/hetznercloudplacementgroup/v1alpha1"
	"github.com/plantonhq/planton/catalog/hetznercloud/hetznercloudplacementgroup/v1alpha1/iac/pulumi/module"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/stackinput"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stackInput := &hetznercloudplacementgroupv1alpha1.HetznerCloudPlacementGroupStackInput{}

		if err := stackinput.LoadStackInput(ctx, stackInput); err != nil {
			return errors.Wrap(err, "failed to load stack-input")
		}

		return module.Resources(ctx, stackInput)
	})
}
