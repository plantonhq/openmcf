package module

import (
	openstackprovider "github.com/plantonhq/planton/catalog/openstack"
	openstacknetworkv1alpha1 "github.com/plantonhq/planton/catalog/openstack/openstacknetwork/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals bundles the data we need throughout the module.
type Locals struct {
	OpenStackProviderConfig *openstackprovider.OpenStackProviderConfig
	OpenStackNetwork        *openstacknetworkv1alpha1.OpenStackNetwork
}

// initializeLocals copies fields from the stack input into Locals.
func initializeLocals(_ *pulumi.Context, stackInput *openstacknetworkv1alpha1.OpenStackNetworkStackInput) *Locals {
	locals := &Locals{}

	locals.OpenStackNetwork = stackInput.Target
	locals.OpenStackProviderConfig = stackInput.ProviderConfig

	return locals
}
