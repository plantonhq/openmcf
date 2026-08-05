# AliCloudNasFileSystem

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `alicloud.planton.dev/v1`

AliCloudNasFileSystemSpec defines the configuration for an Alibaba Cloud
Network Attached Storage (NAS) file system with a bundled access group and
VPC mount target.

NAS provides fully managed, elastic, shared file storage that supports NFS
and SMB protocols. This component bundles the file system, an optional
access group with access rules, and a VPC mount target into a single
deployable unit (per DD07 composite bundling) because a file system without
a mount target is unreachable.

The bundled flow:
  1. Create the NAS file system with the chosen protocol and storage type.
  2. (Conditional) If access_rules are specified, create a custom access
     group with those rules.
  3. Create a mount target in the specified VPC/VSwitch, associated with
     the custom access group or the default VPC group.

When no access_rules are specified, the mount target uses the built-in
DEFAULT_VPC_GROUP_NAME access group, which grants full read-write access
from all IP addresses within the VPC.

Two file system types are supported:
  - standard (default): auto-scaling capacity, general-purpose workloads.
    Storage types: "Performance", "Capacity", "Premium".
  - extreme: dedicated throughput with fixed pre-allocated capacity.
    Storage types: "standard", "advance". Requires zone_id and capacity.

Important: protocol_type, storage_type, and file_system_type are immutable
after creation -- changing any of these requires destroying and recreating
the file system.

Provider resources:
  Terraform: alicloud_nas_file_system + alicloud_nas_access_group
             + alicloud_nas_access_rule + alicloud_nas_mount_target
  Pulumi:    nas.FileSystem + nas.AccessGroup + nas.AccessRule
             + nas.MountTarget

## Example

```yaml
apiVersion: alicloud.planton.dev/v1
kind: AliCloudNasFileSystem
metadata:
  name: alicloudnasfilesystem-demo
spec:
  region: cn-hangzhou
  protocolType: NFS
  storageType: Performance
  description: Demo NAS file system for local testing
  vpcId:
    value: vpc-demo123
  vswitchId:
    value: vsw-demo123
  accessRules:
    - sourceCidrIp: "10.0.0.0/8"
      rwAccessType: RDWR
      userAccessType: no_squash
      priority: 1
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.fileSystemType` | `string` |  | `standard` |  |
| `spec.protocolType` | `string` | yes |  |  |
| `spec.storageType` | `string` | yes |  |  |
| `spec.description` | `string` |  |  |  |
| `spec.encryption` | `AliCloudNasEncryption` |  |  |  |
| `spec.encryption.encryptType` | `int32` | yes |  |  |
| `spec.encryption.kmsKeyId` | `string` |  |  |  |
| `spec.capacity` | `int32` |  |  |  |
| `spec.zoneId` | `string` |  |  |  |
| `spec.vpcId` | `string \| valueFrom` | yes |  | AliCloudVpc (`status.outputs.vpc_id`) |
| `spec.vswitchId` | `string \| valueFrom` | yes |  | AliCloudVswitch (`status.outputs.vswitch_id`) |
| `spec.accessRules` | `[]AliCloudNasAccessRule` |  |  |  |
| `spec.accessRules[].sourceCidrIp` | `string` | yes |  |  |
| `spec.accessRules[].rwAccessType` | `string` |  | `RDWR` |  |
| `spec.accessRules[].userAccessType` | `string` |  | `no_squash` |  |
| `spec.accessRules[].priority` | `int32` |  | `1` |  |
| `spec.resourceGroupId` | `string` |  |  |  |
| `spec.tags` | `map<string, string>` |  |  |  |

## Field Details

### spec.region

`string` · required

Alibaba Cloud region where the file system will be created.
Must match the region of the VPC and VSwitch used for the mount target.
Examples: "cn-hangzhou", "cn-shanghai", "us-west-1", "ap-southeast-1".

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.fileSystemType

`string` · optional (explicit presence)

File system type. Determines the performance model and available
storage types. This is immutable after creation (ForceNew).

  "standard" -- auto-scaling capacity, general-purpose workloads.
  "extreme"  -- dedicated throughput, fixed capacity, requires zone_id.
Default: "standard"

- default: `standard`
- rule: file_system_type must be one of: standard, extreme

### spec.protocolType

`string` · required

NAS protocol type. Determines how clients mount and access the file
system. This is immutable after creation (ForceNew).

  "NFS" -- Network File System, the standard Linux/Unix protocol.
  "SMB" -- Server Message Block, for Windows compatibility.

- rule: protocol_type must be one of: NFS, SMB
- rule: {"required":true}

### spec.storageType

`string` · required

Storage type. Determines the IOPS, throughput, and latency
characteristics. This is immutable after creation (ForceNew).

For standard file systems:
  "Performance" -- SSD-backed, low latency, up to 10 GiB/s throughput
  "Capacity"    -- HDD-backed, cost-effective for warm/cold data
  "Premium"     -- next-gen SSD, enhanced throughput and IOPS

