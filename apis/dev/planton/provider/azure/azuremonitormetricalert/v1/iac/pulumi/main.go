package main

import (
	"github.com/pkg/errors"
	azuremonitormetricalertv1 "github.com/plantonhq/planton/apis/dev/planton/provider/azure/azuremonitormetricalert/v1"
	"github.com/plantonhq/planton/apis/dev/planton/provider/azure/azuremonitormetricalert/v1/iac/pulumi/module"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/stackinput"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stackInput := &azuremonitormetricalertv1.AzureMonitorMetricAlertStackInput{}

		if err := stackinput.LoadStackInput(ctx, stackInput); err != nil {
			return errors.Wrap(err, "failed to load stack-input")
		}

		return module.Resources(ctx, stackInput)
	})
}
