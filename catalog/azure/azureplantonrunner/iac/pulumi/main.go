// Package main provides the Pulumi program entrypoint for the Azure
// Planton Runner appliance: a standing, outbound-only runner on Azure
// Container Apps that executes deploy and cloud operations from inside
// your network perimeter.
package main

import (
	"github.com/pkg/errors"
	"github.com/plantonhq/planton/catalog/azure/azureplantonrunner/iac/pulumi/module"
	azureplantonrunnerv1alpha1 "github.com/plantonhq/planton/catalog/azure/azureplantonrunner/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/stackinput"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stackInput := &azureplantonrunnerv1alpha1.AzurePlantonRunnerStackInput{}

		if err := stackinput.LoadStackInput(ctx, stackInput); err != nil {
			return errors.Wrap(err, "failed to load stack-input")
		}

		return module.Resources(ctx, stackInput)
	})
}
