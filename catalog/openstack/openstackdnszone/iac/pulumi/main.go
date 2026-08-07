package main

import (
	"github.com/plantonhq/planton/catalog/openstack/openstackdnszone/iac/pulumi/module"
	openstackdnszonev1alpha1 "github.com/plantonhq/planton/catalog/openstack/openstackdnszone/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/stackinput"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stackInput := &openstackdnszonev1alpha1.OpenStackDnsZoneStackInput{}

		if err := stackinput.LoadStackInput(ctx, stackInput); err != nil {
			return err
		}

		return module.Resources(ctx, stackInput)
	})
}
