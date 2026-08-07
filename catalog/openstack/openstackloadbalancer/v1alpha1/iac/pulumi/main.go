package main

import (
	openstackloadbalancerv1alpha1 "github.com/plantonhq/planton/catalog/openstack/openstackloadbalancer/v1alpha1"
	"github.com/plantonhq/planton/catalog/openstack/openstackloadbalancer/v1alpha1/iac/pulumi/module"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/stackinput"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stackInput := &openstackloadbalancerv1alpha1.OpenStackLoadBalancerStackInput{}

		if err := stackinput.LoadStackInput(ctx, stackInput); err != nil {
			return err
		}

		return module.Resources(ctx, stackInput)
	})
}
