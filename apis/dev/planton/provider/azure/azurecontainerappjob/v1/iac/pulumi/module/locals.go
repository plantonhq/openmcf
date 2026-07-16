package module

import (
	"strings"

	azurecontainerappjobv1 "github.com/plantonhq/planton/apis/dev/planton/provider/azure/azurecontainerappjob/v1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AzureContainerAppJob *azurecontainerappjobv1.AzureContainerAppJob
	ResourceGroupName    string
	AzureTags            map[string]string
}

// probeTransportStrings maps the probe-transport enum to Azure's
// case-sensitive wire values.
var probeTransportStrings = map[azurecontainerappjobv1.AzureContainerAppJobProbeTransport]string{
	azurecontainerappjobv1.AzureContainerAppJobProbeTransport_TCP_SOCKET: "TCP",
	azurecontainerappjobv1.AzureContainerAppJobProbeTransport_HTTP_GET:   "HTTP",
	azurecontainerappjobv1.AzureContainerAppJobProbeTransport_HTTPS_GET:  "HTTPS",
}

// volumeStorageTypeStrings maps the volume storage-type enum to ARM's
// values. Unspecified deploys EmptyDir (resolved in volume construction).
var volumeStorageTypeStrings = map[azurecontainerappjobv1.AzureContainerAppJobVolumeStorageType]string{
	azurecontainerappjobv1.AzureContainerAppJobVolumeStorageType_EMPTY_DIR:      "EmptyDir",
	azurecontainerappjobv1.AzureContainerAppJobVolumeStorageType_AZURE_FILE:     "AzureFile",
	azurecontainerappjobv1.AzureContainerAppJobVolumeStorageType_NFS_AZURE_FILE: "NfsAzureFile",
	azurecontainerappjobv1.AzureContainerAppJobVolumeStorageType_SECRET:         "Secret",
}

// identityTypeStrings maps the identity-type enum to ARM's values.
var identityTypeStrings = map[azurecontainerappjobv1.AzureContainerAppJobIdentityType]string{
	azurecontainerappjobv1.AzureContainerAppJobIdentityType_SYSTEM_ASSIGNED:          "SystemAssigned",
	azurecontainerappjobv1.AzureContainerAppJobIdentityType_USER_ASSIGNED:            "UserAssigned",
	azurecontainerappjobv1.AzureContainerAppJobIdentityType_SYSTEM_AND_USER_ASSIGNED: "SystemAssigned, UserAssigned",
}

func initializeLocals(ctx *pulumi.Context, stackInput *azurecontainerappjobv1.AzureContainerAppJobStackInput) *Locals {
	locals := &Locals{}

	locals.AzureContainerAppJob = stackInput.Target
	target := stackInput.Target

	locals.ResourceGroupName = target.Spec.ResourceGroup.GetValue()

	// Metadata-derived tags first, then the user's spec tags merged over
	// them: user tags deliberately win so an org's governance conventions
	// (cost center, owner) can override the derived values where they
	// collide.
	locals.AzureTags = map[string]string{
		// PARITY-EXCEPTION: resource_kind here is the lowered
		// CloudResourceKind enum string and resource_id is omitted when
		// metadata.id is empty, while the Terraform module emits the
		// family-wide snake-case literal and falls back to metadata.name.
		// Output-neutral (tags never feed stack outputs); aligning the two
		// shapes is a family-wide convention change, not a per-kind fix.
		"resource":      "true",
		"resource_name": target.Metadata.Name,
		"resource_kind": strings.ToLower(cloudresourcekind.CloudResourceKind_AzureContainerAppJob.String()),
	}

	if target.Metadata.Id != "" {
		locals.AzureTags["resource_id"] = target.Metadata.Id
	}

	if target.Metadata.Org != "" {
		locals.AzureTags["organization"] = target.Metadata.Org
	}

	if target.Metadata.Env != "" {
		locals.AzureTags["environment"] = target.Metadata.Env
	}

	for key, value := range target.Spec.Tags {
		locals.AzureTags[key] = value
	}

	return locals
}
