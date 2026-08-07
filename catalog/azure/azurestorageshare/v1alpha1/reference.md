# AzureStorageShare

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**AzureStorageShareSpec** defines the configuration for creating an
Azure Files share inside an Azure Storage Account: the SMB/NFS file
system unit. VMs, AKS pods, and container apps mount shares for shared
POSIX-style state -- lift-and-shift app data, user profiles (FSLogix),
CI caches, shared content -- and Azure bills, throttles, tiers, and
snapshots at the share level.

Shares are many-per-account with independent lifecycles, which is why
they are a first-class kind referencing the account rather than a list
folded into the account's spec. The parent is fixed at creation: a
share cannot move between accounts.

## Example

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureStorageShare
metadata:
  name: test-storage-share
spec:
  storageAccountId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Storage/storageAccounts/plantonhackstorage
  shareName: hack-team-files
  quotaGb: 500
  # Exercises the protocol enum mapping (SMB is also the default; stated
  # to prove the seam renders).
  enabledProtocol: SMB
  # Exercises the tier enum mapping (Cool = cheapest at rest, priciest
  # per operation).
  accessTier: COOL
  # Exercises the ACL block rendering, including a policy with the full
  # validity window and one leaving the window to the SAS token.
  acls:
    - id: readers
      accessPolicies:
        - permissions: rl
          start: "2026-07-01T00:00:00Z"
          expiry: "2027-07-01T00:00:00Z"
    - id: writers
      accessPolicies:
        - permissions: rwdl
  metadata:
    purpose: hack-validation
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.storageAccountId` | `string \| valueFrom` | yes |  | AzureStorageAccount (`status.outputs.storage_account_id`) |
| `spec.shareName` | `string` | yes |  |  |
| `spec.quotaGb` | `int32` | yes |  |  |
| `spec.enabledProtocol` | `enum` |  |  |  |
| `spec.accessTier` | `enum` |  |  |  |
| `spec.acls` | `[]AzureStorageShareAcl` |  |  |  |
| `spec.acls[].id` | `string` | yes |  |  |
| `spec.acls[].accessPolicies` | `[]AzureStorageShareAclAccessPolicy` |  |  |  |
| `spec.acls[].accessPolicies[].permissions` | `string` | yes |  |  |
| `spec.acls[].accessPolicies[].start` | `string` |  |  |  |
| `spec.acls[].accessPolicies[].expiry` | `string` |  |  |  |
| `spec.metadata` | `map<string, string>` |  |  |  |

## Field Details

### spec.storageAccountId

`string | valueFrom` · required

The storage account the share lives in, by ARM ID. References an
AzureStorageAccount's storage_account_id output so the account and
its shares compose in one manifest set. Fixed at creation. The
account's kind gates what the share can be: NFS shares and the
PREMIUM tier need a FileStorage (premium file) account; quotas above
5120 GB on standard accounts need the account's
large_file_share_enabled.

- references: AzureStorageAccount (`status.outputs.storage_account_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureStorageAccount, name: <that resource's name>, fieldPath: status.outputs.storage_account_id}} -- a bare string does not parse

### spec.shareName

`string` · required

