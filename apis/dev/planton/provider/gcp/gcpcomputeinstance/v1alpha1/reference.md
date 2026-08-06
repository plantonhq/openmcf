# GcpComputeInstance

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `gcp.planton.dev/v1alpha1`

GcpComputeInstanceSpec defines the configuration for a Google Compute
Engine virtual machine — the general-purpose compute node behind
databases-on-VM, stateful application servers, GPU workers, bastion
hosts, and any workload that needs a full OS rather than a managed
runtime.

Important behavioral notes:

  - zone, machine confidential mode, boot-disk source/encryption, NIC
    count and their networks, scratch disks, hostname, and reservation
    affinity are create-time decisions — changing them replaces the VM.
    Many other properties (machine_type, labels, metadata, service
    account, shielded config) are updated by stopping the VM, which the
    provider only does when allow_stopping_for_update is true.
  - Spot VMs (scheduling.provisioning_model = "SPOT") can be preempted
    at any time. Both engines set the API's legacy preemptible flag
    automatically for Spot; automatic restart is forced off.
  - deletion_protection guards against accidental instance deletion but
    does NOT protect data: the boot disk's auto_delete and each attached
    disk's own lifecycle are the data levers. Attached data disks are
    first-class GcpComputeDisk resources referenced here — they survive
    the VM unless their own configuration says otherwise.
  - desired_status stops ("TERMINATED"), suspends ("SUSPENDED"), or
    starts ("RUNNING") the VM without destroying it — compute billing
    stops while terminated; disks keep billing.

## Example

