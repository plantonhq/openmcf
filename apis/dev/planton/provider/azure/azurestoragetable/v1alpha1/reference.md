# AzureStorageTable

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**AzureStorageTableSpec** defines the configuration for creating a
Storage table inside an Azure Storage Account: the serverless NoSQL
key/value store of Azure storage. Applications store schemaless
entities addressed by partition key + row key -- device state, user
profiles, audit trails, IoT telemetry -- at petabyte scale with
single-digit-millisecond point reads and no capacity provisioning.
(Cosmos DB's Table API is the premium sibling -- global distribution,
throughput SLAs -- at a very different price point.)

Tables are many-per-account with independent lifecycles, which is why
they are a first-class kind referencing the account rather than a list
folded into the account's spec. The parent is fixed at creation: a
table cannot move between accounts.

Operational contract: the provider drives table creation and stored
access policies through the table DATA PLANE with shared-key
authorization, so the parent account must keep
shared_access_key_enabled true (Azure's default) for deploys to work.

## Example

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureStorageTable
metadata:
  name: test-storage-table
spec:
  storageAccountId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Storage/storageAccounts/plantonhackstorage
  tableName: hackEntities
  # Exercises the ACL block rendering -- table policies require the full
  # validity window (start + expiry).
  acls:
    - id: readers
      accessPolicies:
        - permissions: r
          start: "2026-07-01T00:00:00Z"
          expiry: "2027-07-01T00:00:00Z"
    - id: writers
      accessPolicies:
        - permissions: raud
          start: "2026-07-01T00:00:00Z"
          expiry: "2027-07-01T00:00:00Z"
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.storageAccountId` | `string \| valueFrom` | yes |  | AzureStorageAccount (`status.outputs.storage_account_id`) |
| `spec.tableName` | `string` | yes |  |  |
| `spec.acls` | `[]AzureStorageTableAcl` |  |  |  |
| `spec.acls[].id` | `string` | yes |  |  |
| `spec.acls[].accessPolicies` | `[]AzureStorageTableAclAccessPolicy` |  |  |  |
| `spec.acls[].accessPolicies[].permissions` | `string` | yes |  |  |
| `spec.acls[].accessPolicies[].start` | `string` | yes |  |  |
| `spec.acls[].accessPolicies[].expiry` | `string` | yes |  |  |

## Field Details

### spec.storageAccountId

`string | valueFrom` · required

The storage account the table lives in, by ARM ID. References an
AzureStorageAccount's storage_account_id output so the account and
its tables compose in one manifest set. Fixed at creation.

- references: AzureStorageAccount (`status.outputs.storage_account_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureStorageAccount, name: <that resource's name>, fieldPath: status.outputs.storage_account_id}} -- a bare string does not parse

### spec.tableName

`string` · required

The table's name: 3-63 alphanumeric characters starting with a
letter -- no hyphens or underscores, and never the literal word
"table" (Azure reserves it). Unique within the account. Changing
the name replaces the table.

- rule: table_name must be 3-63 alphanumeric characters starting with a letter, and cannot be the literal word "table"
- rule: {"required":true}

### spec.acls

`[]AzureStorageTableAcl`

Stored access policies (signed identifiers) for the table. Each
policy anchors shared-access-signature tokens: revoking or
shortening the policy immediately revokes every SAS issued against
it -- the operational reason to prefer policy-anchored SAS over
ad-hoc SAS. At most five policies per table (Azure's limit).

- rule: {"repeated":{"maxItems":"5"}}

### spec.acls[].id

`string` · required

The policy's identifier -- the name SAS tokens reference (1-64
characters). Keep it stable: rotating the id orphans SAS tokens
issued against the old one.

- rule: {"required":true,"string":{"minLen":"1","maxLen":"64"}}

### spec.acls[].accessPolicies

`[]AzureStorageTableAclAccessPolicy`

The policy's validity window and permissions. Unlike blob/share
policies, Azure requires table policies to carry the full window --
start, expiry, and permissions are all mandatory here.

### spec.acls[].accessPolicies[].permissions

`string` · required

The table's data-plane permission letters, in Azure's strict order:
r (read/query), a (add), u (update), d (delete). E.g. "r" for
query-only consumers, "raud" for full entity access.

- rule: permissions must be a non-empty combination of r, a, u, d in that order (e.g. "r", "raud")
- rule: {"required":true}

### spec.acls[].accessPolicies[].start

`string` · required

When the policy becomes valid, RFC 3339 UTC (e.g.
"2026-07-01T00:00:00Z"). Required on table policies.

- rule: {"required":true}

### spec.acls[].accessPolicies[].expiry

`string` · required

When the policy expires, RFC 3339 UTC. Required on table policies
-- the expiry is the revocation lever for every SAS anchored here.

- rule: {"required":true}

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureStorageTable, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.table_id` | `string` | The Azure Resource Manager ID of the table. Role assignments (Storage Table Data Reader/Contributor) scope to it for table-level data access. Format: /subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.Storage/storageAccounts/{account}/tableServices/default/tables/{name} |
| `status.outputs.table_name` | `string` | The table's name -- what SDK clients, Functions table bindings, and app settings reference within the account. |
| `status.outputs.storage_account_name` | `string` | The name of the storage account the table lives in, parsed from the resolved account ID -- saves consumers a second reference when they need the account/table pair. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.storageAccountId` | AzureStorageAccount | `status.outputs.storage_account_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
