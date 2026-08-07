package module

import (
	"strings"

	gcpprovider "github.com/plantonhq/planton/catalog/gcp"
	gcpalloydbinstancev1alpha1 "github.com/plantonhq/planton/catalog/gcp/gcpalloydbinstance/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/gcp/gcplabelkeys"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	GcpProviderConfig  *gcpprovider.GcpProviderConfig
	GcpAlloydbInstance *gcpalloydbinstancev1alpha1.GcpAlloydbInstance
	GcpLabels          map[string]string
}

func initializeLocals(_ *pulumi.Context, stackInput *gcpalloydbinstancev1alpha1.GcpAlloydbInstanceStackInput) *Locals {
	locals := &Locals{}
	locals.GcpAlloydbInstance = stackInput.Target
	locals.GcpLabels = map[string]string{
		gcplabelkeys.Resource:     "true",
		gcplabelkeys.ResourceName: locals.GcpAlloydbInstance.Spec.InstanceId,
		gcplabelkeys.ResourceKind: strings.ToLower(cloudresourcekind.CloudResourceKind_GcpAlloydbInstance.String()),
	}

	if locals.GcpAlloydbInstance.Metadata.Org != "" {
		locals.GcpLabels[gcplabelkeys.Organization] = locals.GcpAlloydbInstance.Metadata.Org
	}
	if locals.GcpAlloydbInstance.Metadata.Env != "" {
		locals.GcpLabels[gcplabelkeys.Environment] = locals.GcpAlloydbInstance.Metadata.Env
	}
	if locals.GcpAlloydbInstance.Metadata.Id != "" {
		locals.GcpLabels[gcplabelkeys.ResourceId] = locals.GcpAlloydbInstance.Metadata.Id
	}

	locals.GcpProviderConfig = stackInput.ProviderConfig
	return locals
}
