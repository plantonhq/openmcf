package module

import (
	"strings"

	gcpprovider "github.com/plantonhq/planton/catalog/gcp"
	gcpbigtableinstancev1alpha1 "github.com/plantonhq/planton/catalog/gcp/gcpbigtableinstance/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/gcp/gcplabelkeys"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	GcpProviderConfig   *gcpprovider.GcpProviderConfig
	GcpBigtableInstance *gcpbigtableinstancev1alpha1.GcpBigtableInstance
	GcpLabels           map[string]string
}

func initializeLocals(_ *pulumi.Context, stackInput *gcpbigtableinstancev1alpha1.GcpBigtableInstanceStackInput) *Locals {
	locals := &Locals{}
	locals.GcpBigtableInstance = stackInput.Target

	// User labels first so platform attribution labels win on key
	// conflicts — identical merge order to the Terraform module.
	locals.GcpLabels = map[string]string{}
	for key, value := range locals.GcpBigtableInstance.Spec.Labels {
		locals.GcpLabels[key] = value
	}
	locals.GcpLabels[gcplabelkeys.Resource] = "true"
	locals.GcpLabels[gcplabelkeys.ResourceName] = locals.GcpBigtableInstance.Spec.InstanceName
	locals.GcpLabels[gcplabelkeys.ResourceKind] = strings.ToLower(cloudresourcekind.CloudResourceKind_GcpBigtableInstance.String())

	if locals.GcpBigtableInstance.Metadata.Org != "" {
		locals.GcpLabels[gcplabelkeys.Organization] = locals.GcpBigtableInstance.Metadata.Org
	}
	if locals.GcpBigtableInstance.Metadata.Env != "" {
		locals.GcpLabels[gcplabelkeys.Environment] = locals.GcpBigtableInstance.Metadata.Env
	}
	if locals.GcpBigtableInstance.Metadata.Id != "" {
		locals.GcpLabels[gcplabelkeys.ResourceId] = locals.GcpBigtableInstance.Metadata.Id
	}

	locals.GcpProviderConfig = stackInput.ProviderConfig
	return locals
}
