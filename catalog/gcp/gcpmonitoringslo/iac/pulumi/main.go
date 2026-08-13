package main

import (
	"github.com/plantonhq/planton/catalog/gcp/gcpmonitoringslo/iac/pulumi/module"
	gcpmonitoringslov1alpha1 "github.com/plantonhq/planton/catalog/gcp/gcpmonitoringslo/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/stackinput"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stackInput := &gcpmonitoringslov1alpha1.GcpMonitoringSloStackInput{}
		if err := stackinput.LoadStackInput(ctx, stackInput); err != nil {
			return err
		}
		return module.Resources(ctx, stackInput)
	})
}
