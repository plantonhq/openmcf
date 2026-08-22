package module

import (
	"github.com/pkg/errors"
	digitaloceanprojectv1alpha1 "github.com/plantonhq/planton/catalog/digitalocean/digitaloceanproject/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/digitalocean/pulumidigitaloceanprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Resources is the module entry point.
func Resources(
	ctx *pulumi.Context,
	stackInput *digitaloceanprojectv1alpha1.DigitalOceanProjectStackInput,
) error {
	// 1. Prepare locals (target handle).
	locals := initializeLocals(ctx, stackInput)

	// 2. Create a DigitalOcean provider from the supplied credential.
	digitalOceanProvider, err := pulumidigitaloceanprovider.Get(
		ctx,
		stackInput.ProviderConfig,
	)
	if err != nil {
		return errors.Wrap(err, "failed to setup digitalocean provider")
	}

	// 3. Create the project.
	if _, err := project(ctx, locals, digitalOceanProvider); err != nil {
		return errors.Wrap(err, "failed to create project")
	}

	return nil
}
