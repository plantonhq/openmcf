# OciFileSystem

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `oci.planton.dev/v1`

OciFileSystemSpec defines the specification for an OCI File Storage
file system -- an NFS-compatible, fully managed network file system
bundled with a dedicated mount target and one or more NFS exports.

The mount target provides the NFS endpoint (IP address) in a subnet,
and exports define the paths at which the file system is accessible.
Export options control per-client NFS access (read/write, identity
squashing, privileged port requirements).

The auto-created export set on the mount target can optionally be
configured to control NFS capacity reporting via statfs.

Excluded from v1:
  - oci_file_storage_snapshot -- operational concern, independent lifecycle
  - oci_file_storage_replication -- cross-region, separate lifecycle
  - oci_file_storage_filesystem_snapshot_policy -- reusable across file systems
  - oci_file_storage_file_system_quota_rule -- advanced admin feature
  - oci_file_storage_outbound_connector -- specialized LDAP integration
  - Kerberos / LDAP ID mapping on mount target -- very low adoption
  - source_snapshot_id -- clone/restore scenario
  - defined_tags, system_tags -- managed by platform
  - freeform_tags -- auto-populated from metadata labels
  - locks -- platform-managed

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.compartmentId` | `string \| valueFrom` | yes |  | OciCompartment (`status.outputs.compartment_id`) |
| `spec.availabilityDomain` | `string` | yes |  |  |
| `spec.displayName` | `string` |  |  |  |
| `spec.kmsKeyId` | `string \| valueFrom` |  |  |  |
| `spec.filesystemSnapshotPolicyId` | `string \| valueFrom` |  |  |  |
| `spec.mountTarget` | `MountTarget` | yes |  |  |
| `spec.mountTarget.subnetId` | `string \| valueFrom` | yes |  | OciSubnet (`status.outputs.subnet_id`) |
| `spec.mountTarget.displayName` | `string` |  |  |  |
| `spec.mountTarget.hostnameLabel` | `string` |  |  |  |
| `spec.mountTarget.ipAddress` | `string` |  |  |  |
| `spec.mountTarget.nsgIds` | `[]string \| valueFrom` |  |  | OciSecurityGroup (`status.outputs.network_security_group_id`) |
| `spec.mountTarget.requestedThroughput` | `int64` |  |  |  |
| `spec.mountTarget.maxFsStatBytes` | `int64` |  |  |  |
| `spec.mountTarget.maxFsStatFiles` | `int64` |  |  |  |
| `spec.exports` | `[]Export` | yes |  |  |
| `spec.exports[].path` | `string` | yes |  |  |
| `spec.exports[].exportOptions` | `[]ExportOption` |  |  |  |
| `spec.exports[].exportOptions[].source` | `string` | yes |  |  |
| `spec.exports[].exportOptions[].access` | `enum` |  |  |  |
| `spec.exports[].exportOptions[].identitySquash` | `enum` |  |  |  |
| `spec.exports[].exportOptions[].requirePrivilegedSourcePort` | `bool` |  |  |  |
| `spec.exports[].exportOptions[].isAnonymousAccessAllowed` | `bool` |  |  |  |
| `spec.exports[].exportOptions[].anonymousUid` | `int64` |  |  |  |
| `spec.exports[].exportOptions[].anonymousGid` | `int64` |  |  |  |

## Field Details

### spec.compartmentId

`string | valueFrom` · required

OCID of the compartment where the file system and mount target
will be created.

- references: OciCompartment (`status.outputs.compartment_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: OciCompartment, name: <that resource's name>, fieldPath: status.outputs.compartment_id}} -- a bare string does not parse

### spec.availabilityDomain

`string` · required

Availability domain where the file system and mount target are created.
Both must reside in the same AD. Example: "Uocm:US-ASHBURN-AD-1".
Changing this forces recreation of all resources.

- rule: {"string":{"minLen":"1"}}

### spec.displayName

`string`

Display name for the file system. When omitted, OCI generates one.

### spec.kmsKeyId

`string | valueFrom`

OCID of a KMS master encryption key for server-side encryption.
When unset, Oracle-managed keys are used.

- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.filesystemSnapshotPolicyId

`string | valueFrom`

OCID of a filesystem snapshot policy to attach for automated snapshots.
The policy must exist in the same availability domain.

- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.mountTarget

`MountTarget` · required

Configuration for the dedicated NFS mount target. The mount target
provides the network endpoint (IP address) that clients use to
mount the file system via NFS.

Note: OCI's default service limit is 2 mount targets per availability
domain. Request a limit increase if deploying more than 2 file systems
in the same AD.

