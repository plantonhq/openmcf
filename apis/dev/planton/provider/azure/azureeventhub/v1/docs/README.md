# AzureEventHub -- Design Research

## The Resource

An event hub (`Microsoft.EventHub/namespaces/eventhubs`) is one
partitioned, replayable event stream inside an Event Hubs namespace. The
component maps onto `azurerm_eventhub` (azurerm v4.x,
`internal/services/eventhub/eventhub_resource.go`), parity-verified
against pulumi-azure v6 (`eventhub.EventHub`).

## Field Mapping (azurerm → spec)

| azurerm | spec | Notes |
|---|---|---|
| `namespace_id` | same | Required, ForceNew; FK → `AzureEventHubNamespace.namespace_id` |
| `name` | `event_hub_name` | Required, ForceNew (replaces the hub and its retained events); 1-256 chars, start/end alphanumeric; also the Kafka topic name |
| `partition_count` | same | 1-1024 in the spec; the tier caps (32 shared, 1024 PREMIUM/dedicated) and the never-decrease / increase-only-on-PREMIUM-or-dedicated contracts are apply-time (see below) |
| `message_retention` | same | Simple days, 1-90; tier caps (1/7/90) apply-time; XOR with retention_description (message CEL) |
| `retention_description` | `retention_description` block | cleanup_policy DELETE/COMPACT (ForceNew), retention_time_in_hours, tombstone_retention_time_in_hours; field-to-policy pairing is a message CEL (see below) |
| `status` | enum | Active (default) / Disabled / SendDisabled; the transitional server states (Creating, Deleting, Renaming, ...) are not knobs and are not modeled |
| `capture_description` | `capture_description` block | enabled (required-explicit), encoding AVRO/AVRO_DEFLATE, interval 60-900s, size 10-500 MB, skip_empty_archives, destination |
| `capture_description.destination.name` | -- | Recorded constant below |
| `capture_description.destination.*` | same | archive_name_format (all-nine-tokens CEL), blob_container_name FK → `AzureStorageContainer.container_name`, storage_account_id FK → `AzureStorageAccount.storage_account_id`, storage_authentication_type + storage_authentication_id (spec-paired) |

Outputs: `event_hub_id` (the hub-level RBAC scope and the parent seam
for consumer groups and hub-scoped SAS rules), `event_hub_name` (the
Kafka topic name), `partition_ids`.

## Decomposition Decisions

- **The hub is a first-class kind, not a namespace bundle.** Hubs are
  many-per-namespace with independent lifecycles, owned by stream teams,
  and individually FK-referenced (consumer groups, hub-scoped SAS rules,
  hub-level role assignments).
- **Capture folds into the hub.** Azure models it as a property of the
  hub with no independent lifecycle, and nothing references it -- the
  inline block is the honest declarative shape.
- **Consumer groups are a first-class kind**
  (AzureEventHubConsumerGroup) -- many per hub, one per consuming
  application, referencing `event_hub_id`.

## The Retention XOR Contract

Exactly one of `message_retention` (simple day count) and
`retention_description` (hour-granular windows and compaction) must be
set -- a message CEL. Within `retention_description`, the hour field is
paired to the policy: DELETE takes `retention_time_in_hours` (and no
tombstone window); COMPACT takes `tombstone_retention_time_in_hours`
(and no retention window). The provider sends whichever fields are
present and **Azure silently ignores the mismatched one** -- a config
that looks meaningful but does nothing. The spec rejects the mismatch up
front instead.

COMPACT is Kafka-style log compaction: the hub keeps the LATEST event
per partition key forever; tombstones (null-value events marking a key
deleted) stay readable for the tombstone window so consumers observe
deletions. COMPACT requires STANDARD or higher, and the cleanup policy
is ForceNew -- fixed at creation.

## The Capture Identity Contract

`storage_authentication_type` picks how capture writes to the storage
account: STORAGE_SAS (Azure's default -- service-managed SAS, no
identity setup; the storage firewall must admit the service),
SYSTEM_ASSIGNED, or USER_ASSIGNED (+ `storage_authentication_id`,
required with USER_ASSIGNED and rejected otherwise -- a message CEL).
For the identity paths, the chosen identity must ride the **namespace's**
identity block (the hub has no identity of its own) and hold Storage
Blob Data Contributor on the account (compose an AzureRoleAssignment).

## Recorded Skips (with reasons)

- **`capture_description.destination.name`** -- a one-value constant,
  not a knob: Azure accepts exactly `EventHubArchive.AzureBlockBlob`
  (the Data Lake variant retired with Gen1). Both modules send it
  unconditionally.
- **Tags** -- ARM does not support tags on Event Hubs entities
  (hubs/consumer groups/rules); the platform's identity tags live on the
  parent namespace.

## Apply-Time Contracts

These depend on the parent namespace's tier, which a reference cannot
see at validation time -- Azure enforces them at apply:

- `partition_count` is capped at 32 on shared BASIC/STANDARD
  namespaces; up to 1024 on PREMIUM or a dedicated cluster.
- `partition_count` can only be INCREASED, and only on PREMIUM or a
  dedicated cluster; shared-namespace hubs keep their creation-time
  count for life.
- Retention is capped by tier: 1 day (BASIC), 7 (STANDARD), 90
  (PREMIUM/dedicated) -- 24/168/2160 in hours.
- The COMPACT cleanup policy requires STANDARD or higher.

## Operational Behavior Worth Knowing

- **Downstream parallelism cannot exceed the partition count** -- each
  partition is an independently consumed, ordered sequence. Guidelines:
  2-4 for low throughput, 8-16 for typical production, 32+ for
  high-throughput ingestion.
- **BASIC namespaces allow a single consumer group per hub** (the
  built-in `$Default`); STANDARD allows 20 -- plan the consuming
  applications against the namespace tier.
- **A capture window closes when EITHER the interval elapses OR the
  size limit accumulates**, whichever comes first. Azure's defaults:
  300 seconds / 300 MB. `skip_empty_archives` false (Azure's default)
  writes empty Avro files, keeping a continuous file cadence for
  downstream batch jobs.
- **SEND_DISABLED is drain mode**: new sends are rejected while
  consumers keep reading -- useful for decommissioning a stream without
  losing in-flight events.

## Composition

- `namespace_id` → `AzureEventHubNamespace.status.outputs.namespace_id`
- `capture_description.destination.blob_container_name` → `AzureStorageContainer.status.outputs.container_name`
- `capture_description.destination.storage_account_id` → `AzureStorageAccount.status.outputs.storage_account_id`
- `capture_description.destination.storage_authentication_id` → `AzureUserAssignedIdentity.status.outputs.identity_id`
- `event_hub_id` output ← AzureEventHubConsumerGroup,
  AzureEventHubAuthorizationRule (hub scope), hub-level data-plane role
  assignments
