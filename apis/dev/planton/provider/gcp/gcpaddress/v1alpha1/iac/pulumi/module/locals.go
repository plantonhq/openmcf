package module

import (
	"strconv"
	"strings"

	gcpprovider "github.com/plantonhq/planton/apis/dev/planton/provider/gcp"
	gcpaddressv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/gcp/gcpaddress/v1alpha1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/gcp/gcplabelkeys"
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
	locals.GcpLabels = map[string]string{
		gcplabelkeys.Resource:     strconv.FormatBool(true),
		gcplabelkeys.ResourceName: locals.GcpAddress.Spec.AddressName,
		gcplabelkeys.ResourceKind: strings.ToLower(cloudresourcekind.CloudResourceKind_GcpAddress.String()),
	}

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
