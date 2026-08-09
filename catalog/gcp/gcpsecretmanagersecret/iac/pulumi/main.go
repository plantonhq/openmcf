package main

import (
	"github.com/plantonhq/planton/catalog/gcp/gcpsecretmanagersecret/iac/pulumi/module"
	gcpsecretmanagersecretv1alpha1 "github.com/plantonhq/planton/catalog/gcp/gcpsecretmanagersecret/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/stackinput"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stackInput := &gcpsecretmanagersecretv1alpha1.GcpSecretManagerSecretStackInput{}
		if err := stackinput.LoadStackInput(ctx, stackInput); err != nil {
			return err
		}
		return module.Resources(ctx, stackInput)
	})
}
