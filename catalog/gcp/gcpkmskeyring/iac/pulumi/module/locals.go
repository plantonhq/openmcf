package module

import (
	gcpprovider "github.com/plantonhq/planton/catalog/gcp"
	gcpkmskeyringv1alpha1 "github.com/plantonhq/planton/catalog/gcp/gcpkmskeyring/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	GcpProviderConfig *gcpprovider.GcpProviderConfig
	GcpKmsKeyRing     *gcpkmskeyringv1alpha1.GcpKmsKeyRing
}

func initializeLocals(_ *pulumi.Context, stackInput *gcpkmskeyringv1alpha1.GcpKmsKeyRingStackInput) *Locals {
	locals := &Locals{}
	locals.GcpKmsKeyRing = stackInput.Target

	// The key ring resource has no labels surface in the Cloud KMS API —
	// no platform attribution labels are computed or stamped, identically
	// on both engines. (Labels live on the crypto keys inside the ring.)

	locals.GcpProviderConfig = stackInput.ProviderConfig
	return locals
}
