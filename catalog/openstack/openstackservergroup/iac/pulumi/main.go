package main

import (
	"github.com/plantonhq/planton/catalog/openstack/openstackservergroup/iac/pulumi/module"
	openstackservergroupv1alpha1 "github.com/plantonhq/planton/catalog/openstack/openstackservergroup/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/stackinput"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stackInput := &openstackservergroupv1alpha1.OpenStackServerGroupStackInput{}

		if err := stackinput.LoadStackInput(ctx, stackInput); err != nil {
			return err
		}

		return module.Resources(ctx, stackInput)
	})
}
