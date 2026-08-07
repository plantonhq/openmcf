package module

import (
	openstackprovider "github.com/plantonhq/planton/catalog/openstack"
	openstackvolumeattachv1alpha1 "github.com/plantonhq/planton/catalog/openstack/openstackvolumeattach/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals bundles the data we need throughout the module.
type Locals struct {
	OpenStackProviderConfig *openstackprovider.OpenStackProviderConfig
	OpenStackVolumeAttach   *openstackvolumeattachv1alpha1.OpenStackVolumeAttach
	// InstanceId is the resolved instance ID from the StringValueOrRef.
	InstanceId string
	// VolumeId is the resolved volume ID from the StringValueOrRef.
	VolumeId string
}

// initializeLocals copies fields from the stack input into Locals.
func initializeLocals(_ *pulumi.Context, stackInput *openstackvolumeattachv1alpha1.OpenStackVolumeAttachStackInput) *Locals {
	locals := &Locals{}

	locals.OpenStackVolumeAttach = stackInput.Target
	locals.OpenStackProviderConfig = stackInput.ProviderConfig

	// Extract instance_id from StringValueOrRef (required field).
	// At runtime, the value is resolved by the FK resolver middleware.
	locals.InstanceId = stackInput.Target.Spec.InstanceId.GetValue()

	// Extract volume_id from StringValueOrRef (required field).
	// At runtime, the value is resolved by the FK resolver middleware.
	locals.VolumeId = stackInput.Target.Spec.VolumeId.GetValue()

	return locals
}
