package main

import (
	"github.com/plantonhq/planton/catalog/gcp/gcpiamdenypolicy/iac/pulumi/module"
	gcpiamdenypolicyv1alpha1 "github.com/plantonhq/planton/catalog/gcp/gcpiamdenypolicy/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/stackinput"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stackInput := &gcpiamdenypolicyv1alpha1.GcpIamDenyPolicyStackInput{}
		if err := stackinput.LoadStackInput(ctx, stackInput); err != nil {
			return err
		}
		return module.Resources(ctx, stackInput)
	})
}
