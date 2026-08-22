// Package main provides the Pulumi program entrypoint for the GCP Planton
// Runner appliance: a standing, outbound-only runner on Cloud Run that
// executes deploy and cloud operations from inside your project's network
// perimeter.
package main

import (
	"github.com/pkg/errors"
	"github.com/plantonhq/planton/catalog/gcp/gcpplantonrunner/iac/pulumi/module"
	gcpplantonrunnerv1alpha1 "github.com/plantonhq/planton/catalog/gcp/gcpplantonrunner/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/stackinput"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stackInput := &gcpplantonrunnerv1alpha1.GcpPlantonRunnerStackInput{}

		if err := stackinput.LoadStackInput(ctx, stackInput); err != nil {
			return errors.Wrap(err, "failed to load stack-input")
		}

		return module.Resources(ctx, stackInput)
	})
}
