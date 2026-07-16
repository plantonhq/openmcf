package module

import (
	azureeventhubv1 "github.com/plantonhq/planton/apis/dev/planton/provider/azure/azureeventhub/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AzureEventHub *azureeventhubv1.AzureEventHub
}

// Event hubs carry no Azure tags: ARM does not support tags on Event Hubs
// entities (hubs/consumer groups/rules), so the platform's identity tags
// live on the parent namespace.

// statusStrings maps the gate-state enum to ARM's wire values. The
// unspecified row deploys Active -- an unmapped enum would send the empty
// string, which the provider rejects.
var statusStrings = map[azureeventhubv1.AzureEventHubEntityStatus]string{
	azureeventhubv1.AzureEventHubEntityStatus_azure_event_hub_entity_status_unspecified: "Active",
	azureeventhubv1.AzureEventHubEntityStatus_ACTIVE:                                    "Active",
	azureeventhubv1.AzureEventHubEntityStatus_DISABLED:                                  "Disabled",
	azureeventhubv1.AzureEventHubEntityStatus_SEND_DISABLED:                             "SendDisabled",
}

// cleanupPolicyStrings maps the retention cleanup policy to ARM's values --
// the spec enum requires an explicit choice when retention_description is
// declared, so no unspecified fallback row is needed.
var cleanupPolicyStrings = map[azureeventhubv1.AzureEventHubCleanupPolicy]string{
	azureeventhubv1.AzureEventHubCleanupPolicy_DELETE:  "Delete",
	azureeventhubv1.AzureEventHubCleanupPolicy_COMPACT: "Compact",
}

// captureEncodingStrings maps the capture encoding to ARM's values --
// required/explicit in the spec, no fallback row needed.
var captureEncodingStrings = map[azureeventhubv1.AzureEventHubCaptureEncoding]string{
	azureeventhubv1.AzureEventHubCaptureEncoding_AVRO:         "Avro",
	azureeventhubv1.AzureEventHubCaptureEncoding_AVRO_DEFLATE: "AvroDeflate",
}

// captureAuthStrings maps the capture storage-auth mode. The unspecified
// row keeps Azure's default (service-managed SAS) -- "StorageSAS" is the
// provider's marker for "send no identity", not an ARM value.
var captureAuthStrings = map[azureeventhubv1.AzureEventHubCaptureStorageAuthenticationType]string{
	azureeventhubv1.AzureEventHubCaptureStorageAuthenticationType_azure_event_hub_capture_storage_authentication_type_unspecified: "StorageSAS",
	azureeventhubv1.AzureEventHubCaptureStorageAuthenticationType_STORAGE_SAS:                                                     "StorageSAS",
	azureeventhubv1.AzureEventHubCaptureStorageAuthenticationType_SYSTEM_ASSIGNED:                                                 "SystemAssigned",
	azureeventhubv1.AzureEventHubCaptureStorageAuthenticationType_USER_ASSIGNED:                                                   "UserAssigned",
}

func initializeLocals(ctx *pulumi.Context, stackInput *azureeventhubv1.AzureEventHubStackInput) *Locals {
	locals := &Locals{}
	locals.AzureEventHub = stackInput.Target
	return locals
}
