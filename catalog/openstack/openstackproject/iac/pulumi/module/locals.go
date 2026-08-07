package module

import (
	openstackprovider "github.com/plantonhq/planton/catalog/openstack"
	openstackprojectv1alpha1 "github.com/plantonhq/planton/catalog/openstack/openstackproject/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals bundles the data we need throughout the module.
type Locals struct {
	OpenStackProviderConfig *openstackprovider.OpenStackProviderConfig
	OpenStackProject        *openstackprojectv1alpha1.OpenStackProject
}

// initializeLocals copies fields from the stack input into Locals.
func initializeLocals(_ *pulumi.Context, stackInput *openstackprojectv1alpha1.OpenStackProjectStackInput) *Locals {
	locals := &Locals{}

	locals.OpenStackProject = stackInput.Target
	locals.OpenStackProviderConfig = stackInput.ProviderConfig

	return locals
}
