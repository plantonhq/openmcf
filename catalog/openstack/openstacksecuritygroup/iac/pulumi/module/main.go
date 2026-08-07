package module

import (
	"github.com/pkg/errors"
	openstacksecuritygroupv1alpha1 "github.com/plantonhq/planton/catalog/openstack/openstacksecuritygroup/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/openstack/pulumiopenstackprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Resources is the entry point called by the Planton CLI.
func Resources(
	ctx *pulumi.Context,
	stackInput *openstacksecuritygroupv1alpha1.OpenStackSecurityGroupStackInput,
) error {
	// 1. Gather handy references.
	locals := initializeLocals(ctx, stackInput)

	// 2. Build a Pulumi OpenStack provider from the supplied credential.
	openstackProvider, err := pulumiopenstackprovider.Get(ctx, stackInput.ProviderConfig)
	if err != nil {
		return errors.Wrap(err, "failed to setup openstack provider")
	}

	// 3. Create the security group and inline rules.
	if err := securityGroup(ctx, locals, openstackProvider); err != nil {
		return errors.Wrap(err, "failed to create openstack security group")
	}

	return nil
}
