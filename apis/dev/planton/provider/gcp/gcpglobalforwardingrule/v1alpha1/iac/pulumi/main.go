package main

import (
	gcpglobalforwardingrulev1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/gcp/gcpglobalforwardingrule/v1alpha1"
	"github.com/plantonhq/planton/apis/dev/planton/provider/gcp/gcpglobalforwardingrule/v1alpha1/iac/pulumi/module"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/stackinput"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stackInput := &gcpglobalforwardingrulev1alpha1.GcpGlobalForwardingRuleStackInput{}
		if err := stackinput.LoadStackInput(ctx, stackInput); err != nil {
			return err
		}
		return module.Resources(ctx, stackInput)
	})
}
