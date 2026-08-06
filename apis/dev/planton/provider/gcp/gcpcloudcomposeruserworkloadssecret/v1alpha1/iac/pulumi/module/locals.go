package module

import (
	gcpprovider "github.com/plantonhq/planton/apis/dev/planton/provider/gcp"
	gcpcloudcomposeruserworkloadssecretv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/gcp/gcpcloudcomposeruserworkloadssecret/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	GcpProviderConfig                   *gcpprovider.GcpProviderConfig
	GcpCloudComposerUserWorkloadsSecret *gcpcloudcomposeruserworkloadssecretv1alpha1.GcpCloudComposerUserWorkloadsSecret
}

func initializeLocals(_ *pulumi.Context, stackInput *gcpcloudcomposeruserworkloadssecretv1alpha1.GcpCloudComposerUserWorkloadsSecretStackInput) *Locals {
	locals := &Locals{}
	locals.GcpCloudComposerUserWorkloadsSecret = stackInput.Target

	// Kubernetes Secrets carry no GCP labels surface — no platform
	// attribution labels are stamped, identically on both engines.

	locals.GcpProviderConfig = stackInput.ProviderConfig
	return locals
}
