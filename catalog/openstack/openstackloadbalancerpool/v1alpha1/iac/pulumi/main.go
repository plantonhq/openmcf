package main

import (
	openstackloadbalancerpoolv1alpha1 "github.com/plantonhq/planton/catalog/openstack/openstackloadbalancerpool/v1alpha1"
	"github.com/plantonhq/planton/catalog/openstack/openstackloadbalancerpool/v1alpha1/iac/pulumi/module"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/stackinput"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stackInput := &openstackloadbalancerpoolv1alpha1.OpenStackLoadBalancerPoolStackInput{}

		if err := stackinput.LoadStackInput(ctx, stackInput); err != nil {
			return err
		}

		return module.Resources(ctx, stackInput)
	})
}
