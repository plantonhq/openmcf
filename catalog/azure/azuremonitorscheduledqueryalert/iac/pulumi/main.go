package main

import (
	"github.com/pkg/errors"
	"github.com/plantonhq/planton/catalog/azure/azuremonitorscheduledqueryalert/iac/pulumi/module"
	azuremonitorscheduledqueryalertv1alpha1 "github.com/plantonhq/planton/catalog/azure/azuremonitorscheduledqueryalert/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/stackinput"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stackInput := &azuremonitorscheduledqueryalertv1alpha1.AzureMonitorScheduledQueryAlertStackInput{}

		if err := stackinput.LoadStackInput(ctx, stackInput); err != nil {
			return errors.Wrap(err, "failed to load stack-input")
		}

		return module.Resources(ctx, stackInput)
	})
}
