package main

import (
	"github.com/plantonhq/planton/catalog/openstack/openstackloadbalancerlistener/iac/pulumi/module"
	openstackloadbalancerlistenerv1alpha1 "github.com/plantonhq/planton/catalog/openstack/openstackloadbalancerlistener/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/stackinput"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stackInput := &openstackloadbalancerlistenerv1alpha1.OpenStackLoadBalancerListenerStackInput{}

		if err := stackinput.LoadStackInput(ctx, stackInput); err != nil {
			return err
		}

		return module.Resources(ctx, stackInput)
	})
}
