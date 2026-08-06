package module

import (
	gcpurlmapv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/gcp/gcpurlmap/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals mirrors the Terraform module's locals {} convention: the resolved
// resource plus any derived values the module needs.
type Locals struct {
	GcpUrlMap *gcpurlmapv1alpha1.GcpUrlMap

	// The cloud-side name defaults to metadata.name when the spec leaves
	// url_map_name empty — the same naming basis every kind uses.
	UrlMapName string
}

func initializeLocals(ctx *pulumi.Context, stackInput *gcpurlmapv1alpha1.GcpUrlMapStackInput) *Locals {
	target := stackInput.Target

	urlMapName := target.Spec.UrlMapName
	if urlMapName == "" {
		urlMapName = target.Metadata.Name
	}

	return &Locals{
		GcpUrlMap:  target,
		UrlMapName: urlMapName,
	}
}
