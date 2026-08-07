package module

import (
	"github.com/pkg/errors"
	openstackinstancev1alpha1 "github.com/plantonhq/planton/catalog/openstack/openstackinstance/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/openstack/pulumiopenstackprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Resources is the entry point called by the Planton CLI.
func Resources(
	ctx *pulumi.Context,
	stackInput *openstackinstancev1alpha1.OpenStackInstanceStackInput,
) error {
	// 1. Gather handy references.
	locals := initializeLocals(ctx, stackInput)

	// 2. Build a Pulumi OpenStack provider from the supplied credential.
	openstackProvider, err := pulumiopenstackprovider.Get(ctx, stackInput.ProviderConfig)
	if err != nil {
		return errors.Wrap(err, "failed to setup openstack provider")
	}

	// 3. Create the compute instance.
	if err := instance(ctx, locals, openstackProvider); err != nil {
		return errors.Wrap(err, "failed to create openstack instance")
	}

	return nil
}
