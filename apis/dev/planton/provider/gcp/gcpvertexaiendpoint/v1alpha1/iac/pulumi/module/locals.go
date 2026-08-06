package module

import (
	"strings"

	gcpprovider "github.com/plantonhq/planton/apis/dev/planton/provider/gcp"
	gcpvertexaiendpointv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/gcp/gcpvertexaiendpoint/v1alpha1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/gcp/gcplabelkeys"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	GcpProviderConfig   *gcpprovider.GcpProviderConfig
	GcpVertexAiEndpoint *gcpvertexaiendpointv1alpha1.GcpVertexAiEndpoint
	GcpLabels           map[string]string
	DisplayName         string
	EndpointName        string
}

func initializeLocals(_ *pulumi.Context, stackInput *gcpvertexaiendpointv1alpha1.GcpVertexAiEndpointStackInput) *Locals {
	locals := &Locals{}
	locals.GcpVertexAiEndpoint = stackInput.Target
	locals.GcpProviderConfig = stackInput.ProviderConfig

	locals.DisplayName = locals.GcpVertexAiEndpoint.Spec.DisplayName

	// Vertex AI requires the endpoint's name to be numeric and the API will
	// not generate one — when the spec omits endpoint_name, derive a stable
	// ID from the resource identity (identical derivation in the Terraform
	// module, so the same manifest yields the same endpoint ID on either
	// engine).
	locals.EndpointName = locals.GcpVertexAiEndpoint.Spec.EndpointName
	if locals.EndpointName == "" {
		locals.EndpointName = deriveEndpointName(
			locals.GcpVertexAiEndpoint.Metadata.Org,
			locals.GcpVertexAiEndpoint.Metadata.Env,
			locals.GcpVertexAiEndpoint.Metadata.Name,
		)
	}

	// User labels first so platform attribution labels win on key
	// conflicts — identical merge order to the Terraform module.
	locals.GcpLabels = map[string]string{}
	for key, value := range locals.GcpVertexAiEndpoint.Spec.Labels {
		locals.GcpLabels[key] = value
	}
	locals.GcpLabels[gcplabelkeys.Resource] = "true"
	locals.GcpLabels[gcplabelkeys.ResourceName] = strings.ToLower(locals.GcpVertexAiEndpoint.Metadata.Name)
	locals.GcpLabels[gcplabelkeys.ResourceKind] = strings.ToLower(cloudresourcekind.CloudResourceKind_GcpVertexAiEndpoint.String())

	if locals.GcpVertexAiEndpoint.Metadata.Org != "" {
		locals.GcpLabels[gcplabelkeys.Organization] = locals.GcpVertexAiEndpoint.Metadata.Org
	}
	if locals.GcpVertexAiEndpoint.Metadata.Env != "" {
		locals.GcpLabels[gcplabelkeys.Environment] = locals.GcpVertexAiEndpoint.Metadata.Env
	}
	if locals.GcpVertexAiEndpoint.Metadata.Id != "" {
		locals.GcpLabels[gcplabelkeys.ResourceId] = locals.GcpVertexAiEndpoint.Metadata.Id
	}

	return locals
}
