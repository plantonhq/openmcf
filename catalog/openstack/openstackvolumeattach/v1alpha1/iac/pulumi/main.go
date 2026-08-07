package main

import (
	openstackvolumeattachv1alpha1 "github.com/plantonhq/planton/catalog/openstack/openstackvolumeattach/v1alpha1"
	"github.com/plantonhq/planton/catalog/openstack/openstackvolumeattach/v1alpha1/iac/pulumi/module"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/stackinput"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stackInput := &openstackvolumeattachv1alpha1.OpenStackVolumeAttachStackInput{}

		if err := stackinput.LoadStackInput(ctx, stackInput); err != nil {
			return err
		}

		return module.Resources(ctx, stackInput)
	})
}
