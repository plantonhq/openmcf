package module

import (
	"strings"

	gcpprovider "github.com/plantonhq/planton/catalog/gcp"
	gcpcloudcomposerenvironmentv1alpha1 "github.com/plantonhq/planton/catalog/gcp/gcpcloudcomposerenvironment/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/gcp/gcplabelkeys"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals holds pre-computed values used across the module.
type Locals struct {
	GcpProviderConfig           *gcpprovider.GcpProviderConfig
	GcpCloudComposerEnvironment *gcpcloudcomposerenvironmentv1alpha1.GcpCloudComposerEnvironment
	GcpLabels                   map[string]string
}

func initializeLocals(_ *pulumi.Context, stackInput *gcpcloudcomposerenvironmentv1alpha1.GcpCloudComposerEnvironmentStackInput) *Locals {
	locals := &Locals{}
	locals.GcpCloudComposerEnvironment = stackInput.Target

	// Determine resource name for labels.
	resourceName := locals.GcpCloudComposerEnvironment.Spec.EnvironmentName
	if resourceName == "" && locals.GcpCloudComposerEnvironment.Metadata != nil {
		resourceName = locals.GcpCloudComposerEnvironment.Metadata.Name
	}

	// User labels first so platform attribution labels win on key
	// conflicts — identical merge order to the Terraform module.
	locals.GcpLabels = map[string]string{}
	for key, value := range locals.GcpCloudComposerEnvironment.Spec.Labels {
		locals.GcpLabels[key] = value
	}
	locals.GcpLabels[gcplabelkeys.Resource] = "true"
	locals.GcpLabels[gcplabelkeys.ResourceName] = resourceName
	locals.GcpLabels[gcplabelkeys.ResourceKind] = strings.ToLower(cloudresourcekind.CloudResourceKind_GcpCloudComposerEnvironment.String())

	if locals.GcpCloudComposerEnvironment.Metadata.Org != "" {
		locals.GcpLabels[gcplabelkeys.Organization] = locals.GcpCloudComposerEnvironment.Metadata.Org
	}
	if locals.GcpCloudComposerEnvironment.Metadata.Env != "" {
		locals.GcpLabels[gcplabelkeys.Environment] = locals.GcpCloudComposerEnvironment.Metadata.Env
	}
	if locals.GcpCloudComposerEnvironment.Metadata.Id != "" {
		locals.GcpLabels[gcplabelkeys.ResourceId] = locals.GcpCloudComposerEnvironment.Metadata.Id
	}

	locals.GcpProviderConfig = stackInput.ProviderConfig
	return locals
}
