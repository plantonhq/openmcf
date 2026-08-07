# AzureStorageDataLakeGen2Filesystem

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**AzureStorageDataLakeGen2FilesystemSpec** defines the configuration for
creating a Data Lake Storage Gen2 filesystem inside an Azure Storage
Account: the namespace unit of an analytics data lake. A filesystem is
where hierarchical namespace (HNS) data lives -- Spark, Databricks,
Synapse, and the abfss:// driver all address data as
abfss://{filesystem}@{account}.dfs.core.windows.net/path. Analytics
estates conventionally provision one filesystem per data-lake zone
(raw, curated, gold) so each zone carries its own access control.

Filesystems are many-per-account with independent lifecycles and are
the grant boundary for data-plane RBAC and POSIX ACLs, which is why
they are a first-class kind rather than a list folded into the
account's spec. The parent is fixed at creation: a filesystem cannot
move between accounts.

**The account must have hierarchical namespace enabled** (the storage
account's is_hns_enabled) for POSIX access control -- Azure rejects
ACL, owner, and group settings on a flat-namespace account at apply
time. A filesystem on a flat account is just a blob container wearing
a dfs endpoint; create it as an AzureStorageContainer instead.

## Example

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureStorageDataLakeGen2Filesystem
metadata:
  name: test-data-lake-filesystem
spec:
  storageAccountId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Storage/storageAccounts/plantonhacklake
  filesystemName: raw-zone
  # Exercises the encryption-scope seam (sent only when set).
  defaultEncryptionScope:
    value: regulatedScope
  # Exercises root ownership (object ID and $superuser literal forms;
  # both fields also take a valueFrom reference to a user-assigned
  # identity's principal_id).
  owner:
    value: 11111111-2222-3333-4444-555555555555
  group:
    value: $superuser
  # Exercises every ACE shape: both scopes, all four types, the
  # qualified USER/GROUP entries, and the rwx vocabulary. objectId
  # takes the value/valueFrom wrapper -- a workload identity's zone
  # access references its principal_id output.
  aces:
    - type: USER
      permissions: rwx
    - type: USER
      objectId:
        value: 11111111-2222-3333-4444-555555555555
      permissions: r-x
    - scope: DEFAULT
      type: GROUP
      objectId:
        value: 66666666-7777-8888-9999-000000000000
      permissions: r-x
    - type: MASK
      permissions: r-x
    - type: OTHER
      permissions: "---"
  # Exercises the properties map (values base64-encoded per Azure's
  # requirement).
  properties:
    environment: cHJvZHVjdGlvbg==
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.storageAccountId` | `string \| valueFrom` | yes |  | AzureStorageAccount (`status.outputs.storage_account_id`) |
| `spec.filesystemName` | `string` | yes |  |  |
| `spec.defaultEncryptionScope` | `string \| valueFrom` |  |  | AzureStorageEncryptionScope (`status.outputs.encryption_scope_name`) |
| `spec.owner` | `string \| valueFrom` |  |  | AzureUserAssignedIdentity (`status.outputs.principal_id`) |
| `spec.group` | `string \| valueFrom` |  |  | AzureUserAssignedIdentity (`status.outputs.principal_id`) |
| `spec.aces` | `[]AzureStorageDataLakeGen2FilesystemAce` |  |  |  |
| `spec.aces[].scope` | `enum` |  |  |  |
| `spec.aces[].type` | `enum` |  |  |  |
| `spec.aces[].objectId` | `string \| valueFrom` |  |  | AzureUserAssignedIdentity (`status.outputs.principal_id`) |
| `spec.aces[].permissions` | `string` | yes |  |  |
| `spec.properties` | `map<string, string>` |  |  |  |

## Field Details

### spec.storageAccountId

`string | valueFrom` · required

The storage account the filesystem lives in, by ARM ID. References
an AzureStorageAccount's storage_account_id output so the account
and its filesystems compose in one manifest set. The account should
carry is_hns_enabled: true -- POSIX access control (owner, group,
aces) is rejected on flat-namespace accounts. Fixed at creation.

- references: AzureStorageAccount (`status.outputs.storage_account_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureStorageAccount, name: <that resource's name>, fieldPath: status.outputs.storage_account_id}} -- a bare string does not parse

### spec.filesystemName

`string` · required

The filesystem's name: 3-63 lowercase letters, digits, and hyphens,
not starting with a hyphen ("$root" is the one special name Azure
also accepts). Unique within the account; it becomes the container
segment of every abfss:// and dfs URL. Changing the name replaces
the filesystem -- and everything stored in it.

- rule: filesystem_name must be 3-63 lowercase letters, digits, and hyphens, not starting with a hyphen (or the special name $root)
- rule: {"required":true}

### spec.defaultEncryptionScope

`string | valueFrom`

The encryption scope applied to data that doesn't name its own --
sub-account key isolation for just this filesystem (e.g. a
regulated zone encrypting under a customer-managed key while the
rest of the lake uses platform keys). References an
AzureStorageEncryptionScope's name output; the scope must live on
the SAME account as the filesystem. Fixed at creation.

- references: AzureStorageEncryptionScope (`status.outputs.encryption_scope_name`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureStorageEncryptionScope, name: <that resource's name>, fieldPath: status.outputs.encryption_scope_name}} -- a bare string does not parse

### spec.owner

`string | valueFrom`

The Entra principal that OWNS the filesystem's root path ("/"). The
owning user always holds the root's user-class permissions
regardless of the ACL. Takes an Entra object ID (a GUID -- for a
managed identity this is the PRINCIPAL id, not the client id) or
the special literal $superuser. References an
AzureUserAssignedIdentity's principal_id output so a workload
identity can own its zone; Entra users, groups, and service
principals are granted by literal object ID. Unset leaves Azure's
default owner ($superuser -- the account's shared-key principal).

- references: AzureUserAssignedIdentity (`status.outputs.principal_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureUserAssignedIdentity, name: <that resource's name>, fieldPath: status.outputs.principal_id}} -- a bare string does not parse

### spec.group

`string | valueFrom`

The Entra principal that owns the filesystem's root path ("/") as
its GROUP. Group-class ACL entries without an explicit object ID
evaluate against this group. Takes an Entra object ID (a GUID) or
the special literal $superuser; references an
AzureUserAssignedIdentity's principal_id output when a workload
identity plays the owning-group role, while Entra security groups
are granted by literal object ID. Unset leaves Azure's default
owning group.

- references: AzureUserAssignedIdentity (`status.outputs.principal_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureUserAssignedIdentity, name: <that resource's name>, fieldPath: status.outputs.principal_id}} -- a bare string does not parse

### spec.aces

`[]AzureStorageDataLakeGen2FilesystemAce`

The POSIX access control list applied to the filesystem's ROOT path
("/"). Access entries gate operations on the root itself; default
entries are the template newly created children inherit. Requires
hierarchical namespace on the account -- Azure rejects ACLs on
flat-namespace accounts at apply time. Directories deeper in the
tree manage their own ACLs (via SDKs or Storage Explorer); the
filesystem kind owns only the root.

- rule: object_id can only be set on USER and GROUP entries -- MASK and OTHER entries apply to every caller, not a specific principal

### spec.aces[].scope

`enum`

Whether this is an ACCESS entry (gates operations on the root path
itself) or a DEFAULT entry (the template newly created children
inherit). Unspecified means ACCESS -- Azure's default.

- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `azure_storage_data_lake_gen2_filesystem_ace_scope_unspecified` -- Not specified: ACCESS -- Azure's default.
- `ACCESS` -- Gates operations on the root path itself.
- `DEFAULT` -- The template newly created children inherit -- how a zone's permission posture propagates to files landing in it.

### spec.aces[].type

`enum`

The POSIX entry class. USER and GROUP entries may name a specific
Entra principal via object_id (or apply to the owning user/group
when unqualified); MASK caps the effective permissions of every
named entry; OTHER covers callers matched by no other entry.

- rule: {"enum":{"definedOnly":true,"notIn":[0]}}

Allowed values (use exactly as shown):

- `azure_storage_data_lake_gen2_filesystem_ace_type_unspecified` -- Not specified -- invalid; choose an explicit entry class.
- `USER` -- A user entry: the owning user when unqualified, or the principal named by object_id.
- `GROUP` -- A group entry: the owning group when unqualified, or the group named by object_id.
- `MASK` -- The mask entry: caps the EFFECTIVE permissions of every named user, named group, and owning-group entry.
- `OTHER` -- The other entry: callers matched by no other entry.

### spec.aces[].objectId

`string | valueFrom`

The Entra object ID (a GUID) of the user, group, service principal,
or managed identity this entry applies to. Only valid for USER and
GROUP entries; leave unset to address the root path's OWNING user
or group. References an AzureUserAssignedIdentity's principal_id
output so a workload identity's zone access composes by reference;
Entra users and security groups are granted by literal object ID.
NOTE: for a managed identity this is the PRINCIPAL id, not the
client id -- an ACL naming the client id silently never matches.

- references: AzureUserAssignedIdentity (`status.outputs.principal_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureUserAssignedIdentity, name: <that resource's name>, fieldPath: status.outputs.principal_id}} -- a bare string does not parse

### spec.aces[].permissions

`string` · required

The entry's permissions in POSIX short form: exactly three
characters -- r or - (read), w or - (write), x or - (execute, i.e.
traverse for directories). Examples: rwx (full), r-x (read and
list), --- (none).

- rule: permissions must be the three-character POSIX form [r-][w-][x-], e.g. rwx, r-x, or ---
- rule: {"required":true}

### spec.properties

`map<string, string>`

Free-form properties stored on the filesystem. Azure requires the
VALUES to be base64-encoded strings (the keys stay plain) -- e.g.
environment: cHJvZHVjdGlvbg== -- and returns them base64-encoded on
read. Visible to anyone who can read filesystem properties; not for
secrets.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureStorageDataLakeGen2Filesystem, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.filesystem_id` | `string` | The filesystem's Azure Resource Manager ID -- ADLS filesystems surface in ARM as blob containers, so this is the container-proxy ID and the scope data-plane role assignments (Storage Blob Data Reader/Contributor/Owner) target for filesystem-level access. Format: /subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.Storage/storageAccounts/{account}/blobServices/default/containers/{name} |
| `status.outputs.filesystem_name` | `string` | The filesystem's name -- the container segment of every abfss:// and dfs URL (abfss://{name}@{account}.dfs.core.windows.net/). |
| `status.outputs.storage_account_name` | `string` | The name of the storage account the filesystem lives in, parsed from the resolved account ID -- saves consumers a second reference when they need the account/filesystem pair (abfss URLs, Spark configs, mount definitions). |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.storageAccountId` | AzureStorageAccount | `status.outputs.storage_account_id` |
| `spec.defaultEncryptionScope` | AzureStorageEncryptionScope | `status.outputs.encryption_scope_name` |
| `spec.owner` | AzureUserAssignedIdentity | `status.outputs.principal_id` |
| `spec.group` | AzureUserAssignedIdentity | `status.outputs.principal_id` |
| `spec.aces[].objectId` | AzureUserAssignedIdentity | `status.outputs.principal_id` |

## See Also

- [Overview](../README.md)
