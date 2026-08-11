package main

import (
	"github.com/plantonhq/planton/catalog/gcp/gcpidentityplatformconfig/iac/pulumi/module"
	gcpidentityplatformconfigv1alpha1 "github.com/plantonhq/planton/catalog/gcp/gcpidentityplatformconfig/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/stackinput"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stackInput := &gcpidentityplatformconfigv1alpha1.GcpIdentityPlatformConfigStackInput{}
		if err := stackinput.LoadStackInput(ctx, stackInput); err != nil {
			return err
		}
		return module.Resources(ctx, stackInput)
	})
}
