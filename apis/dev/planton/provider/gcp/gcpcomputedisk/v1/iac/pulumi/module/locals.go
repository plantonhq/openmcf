package module

import (
	"strings"

	gcpprovider "github.com/plantonhq/planton/apis/dev/planton/provider/gcp"
	gcpcomputediskv1 "github.com/plantonhq/planton/apis/dev/planton/provider/gcp/gcpcomputedisk/v1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/gcp/gcplabelkeys"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	GcpProviderConfig *gcpprovider.GcpProviderConfig
	GcpComputeDisk    *gcpcomputediskv1.GcpComputeDisk
	GcpLabels         map[string]string
	// DiskName is the cloud-side name: spec.disk_name when set,
	// metadata.name otherwise — the same explicit conditional as the
	// Terraform module, so both engines derive the identical name.
	DiskName string
}

func initializeLocals(_ *pulumi.Context, stackInput *gcpcomputediskv1.GcpComputeDiskStackInput) *Locals {
	locals := &Locals{}
	locals.GcpComputeDisk = stackInput.Target
	locals.GcpProviderConfig = stackInput.ProviderConfig

	locals.DiskName = stackInput.Target.Spec.DiskName
	if locals.DiskName == "" {
		locals.DiskName = stackInput.Target.Metadata.Name
	}

	// User labels first so platform attribution labels win on key
	// conflicts — identical merge order to the Terraform module.
	locals.GcpLabels = map[string]string{}
	for key, value := range stackInput.Target.Spec.Labels {
		locals.GcpLabels[key] = value
	}
	locals.GcpLabels[gcplabelkeys.Resource] = "true"
	locals.GcpLabels[gcplabelkeys.ResourceName] = locals.DiskName
	locals.GcpLabels[gcplabelkeys.ResourceKind] = strings.ToLower(cloudresourcekind.CloudResourceKind_GcpComputeDisk.String())

	if locals.GcpComputeDisk.Metadata.Org != "" {
		locals.GcpLabels[gcplabelkeys.Organization] = locals.GcpComputeDisk.Metadata.Org
	}
	if locals.GcpComputeDisk.Metadata.Env != "" {
		locals.GcpLabels[gcplabelkeys.Environment] = locals.GcpComputeDisk.Metadata.Env
	}
	if locals.GcpComputeDisk.Metadata.Id != "" {
		locals.GcpLabels[gcplabelkeys.ResourceId] = locals.GcpComputeDisk.Metadata.Id
	}

	return locals
}
