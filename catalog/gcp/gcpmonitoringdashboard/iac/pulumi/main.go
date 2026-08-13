package main

import (
	"github.com/plantonhq/planton/catalog/gcp/gcpmonitoringdashboard/iac/pulumi/module"
	gcpmonitoringdashboardv1alpha1 "github.com/plantonhq/planton/catalog/gcp/gcpmonitoringdashboard/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/stackinput"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stackInput := &gcpmonitoringdashboardv1alpha1.GcpMonitoringDashboardStackInput{}
		if err := stackinput.LoadStackInput(ctx, stackInput); err != nil {
			return err
		}
		return module.Resources(ctx, stackInput)
	})
}
