# AzureEventHub

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1`

**AzureEventHubSpec** defines the configuration for creating an event hub
inside an Azure Event Hubs namespace: one partitioned, replayable event
stream.

An event hub is the core streaming entity. Producers append events to
partitions; consumers read them through consumer groups
(AzureEventHubConsumerGroup), each keeping its own offset -- so the same
stream feeds real-time processing, batch analytics, and archival
independently. Event hubs are many-per-namespace with independent
lifecycles, which is why they are a first-class kind referencing the
namespace rather than a list folded into the namespace's spec. The hub's
ARM ID is also a data-plane RBAC scope (grant "Azure Event Hubs Data
Receiver/Sender" on exactly this hub) and the target for hub-scoped SAS
credentials (AzureEventHubAuthorizationRule with event_hub_id).

**Capture** archives every event to Azure Blob Storage in Avro format on
a size-or-interval cadence -- the built-in bridge from streaming to batch
(cold storage, replay beyond retention, audit). It folds into this spec
because Azure models it as a property of the hub with no independent
lifecycle.

**Contracts enforced by Azure at apply time** (they depend on the parent
namespace's tier, which a reference cannot see at validation time):
- partition_count is capped at 32 on shared BASIC/STANDARD namespaces;
  up to 1024 on PREMIUM or a dedicated cluster.
- partition_count can only be INCREASED, and only on PREMIUM or a
  dedicated cluster; shared-namespace hubs keep their creation-time
  count for life.
- message_retention is capped at 1 day on BASIC, 7 on STANDARD, 90 on
  PREMIUM/dedicated.
- The Compact cleanup policy requires STANDARD or higher.

**ForceNew fields** (changing these replaces the hub and its retained
events): `event_hub_name`, `retention_description.cleanup_policy`.

## Example

```yaml
# Offline-plan manifest: a hub exercising every deep seam at once -- the
# rich retention model (Kafka-style compaction takes its own manifest;
# this one renders the DELETE window), the SendDisabled gate, and the
# full capture block with user-assigned-identity storage auth.
apiVersion: azure.planton.dev/v1
kind: AzureEventHub
metadata:
  name: test-event-hub
spec:
  namespaceId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.EventHub/namespaces/hack-eventhubs-ns
  eventHubName: telemetry
  partitionCount: 8
  retentionDescription:
    cleanupPolicy: DELETE
    retentionTimeInHours: 168
  status: SEND_DISABLED
  captureDescription:
    enabled: true
    encoding: AVRO_DEFLATE
    intervalInSeconds: 600
    sizeLimitInBytes: 104857600
    skipEmptyArchives: true
    destination:
      archiveNameFormat: "{Namespace}/{EventHub}/{PartitionId}/{Year}/{Month}/{Day}/{Hour}/{Minute}/{Second}"
      blobContainerName:
        value: capture-archives
      storageAccountId:
        value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Storage/storageAccounts/hackcapturesa
      storageAuthenticationType: USER_ASSIGNED
      storageAuthenticationId:
        value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.ManagedIdentity/userAssignedIdentities/hack-uai
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.namespaceId` | `string \| valueFrom` | yes |  | AzureEventHubNamespace (`status.outputs.namespace_id`) |
| `spec.eventHubName` | `string` | yes |  |  |
| `spec.partitionCount` | `int32` | yes |  |  |
| `spec.messageRetention` | `int32` |  |  |  |
| `spec.retentionDescription` | `AzureEventHubRetentionDescription` |  |  |  |
| `spec.retentionDescription.cleanupPolicy` | `enum` |  |  |  |
| `spec.retentionDescription.retentionTimeInHours` | `int32` |  |  |  |
| `spec.retentionDescription.tombstoneRetentionTimeInHours` | `int32` |  |  |  |
| `spec.status` | `enum` |  |  |  |
| `spec.captureDescription` | `AzureEventHubCaptureDescription` |  |  |  |
| `spec.captureDescription.enabled` | `bool` | yes |  |  |
| `spec.captureDescription.encoding` | `enum` |  |  |  |
| `spec.captureDescription.intervalInSeconds` | `int32` |  |  |  |
| `spec.captureDescription.sizeLimitInBytes` | `int32` |  |  |  |
| `spec.captureDescription.skipEmptyArchives` | `bool` |  |  |  |
| `spec.captureDescription.destination` | `AzureEventHubCaptureDestination` | yes |  |  |
| `spec.captureDescription.destination.archiveNameFormat` | `string` | yes |  |  |
| `spec.captureDescription.destination.blobContainerName` | `string \| valueFrom` | yes |  | AzureStorageContainer (`status.outputs.container_name`) |
| `spec.captureDescription.destination.storageAccountId` | `string \| valueFrom` | yes |  | AzureStorageAccount (`status.outputs.storage_account_id`) |
| `spec.captureDescription.destination.storageAuthenticationType` | `enum` |  |  |  |
| `spec.captureDescription.destination.storageAuthenticationId` | `string \| valueFrom` |  |  | AzureUserAssignedIdentity (`status.outputs.identity_id`) |

## Field Details

### spec.namespaceId

`string | valueFrom` · required

The Event Hubs namespace the hub lives in, by ARM ID. References an
AzureEventHubNamespace's namespace_id output so the namespace and its
hubs compose in one manifest set. Fixed at creation.

- references: AzureEventHubNamespace (`status.outputs.namespace_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureEventHubNamespace, name: <that resource's name>, fieldPath: status.outputs.namespace_id}} -- a bare string does not parse

