package module

import (
	"strings"

	gcpprovider "github.com/plantonhq/planton/apis/dev/planton/provider/gcp"
	gcpmemorystoreinstancev1 "github.com/plantonhq/planton/apis/dev/planton/provider/gcp/gcpmemorystoreinstance/v1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/gcp/gcplabelkeys"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	GcpProviderConfig      *gcpprovider.GcpProviderConfig
	GcpMemorystoreInstance *gcpmemorystoreinstancev1.GcpMemorystoreInstance
	GcpLabels              map[string]string
}

func initializeLocals(_ *pulumi.Context, stackInput *gcpmemorystoreinstancev1.GcpMemorystoreInstanceStackInput) *Locals {
	locals := &Locals{}
	locals.GcpMemorystoreInstance = stackInput.Target

	// User labels first so platform attribution labels win on key
	// conflicts — identical merge order to the Terraform module.
	locals.GcpLabels = map[string]string{}
	for key, value := range locals.GcpMemorystoreInstance.Spec.Labels {
		locals.GcpLabels[key] = value
	}
	locals.GcpLabels[gcplabelkeys.Resource] = "true"
	locals.GcpLabels[gcplabelkeys.ResourceName] = locals.GcpMemorystoreInstance.Spec.InstanceName
	locals.GcpLabels[gcplabelkeys.ResourceKind] = strings.ToLower(cloudresourcekind.CloudResourceKind_GcpMemorystoreInstance.String())

	if locals.GcpMemorystoreInstance.Metadata.Org != "" {
		locals.GcpLabels[gcplabelkeys.Organization] = locals.GcpMemorystoreInstance.Metadata.Org
	}
	if locals.GcpMemorystoreInstance.Metadata.Env != "" {
		locals.GcpLabels[gcplabelkeys.Environment] = locals.GcpMemorystoreInstance.Metadata.Env
	}
	if locals.GcpMemorystoreInstance.Metadata.Id != "" {
		locals.GcpLabels[gcplabelkeys.ResourceId] = locals.GcpMemorystoreInstance.Metadata.Id
	}

	locals.GcpProviderConfig = stackInput.ProviderConfig
	return locals
}
