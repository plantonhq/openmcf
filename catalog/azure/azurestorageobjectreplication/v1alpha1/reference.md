# AzureStorageObjectReplication

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**AzureStorageObjectReplicationSpec** defines the configuration for an
object replication policy between TWO Azure Storage Accounts:
asynchronous, rule-driven copying of block blobs from containers on a
source account to containers on a destination account. This is the
storage-level answer to cross-region DR (replicate to an account in a
paired region), data distribution (fan out to a read-local copy), and
tenant offboarding/archival -- without any application-side copy jobs.

One policy spans exactly one account pair; Azure materializes it on
BOTH accounts (the destination holds the authoritative copy that
assigns rule IDs, the source holds the mirror), which the modules
handle as one unit -- this kind IS the pair. Both accounts, and both
sides of every rule, are fixed at creation.

**Both accounts must be prepared for replication**: blob versioning
AND change feed enabled (the account spec's blob_properties) on the
source, and blob versioning on the destination -- Azure rejects the
policy at apply time otherwise. Versioning rules out
hierarchical-namespace accounts (ADLS Gen2 does not support it).
Replication is asynchronous with no RPO guarantee unless the account
opts into metrics and monitors them.

## Example

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureStorageObjectReplication
metadata:
  name: test-storage-object-replication
spec:
  sourceStorageAccountId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Storage/storageAccounts/plantonhackorsrc
  destinationStorageAccountId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Storage/storageAccounts/plantonhackordst
  rules:
    # Exercises the backfill-everything special value and prefix
    # filters.
    - sourceContainerName:
        value: invoices
      destinationContainerName:
        value: invoices-replica
      copyBlobsCreatedAfter: Everything
      prefixMatch:
        - invoices/2026
        - receipts/
    # Exercises the RFC 3339 instant form alongside the provider-default
    # rule shape.
    - sourceContainerName:
        value: exports
      destinationContainerName:
        value: exports-replica
      copyBlobsCreatedAfter: "2026-01-01T00:00:00Z"
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.sourceStorageAccountId` | `string \| valueFrom` | yes |  | AzureStorageAccount (`status.outputs.storage_account_id`) |
| `spec.destinationStorageAccountId` | `string \| valueFrom` | yes |  | AzureStorageAccount (`status.outputs.storage_account_id`) |
| `spec.rules` | `[]AzureStorageObjectReplicationRule` | yes |  |  |
| `spec.rules[].sourceContainerName` | `string \| valueFrom` | yes |  | AzureStorageContainer (`status.outputs.container_name`) |
| `spec.rules[].destinationContainerName` | `string \| valueFrom` | yes |  | AzureStorageContainer (`status.outputs.container_name`) |
| `spec.rules[].copyBlobsCreatedAfter` | `string` |  | `OnlyNewObjects` |  |
| `spec.rules[].prefixMatch` | `[]string` |  |  |  |

## Field Details

### spec.sourceStorageAccountId

`string | valueFrom` · required

The account blobs are copied FROM, by ARM ID. References an
AzureStorageAccount's storage_account_id output. Needs blob
versioning AND change feed enabled (blob_properties on the account
spec). Fixed at creation.

- references: AzureStorageAccount (`status.outputs.storage_account_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureStorageAccount, name: <that resource's name>, fieldPath: status.outputs.storage_account_id}} -- a bare string does not parse

### spec.destinationStorageAccountId

`string | valueFrom` · required

The account blobs are copied TO, by ARM ID. References an
AzureStorageAccount's storage_account_id output. Needs blob
versioning enabled. May live in any region -- a paired-region
destination is the DR pattern; same-region is the
distribution/archival pattern. Fixed at creation.

- references: AzureStorageAccount (`status.outputs.storage_account_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureStorageAccount, name: <that resource's name>, fieldPath: status.outputs.storage_account_id}} -- a bare string does not parse

### spec.rules

`[]AzureStorageObjectReplicationRule` · required

What gets replicated: each rule maps ONE source container to ONE
destination container, with optional filters. A policy carries up
to 1000 rules; a container pair appears in at most one rule.

- rule: {"repeated":{"minItems":"1","maxItems":"1000"}}

### spec.rules[].sourceContainerName

`string | valueFrom` · required

The container on the SOURCE account blobs are copied from, by name.
References an AzureStorageContainer's container_name output -- the
container must live on the policy's source account.

- references: AzureStorageContainer (`status.outputs.container_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureStorageContainer, name: <that resource's name>, fieldPath: status.outputs.container_name}} -- a bare string does not parse

### spec.rules[].destinationContainerName

`string | valueFrom` · required

The container on the DESTINATION account blobs are copied to, by
name. References an AzureStorageContainer's container_name output --
the container must live on the policy's destination account.

- references: AzureStorageContainer (`status.outputs.container_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureStorageContainer, name: <that resource's name>, fieldPath: status.outputs.container_name}} -- a bare string does not parse

### spec.rules[].copyBlobsCreatedAfter

`string` · optional (explicit presence)

Which existing blobs join the copy, beyond everything created after
the rule: OnlyNewObjects (the default -- no backfill), Everything
(backfill the whole container), or an RFC 3339 UTC instant
(backfill blobs created after that moment, e.g.
2026-01-01T00:00:00Z). Backfilling a large container takes time
proportional to its size.

- default: `OnlyNewObjects`
- rule: copy_blobs_created_after must be OnlyNewObjects, Everything, or an RFC 3339 UTC instant like 2026-01-01T00:00:00Z

### spec.rules[].prefixMatch

`[]string`

Replicate ONLY blobs whose names start with one of these prefixes
(e.g. invoices/2026) -- an INCLUDE filter, matching ARM's own
prefixMatch semantics. Empty means every blob in the container.

- rule: {"repeated":{"items":{"string":{"minLen":"1"}}}}

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureStorageObjectReplication, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.source_object_replication_id` | `string` | The policy's ARM ID on the SOURCE account. Format: /subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.Storage/storageAccounts/{src}/objectReplicationPolicies/{policyId} |
| `status.outputs.destination_object_replication_id` | `string` | The policy's ARM ID on the DESTINATION account -- the authoritative copy (Azure assigns rule IDs here first). Format: /subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.Storage/storageAccounts/{dst}/objectReplicationPolicies/{policyId} |
| `status.outputs.policy_id` | `string` | The server-assigned policy GUID shared by both sides -- what `az storage account or-policy show --policy-id` and the monitoring surfaces key on. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.sourceStorageAccountId` | AzureStorageAccount | `status.outputs.storage_account_id` |
| `spec.destinationStorageAccountId` | AzureStorageAccount | `status.outputs.storage_account_id` |
| `spec.rules[].sourceContainerName` | AzureStorageContainer | `status.outputs.container_name` |
| `spec.rules[].destinationContainerName` | AzureStorageContainer | `status.outputs.container_name` |

## See Also

- [Overview](../README.md)
