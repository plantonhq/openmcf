package module

import (
	gcpprovider "github.com/plantonhq/planton/apis/dev/planton/provider/gcp"
	gcpcloudcomposeruserworkloadsconfigmapv1 "github.com/plantonhq/planton/apis/dev/planton/provider/gcp/gcpcloudcomposeruserworkloadsconfigmap/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	GcpProviderConfig                      *gcpprovider.GcpProviderConfig
	GcpCloudComposerUserWorkloadsConfigMap *gcpcloudcomposeruserworkloadsconfigmapv1.GcpCloudComposerUserWorkloadsConfigMap
}

func initializeLocals(_ *pulumi.Context, stackInput *gcpcloudcomposeruserworkloadsconfigmapv1.GcpCloudComposerUserWorkloadsConfigMapStackInput) *Locals {
	locals := &Locals{}
	locals.GcpCloudComposerUserWorkloadsConfigMap = stackInput.Target

	// Kubernetes ConfigMaps carry no GCP labels surface — no platform
	// attribution labels are stamped, identically on both engines.

	locals.GcpProviderConfig = stackInput.ProviderConfig
	return locals
}
