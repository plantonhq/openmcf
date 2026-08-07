package main

import (
	"github.com/plantonhq/planton/catalog/openstack/openstackloadbalancermember/iac/pulumi/module"
	openstackloadbalancermemberv1alpha1 "github.com/plantonhq/planton/catalog/openstack/openstackloadbalancermember/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/stackinput"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stackInput := &openstackloadbalancermemberv1alpha1.OpenStackLoadBalancerMemberStackInput{}

		if err := stackinput.LoadStackInput(ctx, stackInput); err != nil {
			return err
		}

		return module.Resources(ctx, stackInput)
	})
}
