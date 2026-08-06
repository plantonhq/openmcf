package module

import (
	gcpprovider "github.com/plantonhq/planton/apis/dev/planton/provider/gcp"
	gcpvertexaideployedindexv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/gcp/gcpvertexaideployedindex/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	GcpProviderConfig        *gcpprovider.GcpProviderConfig
	GcpVertexAiDeployedIndex *gcpvertexaideployedindexv1alpha1.GcpVertexAiDeployedIndex
}

func initializeLocals(_ *pulumi.Context, stackInput *gcpvertexaideployedindexv1alpha1.GcpVertexAiDeployedIndexStackInput) *Locals {
	locals := &Locals{}
	locals.GcpVertexAiDeployedIndex = stackInput.Target
	locals.GcpProviderConfig = stackInput.ProviderConfig

	// This resource class carries NO labels and NO project field in the
	// GCP API — the deployment lives inside the index endpoint resource
	// and inherits its project — so there is no label merge here:
	// platform attribution is impossible on a DeployedIndex and none is
	// faked.

	return locals
}
