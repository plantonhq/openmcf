package main

import (
	"github.com/plantonhq/planton/catalog/openstack/openstacknetworkport/iac/pulumi/module"
	openstacknetworkportv1alpha1 "github.com/plantonhq/planton/catalog/openstack/openstacknetworkport/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/stackinput"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stackInput := &openstacknetworkportv1alpha1.OpenStackNetworkPortStackInput{}

		if err := stackinput.LoadStackInput(ctx, stackInput); err != nil {
			return err
		}

		return module.Resources(ctx, stackInput)
	})
}
