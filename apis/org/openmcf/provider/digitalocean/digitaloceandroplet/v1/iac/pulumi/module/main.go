package module

import (
	"github.com/pkg/errors"
	digitaloceandropletv1 "github.com/plantonhq/openmcf/apis/org/openmcf/provider/digitalocean/digitaloceandroplet/v1"
	"github.com/plantonhq/openmcf/pkg/iac/pulumi/pulumimodule/provider/digitalocean/pulumidigitaloceanprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Resources is the module entry point—keeps symmetry with other Planton modules.
func Resources(
	ctx *pulumi.Context,
	stackInput *digitaloceandropletv1.DigitalOceanDropletStackInput,
) error {
	// 1. Prepare locals (metadata, labels, credentials, etc.).
	locals := initializeLocals(ctx, stackInput)

	// 2. DigitalOcean provider from supplied credential.
	digitalOceanProvider, err := pulumidigitaloceanprovider.Get(
		ctx,
		stackInput.ProviderConfig,
	)
	if err != nil {
		return errors.Wrap(err, "failed to setup digitalocean provider")
	}

	// 3. Create the Droplet.
	if _, err := droplet(ctx, locals, digitalOceanProvider); err != nil {
		return errors.Wrap(err, "failed to create droplet")
	}

	return nil
}
