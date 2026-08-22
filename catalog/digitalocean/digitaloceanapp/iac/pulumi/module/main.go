package module

import (
	"github.com/pkg/errors"
	digitaloceanappv1alpha1 "github.com/plantonhq/planton/catalog/digitalocean/digitaloceanapp/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/digitalocean/pulumidigitaloceanprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Resources(
	ctx *pulumi.Context,
	stackInput *digitaloceanappv1alpha1.DigitalOceanAppStackInput,
) error {
	locals := initializeLocals(ctx, stackInput)

	digitalOceanProvider, err := pulumidigitaloceanprovider.Get(
		ctx,
		stackInput.ProviderConfig,
	)
	if err != nil {
		return errors.Wrap(err, "failed to setup digitalocean provider")
	}

	if _, err := app(ctx, locals, digitalOceanProvider); err != nil {
		return errors.Wrap(err, "failed to create digitalocean app")
	}

	return nil
}
