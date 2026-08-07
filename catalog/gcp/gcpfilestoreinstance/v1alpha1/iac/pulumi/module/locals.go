package module

import (
	"strings"

	gcpprovider "github.com/plantonhq/planton/catalog/gcp"
	gcpfilestoreinstancev1alpha1 "github.com/plantonhq/planton/catalog/gcp/gcpfilestoreinstance/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/gcp/gcplabelkeys"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	GcpProviderConfig    *gcpprovider.GcpProviderConfig
	GcpFilestoreInstance *gcpfilestoreinstancev1alpha1.GcpFilestoreInstance
	GcpLabels            map[string]string
	// InstanceName is the cloud-side name: spec.instance_name when set,
	// metadata.name otherwise — the same explicit conditional as the
	// Terraform module, so both engines derive the identical name.
	InstanceName string
}

func initializeLocals(_ *pulumi.Context, stackInput *gcpfilestoreinstancev1alpha1.GcpFilestoreInstanceStackInput) *Locals {
	locals := &Locals{}
	locals.GcpFilestoreInstance = stackInput.Target
	locals.GcpProviderConfig = stackInput.ProviderConfig

	locals.InstanceName = stackInput.Target.Spec.InstanceName
	if locals.InstanceName == "" {
		locals.InstanceName = stackInput.Target.Metadata.Name
	}

	// User labels first so platform attribution labels win on key
	// conflicts — identical merge order to the Terraform module.
	locals.GcpLabels = map[string]string{}
	for key, value := range stackInput.Target.Spec.Labels {
		locals.GcpLabels[key] = value
	}
	locals.GcpLabels[gcplabelkeys.Resource] = "true"
	locals.GcpLabels[gcplabelkeys.ResourceName] = locals.InstanceName
	locals.GcpLabels[gcplabelkeys.ResourceKind] = strings.ToLower(cloudresourcekind.CloudResourceKind_GcpFilestoreInstance.String())

	if locals.GcpFilestoreInstance.Metadata.Org != "" {
		locals.GcpLabels[gcplabelkeys.Organization] = locals.GcpFilestoreInstance.Metadata.Org
	}
	if locals.GcpFilestoreInstance.Metadata.Env != "" {
		locals.GcpLabels[gcplabelkeys.Environment] = locals.GcpFilestoreInstance.Metadata.Env
	}
	if locals.GcpFilestoreInstance.Metadata.Id != "" {
		locals.GcpLabels[gcplabelkeys.ResourceId] = locals.GcpFilestoreInstance.Metadata.Id
	}

	return locals
}
