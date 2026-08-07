package main

import (
	"github.com/plantonhq/planton/catalog/openstack/openstackroleassignment/iac/pulumi/module"
	openstackrav1 "github.com/plantonhq/planton/catalog/openstack/openstackroleassignment/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/stackinput"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stackInput := &openstackrav1.OpenStackRoleAssignmentStackInput{}

		if err := stackinput.LoadStackInput(ctx, stackInput); err != nil {
			return err
		}

		return module.Resources(ctx, stackInput)
	})
}
