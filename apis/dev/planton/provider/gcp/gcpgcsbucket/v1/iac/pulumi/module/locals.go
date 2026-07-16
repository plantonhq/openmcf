package module

import (
	"strings"

	gcpprovider "github.com/plantonhq/planton/apis/dev/planton/provider/gcp"
	gcpgcsbucketv1 "github.com/plantonhq/planton/apis/dev/planton/provider/gcp/gcpgcsbucket/v1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/gcp/gcplabelkeys"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	GcpProviderConfig *gcpprovider.GcpProviderConfig
	GcpGcsBucket      *gcpgcsbucketv1.GcpGcsBucket
	GcpLabels         map[string]string
}

func initializeLocals(_ *pulumi.Context, stackInput *gcpgcsbucketv1.GcpGcsBucketStackInput) *Locals {
	locals := &Locals{}
	locals.GcpGcsBucket = stackInput.Target

	// User labels first so platform attribution labels win on key
	// conflicts — identical merge order to the Terraform module.
	locals.GcpLabels = map[string]string{}
	for key, value := range locals.GcpGcsBucket.Spec.Labels {
		locals.GcpLabels[key] = value
	}
	locals.GcpLabels[gcplabelkeys.Resource] = "true"
	locals.GcpLabels[gcplabelkeys.ResourceName] = locals.GcpGcsBucket.Spec.BucketName
	locals.GcpLabels[gcplabelkeys.ResourceKind] = strings.ToLower(cloudresourcekind.CloudResourceKind_GcpGcsBucket.String())

	if locals.GcpGcsBucket.Metadata.Org != "" {
		locals.GcpLabels[gcplabelkeys.Organization] = locals.GcpGcsBucket.Metadata.Org
	}
	if locals.GcpGcsBucket.Metadata.Env != "" {
		locals.GcpLabels[gcplabelkeys.Environment] = locals.GcpGcsBucket.Metadata.Env
	}
	if locals.GcpGcsBucket.Metadata.Id != "" {
		locals.GcpLabels[gcplabelkeys.ResourceId] = locals.GcpGcsBucket.Metadata.Id
	}

	locals.GcpProviderConfig = stackInput.ProviderConfig
	return locals
}
