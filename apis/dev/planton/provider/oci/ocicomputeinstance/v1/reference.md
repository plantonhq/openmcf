# OciComputeInstance

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `oci.planton.dev/v1`

OciComputeInstanceSpec defines the specification for an Oracle Cloud
Infrastructure compute instance.

A compute instance is a virtual machine or bare metal host running on OCI
infrastructure. Flex shapes allow configuring the exact number of OCPUs and
memory; standard shapes use fixed configurations. The instance's primary
VNIC is configured inline via create_vnic_details, which determines the
subnet, network security groups, and public IP assignment.

Cloud-init scripts and SSH keys are passed through the metadata map using
the well-known keys "user_data" (base64-encoded) and "ssh_authorized_keys".

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.compartmentId` | `string \| valueFrom` | yes |  | OciCompartment (`status.outputs.compartment_id`) |
| `spec.availabilityDomain` | `string` | yes |  |  |
| `spec.shape` | `string` | yes |  |  |
| `spec.displayName` | `string` |  |  |  |
| `spec.shapeConfig` | `ShapeConfig` |  |  |  |
| `spec.shapeConfig.ocpus` | `float` |  |  |  |
| `spec.shapeConfig.memoryInGbs` | `float` |  |  |  |
| `spec.shapeConfig.baselineOcpuUtilization` | `string` |  |  |  |
| `spec.shapeConfig.nvmes` | `int32` |  |  |  |
| `spec.sourceDetails` | `SourceDetails` | yes |  |  |
| `spec.sourceDetails.sourceType` | `enum` |  |  |  |
| `spec.sourceDetails.sourceId` | `string` | yes |  |  |
| `spec.sourceDetails.bootVolumeSizeInGbs` | `int64` |  |  |  |
| `spec.sourceDetails.bootVolumeVpusPerGb` | `int64` |  |  |  |
| `spec.sourceDetails.kmsKeyId` | `string \| valueFrom` |  |  | OciKmsKey (`status.outputs.key_id`) |
| `spec.createVnicDetails` | `CreateVnicDetails` | yes |  |  |
| `spec.createVnicDetails.subnetId` | `string \| valueFrom` | yes |  | OciSubnet (`status.outputs.subnet_id`) |
| `spec.createVnicDetails.nsgIds` | `[]string \| valueFrom` |  |  | OciSecurityGroup (`status.outputs.network_security_group_id`) |
| `spec.createVnicDetails.assignPublicIp` | `bool` |  |  |  |
| `spec.createVnicDetails.displayName` | `string` |  |  |  |
| `spec.createVnicDetails.hostnameLabel` | `string` |  |  |  |
| `spec.createVnicDetails.privateIp` | `string` |  |  |  |
| `spec.createVnicDetails.skipSourceDestCheck` | `bool` |  |  |  |
| `spec.createVnicDetails.assignPrivateDnsRecord` | `bool` |  |  |  |
| `spec.metadata` | `map<string, string>` |  |  |  |
| `spec.faultDomain` | `string` |  |  |  |
| `spec.isPvEncryptionInTransitEnabled` | `bool` |  |  |  |
| `spec.agentConfig` | `AgentConfig` |  |  |  |
| `spec.agentConfig.areAllPluginsDisabled` | `bool` |  |  |  |
| `spec.agentConfig.isManagementDisabled` | `bool` |  |  |  |
| `spec.agentConfig.isMonitoringDisabled` | `bool` |  |  |  |
| `spec.agentConfig.pluginsConfig` | `[]PluginConfig` |  |  |  |
| `spec.agentConfig.pluginsConfig[].name` | `string` | yes |  |  |
| `spec.agentConfig.pluginsConfig[].desiredState` | `enum` |  |  |  |
| `spec.availabilityConfig` | `AvailabilityConfig` |  |  |  |
| `spec.availabilityConfig.isLiveMigrationPreferred` | `bool` |  |  |  |
| `spec.availabilityConfig.recoveryAction` | `enum` |  |  |  |
| `spec.launchOptions` | `LaunchOptions` |  |  |  |
| `spec.launchOptions.bootVolumeType` | `string` |  |  |  |
| `spec.launchOptions.networkType` | `string` |  |  |  |
| `spec.launchOptions.firmware` | `enum` |  |  |  |
| `spec.launchOptions.isPvEncryptionInTransitEnabled` | `bool` |  |  |  |
| `spec.launchOptions.isConsistentVolumeNamingEnabled` | `bool` |  |  |  |
| `spec.instanceOptions` | `InstanceOptions` |  |  |  |
| `spec.instanceOptions.areLegacyImdsEndpointsDisabled` | `bool` |  |  |  |
| `spec.preemptibleInstanceConfig` | `PreemptibleInstanceConfig` |  |  |  |
| `spec.preemptibleInstanceConfig.preserveBootVolume` | `bool` |  |  |  |
| `spec.capacityReservationId` | `string \| valueFrom` |  |  |  |
| `spec.dedicatedVmHostId` | `string \| valueFrom` |  |  |  |
| `spec.platformConfig` | `PlatformConfig` |  |  |  |
| `spec.platformConfig.type` | `enum` |  |  |  |
| `spec.platformConfig.isSecureBootEnabled` | `bool` |  |  |  |
| `spec.platformConfig.isMeasuredBootEnabled` | `bool` |  |  |  |
| `spec.platformConfig.isTrustedPlatformModuleEnabled` | `bool` |  |  |  |
| `spec.platformConfig.isMemoryEncryptionEnabled` | `bool` |  |  |  |
| `spec.platformConfig.isSymmetricMultiThreadingEnabled` | `bool` |  |  |  |
| `spec.platformConfig.areVirtualInstructionsEnabled` | `bool` |  |  |  |
| `spec.platformConfig.isAccessControlServiceEnabled` | `bool` |  |  |  |
| `spec.platformConfig.isInputOutputMemoryManagementUnitEnabled` | `bool` |  |  |  |
| `spec.platformConfig.numaNodesPerSocket` | `string` |  |  |  |
| `spec.platformConfig.percentageOfCoresEnabled` | `int32` |  |  |  |

## Field Details

### spec.compartmentId

`string | valueFrom` · required

OCID of the compartment where the instance will be created.

- references: OciCompartment (`status.outputs.compartment_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: OciCompartment, name: <that resource's name>, fieldPath: status.outputs.compartment_id}} -- a bare string does not parse