### spec.eventHubName

`string` · required

The hub's name -- unique within the namespace, 1-256 characters.
Starts and ends with a letter or number; letters, numbers, periods,
hyphens, and underscores in between.

**ForceNew**: changing the name replaces the hub and its retained
events.

- rule: event_hub_name must start and end with a letter or number and may contain letters, numbers, periods, hyphens, and underscores (max 256 characters)
- rule: {"required":true,"string":{"minLen":"1","maxLen":"256"}}

### spec.partitionCount

`int32` · required

The number of partitions -- the hub's unit of parallelism and
ordering. Each partition is an independently consumed, ordered event
sequence; downstream parallelism cannot exceed the partition count.

1-32 on shared BASIC/STANDARD namespaces; up to 1024 on PREMIUM or a
dedicated cluster (Azure enforces the tier cap at apply time).
The count can never be DECREASED, and can only be increased on
PREMIUM/dedicated -- on shared namespaces, size for peak up front.

Guidelines: 2-4 for low throughput; 8-16 for typical production;
32+ for high-throughput ingestion with many parallel consumers.

- rule: {"required":true,"int32":{"lte":1024,"gte":1}}

### spec.messageRetention

`int32` · optional (explicit presence)

How long events stay replayable, in days. Caps by tier: 1 (BASIC),
7 (STANDARD), 90 (PREMIUM/dedicated) -- Azure enforces the cap at
apply time. Unset lets Azure default it (1 day). Exactly one of
message_retention and retention_description must be set: this simple
day count, or the richer hours/compaction block below.

- rule: {"int32":{"lte":90,"gte":1}}

### spec.retentionDescription

`AzureEventHubRetentionDescription`

The richer retention model: hour-granular windows and Kafka-style
log compaction. Exactly one of message_retention and
retention_description must be set.

- rule: DELETE takes retention_time_in_hours (and no tombstone window); COMPACT takes tombstone_retention_time_in_hours (and no retention window) -- Azure ignores the mismatched field silently, so the spec rejects it up front

### spec.retentionDescription.cleanupPolicy

`enum`

How the hub reclaims space: DELETE removes events past the retention
window (the classic streaming model); COMPACT keeps the LATEST event
per partition key forever, removing older values -- Kafka-style
compacted topics for changelog/table semantics. COMPACT requires
STANDARD tier or higher (Azure enforces at apply time).

**ForceNew**: the cleanup policy is fixed at creation -- changing it
replaces the hub.

- rule: {"enum":{"definedOnly":true,"notIn":[0]}}

Allowed values (use exactly as shown):

- `azure_event_hub_cleanup_policy_unspecified` -- Not specified -- invalid; choose DELETE or COMPACT explicitly.
- `DELETE` -- Remove events older than the retention window -- the classic streaming model.
- `COMPACT` -- Keep the latest event per partition key forever (Kafka-style compaction) -- changelog/table semantics. STANDARD tier or higher.

