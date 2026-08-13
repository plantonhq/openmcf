package main

import (
	"github.com/pkg/errors"
	"github.com/plantonhq/planton/catalog/azure/azurenetworkwatcherflowlog/iac/pulumi/module"
	azurenetworkwatcherflowlogv1alpha1 "github.com/plantonhq/planton/catalog/azure/azurenetworkwatcherflowlog/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/stackinput"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stackInput := &azurenetworkwatcherflowlogv1alpha1.AzureNetworkWatcherFlowLogStackInput{}

		if err := stackinput.LoadStackInput(ctx, stackInput); err != nil {
			return errors.Wrap(err, "failed to load stack-input")
		}

		return module.Resources(ctx, stackInput)
	})
}