```yaml
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpComputeInstance
metadata:
  name: test-vm
spec:
  projectId:
    value: planton-testing
  zone: us-central1-a
  machineType: e2-micro
  bootDisk:
    image: debian-cloud/debian-12
    sizeGb: 10
    type: pd-balanced
  networkInterfaces:
    - network:
        value: default
      accessConfigs:
        - networkTier: PREMIUM
  scheduling:
    provisioningModel: SPOT
    onHostMaintenance: TERMINATE
    instanceTerminationAction: DELETE
  shieldedInstanceConfig:
    enableSecureBoot: true
  labels:
    env: testing
  tags:
    - http-server
  metadata:
    enable-oslogin: "TRUE"
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.projectId` | `string \| valueFrom` |  |  | GcpProject (`status.outputs.project_id`) |
| `spec.instanceName` | `string` |  |  |  |
| `spec.zone` | `string` | yes |  |  |
| `spec.machineType` | `string` | yes |  |  |
| `spec.description` | `string` |  |  |  |
| `spec.hostname` | `string` |  |  |  |
| `spec.bootDisk` | `GcpComputeInstanceBootDisk` | yes |  |  |
| `spec.bootDisk.image` | `string` |  |  |  |
| `spec.bootDisk.sourceSnapshot` | `string` |  |  |  |
| `spec.bootDisk.sourceDisk` | `string \| valueFrom` |  |  | GcpComputeDisk (`status.outputs.self_link`) |
| `spec.bootDisk.sizeGb` | `int32` |  |  |  |
| `spec.bootDisk.type` | `string` |  |  |  |
| `spec.bootDisk.autoDelete` | `bool` |  | `true` |  |
| `spec.bootDisk.deviceName` | `string` |  |  |  |
| `spec.bootDisk.kmsKey` | `string \| valueFrom` |  |  | GcpKmsKey (`status.outputs.key_id`) |
| `spec.bootDisk.diskLabels` | `map<string, string>` |  |  |  |
| `spec.bootDisk.provisionedIops` | `int64` |  |  |  |
| `spec.bootDisk.provisionedThroughput` | `int64` |  |  |  |
| `spec.bootDisk.architecture` | `string` |  |  |  |
| `spec.bootDisk.enableConfidentialCompute` | `bool` |  |  |  |
| `spec.bootDisk.resourcePolicies` | `[]string` |  |  |  |
| `spec.bootDisk.storagePool` | `string` |  |  |  |
| `spec.attachedDisks` | `[]GcpComputeInstanceAttachedDisk` |  |  |  |
| `spec.attachedDisks[].source` | `string \| valueFrom` | yes |  | GcpComputeDisk (`status.outputs.self_link`) |
| `spec.attachedDisks[].deviceName` | `string` |  |  |  |
| `spec.attachedDisks[].mode` | `string` |  |  |  |
| `spec.attachedDisks[].kmsKey` | `string \| valueFrom` |  |  | GcpKmsKey (`status.outputs.key_id`) |
| `spec.scratchDisks` | `[]GcpComputeInstanceScratchDisk` |  |  |  |
| `spec.scratchDisks[].interface` | `string` | yes |  |  |
| `spec.scratchDisks[].sizeGb` | `int32` |  |  |  |
| `spec.scratchDisks[].deviceName` | `string` |  |  |  |
| `spec.networkInterfaces` | `[]GcpComputeInstanceNetworkInterface` | yes |  |  |
| `spec.networkInterfaces[].network` | `string \| valueFrom` |  |  | GcpVpcNetwork (`status.outputs.network_self_link`) |
| `spec.networkInterfaces[].subnetwork` | `string \| valueFrom` |  |  | GcpSubnetwork (`status.outputs.subnetwork_self_link`) |
| `spec.networkInterfaces[].subnetworkProject` | `string` |  |  |  |
| `spec.networkInterfaces[].networkIp` | `string \| valueFrom` |  |  | GcpAddress (`status.outputs.address`) |
| `spec.networkInterfaces[].accessConfigs` | `[]GcpComputeInstanceAccessConfig` |  |  |  |
| `spec.networkInterfaces[].accessConfigs[].natIp` | `string \| valueFrom` |  |  | GcpAddress (`status.outputs.address`) |
| `spec.networkInterfaces[].accessConfigs[].networkTier` | `string` |  |  |  |
| `spec.networkInterfaces[].accessConfigs[].publicPtrDomainName` | `string` |  |  |  |
| `spec.networkInterfaces[].ipv6AccessConfigs` | `[]GcpComputeInstanceIpv6AccessConfig` |  |  |  |
| `spec.networkInterfaces[].ipv6AccessConfigs[].networkTier` | `string` | yes |  |  |
| `spec.networkInterfaces[].ipv6AccessConfigs[].publicPtrDomainName` | `string` |  |  |  |
| `spec.networkInterfaces[].stackType` | `string` |  |  |  |
| `spec.networkInterfaces[].nicType` | `string` |  |  |  |
| `spec.networkInterfaces[].queueCount` | `int32` |  |  |  |
| `spec.networkInterfaces[].aliasIpRanges` | `[]GcpComputeInstanceAliasIpRange` |  |  |  |
| `spec.networkInterfaces[].aliasIpRanges[].ipCidrRange` | `string` | yes |  |  |
| `spec.networkInterfaces[].aliasIpRanges[].subnetworkRangeName` | `string` |  |  |  |
| `spec.serviceAccount` | `GcpComputeInstanceServiceAccount` |  |  |  |
| `spec.serviceAccount.email` | `string \| valueFrom` |  |  | GcpServiceAccount (`status.outputs.email`) |
| `spec.serviceAccount.scopes` | `[]string` | yes |  |  |
| `spec.scheduling` | `GcpComputeInstanceScheduling` |  |  |  |
| `spec.scheduling.provisioningModel` | `string` |  |  |  |
| `spec.scheduling.automaticRestart` | `bool` |  | `true` |  |
| `spec.scheduling.onHostMaintenance` | `string` |  |  |  |
| `spec.scheduling.instanceTerminationAction` | `string` |  |  |  |
| `spec.scheduling.maxRunDurationSeconds` | `int64` |  |  |  |
| `spec.scheduling.terminationTime` | `string` |  |  |  |
| `spec.scheduling.discardLocalSsdsOnStop` | `bool` |  |  |  |
| `spec.scheduling.availabilityDomain` | `int32` |  |  |  |
| `spec.scheduling.minNodeCpus` | `int32` |  |  |  |
| `spec.scheduling.nodeAffinities` | `[]GcpComputeInstanceNodeAffinity` |  |  |  |
| `spec.scheduling.nodeAffinities[].key` | `string` | yes |  |  |
| `spec.scheduling.nodeAffinities[].operator` | `string` | yes |  |  |
| `spec.scheduling.nodeAffinities[].values` | `[]string` | yes |  |  |
| `spec.scheduling.localSsdRecoveryTimeoutSeconds` | `int64` |  |  |  |
| `spec.shieldedInstanceConfig` | `GcpComputeInstanceShieldedConfig` |  |  |  |
| `spec.shieldedInstanceConfig.enableSecureBoot` | `bool` |  |  |  |
| `spec.shieldedInstanceConfig.enableVtpm` | `bool` |  | `true` |  |
| `spec.shieldedInstanceConfig.enableIntegrityMonitoring` | `bool` |  | `true` |  |
| `spec.confidentialInstanceConfig` | `GcpComputeInstanceConfidentialConfig` |  |  |  |
| `spec.confidentialInstanceConfig.confidentialInstanceType` | `string` | yes |  |  |
| `spec.advancedMachineFeatures` | `GcpComputeInstanceAdvancedMachineFeatures` |  |  |  |
| `spec.advancedMachineFeatures.enableNestedVirtualization` | `bool` |  |  |  |
| `spec.advancedMachineFeatures.threadsPerCore` | `int32` |  |  |  |
| `spec.advancedMachineFeatures.visibleCoreCount` | `int32` |  |  |  |
| `spec.advancedMachineFeatures.enableUefiNetworking` | `bool` |  |  |  |
| `spec.advancedMachineFeatures.performanceMonitoringUnit` | `string` |  |  |  |
| `spec.advancedMachineFeatures.turboMode` | `string` |  |  |  |
| `spec.guestAccelerators` | `[]GcpComputeInstanceGuestAccelerator` |  |  |  |
| `spec.guestAccelerators[].type` | `string` | yes |  |  |
| `spec.guestAccelerators[].count` | `int32` | yes |  |  |
| `spec.reservationAffinity` | `GcpComputeInstanceReservationAffinity` |  |  |  |
| `spec.reservationAffinity.type` | `string` | yes |  |  |
| `spec.reservationAffinity.specificReservation` | `GcpComputeInstanceSpecificReservation` |  |  |  |
| `spec.reservationAffinity.specificReservation.key` | `string` | yes |  |  |
| `spec.reservationAffinity.specificReservation.values` | `[]string` | yes |  |  |
| `spec.totalEgressBandwidthTier` | `string` |  |  |  |
| `spec.metadata` | `map<string, string>` |  |  |  |
| `spec.startupScript` | `string` |  |  |  |
| `spec.sshKeys` | `[]string` |  |  |  |
| `spec.labels` | `map<string, string>` |  |  |  |
| `spec.tags` | `[]string` |  |  |  |
| `spec.resourceManagerTags` | `map<string, string>` |  |  |  |
| `spec.resourcePolicies` | `[]string` |  |  |  |
| `spec.minCpuPlatform` | `string` |  |  |  |
| `spec.canIpForward` | `bool` |  |  |  |
| `spec.enableDisplay` | `bool` |  |  |  |
| `spec.deletionProtection` | `bool` |  |  |  |
| `spec.desiredStatus` | `string` |  |  |  |
| `spec.allowStoppingForUpdate` | `bool` |  | `true` |  |
| `spec.keyRevocationActionType` | `string` |  |  |  |

