package module

import (
	openstackprovider "github.com/plantonhq/planton/catalog/openstack"
	openstackrav1 "github.com/plantonhq/planton/catalog/openstack/openstackroleassignment/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals bundles the data we need throughout the module.
type Locals struct {
	OpenStackProviderConfig *openstackprovider.OpenStackProviderConfig
	OpenStackRoleAssignment *openstackrav1.OpenStackRoleAssignment
}

// initializeLocals copies fields from the stack input into Locals.
func initializeLocals(_ *pulumi.Context, stackInput *openstackrav1.OpenStackRoleAssignmentStackInput) *Locals {
	locals := &Locals{}

	locals.OpenStackRoleAssignment = stackInput.Target
	locals.OpenStackProviderConfig = stackInput.ProviderConfig

	return locals
}