For extreme file systems:
  "standard"    -- baseline extreme NAS
  "advance"     -- higher throughput extreme NAS

- rule: storage_type must be one of: Performance, Capacity, Premium, standard, advance
- rule: {"required":true}

### spec.description

`string`

Human-readable description of the file system's purpose.

### spec.encryption

`AliCloudNasEncryption`

Encryption configuration. When omitted, the file system is not encrypted.
Encryption type is immutable after creation (ForceNew).

### spec.encryption.encryptType

`int32` · required

Encryption type.
  1 = NAS-managed key: Alibaba Cloud manages the encryption key.
      No additional configuration needed.
  2 = Customer-managed KMS key: You manage the key via KMS.
      Requires kms_key_id to be set.

- rule: encrypt_type must be 1 (NAS-managed) or 2 (KMS customer-managed)
- rule: {"required":true}

### spec.encryption.kmsKeyId

`string`

KMS key ID for customer-managed encryption.
Required when encrypt_type is 2. Ignored when encrypt_type is 1.
Create one with AliCloudKmsKey and reference its key_id output.

### spec.capacity

`int32`

File system capacity in GiB.
Required for extreme NAS (minimum 100 GiB). Ignored for standard NAS,
which auto-scales capacity based on stored data.

### spec.zoneId

`string`

Availability zone for the file system.
Required for extreme NAS. Optional for standard NAS (the service
auto-assigns a zone when omitted).
Format: "{region}-{letter}", e.g. "cn-hangzhou-a", "cn-hangzhou-b".

### spec.vpcId

`string | valueFrom` · required

VPC ID for the mount target.
The mount target is created in this VPC, making the file system
accessible to resources within it.

- references: AliCloudVpc (`status.outputs.vpc_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AliCloudVpc, name: <that resource's name>, fieldPath: status.outputs.vpc_id}} -- a bare string does not parse

### spec.vswitchId

`string | valueFrom` · required

VSwitch ID for the mount target.
Must be in the VPC specified by vpc_id. Determines which availability
zone the mount target is placed in.

- references: AliCloudVswitch (`status.outputs.vswitch_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AliCloudVswitch, name: <that resource's name>, fieldPath: status.outputs.vswitch_id}} -- a bare string does not parse

### spec.accessRules

`[]AliCloudNasAccessRule`

Access rules that control which IP ranges can mount the file system and
with what permissions.

When one or more access_rules are specified, a custom access group is
auto-created and associated with the mount target. This gives fine-grained
control over which subnets or hosts can access the file system.

When omitted (empty list), the mount target uses the built-in
DEFAULT_VPC_GROUP_NAME access group, which allows full read-write access
from all IP addresses within the VPC.

### spec.accessRules[].sourceCidrIp

`string` · required

Source CIDR IP address or block.
Allows traffic from this network range to access the file system.
Examples: "0.0.0.0/0" for all IPs, "10.0.0.0/8" for a private range,
"10.0.1.5" for a single host.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.accessRules[].rwAccessType

`string` · optional (explicit presence)

Read-write access type.
  "RDWR"   -- full read-write access (default)
  "RDONLY" -- read-only access
Default: "RDWR"

- default: `RDWR`
- rule: rw_access_type must be one of: RDWR, RDONLY

### spec.accessRules[].userAccessType

`string` · optional (explicit presence)

User identity mapping for NFS root access.
  "no_squash"    -- preserve original user identity (default)
  "root_squash"  -- map root (uid 0) to anonymous user
  "all_squash"   -- map all users to anonymous user
Default: "no_squash"

- default: `no_squash`
- rule: user_access_type must be one of: no_squash, root_squash, all_squash

### spec.accessRules[].priority

`int32` · optional (explicit presence)

Priority of this access rule. Lower values have higher precedence.
Range: 1-100.
Default: 1

- default: `1`
- rule: priority must be between 1 and 100

### spec.resourceGroupId

`string`

Alibaba Cloud resource group ID for organizational grouping (per DD05).
If omitted, the file system is placed in the account's default resource
group.

### spec.tags

`map<string, string>`

Tags to apply to the file system resource.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AliCloudNasFileSystem, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.file_system_id` | `string` | The NAS file system ID assigned by Alibaba Cloud (e.g., "1ca404a348"). Used to reference this file system from other resources such as snapshot policies or lifecycle policies. |
| `status.outputs.mount_target_domain` | `string` | The mount target domain name. This is the endpoint used to mount the file system from ECS instances, containers, or any resource within the VPC. NFS mount example: mount -t nfs -o vers=4,minorversion=0,noresvport <domain>:/ /mnt/nas SMB mount example (Windows): net use Z: \\<domain>\myshare |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.vpcId` | AliCloudVpc | `status.outputs.vpc_id` |
| `spec.vswitchId` | AliCloudVswitch | `status.outputs.vswitch_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
