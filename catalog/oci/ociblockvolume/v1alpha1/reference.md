# OciBlockVolume

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `oci.planton.dev/v1alpha1`

OciBlockVolumeSpec defines the specification for an OCI Block Volume --
a high-performance, durable block storage device that can be attached to
compute instances, with configurable performance tiers (VPUs/GB),
automatic performance tuning, cross-region replication for DR, and an
optional backup policy assignment for scheduled backups.

Performance tiers (vpus_per_gb):
  - 0     = Lower Cost       (2 IOPS/GB, 240 KB/s/GB)
  - 10    = Balanced          (60 IOPS/GB, 480 KB/s/GB)  [OCI default]
  - 20    = Higher Performance (75 IOPS/GB, 600 KB/s/GB)
  - 30-120 = Ultra High Performance (in increments of 10)

Excluded from v1:
  - source_details -- clone/restore from volume, backup, or replica
  - cluster_placement_group_id -- very niche placement constraint
  - defined_tags, system_tags -- managed by platform
  - freeform_tags -- auto-populated from metadata labels
  - is_auto_tune_enabled -- deprecated, use autotune_policies
  - backup_policy_id on volume -- deprecated, using assignment resource
  - size_in_mbs -- deprecated, use size_in_gbs
  - volume_backup_id -- deprecated, use source_details

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.compartmentId` | `string \| valueFrom` | yes |  | OciCompartment (`status.outputs.compartment_id`) |
| `spec.availabilityDomain` | `string` | yes |  |  |
| `spec.displayName` | `string` |  |  |  |
| `spec.sizeInGbs` | `int32` |  |  |  |
| `spec.vpusPerGb` | `int32` |  |  |  |
| `spec.kmsKeyId` | `string \| valueFrom` |  |  |  |
| `spec.isReservationsEnabled` | `bool` |  |  |  |
| `spec.autotunePolicies` | `[]AutotunePolicy` |  |  |  |
| `spec.autotunePolicies[].autotuneType` | `enum` |  |  |  |
| `spec.autotunePolicies[].maxVpusPerGb` | `int32` |  |  |  |
| `spec.blockVolumeReplicas` | `[]BlockVolumeReplica` |  |  |  |
| `spec.blockVolumeReplicas[].availabilityDomain` | `string` | yes |  |  |
| `spec.blockVolumeReplicas[].displayName` | `string` |  |  |  |
| `spec.blockVolumeReplicas[].xrrKmsKeyId` | `string \| valueFrom` |  |  |  |
| `spec.backupPolicyId` | `string \| valueFrom` |  |  |  |
| `spec.xrcKmsKeyId` | `string \| valueFrom` |  |  |  |

## Field Details

### spec.compartmentId

`string | valueFrom` · required

OCID of the compartment where the block volume will be created.

- references: OciCompartment (`status.outputs.compartment_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: OciCompartment, name: <that resource's name>, fieldPath: status.outputs.compartment_id}} -- a bare string does not parse

### spec.availabilityDomain

`string` · required

Availability domain where the volume is placed (e.g., "Uocm:US-ASHBURN-AD-1").
The volume and any attached compute instance must be in the same AD.
Changing this forces recreation.

- rule: {"string":{"minLen":"1"}}

### spec.displayName

`string`

Display name for the volume. When omitted, the metadata name is used.

### spec.sizeInGbs

`int32`

Size of the volume in gigabytes. Valid range: 50-32768 (50 GB to 32 TB).
Must be specified explicitly to prevent accidental creation at OCI's
1 TB default.

### spec.vpusPerGb

`int32` · optional (explicit presence)

Volume Performance Units per GB. Controls IOPS and throughput.
When unset, OCI defaults to 10 (Balanced). Set to 0 for Lower Cost.
Valid values: 0, 10, 20, 30-120 (in increments of 10).
Stored as Int64 string in OCI API.

### spec.kmsKeyId

`string | valueFrom`

OCID of a KMS master encryption key for volume encryption.
When unset, Oracle-managed keys are used.

- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.isReservationsEnabled

`bool`

Enables SCSI persistent reservation support on the volume.
Required for shared-storage clustering scenarios (e.g., Oracle RAC).

### spec.autotunePolicies

`[]AutotunePolicy`

Autotune policies that automatically adjust volume performance
(VPUs/GB) based on detachment state or workload patterns.

- rule: max_vpus_per_gb must be > 0 when autotune_type is performance_based

### spec.autotunePolicies[].autotuneType

`enum`

Type of autotune policy to apply.

- rule: {"enum":{"definedOnly":true,"notIn":[0]}}

Allowed values (use exactly as shown):

- `unspecified`
- `detached_volume` -- Automatically adjusts VPUs to 0 (Lower Cost) when volume is detached, and restores previous VPUs when re-attached.
- `performance_based` -- Dynamically adjusts VPUs based on workload demand, up to max_vpus_per_gb.

### spec.autotunePolicies[].maxVpusPerGb

`int32`

Maximum VPUs/GB for performance-based autotune.
Required when autotune_type is performance_based.

### spec.blockVolumeReplicas

`[]BlockVolumeReplica`

Cross-region block volume replicas for disaster recovery.
Each replica is placed in a target availability domain (can be
in a different region) and asynchronously replicated.

### spec.blockVolumeReplicas[].availabilityDomain

`string` · required

Availability domain for the replica (e.g., "Uocm:US-PHOENIX-AD-1").
Can be in a different region from the source volume.

- rule: {"string":{"minLen":"1"}}

### spec.blockVolumeReplicas[].displayName

`string`

Display name for the replica. When omitted, OCI generates one.

### spec.blockVolumeReplicas[].xrrKmsKeyId

`string | valueFrom`

OCID of a KMS key for encrypting the cross-region replica.

- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.backupPolicyId

`string | valueFrom`

OCID of a backup policy to assign to this volume.
OCI provides Oracle-defined policies (Gold, Silver, Bronze) or
custom user-defined policies. When set, the component creates an
oci_core_volume_backup_policy_assignment sub-resource.

- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.xrcKmsKeyId

`string | valueFrom`

OCID of a KMS key used to encrypt cross-region volume backups.
Only relevant when a backup policy with cross-region copy is assigned.
Changing this forces recreation.

- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

## Validation Rules

- `size_in_gbs_minimum`: size_in_gbs must be at least 50 (the OCI minimum block volume size in GB)

## Outputs

Reference an output from another manifest as `valueFrom: {kind: OciBlockVolume, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.volume_id` | `string` | OCID of the block volume. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.compartmentId` | OciCompartment | `status.outputs.compartment_id` |

## See Also

- [Overview](../README.md)
