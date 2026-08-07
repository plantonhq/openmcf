package module

import (
	openstackprovider "github.com/plantonhq/planton/catalog/openstack"
	openstackkeypairv1alpha1 "github.com/plantonhq/planton/catalog/openstack/openstackkeypair/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals bundles the data we need throughout the module.
type Locals struct {
	OpenStackProviderConfig *openstackprovider.OpenStackProviderConfig
	OpenStackKeypair        *openstackkeypairv1alpha1.OpenStackKeypair
}

// initializeLocals copies fields from the stack input into Locals.
func initializeLocals(_ *pulumi.Context, stackInput *openstackkeypairv1alpha1.OpenStackKeypairStackInput) *Locals {
	locals := &Locals{}

	locals.OpenStackKeypair = stackInput.Target
	locals.OpenStackProviderConfig = stackInput.ProviderConfig

	return locals
}