### spec.availabilityDomain

`string` · required

Availability domain where the instance will be placed
(e.g. "Ixxj:US-ASHBURN-AD-1"). Changing this forces recreation.

- rule: {"string":{"minLen":"1"}}

### spec.shape

`string` · required

Compute shape that determines the base hardware profile
(e.g. "VM.Standard.E4.Flex", "VM.Standard.A1.Flex", "BM.Standard3.64").
Flex shapes require shape_config to specify OCPUs and memory.

- rule: {"string":{"minLen":"1"}}

### spec.displayName

`string`

Human-readable name shown in the OCI Console.
Falls back to metadata.name if not provided.

### spec.shapeConfig

`ShapeConfig`

Shape configuration for flexible shapes. Required when the chosen shape
name contains "Flex". Ignored for fixed shapes.

### spec.shapeConfig.ocpus

`float` · optional (explicit presence)

Number of OCPUs (Oracle CPUs). Each OCPU maps to a physical core
with simultaneous multi-threading.

### spec.shapeConfig.memoryInGbs

`float` · optional (explicit presence)

Amount of memory in GiB. Flex shapes allow a range per OCPU; consult
the shape documentation for valid ratios.

### spec.shapeConfig.baselineOcpuUtilization

`string`

Baseline OCPU utilization for burstable instances.
Valid values: "BASELINE_1_8", "BASELINE_1_2", "BASELINE_1_1".
Only applicable to burstable shapes.

### spec.shapeConfig.nvmes

`int32` · optional (explicit presence)

Number of NVMe drives for dense-IO shapes.

### spec.sourceDetails

`SourceDetails` · required

Boot source for the instance: either an image OCID or a boot volume OCID.

- rule: {"required":true}

### spec.sourceDetails.sourceType

`enum`

Whether the instance boots from an image or an existing boot volume.

- rule: {"enum":{"notIn":[0]}}

Allowed values (use exactly as shown):

- `source_type_unspecified`
- `image`
- `boot_volume`

### spec.sourceDetails.sourceId

`string` · required

OCID of the image or boot volume to launch from.

- rule: {"string":{"minLen":"1"}}

### spec.sourceDetails.bootVolumeSizeInGbs

`int64` · optional (explicit presence)

Size of the boot volume in GiB. When launching from an image, defaults
to the image's minimum size. Must be >= the image minimum.

### spec.sourceDetails.bootVolumeVpusPerGb

`int64` · optional (explicit presence)

VPUs per GB for the boot volume (10=Balanced, 20=Higher Performance,
30-120=Ultra High Performance). Defaults to 10.

### spec.sourceDetails.kmsKeyId

`string | valueFrom`

OCID of a KMS key to encrypt the boot volume at rest.