### spec.retentionDescription.retentionTimeInHours

`int32` · optional (explicit presence)

The retention window in hours, for the DELETE policy (e.g. 168 =
7 days). Tier caps apply as with message_retention (24h BASIC, 168h
STANDARD, 2160h PREMIUM/dedicated; Azure enforces at apply time).
Required with DELETE; must be absent with COMPACT (where events
never expire by time).

- rule: {"int32":{"gte":1}}

### spec.retentionDescription.tombstoneRetentionTimeInHours

`int32` · optional (explicit presence)

How long a tombstone (a null-value event marking a key as deleted)
stays readable before compaction removes it, in hours -- consumers
must replay within this window to observe deletions. Required with
COMPACT; must be absent with DELETE.

- rule: {"int32":{"gte":1}}

### spec.status

`enum`

The hub's gate state: ACTIVE (normal), DISABLED (sends and receives
rejected; events retained), or SEND_DISABLED (receive-only drain
mode). Unspecified deploys ACTIVE. The transitional server states
(Creating, Deleting, Renaming, ...) are not knobs and are not
modeled.

- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `azure_event_hub_entity_status_unspecified` -- Not specified -- deploys ACTIVE.
- `ACTIVE` -- Normal operation: sends and receives flow.
- `DISABLED` -- Sends and receives are rejected; retained events stay stored.
- `SEND_DISABLED` -- Receive-only drain mode: new sends are rejected, consumers keep reading.

### spec.captureDescription

`AzureEventHubCaptureDescription`

Capture: archive every event to Azure Blob Storage in Avro format on
a size-or-interval cadence. The built-in streaming-to-batch bridge --
cold storage, replay beyond the retention window, audit trails --
with no consumer application to run.

### spec.captureDescription.enabled

`bool` · required · optional (explicit presence)

Whether capture is running. Azure requires the flag stated explicitly
when the block is declared (presence-tracked so an explicit false is
legal); keeping the block with enabled=false preserves the
configuration while pausing archival.

- rule: {"required":true}

### spec.captureDescription.encoding

`enum`

The archive encoding. AVRO is the standard, queryable-by-everything
choice; AVRO_DEFLATE trades CPU for smaller archives.

- rule: {"enum":{"definedOnly":true,"notIn":[0]}}

Allowed values (use exactly as shown):

- `azure_event_hub_capture_encoding_unspecified` -- Not specified -- invalid; choose an encoding explicitly.
- `AVRO` -- Apache Avro -- the standard capture format, readable by every analytics engine.
- `AVRO_DEFLATE` -- Deflate-compressed Avro -- smaller archives at some CPU cost.

### spec.captureDescription.intervalInSeconds

`int32` · optional (explicit presence)

How often an archive window closes, in seconds (60-900). A window
closes when EITHER this interval elapses OR size_limit_in_bytes
accumulates, whichever comes first. Azure's default: 300.

- rule: {"int32":{"lte":900,"gte":60}}

### spec.captureDescription.sizeLimitInBytes

`int32` · optional (explicit presence)

How much data closes an archive window, in bytes (10485760 =
10 MB to 524288000 = 500 MB). Azure's default: 314572800 (300 MB).

- rule: {"int32":{"lte":524288000,"gte":10485760}}

### spec.captureDescription.skipEmptyArchives

`bool` · optional (explicit presence)

Whether to skip writing archive files for windows with no events.
Azure's default (false) writes empty Avro files, which keeps a
continuous file cadence for downstream batch jobs; true saves
storage on sparse streams.

### spec.captureDescription.destination

`AzureEventHubCaptureDestination` · required

Where the archives land.

- rule: {"required":true}
- rule: storage_authentication_id is required with USER_ASSIGNED and must be absent otherwise -- Azure requires the identity for the user-assigned path and silently ignores it on the others

### spec.captureDescription.destination.archiveNameFormat

`string` · required

The naming pattern for archive blobs. Must contain ALL of the
placeholders {Namespace}, {EventHub}, {PartitionId}, {Year}, {Month},
{Day}, {Hour}, {Minute}, {Second} -- order and surrounding text are
free. Example:
"{Namespace}/{EventHub}/{PartitionId}/{Year}/{Month}/{Day}/{Hour}/{Minute}/{Second}"

