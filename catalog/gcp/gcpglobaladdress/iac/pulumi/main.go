package main

import (
	"github.com/plantonhq/planton/catalog/gcp/gcpglobaladdress/iac/pulumi/module"
	gcpglobaladdressv1alpha1 "github.com/plantonhq/planton/catalog/gcp/gcpglobaladdress/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/stackinput"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stackInput := &gcpglobaladdressv1alpha1.GcpGlobalAddressStackInput{}
		if err := stackinput.LoadStackInput(ctx, stackInput); err != nil {
			return err
		}
		return module.Resources(ctx, stackInput)
	})
}
