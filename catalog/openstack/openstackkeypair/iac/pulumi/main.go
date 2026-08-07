package main

import (
	"github.com/plantonhq/planton/catalog/openstack/openstackkeypair/iac/pulumi/module"
	openstackkeypairv1alpha1 "github.com/plantonhq/planton/catalog/openstack/openstackkeypair/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/stackinput"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stackInput := &openstackkeypairv1alpha1.OpenStackKeypairStackInput{}

		if err := stackinput.LoadStackInput(ctx, stackInput); err != nil {
			return err
		}

		return module.Resources(ctx, stackInput)
	})
}