The share's name: 3-63 lowercase letters, digits, and hyphens;
starts and ends with a letter or digit; no consecutive hyphens.
Unique within the account (it becomes the mount path segment:
//{account}.file.core.windows.net/{name}). Changing the name
replaces the share.

- rule: share_name must be 3-63 lowercase letters, digits, and hyphens, starting and ending with a letter or digit, with no consecutive hyphens
- rule: {"required":true,"string":{"minLen":"3","maxLen":"63"}}

### spec.quotaGb

`int32` · required

The share's maximum size in gigabytes -- the provisioned quota SMB
clients see as the drive size and Azure enforces on writes. Standard
accounts allow 1-5120 GB (up to 102400 GB when the account enables
large_file_share_enabled); premium FileStorage accounts require at
least 100 GB and bill provisioned capacity whether used or not.
Grows in place; shrinking below used capacity fails.

- rule: {"required":true,"int32":{"lte":102400,"gte":1}}

### spec.enabledProtocol

`enum`

The file-sharing protocol. Unspecified means SMB -- what Windows
mounts and most Linux distributions mount via cifs. NFS (v4.1) is
for Linux workloads that need POSIX semantics (hard links, chmod)
and requires a premium FileStorage account; NFS shares are reachable
only over private network paths (no public-endpoint mounts).
Fixed at creation.

- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `azure_storage_share_protocol_unspecified` -- Not specified: SMB.
- `SMB` -- SMB 3.x/2.1 plus the REST API -- what Windows and most Linux mounts use. Works on every account kind.
- `NFS` -- NFS v4.1 -- POSIX semantics for Linux workloads. Requires a premium FileStorage account and private network reachability.

### spec.accessTier

`enum`

The share's performance/billing tier. Unspecified lets Azure pick
its default (TRANSACTION_OPTIMIZED on standard accounts, PREMIUM on
FileStorage accounts). HOT and COOL trade lower at-rest prices for
per-operation charges -- COOL suits shares that are mostly read
rarely; PREMIUM (SSD, provisioned) is required -- and the only legal
value -- on FileStorage accounts.

- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `azure_storage_share_access_tier_unspecified` -- Not specified: Azure's default for the account kind (TransactionOptimized on standard, Premium on FileStorage).
- `TRANSACTION_OPTIMIZED` -- Standard HDD-backed, optimized for transaction-heavy workloads -- the standard-account default.
- `HOT` -- Lower per-operation cost than TransactionOptimized with moderate at-rest cost -- general-purpose file serving.
- `COOL` -- Cheapest at rest, priciest per operation -- shares that are mostly archived and read rarely.
- `PREMIUM` -- SSD-backed provisioned performance -- required (and the only legal tier) on premium FileStorage accounts.

### spec.acls

`[]AzureStorageShareAcl`

Stored access policies (signed identifiers) for the share. Each
policy anchors shared-access-signature tokens: revoking or
shortening the policy immediately revokes every SAS issued against
it -- the operational reason to prefer policy-anchored SAS over
ad-hoc SAS. At most five policies per share (Azure's limit).

- rule: {"repeated":{"maxItems":"5"}}

### spec.acls[].id

`string` · required

The policy's identifier -- the name SAS tokens reference (1-64
characters). Keep it stable: rotating the id orphans SAS tokens
issued against the old one.

- rule: {"required":true,"string":{"minLen":"1","maxLen":"64"}}

### spec.acls[].accessPolicies

`[]AzureStorageShareAclAccessPolicy`

The policy's validity window and permissions. Azure allows the
window fields to live on the SAS token instead, so both are
optional here; permissions must be declared on the policy.

### spec.acls[].accessPolicies[].permissions

`string` · required

The share's data-plane permission letters, in Azure's strict order:
r (read), w (write), d (delete), l (list). E.g. "rl" for
read-and-list consumers, "rwdl" for full data access.

- rule: permissions must be a non-empty combination of r, w, d, l in that order (e.g. "rl", "rwdl")
- rule: {"required":true}

### spec.acls[].accessPolicies[].start

`string`

When the policy becomes valid, RFC 3339 UTC (e.g.
"2026-07-01T00:00:00Z"). Omit to leave the start open (valid
immediately or governed by the SAS token's own start).

### spec.acls[].accessPolicies[].expiry

`string`

When the policy expires, RFC 3339 UTC. Omit to leave expiry to the
SAS token -- but a policy-level expiry is the revocation lever, so
production policies should set it.

### spec.metadata

`map<string, string>`

Free-form metadata key/value pairs stored on the share (visible to
anyone who can read share properties -- not for secrets). Keys must
be valid C# identifiers per Azure's rule; lowercase is canonical
(Azure lowercases keys on read).

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureStorageShare, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.share_id` | `string` | The Azure Resource Manager ID of the share -- the management-plane identity ARM reads and policy target. Format: /subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.Storage/storageAccounts/{account}/fileServices/default/shares/{name} |
| `status.outputs.rbac_scope_id` | `string` | The scope data-plane role assignments target for share-level file access (Storage File Data SMB Share Reader/Contributor/Elevated Contributor). Azure Files RBAC deliberately uses a DIFFERENT segment than the management ID -- .../fileServices/default/fileshares/{name} -- so this output exists precisely so grants never have to rewrite the management ID by hand. |
| `status.outputs.share_name` | `string` | The share's name -- what mount commands, CSI volume definitions, and app settings reference within the account. |
| `status.outputs.storage_account_name` | `string` | The name of the storage account the share lives in, parsed from the resolved account ID -- saves consumers a second reference when they need the account/share pair. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.storageAccountId` | AzureStorageAccount | `status.outputs.storage_account_id` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AzureContainerAppEnvironmentStorage | `spec.shareName` | `status.outputs.share_name` |

## See Also

- [Overview](../README.md)
