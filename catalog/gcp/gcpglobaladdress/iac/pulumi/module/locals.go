package module

import (
	"strconv"
	"strings"

	gcpprovider "github.com/plantonhq/planton/catalog/gcp"
	gcpglobaladdressv1alpha1 "github.com/plantonhq/planton/catalog/gcp/gcpglobaladdress/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/gcp/gcplabelkeys"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	GcpProviderConfig *gcpprovider.GcpProviderConfig
	GcpGlobalAddress  *gcpglobaladdressv1alpha1.GcpGlobalAddress
	GcpLabels         map[string]string
}

func initializeLocals(_ *pulumi.Context, stackInput *gcpglobaladdressv1alpha1.GcpGlobalAddressStackInput) *Locals {
	locals := &Locals{}
	locals.GcpGlobalAddress = stackInput.Target
	// User labels first so platform attribution labels win on key
	// conflicts — identical merge order to the Terraform module.
	locals.GcpLabels = map[string]string{}
	for key, value := range locals.GcpGlobalAddress.Spec.Labels {
		locals.GcpLabels[key] = value
	}
	locals.GcpLabels[gcplabelkeys.Resource] = strconv.FormatBool(true)
	locals.GcpLabels[gcplabelkeys.ResourceName] = locals.GcpGlobalAddress.Spec.AddressName
	locals.GcpLabels[gcplabelkeys.ResourceKind] = strings.ToLower(cloudresourcekind.CloudResourceKind_GcpGlobalAddress.String())

	if locals.GcpGlobalAddress.Metadata.Org != "" {
		locals.GcpLabels[gcplabelkeys.Organization] = locals.GcpGlobalAddress.Metadata.Org
	}
	if locals.GcpGlobalAddress.Metadata.Env != "" {
		locals.GcpLabels[gcplabelkeys.Environment] = locals.GcpGlobalAddress.Metadata.Env
	}
	if locals.GcpGlobalAddress.Metadata.Id != "" {
		locals.GcpLabels[gcplabelkeys.ResourceId] = locals.GcpGlobalAddress.Metadata.Id
	}

	locals.GcpProviderConfig = stackInput.ProviderConfig
	return locals
}