- rule: archive_name_format must contain all of {Namespace}, {EventHub}, {PartitionId}, {Year}, {Month}, {Day}, {Hour}, {Minute}, {Second}
- rule: {"required":true}

### spec.captureDescription.destination.blobContainerName

`string | valueFrom` · required

The blob container archives are written into. References an
AzureStorageContainer's container_name output; the container must
live in the storage account below.

- references: AzureStorageContainer (`status.outputs.container_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureStorageContainer, name: <that resource's name>, fieldPath: status.outputs.container_name}} -- a bare string does not parse

### spec.captureDescription.destination.storageAccountId

`string | valueFrom` · required

The storage account holding the container, by ARM ID. References an
AzureStorageAccount's storage_account_id output.

- references: AzureStorageAccount (`status.outputs.storage_account_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureStorageAccount, name: <that resource's name>, fieldPath: status.outputs.storage_account_id}} -- a bare string does not parse

### spec.captureDescription.destination.storageAuthenticationType

`enum`

How capture authenticates to the storage account. STORAGE_SAS
(Azure's default) uses service-managed SAS; SYSTEM_ASSIGNED uses the
namespace's system-assigned identity; USER_ASSIGNED uses the
identity below. For the identity paths, grant the chosen identity
"Storage Blob Data Contributor" on the account (compose an
AzureRoleAssignment) and attach it via the namespace's identity
block.

- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `azure_event_hub_capture_storage_authentication_type_unspecified` -- Not specified -- Azure's default: service-managed SAS.
- `STORAGE_SAS` -- Service-managed shared-access signatures (Azure's default). No identity setup needed; the storage firewall must admit the service.
- `SYSTEM_ASSIGNED` -- The namespace's system-assigned managed identity -- keyless; grant it Storage Blob Data Contributor on the account.
- `USER_ASSIGNED` -- A user-assigned managed identity (storage_authentication_id) -- keyless and pre-grantable before the namespace exists.

### spec.captureDescription.destination.storageAuthenticationId

`string | valueFrom`

The user-assigned identity capture writes with, by ARM ID --
required with (and only meaningful with) USER_ASSIGNED. References
an AzureUserAssignedIdentity's identity_id output; the identity must
also be attached via the namespace's identity block.

- references: AzureUserAssignedIdentity (`status.outputs.identity_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureUserAssignedIdentity, name: <that resource's name>, fieldPath: status.outputs.identity_id}} -- a bare string does not parse

## Validation Rules

- `event_hub_retention_exactly_one_model`: set exactly one retention model: message_retention (simple day count) or retention_description (hour-granular windows and log compaction)

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureEventHub, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.event_hub_id` | `string` | The Azure Resource Manager ID of the event hub. Format: /subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.EventHub/namespaces/{ns}/eventhubs/{name} The parent reference for consumer groups and hub-scoped authorization rules, and the scope for hub-level data-plane RBAC. |
| `status.outputs.event_hub_name` | `string` | The hub's name -- what producers and consumers address within the namespace (the Kafka topic name on the Kafka endpoint). |
| `status.outputs.partition_ids` | `[]string` | The identifiers of the hub's partitions (e.g. "0", "1", ...), as Azure assigned them -- what partition-aware consumers enumerate. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.namespaceId` | AzureEventHubNamespace | `status.outputs.namespace_id` |
| `spec.captureDescription.destination.blobContainerName` | AzureStorageContainer | `status.outputs.container_name` |
| `spec.captureDescription.destination.storageAccountId` | AzureStorageAccount | `status.outputs.storage_account_id` |
| `spec.captureDescription.destination.storageAuthenticationId` | AzureUserAssignedIdentity | `status.outputs.identity_id` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AzureEventHubAuthorizationRule | `spec.eventHubId` | `status.outputs.event_hub_id` |
| AzureEventHubConsumerGroup | `spec.eventHubId` | `status.outputs.event_hub_id` |
| AzureMonitorActionGroup | `spec.eventHubReceivers[].eventHubName` | `status.outputs.event_hub_name` |
| AzureMonitorDiagnosticSetting | `spec.eventhubName` | `status.outputs.event_hub_name` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
