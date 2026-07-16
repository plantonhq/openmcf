package module

import (
	"strings"

	gcpprovider "github.com/plantonhq/planton/apis/dev/planton/provider/gcp"
	gcpbigquerytablev1 "github.com/plantonhq/planton/apis/dev/planton/provider/gcp/gcpbigquerytable/v1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/gcp/gcplabelkeys"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	GcpProviderConfig *gcpprovider.GcpProviderConfig
	GcpBigQueryTable  *gcpbigquerytablev1.GcpBigQueryTable
	GcpLabels         map[string]string
}

func initializeLocals(_ *pulumi.Context, stackInput *gcpbigquerytablev1.GcpBigQueryTableStackInput) *Locals {
	locals := &Locals{}
	locals.GcpBigQueryTable = stackInput.Target

	// User labels first so platform attribution labels win on key
	// conflicts — identical merge order to the Terraform module.
	locals.GcpLabels = map[string]string{}
	for key, value := range locals.GcpBigQueryTable.Spec.Labels {
		locals.GcpLabels[key] = value
	}
	locals.GcpLabels[gcplabelkeys.Resource] = "true"
	locals.GcpLabels[gcplabelkeys.ResourceName] = locals.GcpBigQueryTable.Spec.TableId
	locals.GcpLabels[gcplabelkeys.ResourceKind] = strings.ToLower(cloudresourcekind.CloudResourceKind_GcpBigQueryTable.String())

	if locals.GcpBigQueryTable.Metadata.Org != "" {
		locals.GcpLabels[gcplabelkeys.Organization] = locals.GcpBigQueryTable.Metadata.Org
	}
	if locals.GcpBigQueryTable.Metadata.Env != "" {
		locals.GcpLabels[gcplabelkeys.Environment] = locals.GcpBigQueryTable.Metadata.Env
	}
	if locals.GcpBigQueryTable.Metadata.Id != "" {
		locals.GcpLabels[gcplabelkeys.ResourceId] = locals.GcpBigQueryTable.Metadata.Id
	}

	locals.GcpProviderConfig = stackInput.ProviderConfig
	return locals
}