- references: OciKmsKey (`status.outputs.key_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: OciKmsKey, name: <that resource's name>, fieldPath: status.outputs.key_id}} -- a bare string does not parse

### spec.createVnicDetails

`CreateVnicDetails` · required

Primary VNIC configuration that determines the instance's network placement.

- rule: {"required":true}

### spec.createVnicDetails.subnetId

`string | valueFrom` · required

OCID of the subnet for the primary VNIC.

- references: OciSubnet (`status.outputs.subnet_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: OciSubnet, name: <that resource's name>, fieldPath: status.outputs.subnet_id}} -- a bare string does not parse

### spec.createVnicDetails.nsgIds

`[]string | valueFrom`

OCIDs of network security groups to associate with the VNIC.

- references: OciSecurityGroup (`status.outputs.network_security_group_id`)
- rule: {"repeated":{"maxItems":"5"}}
- rule: write as {value: <literal>} or {valueFrom: {kind: OciSecurityGroup, name: <that resource's name>, fieldPath: status.outputs.network_security_group_id}} -- a bare string does not parse

### spec.createVnicDetails.assignPublicIp

`bool` · optional (explicit presence)

Whether to assign a public IP. When unset, uses the subnet default
(public subnets assign; private subnets do not).

### spec.createVnicDetails.displayName

`string`

Display name for the VNIC in the OCI Console.

### spec.createVnicDetails.hostnameLabel

`string`

Hostname label for DNS within the subnet's DNS domain.
Must be alphanumeric, start with a letter, max 63 chars.

### spec.createVnicDetails.privateIp

`string`

Specific private IP address to assign. Must be available within the
subnet's CIDR. When omitted, OCI auto-assigns.

### spec.createVnicDetails.skipSourceDestCheck

`bool` · optional (explicit presence)

When true, disables source/destination checking on the VNIC.
Required for NAT instances or virtual routers.

### spec.createVnicDetails.assignPrivateDnsRecord

`bool` · optional (explicit presence)

Whether to register a private DNS record for the VNIC.

### spec.metadata

`map<string, string>`

Instance metadata key-value pairs. Common keys:
  - "ssh_authorized_keys": newline-separated public SSH keys
  - "user_data": base64-encoded cloud-init script

### spec.faultDomain

`string`

Fault domain within the availability domain (e.g. "FAULT-DOMAIN-1").
OCI distributes instances across fault domains for HA when unspecified.

### spec.isPvEncryptionInTransitEnabled

`bool` · optional (explicit presence)

When true, enables in-transit encryption for the paravirtualized boot
and data volume attachments. Changing this forces recreation.

### spec.agentConfig

`AgentConfig`

Oracle Cloud Agent configuration controlling monitoring, management,
and plugin behavior on the instance.

### spec.agentConfig.areAllPluginsDisabled

`bool` · optional (explicit presence)

When true, disables all Oracle Cloud Agent plugins.

### spec.agentConfig.isManagementDisabled

`bool` · optional (explicit presence)

When true, disables the management agent (OS Management, etc.).

### spec.agentConfig.isMonitoringDisabled

`bool` · optional (explicit presence)

When true, disables the monitoring agent (Compute Instance Monitoring).

### spec.agentConfig.pluginsConfig

`[]PluginConfig`

Per-plugin overrides.

### spec.agentConfig.pluginsConfig[].name

`string` · required

Plugin name (e.g. "Vulnerability Scanning", "OS Management Service Agent").

- rule: {"string":{"minLen":"1"}}

### spec.agentConfig.pluginsConfig[].desiredState

`enum`

Whether the plugin should be enabled or disabled.

- rule: {"enum":{"notIn":[0]}}

Allowed values (use exactly as shown):

- `desired_state_unspecified`
- `enabled`
- `disabled`

### spec.availabilityConfig

`AvailabilityConfig`

Availability and recovery behavior for infrastructure maintenance events.

### spec.availabilityConfig.isLiveMigrationPreferred

`bool` · optional (explicit presence)

When true, OCI prefers live migration over reboot during maintenance.

### spec.availabilityConfig.recoveryAction

`enum`

Action to take when the underlying host has an unplanned failure.

Allowed values (use exactly as shown):

- `recovery_action_unspecified`
- `restore_instance`
- `stop_instance`

### spec.launchOptions

`LaunchOptions`

Low-level launch options for boot volume type, network type, and firmware.
Most users can omit this; the defaults chosen by OCI based on the image
and shape are appropriate for the vast majority of workloads.

### spec.launchOptions.bootVolumeType

`string`

Emulation type for the boot volume attachment.
Valid values: "ISCSI", "SCSI", "IDE", "VFIO", "PARAVIRTUALIZED".

### spec.launchOptions.networkType

`string`

Emulation type for the primary network interface.
Valid values: "E1000", "VFIO", "PARAVIRTUALIZED".

### spec.launchOptions.firmware

`enum`

Firmware type for the instance.

Allowed values (use exactly as shown):

- `firmware_unspecified`
- `bios`
- `uefi_64`

### spec.launchOptions.isPvEncryptionInTransitEnabled

`bool` · optional (explicit presence)

In-transit encryption for the boot volume attachment (within launch options).

### spec.launchOptions.isConsistentVolumeNamingEnabled

`bool` · optional (explicit presence)

Consistent naming for attached volumes (e.g. /dev/oracleoci/...).

### spec.instanceOptions

`InstanceOptions`

Instance Metadata Service (IMDS) endpoint configuration.

### spec.instanceOptions.areLegacyImdsEndpointsDisabled

`bool` · optional (explicit presence)

When true, disables the legacy IMDSv1 endpoints. Recommended for
security: use only the IMDSv2 token-based endpoint.

### spec.preemptibleInstanceConfig

`PreemptibleInstanceConfig`

Configures the instance as preemptible (spot-like), allowing OCI to
reclaim it when capacity is needed. Significantly reduces cost for
fault-tolerant workloads.

### spec.preemptibleInstanceConfig.preserveBootVolume

`bool` · optional (explicit presence)

When true, the boot volume is preserved when the instance is preempted.
When false, both instance and boot volume are terminated.

### spec.capacityReservationId

`string | valueFrom`

OCID of a capacity reservation to launch this instance against.

- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.dedicatedVmHostId

`string | valueFrom`

OCID of a dedicated VM host to place this instance on.
Used for workloads requiring physical isolation for compliance or licensing.

- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.platformConfig

`PlatformConfig`

Platform-level security and hardware configuration. The applicable fields
depend on the platform type: VM shapes (amd_vm, intel_vm) support secure
boot, measured boot, TPM, and memory encryption; bare metal shapes support
additional hardware tuning (NUMA, core percentage, SMT, etc.).

### spec.platformConfig.type

`enum`

Platform type. Must match the instance shape family.

- rule: {"enum":{"notIn":[0]}}

Allowed values (use exactly as shown):

- `platform_type_unspecified`
- `amd_milan_bm`
- `amd_milan_bm_gpu`
- `amd_rome_bm`
- `amd_rome_bm_gpu`
- `amd_vm`
- `generic_bm`
- `intel_icelake_bm`
- `intel_skylake_bm`
- `intel_vm`

### spec.platformConfig.isSecureBootEnabled

`bool` · optional (explicit presence)

Enable Secure Boot to verify boot software signatures (VM + BM).

### spec.platformConfig.isMeasuredBootEnabled

`bool` · optional (explicit presence)

Enable Measured Boot for integrity measurements stored in the TPM (VM + BM).

### spec.platformConfig.isTrustedPlatformModuleEnabled

`bool` · optional (explicit presence)

Enable the Trusted Platform Module for secure key storage (VM + BM).

### spec.platformConfig.isMemoryEncryptionEnabled

`bool` · optional (explicit presence)

Enable AMD SEV or Intel TME memory encryption (VM + BM).

### spec.platformConfig.isSymmetricMultiThreadingEnabled

`bool` · optional (explicit presence)

Enable Symmetric Multi-Threading (SMT/Hyperthreading). BM shapes only.

### spec.platformConfig.areVirtualInstructionsEnabled

`bool` · optional (explicit presence)

Enable nested virtualization (AMD-V or VT-x). BM shapes only.

### spec.platformConfig.isAccessControlServiceEnabled

`bool` · optional (explicit presence)

Enable Access Control Service for PCI passthrough isolation. BM shapes only.

### spec.platformConfig.isInputOutputMemoryManagementUnitEnabled

`bool` · optional (explicit presence)

Enable IOMMU for device memory protection. BM shapes only.

### spec.platformConfig.numaNodesPerSocket

`string`

NUMA nodes per socket configuration. BM shapes only.
Valid values: "NPS0", "NPS1", "NPS2", "NPS4".

### spec.platformConfig.percentageOfCoresEnabled

`int32` · optional (explicit presence)

Percentage of cores enabled on the instance. BM shapes only.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: OciComputeInstance, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.instance_id` | `string` | OCID of the compute instance. |
| `status.outputs.private_ip` | `string` | Private IP address of the primary VNIC. |
| `status.outputs.public_ip` | `string` | Public IP address of the primary VNIC. Empty when no public IP is assigned. |
| `status.outputs.boot_volume_id` | `string` | OCID of the boot volume attached to the instance. |
| `status.outputs.availability_domain` | `string` | Availability domain where the instance was placed. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.compartmentId` | OciCompartment | `status.outputs.compartment_id` |
| `spec.sourceDetails.kmsKeyId` | OciKmsKey | `status.outputs.key_id` |
| `spec.createVnicDetails.subnetId` | OciSubnet | `status.outputs.subnet_id` |
| `spec.createVnicDetails.nsgIds` | OciSecurityGroup | `status.outputs.network_security_group_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
