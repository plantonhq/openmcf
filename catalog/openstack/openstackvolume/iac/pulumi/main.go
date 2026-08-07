package main

import (
	"github.com/plantonhq/planton/catalog/openstack/openstackvolume/iac/pulumi/module"
	openstackvolumev1alpha1 "github.com/plantonhq/planton/catalog/openstack/openstackvolume/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/stackinput"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stackInput := &openstackvolumev1alpha1.OpenStackVolumeStackInput{}

		if err := stackinput.LoadStackInput(ctx, stackInput); err != nil {
			return err
		}

		return module.Resources(ctx, stackInput)
	})
}
