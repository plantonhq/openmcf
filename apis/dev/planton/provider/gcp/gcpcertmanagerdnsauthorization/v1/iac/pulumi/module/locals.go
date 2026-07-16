package module

import (
	"strings"

	gcpprovider "github.com/plantonhq/planton/apis/dev/planton/provider/gcp"
	gcpcertmanagerdnsauthorizationv1 "github.com/plantonhq/planton/apis/dev/planton/provider/gcp/gcpcertmanagerdnsauthorization/v1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/gcp/gcplabelkeys"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	GcpProviderConfig              *gcpprovider.GcpProviderConfig
	GcpCertManagerDnsAuthorization *gcpcertmanagerdnsauthorizationv1.GcpCertManagerDnsAuthorization
	GcpLabels                      map[string]string

	// ProjectId is empty when the manifest omits it — the provider's default
	// project then applies (the same ambient contract the Terraform module
	// honors by passing null).
	ProjectId string

	// AuthorizationName falls back to metadata.name — explicit conditional,
	// so both engines derive the identical cloud-side name.
	AuthorizationName string
}

func initializeLocals(_ *pulumi.Context, stackInput *gcpcertmanagerdnsauthorizationv1.GcpCertManagerDnsAuthorizationStackInput) *Locals {
	locals := &Locals{}
	locals.GcpCertManagerDnsAuthorization = stackInput.Target
	locals.GcpProviderConfig = stackInput.ProviderConfig

	locals.ProjectId = stackInput.Target.Spec.ProjectId.GetValue()

	locals.AuthorizationName = stackInput.Target.Spec.AuthorizationName
	if locals.AuthorizationName == "" {
		locals.AuthorizationName = stackInput.Target.Metadata.Name
	}

	// User labels first so platform attribution labels win on key
	// conflicts — identical merge order to the Terraform module.
	locals.GcpLabels = map[string]string{}
	for key, value := range stackInput.Target.Spec.Labels {
		locals.GcpLabels[key] = value
	}
	locals.GcpLabels[gcplabelkeys.Resource] = "true"
	locals.GcpLabels[gcplabelkeys.ResourceName] = locals.AuthorizationName
	locals.GcpLabels[gcplabelkeys.ResourceKind] = strings.ToLower(cloudresourcekind.CloudResourceKind_GcpCertManagerDnsAuthorization.String())

	if stackInput.Target.Metadata.Org != "" {
		locals.GcpLabels[gcplabelkeys.Organization] = stackInput.Target.Metadata.Org
	}
	if stackInput.Target.Metadata.Env != "" {
		locals.GcpLabels[gcplabelkeys.Environment] = stackInput.Target.Metadata.Env
	}
	if stackInput.Target.Metadata.Id != "" {
		locals.GcpLabels[gcplabelkeys.ResourceId] = stackInput.Target.Metadata.Id
	}

	return locals
}