## Field Details

### spec.projectId

`string | valueFrom`

The GCP project that owns the instance.
Can be a literal project ID or a reference to a GcpProject resource.
If omitted, the provider's default project is used.
Immutable after creation.

- references: GcpProject (`status.outputs.project_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: GcpProject, name: <that resource's name>, fieldPath: status.outputs.project_id}} -- a bare string does not parse

### spec.instanceName

`string`

Name of the Compute Engine instance. 1-63 characters, lowercase
letters, numbers, and hyphens; must start with a letter and cannot end
with a hyphen. When omitted, metadata.name is used.
Immutable after creation.

- rule: instance_name must be 1-63 characters of lowercase letters, numbers, and hyphens, starting with a letter and not ending with a hyphen

### spec.zone

`string` · required

Zone where the instance runs, e.g. "us-central1-a".
Immutable after creation (moving zones replaces the VM).

- rule: {"required":true,"string":{"pattern":"^[a-z]+-[a-z]+[0-9]+-[a-z]$"}}

### spec.machineType

`string` · required

Machine type, e.g. "e2-medium", "n2-standard-4", "c3-highcpu-8", or a
custom shape like "custom-6-20480". Mutable — changing it stops and
restarts the VM (requires allow_stopping_for_update).

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.description

`string`

Human-readable description of the instance.

### spec.hostname

`string`

Custom fully-qualified DNS hostname, e.g. "db-1.prod.internal".
RFC-1035 labels separated by dots; when unset GCP derives
"<name>.c.<project>.internal". Immutable after creation.

- rule: hostname must be a fully qualified DNS name of RFC-1035 labels, e.g. host.example.internal

### spec.bootDisk

`GcpComputeInstanceBootDisk` · required

Boot disk configuration — the disk the OS boots from.

- rule: {"required":true}
- rule: exactly one boot source is required: image (fresh install), source_snapshot (restore), or source_disk (pre-created bootable disk)

### spec.bootDisk.image

`string`

Source image for a fresh boot disk. Accepts an image family
("debian-cloud/debian-12", "ubuntu-os-cloud/ubuntu-2404-lts-amd64") or
a specific image self link. Families resolve to the newest image at
create time. Create-time only.

### spec.bootDisk.sourceSnapshot

`string`

Source snapshot to restore the boot disk from (name or self link).
Create-time only.

### spec.bootDisk.sourceDisk

`string | valueFrom`

Existing bootable disk to boot from, referenced as a GcpComputeDisk.
The disk must live in the instance's zone. When booting from an
existing disk, size/type/encryption below are ignored — the disk
already owns them.

- references: GcpComputeDisk (`status.outputs.self_link`)
- rule: write as {value: <literal>} or {valueFrom: {kind: GcpComputeDisk, name: <that resource's name>, fieldPath: status.outputs.self_link}} -- a bare string does not parse

### spec.bootDisk.sizeGb

`int32`

Size of the boot disk in GB. When omitted, the image or snapshot size
is used. Grows in place; never shrinks.

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","int32":{"lte":65536,"gte":10}}

### spec.bootDisk.type

`string`

Disk type: "pd-standard" (HDD), "pd-balanced" (default in GCP for most
machine shapes and the sensible choice), "pd-ssd" (high IOPS), or a
hyperdisk type on supported machine families ("hyperdisk-balanced").

### spec.bootDisk.autoDelete

`bool` · optional (explicit presence)

Delete the boot disk automatically when the instance is deleted.
Defaults to true (matching GCP). Set false to keep the OS disk for
forensics or re-attachment.

- default: `true`

### spec.bootDisk.deviceName

`string`

Device name exposed under /dev/disk/by-id/google-*. When omitted GCP
assigns one.

### spec.bootDisk.kmsKey

`string | valueFrom`

Customer-managed encryption key (CMEK) for the boot disk, referenced
as a GcpKmsKey. The Compute Engine service agent
(service-<project-number>@compute-system.iam.gserviceaccount.com) must
hold roles/cloudkms.cryptoKeyEncrypterDecrypter on the key.
Create-time only.

- references: GcpKmsKey (`status.outputs.key_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: GcpKmsKey, name: <that resource's name>, fieldPath: status.outputs.key_id}} -- a bare string does not parse

### spec.bootDisk.diskLabels

`map<string, string>`

Labels applied to the boot disk itself (distinct from instance
labels) — useful for disk-level cost attribution and snapshot
policies.

### spec.bootDisk.provisionedIops

`int64` · optional (explicit presence)

Provisioned IOPS for hyperdisk types that support tuning
(e.g. hyperdisk-extreme, hyperdisk-balanced). Leave unset for pd-*
types.

- rule: {"int64":{"gt":"0"}}

### spec.bootDisk.provisionedThroughput

`int64` · optional (explicit presence)

Provisioned throughput in MB/s for hyperdisk types that support
tuning (e.g. hyperdisk-throughput, hyperdisk-balanced). Leave unset
for pd-* types.

- rule: {"int64":{"gt":"0"}}

### spec.bootDisk.architecture

`string`

CPU architecture of the disk/image: "X86_64" or "ARM64" (e.g. for
Tau T2A/Axion machine types). Normally inferred from the image.

- rule: architecture must be X86_64 or ARM64

### spec.bootDisk.enableConfidentialCompute

`bool`

Create the boot disk in confidential-compute mode (hyperdisk SKUs
only; requires kms_key).

### spec.bootDisk.resourcePolicies

`[]string`

Self links of resource policies to attach to the boot disk at create
time (e.g. a snapshot schedule). GCP currently allows at most one.
Changing this replaces the instance.

- rule: {"repeated":{"maxItems":"1"}}

### spec.bootDisk.storagePool

`string`

URL of the storage pool to create the boot disk in (hyperdisk storage
pools).

### spec.attachedDisks

`[]GcpComputeInstanceAttachedDisk`

Additional persistent data disks attached to the instance. Each disk
is a first-class GcpComputeDisk resource referenced by self link — the
disk has its own lifecycle and survives this VM unless its own
configuration says otherwise.

### spec.attachedDisks[].source

`string | valueFrom` · required

The disk to attach, referenced as a GcpComputeDisk (or a literal disk
name/self link). The disk must live in the instance's zone.

- references: GcpComputeDisk (`status.outputs.self_link`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: GcpComputeDisk, name: <that resource's name>, fieldPath: status.outputs.self_link}} -- a bare string does not parse

### spec.attachedDisks[].deviceName

`string`

Device name exposed under /dev/disk/by-id/google-*. When omitted GCP
assigns "persistent-disk-N".

### spec.attachedDisks[].mode

`string`

Attachment mode: "READ_WRITE" (default) or "READ_ONLY" (lets many VMs
attach the same disk simultaneously).

- rule: mode must be READ_WRITE or READ_ONLY

### spec.attachedDisks[].kmsKey

`string | valueFrom`

The CMEK key protecting the attached disk, referenced as a GcpKmsKey.
Required only when the disk is CMEK-encrypted — the attachment must
present the same key the disk was created with.

- references: GcpKmsKey (`status.outputs.key_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: GcpKmsKey, name: <that resource's name>, fieldPath: status.outputs.key_id}} -- a bare string does not parse

### spec.scratchDisks

`[]GcpComputeInstanceScratchDisk`

Ephemeral local-SSD scratch disks physically attached to the host.
Contents are lost when the VM stops or is preempted — use only for
caches, temp space, and high-IOPS scratch data. Create-time only.

### spec.scratchDisks[].interface

`string` · required

Disk interface: "NVME" (recommended; highest performance) or "SCSI".

- rule: {"required":true,"string":{"in":["NVME","SCSI"]}}

### spec.scratchDisks[].sizeGb

`int32`

Size in GB. Local SSDs come in fixed 375 GB units (or 3000 GB on
supported Z3 shapes). When omitted, 375 is used.

- rule: scratch disk size_gb must be 375 (standard local SSD unit) or 3000 (Z3 shapes)

### spec.scratchDisks[].deviceName

`string`

Device name exposed under /dev/disk/by-id/google-*.

### spec.networkInterfaces

`[]GcpComputeInstanceNetworkInterface` · required

Network interfaces. At least one is required; multiple NICs must each
attach to a different VPC network. NIC count and their target
networks are immutable after creation.

- rule: {"required":true,"repeated":{"minItems":"1"}}
- rule: each network interface needs a network (auto-mode VPC) or a subnetwork (custom-mode VPC) — set at least one

### spec.networkInterfaces[].network

`string | valueFrom`

VPC network for this interface, referenced as a GcpVpcNetwork.
Sufficient alone only for auto-mode VPCs; custom-mode VPCs need
subnetwork.

- references: GcpVpcNetwork (`status.outputs.network_self_link`)
- rule: write as {value: <literal>} or {valueFrom: {kind: GcpVpcNetwork, name: <that resource's name>, fieldPath: status.outputs.network_self_link}} -- a bare string does not parse

### spec.networkInterfaces[].subnetwork

`string | valueFrom`

Subnetwork for this interface, referenced as a GcpSubnetwork. The
subnetwork's region must contain the instance's zone.

- references: GcpSubnetwork (`status.outputs.subnetwork_self_link`)
- rule: write as {value: <literal>} or {valueFrom: {kind: GcpSubnetwork, name: <that resource's name>, fieldPath: status.outputs.subnetwork_self_link}} -- a bare string does not parse

### spec.networkInterfaces[].subnetworkProject

`string`

Project owning the subnetwork — set when attaching to a Shared VPC
host project's subnetwork from a service project.

### spec.networkInterfaces[].networkIp

`string | valueFrom`

Static internal IP for this interface. Accepts a literal IP or a
reference to a reserved INTERNAL GcpAddress. When omitted, GCP
assigns an ephemeral internal IP from the subnetwork range.

- references: GcpAddress (`status.outputs.address`)
- rule: write as {value: <literal>} or {valueFrom: {kind: GcpAddress, name: <that resource's name>, fieldPath: status.outputs.address}} -- a bare string does not parse

### spec.networkInterfaces[].accessConfigs

`[]GcpComputeInstanceAccessConfig`

External IPv4 access configs. Empty means no external IP (private
VM — pair with Cloud NAT for egress). GCP supports at most one
access config per interface.

- rule: {"repeated":{"maxItems":"1"}}

### spec.networkInterfaces[].accessConfigs[].natIp

`string | valueFrom`

Static external IP, as a literal or a reference to a reserved
EXTERNAL GcpAddress. When omitted, GCP assigns an ephemeral external
IP.

- references: GcpAddress (`status.outputs.address`)
- rule: write as {value: <literal>} or {valueFrom: {kind: GcpAddress, name: <that resource's name>, fieldPath: status.outputs.address}} -- a bare string does not parse

### spec.networkInterfaces[].accessConfigs[].networkTier

`string`

Network service tier for this IP: "PREMIUM" (default; Google's global
backbone) or "STANDARD" (regional, cheaper).

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"in":["PREMIUM","STANDARD"]}}

### spec.networkInterfaces[].accessConfigs[].publicPtrDomainName

`string`

Domain name for the public PTR (reverse DNS) record of this IP.

### spec.networkInterfaces[].ipv6AccessConfigs

`[]GcpComputeInstanceIpv6AccessConfig`

External IPv6 access configs. Requires stack_type "IPV4_IPV6" and a
subnetwork with an external IPv6 range. At most one per interface.

- rule: {"repeated":{"maxItems":"1"}}

### spec.networkInterfaces[].ipv6AccessConfigs[].networkTier

`string` · required

Network service tier for IPv6 traffic. Only "PREMIUM" is valid.

- rule: {"required":true,"string":{"in":["PREMIUM"]}}

### spec.networkInterfaces[].ipv6AccessConfigs[].publicPtrDomainName

`string`

Domain name for the public PTR (reverse DNS) record of the external
IPv6 range.

### spec.networkInterfaces[].stackType

`string`

IP stack of the interface:
  ""            -- same as "IPV4_ONLY" (GCP default)
  "IPV4_ONLY"   -- IPv4 addresses only
  "IPV4_IPV6"   -- dual stack (subnetwork must have an IPv6 range)
  "IPV6_ONLY"   -- IPv6 only (supported on IPv6-enabled subnetworks)

- rule: stack_type must be IPV4_ONLY, IPV4_IPV6, or IPV6_ONLY

### spec.networkInterfaces[].nicType

`string`

vNIC type: "" (GCP picks), "GVNIC" (recommended on modern machine
families; required for TIER_1 bandwidth), "VIRTIO_NET" (legacy), or
the RDMA types "IDPF", "MRDMA", "IRDMA" on specialized shapes.

- rule: nic_type must be one of GVNIC, VIRTIO_NET, IDPF, MRDMA, IRDMA

### spec.networkInterfaces[].queueCount

`int32` · optional (explicit presence)

Networking queue count for Rx and Tx (1-32). When omitted GCP sizes
queues from vCPU count.

- rule: {"int32":{"lte":32,"gte":1}}

### spec.networkInterfaces[].aliasIpRanges

`[]GcpComputeInstanceAliasIpRange`

Alias IP ranges served by this interface — the mechanism behind
per-pod/per-container IPs and multi-IP VMs.

### spec.networkInterfaces[].aliasIpRanges[].ipCidrRange

`string` · required

The alias range: a CIDR ("10.1.2.0/24"), a single IP ("10.1.2.3"), or
a netmask ("/24") to auto-allocate from the range.

- rule: {"required":true}

### spec.networkInterfaces[].aliasIpRanges[].subnetworkRangeName

`string`

Secondary range name on the subnetwork to allocate from. When omitted
the primary range is used.

### spec.serviceAccount

`GcpComputeInstanceServiceAccount`

Service account the VM's workloads authenticate as. When omitted, the
Compute Engine default service account is used with its default
scopes — prefer a dedicated least-privilege account for production.

### spec.serviceAccount.email

`string | valueFrom`

Service account email, referenced as a GcpServiceAccount or a literal
email. Changing it stops and restarts the VM (requires
allow_stopping_for_update).

- references: GcpServiceAccount (`status.outputs.email`)
- rule: write as {value: <literal>} or {valueFrom: {kind: GcpServiceAccount, name: <that resource's name>, fieldPath: status.outputs.email}} -- a bare string does not parse

### spec.serviceAccount.scopes

`[]string` · required

OAuth scopes for the attached account. The modern practice is a
single "https://www.googleapis.com/auth/cloud-platform" scope with
access controlled entirely by IAM roles; narrower legacy scopes
remain supported. Required when this block is set.

- rule: {"repeated":{"minItems":"1"}}

### spec.scheduling

`GcpComputeInstanceScheduling`

Scheduling policy: Spot vs standard provisioning, maintenance
behavior, run-duration limits, and sole-tenant node placement.

### spec.scheduling.provisioningModel

`string`

Provisioning model:
  ""         -- same as "STANDARD"
  "STANDARD" -- on-demand capacity
  "SPOT"     -- deeply discounted preemptible capacity; GCP may
                reclaim the VM at any time
Create-time only.

- rule: provisioning_model must be STANDARD or SPOT

### spec.scheduling.automaticRestart

`bool` · optional (explicit presence)

Restart the VM automatically when Compute Engine (not a user) stops
it. Defaults to true for standard VMs; must be false (or unset) for
Spot.

- default: `true`

### spec.scheduling.onHostMaintenance

`string`

Host maintenance behavior: "" (GCP default "MIGRATE"), "MIGRATE"
(live-migrate; zero downtime), or "TERMINATE" (stop during
maintenance — required for GPUs and confidential VMs).

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"in":["MIGRATE","TERMINATE"]}}

### spec.scheduling.instanceTerminationAction

`string`

What GCP does when a Spot VM is preempted or its run duration
expires: "STOP" (keep the stopped VM and disks) or "DELETE" (remove
the VM). Spot only.

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"in":["STOP","DELETE"]}}

### spec.scheduling.maxRunDurationSeconds

`int64` · optional (explicit presence)

Maximum run duration in seconds, after which
instance_termination_action is executed. The duration clock starts at
every VM start. Mutually exclusive with termination_time.

- rule: {"int64":{"lte":"315576000000","gte":"1"}}

### spec.scheduling.terminationTime

`string`

Absolute timestamp (RFC 3339) at which the VM is terminated. Mutually
exclusive with max_run_duration_seconds.

### spec.scheduling.discardLocalSsdsOnStop

`bool` · optional (explicit presence)

Discard local-SSD contents when the VM is stopped by a lifetime limit
(max_run_duration/termination_time) instead of preserving them.

### spec.scheduling.availabilityDomain

`int32` · optional (explicit presence)

Availability domain for spread-placement within the zone (used with
spread placement policies). 0 means unset.

- rule: {"int32":{"gte":1}}

### spec.scheduling.minNodeCpus

`int32` · optional (explicit presence)

Minimum vCPUs on the sole-tenant node this VM can be scheduled onto.
Sole-tenancy only.

- rule: {"int32":{"gte":1}}

### spec.scheduling.nodeAffinities

`[]GcpComputeInstanceNodeAffinity`

Sole-tenant node affinities selecting which node groups this VM may
run on. Setting any affinity places the VM on sole-tenant hardware.

### spec.scheduling.nodeAffinities[].key

`string` · required

Node-group label key to match (e.g.
"compute.googleapis.com/node-group-name").

- rule: {"required":true}

### spec.scheduling.nodeAffinities[].operator

`string` · required

Match operator: "IN" or "NOT_IN".

- rule: {"required":true,"string":{"in":["IN","NOT_IN"]}}

### spec.scheduling.nodeAffinities[].values

`[]string` · required

Label values to match.

- rule: {"repeated":{"minItems":"1"}}

### spec.scheduling.localSsdRecoveryTimeoutSeconds

`int64` · optional (explicit presence)

How long Compute Engine waits for a local-SSD-preserving recovery
when the host fails, in seconds, before falling back to default
recovery.

- rule: {"int64":{"lte":"604800","gte":"0"}}

### spec.shieldedInstanceConfig

`GcpComputeInstanceShieldedConfig`

Shielded VM configuration (secure boot, vTPM, integrity monitoring).
Requires an image with Shielded VM support (all recent Google-provided
images qualify).

### spec.shieldedInstanceConfig.enableSecureBoot

`bool` · optional (explicit presence)

Verify the boot loader's signature chain; blocks boot on tampering.
GCP default is false because some third-party images are unsigned.

### spec.shieldedInstanceConfig.enableVtpm

`bool` · optional (explicit presence)

Virtual Trusted Platform Module. GCP default is true.

- default: `true`

### spec.shieldedInstanceConfig.enableIntegrityMonitoring

`bool` · optional (explicit presence)

Boot-integrity measurement and monitoring via the vTPM. GCP default
is true.

- default: `true`

### spec.confidentialInstanceConfig

`GcpComputeInstanceConfidentialConfig`

Confidential VM configuration — hardware memory encryption (AMD SEV /
SEV-SNP or Intel TDX). Requires a supported machine family (e.g. N2D,
C2D, C3) and on_host_maintenance = "TERMINATE". Create-time only.

### spec.confidentialInstanceConfig.confidentialInstanceType

`string` · required

Confidential computing technology:
  "SEV"     -- AMD Secure Encrypted Virtualization (N2D/C2D/C3D)
  "SEV_SNP" -- AMD SEV Secure Nested Paging (requires an AMD Milan+
               min_cpu_platform)
  "TDX"     -- Intel Trust Domain Extensions (C3)

- rule: {"required":true,"string":{"in":["SEV","SEV_SNP","TDX"]}}

### spec.advancedMachineFeatures

`GcpComputeInstanceAdvancedMachineFeatures`

Advanced machine features: nested virtualization, SMT control, visible
core count, UEFI networking, performance monitoring unit, and turbo
mode.

### spec.advancedMachineFeatures.enableNestedVirtualization

`bool` · optional (explicit presence)

Expose nested virtualization support (VMX) to the guest — run VMs
inside this VM.

### spec.advancedMachineFeatures.threadsPerCore

`int32` · optional (explicit presence)

Threads per physical core: 1 disables simultaneous multithreading
(SMT) — common for licensing and security isolation; 2 is the
hardware default.

- rule: {"int32":{"in":[1,2]}}

### spec.advancedMachineFeatures.visibleCoreCount

`int32` · optional (explicit presence)

Number of physical cores exposed to the guest (core visibility for
per-core licensing). When unset all cores are visible.

- rule: {"int32":{"gte":1}}

### spec.advancedMachineFeatures.enableUefiNetworking

`bool` · optional (explicit presence)

Enable UEFI networking in the guest firmware. Create-time only.

### spec.advancedMachineFeatures.performanceMonitoringUnit

`string`

Performance monitoring unit exposure level: "STANDARD", "ENHANCED",
or "ARCHITECTURAL".

- rule: performance_monitoring_unit must be STANDARD, ENHANCED, or ARCHITECTURAL

### spec.advancedMachineFeatures.turboMode

`string`

Turbo frequency mode. "ALL_CORE_MAX" runs all cores at maximum turbo
frequency (supported machine families only).

- rule: turbo_mode must be ALL_CORE_MAX

### spec.guestAccelerators

`[]GcpComputeInstanceGuestAccelerator`

GPU accelerator cards attached to the instance. Requires a
GPU-capable zone and on_host_maintenance = "TERMINATE".

### spec.guestAccelerators[].type

`string` · required

Accelerator type available in the instance's zone, e.g.
"nvidia-tesla-t4", "nvidia-l4", "nvidia-a100-80gb".

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.guestAccelerators[].count

`int32` · required

Number of cards of this type.

- rule: {"required":true,"int32":{"gte":1}}

### spec.reservationAffinity

`GcpComputeInstanceReservationAffinity`

Reservation affinity — whether this VM consumes capacity from any
matching reservation, a specific reservation, or none.

- rule: specific_reservation must be set when (and only when) type is SPECIFIC_RESERVATION

### spec.reservationAffinity.type

`string` · required

Reservation consumption mode:
  "ANY_RESERVATION"      -- consume any matching reservation (GCP
                            default)
  "SPECIFIC_RESERVATION" -- consume only the named reservation
  "NO_RESERVATION"       -- never consume reserved capacity

- rule: {"required":true,"string":{"in":["ANY_RESERVATION","SPECIFIC_RESERVATION","NO_RESERVATION"]}}

### spec.reservationAffinity.specificReservation

`GcpComputeInstanceSpecificReservation`

The specific reservation to consume (type SPECIFIC_RESERVATION).

### spec.reservationAffinity.specificReservation.key

`string` · required

Reservation label key — use
"compute.googleapis.com/reservation-name" to target by name.

- rule: {"required":true}

### spec.reservationAffinity.specificReservation.values

`[]string` · required

Reservation label values (the reservation name when using the
reservation-name key).

- rule: {"repeated":{"minItems":"1"}}

### spec.totalEgressBandwidthTier

`string`

Per-VM egress bandwidth tier. "TIER_1" raises the bandwidth cap on
supported machine shapes (N2/N2D/C2/C3 with >= 30 vCPUs and gVNIC);
"DEFAULT" is the standard cap.

- rule: total_egress_bandwidth_tier must be DEFAULT or TIER_1

### spec.metadata

`map<string, string>`

Custom metadata key/value pairs made available to the guest OS via the
metadata server. Well-known keys configure agents and features (e.g.
"enable-oslogin", "startup-script-url").

### spec.startupScript

`string`

Startup script executed by the guest agent on every boot. Maps to the
metadata_startup_script surface, which keeps it distinct from user
metadata and re-runs it on each start.

### spec.sshKeys

`[]string`

SSH public keys in "username:ssh-rsa AAAA... user" format. Folded into
the instance metadata "ssh-keys" key (newline-joined) identically on
both engines. Ignored by VMs using OS Login.

### spec.labels

`map<string, string>`

User labels merged with Planton attribution labels (which win on key
conflicts). Keys and values must match GCP label constraints:
lowercase letters, numbers, hyphens, underscores; keys start with a
letter, max 63 characters.

### spec.tags

`[]string`

Network tags used by firewall rules and network routes to select this
instance.

### spec.resourceManagerTags

`map<string, string>`

Resource Manager tags bound to the instance for org-policy and IAM
conditions. Keys in the form "tagKeys/{id}", values "tagValues/{id}".
Create-time only.

### spec.resourcePolicies

`[]string`

Self links of compute resource policies to attach to the instance
(e.g. an instance schedule that starts/stops the VM on a calendar).
GCP currently allows at most one policy per instance.

- rule: {"repeated":{"maxItems":"1"}}

### spec.minCpuPlatform

`string`

Minimum CPU platform for the VM, e.g. "Intel Ice Lake" or
"AMD Milan". Constrains scheduling to hosts with at least this
platform.

### spec.canIpForward

`bool`

Allow sending/receiving packets with source or destination IPs that do
not match the instance's own addresses — required for VMs acting as
routers, NAT gateways, or VPN endpoints. Create-time only.

### spec.enableDisplay

`bool`

Enable the virtual display device (needed by some remote-desktop and
screen-capture tooling on headless VMs).

### spec.deletionProtection

`bool`

Protect the instance against accidental deletion. Deleting a protected
instance fails until this is set back to false. Defaults to false: a
VM is a compute node whose data levers are its disks — the boot disk's
auto_delete and each GcpComputeDisk's own lifecycle guard the data.

### spec.desiredStatus

`string`

Desired lifecycle status of the VM:
  ""           -- same as "RUNNING"
  "RUNNING"    -- started
  "SUSPENDED"  -- suspended to disk (fast resume; memory persisted)
  "TERMINATED" -- stopped (compute billing stops; disks keep billing)
Changing this starts/suspends/stops the VM in place.

- rule: desired_status must be RUNNING, SUSPENDED, or TERMINATED

### spec.allowStoppingForUpdate

`bool` · optional (explicit presence)

Allow the provider to stop and restart the VM when an update requires
it (machine type, service account, network interface changes, ...).
Without it those updates fail instead of causing downtime. Recommended
true for instances whose brief restart is acceptable.

- default: `true`

### spec.keyRevocationActionType

`string`

Action GCP takes on the VM when a Cloud KMS key protecting it is
revoked: "NONE" (default) or "STOP".

- rule: key_revocation_action_type must be NONE or STOP

## Validation Rules

- `spot_requires_no_automatic_restart`: Spot VMs cannot automatically restart after preemption — leave scheduling.automatic_restart unset or set it to false when provisioning_model is SPOT
- `termination_action_requires_spot`: scheduling.instance_termination_action applies only to Spot VMs — set scheduling.provisioning_model to SPOT or remove the termination action
- `confidential_requires_terminate_maintenance`: confidential VMs cannot live-migrate — set scheduling.on_host_maintenance to TERMINATE when confidential_instance_config is configured
- `accelerators_require_terminate_maintenance`: VMs with guest accelerators (GPUs) cannot live-migrate — set scheduling.on_host_maintenance to TERMINATE when guest_accelerators are attached
- `max_run_duration_conflicts_with_termination_time`: max_run_duration_seconds and termination_time both bound the VM's lifetime — set at most one

## Outputs

Reference an output from another manifest as `valueFrom: {kind: GcpComputeInstance, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.instance_name` | `string` | Name of the Compute Engine instance. |
| `status.outputs.instance_id` | `string` | Instance ID (unique numeric identifier). |
| `status.outputs.self_link` | `string` | Self link URL of the instance. |
| `status.outputs.internal_ip` | `string` | Internal (private) IP address of the instance. |
| `status.outputs.external_ip` | `string` | External (public) IP address of the instance (if configured). |
| `status.outputs.status` | `string` | Current status of the instance (RUNNING, STOPPED, etc.). |
| `status.outputs.zone` | `string` | Zone where the instance is located. |
| `status.outputs.machine_type` | `string` | Machine type of the instance. |
| `status.outputs.cpu_platform` | `string` | CPU platform of the instance. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.projectId` | GcpProject | `status.outputs.project_id` |
| `spec.bootDisk.sourceDisk` | GcpComputeDisk | `status.outputs.self_link` |
| `spec.bootDisk.kmsKey` | GcpKmsKey | `status.outputs.key_id` |
| `spec.attachedDisks[].source` | GcpComputeDisk | `status.outputs.self_link` |
| `spec.attachedDisks[].kmsKey` | GcpKmsKey | `status.outputs.key_id` |
| `spec.networkInterfaces[].network` | GcpVpcNetwork | `status.outputs.network_self_link` |
| `spec.networkInterfaces[].subnetwork` | GcpSubnetwork | `status.outputs.subnetwork_self_link` |
| `spec.networkInterfaces[].networkIp` | GcpAddress | `status.outputs.address` |
| `spec.networkInterfaces[].accessConfigs[].natIp` | GcpAddress | `status.outputs.address` |
| `spec.serviceAccount.email` | GcpServiceAccount | `status.outputs.email` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
