package module

import (
	"strconv"
	"strings"

	gcpprovider "github.com/plantonhq/planton/catalog/gcp"
	gcpaddressv1alpha1 "github.com/plantonhq/planton/catalog/gcp/gcpaddress/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/gcp/gcplabelkeys"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	GcpProviderConfig *gcpprovider.GcpProviderConfig
	GcpAddress        *gcpaddressv1alpha1.GcpAddress
	GcpLabels         map[string]string
}

func initializeLocals(_ *pulumi.Context, stackInput *gcpaddressv1alpha1.GcpAddressStackInput) *Locals {
	locals := &Locals{}
	locals.GcpAddress = stackInput.Target
	// User labels first so platform attribution labels win on key
	// conflicts — identical merge order to the Terraform module.
	locals.GcpLabels = map[string]string{}
	for key, value := range locals.GcpAddress.Spec.Labels {
		locals.GcpLabels[key] = value
	}
	locals.GcpLabels[gcplabelkeys.Resource] = strconv.FormatBool(true)
	locals.GcpLabels[gcplabelkeys.ResourceName] = locals.GcpAddress.Spec.AddressName
	locals.GcpLabels[gcplabelkeys.ResourceKind] = strings.ToLower(cloudresourcekind.CloudResourceKind_GcpAddress.String())

	if locals.GcpAddress.Metadata.Org != "" {
		locals.GcpLabels[gcplabelkeys.Organization] = locals.GcpAddress.Metadata.Org
	}
	if locals.GcpAddress.Metadata.Env != "" {
		locals.GcpLabels[gcplabelkeys.Environment] = locals.GcpAddress.Metadata.Env
	}
	if locals.GcpAddress.Metadata.Id != "" {
		locals.GcpLabels[gcplabelkeys.ResourceId] = locals.GcpAddress.Metadata.Id
	}

	locals.GcpProviderConfig = stackInput.ProviderConfig
	return locals
}