- rule: {"required":true}

### spec.mountTarget.subnetId

`string | valueFrom` · required

OCID of the subnet where the mount target will be created.
The subnet determines the VCN and availability domain for
NFS access. Changing this forces recreation.

- references: OciSubnet (`status.outputs.subnet_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: OciSubnet, name: <that resource's name>, fieldPath: status.outputs.subnet_id}} -- a bare string does not parse

### spec.mountTarget.displayName

`string`

Display name for the mount target. When omitted, OCI generates one.

### spec.mountTarget.hostnameLabel

`string`

DNS hostname label for the mount target within the VCN's DNS.
When combined with the subnet and VCN DNS labels, produces an
FQDN like <hostname>.<subnet>.<vcn>.oraclevcn.com.
Changing this forces recreation.

### spec.mountTarget.ipAddress

`string`

Specific private IP address to assign. Must be available in the
subnet's CIDR. When omitted, OCI auto-assigns from the subnet.
Changing this forces recreation.

### spec.mountTarget.nsgIds

`[]string | valueFrom`

OCIDs of network security groups to associate with the mount target.
Controls NFS traffic (port 2049/TCP, 111/TCP for portmapper).

- references: OciSecurityGroup (`status.outputs.network_security_group_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: OciSecurityGroup, name: <that resource's name>, fieldPath: status.outputs.network_security_group_id}} -- a bare string does not parse

### spec.mountTarget.requestedThroughput

`int64`

Requested throughput in Mbps for the mount target. When omitted,
OCI uses the default throughput tier.

### spec.mountTarget.maxFsStatBytes

`int64`

Maximum NFS capacity in bytes reported to clients via statfs.
Configures the auto-created export set. When omitted, the actual
file system metered size is reported.

### spec.mountTarget.maxFsStatFiles

`int64`

Maximum file count reported to clients via statfs.
Configures the auto-created export set. When omitted, the actual
count is reported.

### spec.exports

`[]Export` · required

NFS export paths. Each export makes the file system accessible at
a specific path on the mount target. At least one export is required.

- rule: {"repeated":{"minItems":"1"}}
- rule: export path must start with '/'

### spec.exports[].path

`string` · required

NFS export path. Must start with '/' and be unique within the
mount target's export set. Example: "/shared-data".
Changing this forces recreation.

- rule: {"string":{"minLen":"1"}}

### spec.exports[].exportOptions

`[]ExportOption`

NFS access control rules. Each rule specifies a source CIDR
and its permissions. When omitted, OCI applies default access.

### spec.exports[].exportOptions[].source

`string` · required

Source IP address or CIDR block allowed to access this export.
Use "0.0.0.0/0" for unrestricted access within the network.

- rule: {"string":{"minLen":"1"}}

### spec.exports[].exportOptions[].access

`enum`

NFS access level for this source.

Allowed values (use exactly as shown):

- `access_unspecified`
- `read_write`
- `read_only`

### spec.exports[].exportOptions[].identitySquash

`enum`

Identity squashing mode for NFS requests from this source.

Allowed values (use exactly as shown):

- `identity_squash_unspecified`
- `no_squash`
- `root_squash`
- `all_squash`

### spec.exports[].exportOptions[].requirePrivilegedSourcePort

`bool`

When true, only connections from privileged source ports
(< 1024) are allowed. Standard for UNIX NFS clients.

### spec.exports[].exportOptions[].isAnonymousAccessAllowed

`bool`

When true, anonymous (unauthenticated) access is allowed.

### spec.exports[].exportOptions[].anonymousUid

`int64`

UNIX UID to map anonymous or squashed users to.
Typically 65534 (nobody). Stored as Int64 string in OCI API.

### spec.exports[].exportOptions[].anonymousGid

`int64`

UNIX GID to map anonymous or squashed users to.
Typically 65534 (nogroup). Stored as Int64 string in OCI API.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: OciFileSystem, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.file_system_id` | `string` | OCID of the file system. |
| `status.outputs.mount_target_id` | `string` | OCID of the mount target. |
| `status.outputs.mount_target_ip_address` | `string` | Private IP address of the mount target. Used in NFS mount commands: mount -t nfs <ip>:<export_path> /local/mount/point |
| `status.outputs.export_set_id` | `string` | OCID of the export set associated with the mount target. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.compartmentId` | OciCompartment | `status.outputs.compartment_id` |
| `spec.mountTarget.subnetId` | OciSubnet | `status.outputs.subnet_id` |
| `spec.mountTarget.nsgIds` | OciSecurityGroup | `status.outputs.network_security_group_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
