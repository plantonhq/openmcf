package module

import (
	openstackprovider "github.com/plantonhq/planton/catalog/openstack"
	openstackservergroupv1alpha1 "github.com/plantonhq/planton/catalog/openstack/openstackservergroup/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals bundles the data we need throughout the module.
type Locals struct {
	OpenStackProviderConfig *openstackprovider.OpenStackProviderConfig
	OpenStackServerGroup    *openstackservergroupv1alpha1.OpenStackServerGroup
}

// initializeLocals copies fields from the stack input into Locals.
func initializeLocals(_ *pulumi.Context, stackInput *openstackservergroupv1alpha1.OpenStackServerGroupStackInput) *Locals {
	locals := &Locals{}

	locals.OpenStackServerGroup = stackInput.Target
	locals.OpenStackProviderConfig = stackInput.ProviderConfig

	return locals
}
