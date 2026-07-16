package module

import (
	"strings"

	gcpprovider "github.com/plantonhq/planton/apis/dev/planton/provider/gcp"
	gcpvertexainotebookv1 "github.com/plantonhq/planton/apis/dev/planton/provider/gcp/gcpvertexainotebook/v1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/gcp/gcplabelkeys"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	GcpProviderConfig   *gcpprovider.GcpProviderConfig
	GcpVertexAiNotebook *gcpvertexainotebookv1.GcpVertexAiNotebook
	GcpLabels           map[string]string
	InstanceName        string
}

func initializeLocals(_ *pulumi.Context, stackInput *gcpvertexainotebookv1.GcpVertexAiNotebookStackInput) *Locals {
	locals := &Locals{}
	locals.GcpVertexAiNotebook = stackInput.Target
	locals.GcpProviderConfig = stackInput.ProviderConfig

	// Determine the instance name: explicit instance_name or fall back to metadata.name.
	locals.InstanceName = locals.GcpVertexAiNotebook.Spec.InstanceName
	if locals.InstanceName == "" {
		locals.InstanceName = locals.GcpVertexAiNotebook.Metadata.Name
	}

	// User labels first so platform attribution labels win on key
	// conflicts — identical merge order to the Terraform module.
	locals.GcpLabels = map[string]string{}
	for key, value := range locals.GcpVertexAiNotebook.Spec.Labels {
		locals.GcpLabels[key] = value
	}
	locals.GcpLabels[gcplabelkeys.Resource] = "true"
	locals.GcpLabels[gcplabelkeys.ResourceName] = locals.InstanceName
	locals.GcpLabels[gcplabelkeys.ResourceKind] = strings.ToLower(cloudresourcekind.CloudResourceKind_GcpVertexAiNotebook.String())

	if locals.GcpVertexAiNotebook.Metadata.Org != "" {
		locals.GcpLabels[gcplabelkeys.Organization] = locals.GcpVertexAiNotebook.Metadata.Org
	}
	if locals.GcpVertexAiNotebook.Metadata.Env != "" {
		locals.GcpLabels[gcplabelkeys.Environment] = locals.GcpVertexAiNotebook.Metadata.Env
	}
	if locals.GcpVertexAiNotebook.Metadata.Id != "" {
		locals.GcpLabels[gcplabelkeys.ResourceId] = locals.GcpVertexAiNotebook.Metadata.Id
	}

	return locals
}
