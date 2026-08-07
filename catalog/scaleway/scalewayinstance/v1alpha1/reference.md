# ScalewayInstance

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `scaleway.planton.dev/v1alpha1`

ScalewayInstanceSpec defines the specification for a Scaleway compute instance.

A Scaleway Instance is a virtual machine running on Scaleway's cloud platform.
This is a **composite resource** that bundles the following Scaleway resources
into a single declarative unit:

  1. The instance server (compute + root volume).
  2. An optional dedicated Flexible IP (public IPv4 address).
  3. Optional additional local volumes (l_ssd, scratch) created and attached.
  4. An optional Private Network attachment (inline NIC on the server).

Instances are **zonal** resources (e.g., "fr-par-1"). The zone must match
the zone of any Private Network attachment and security group.

**Composition pattern**: An instance typically sits at Layer 2 in an infra-chart
DAG: VPC -> PrivateNetwork -> Instance. It accepts optional `StringValueOrRef`
references to ScalewayPrivateNetwork and ScalewayInstanceSecurityGroup. Downstream
resources like ScalewayDnsRecord can reference `status.outputs.public_ip_address`,
and ScalewayLoadBalancer backends can reference `status.outputs.private_ip_address`.

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.zone` | `string` | yes |  |  |
| `spec.type` | `string` | yes |  |  |
| `spec.image` | `string` | yes |  |  |
| `spec.publicIp` | `ScalewayInstancePublicIp` |  |  |  |
| `spec.publicIp.reverseDns` | `string` |  |  |  |
| `spec.securityGroupId` | `string \| valueFrom` |  |  | ScalewayInstanceSecurityGroup (`status.outputs.security_group_id`) |
| `spec.privateNetworkId` | `string \| valueFrom` |  |  | ScalewayPrivateNetwork (`status.outputs.private_network_id`) |
| `spec.rootVolume` | `ScalewayInstanceRootVolume` |  |  |  |
| `spec.rootVolume.sizeInGb` | `int32` |  |  |  |
| `spec.rootVolume.volumeType` | `string` |  |  |  |
| `spec.rootVolume.deleteOnTermination` | `bool` |  |  |  |
| `spec.rootVolume.sbsIops` | `int32` |  |  |  |
| `spec.additionalVolumes` | `[]ScalewayInstanceAdditionalVolume` |  |  |  |
| `spec.additionalVolumes[].name` | `string` |  |  |  |
| `spec.additionalVolumes[].volumeType` | `string` | yes | `l_ssd` |  |
| `spec.additionalVolumes[].sizeInGb` | `int32` | yes |  |  |
| `spec.cloudInit` | `string` |  |  |  |
| `spec.state` | `string` |  | `started` |  |
| `spec.protected` | `bool` |  |  |  |

## Field Details

### spec.zone

`string` · required

The Scaleway zone where the instance will be created.
Examples: "fr-par-1", "nl-ams-1", "pl-waw-1"

Instances are zonal resources. The zone determines which physical
data center the instance runs in. Choose a zone in the same region
as the Private Network and other connected resources.

This field is required and cannot be changed after creation.

- rule: {"required":true}

### spec.type

`string` · required

Instance commercial type (required).

Determines CPU cores, RAM, and local storage allocation. Common types:
  - Development:  "DEV1-S" (2 vCPU, 2 GB), "DEV1-M" (3 vCPU, 4 GB)
  - General:      "GP1-S" (8 vCPU, 32 GB), "GP1-M" (16 vCPU, 64 GB)
  - Production:   "PRO2-S" (2 vCPU, 8 GB), "PRO2-M" (4 vCPU, 16 GB)

See Scaleway pricing page for the full list of available types and their
specifications. Type can be changed after creation (the instance will be
stopped, migrated, and restarted automatically).

- rule: {"required":true}

### spec.image

`string` · required

Base image UUID or label (required).

The image determines the operating system and initial disk contents.
Can be a full UUID (e.g., "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx") or
a human-friendly label (e.g., "ubuntu_focal", "debian_bullseye",
"ubuntu_jammy", "centos_stream_9").

Use labels for quick setup and UUIDs for reproducible deployments.
The available images depend on the instance zone.

- rule: {"required":true}

### spec.publicIp

`ScalewayInstancePublicIp`

Public IP configuration.

If set, a dedicated Flexible IPv4 address is created and attached to
the instance. The instance will be reachable from the internet on this
IP. If omitted, the instance has **no public IP** and is reachable only
via the Private Network (through a Public Gateway bastion, VPN, or
Load Balancer).

**When to set**: Development instances, bastion hosts, standalone servers.
**When to omit**: Production workloads behind a Load Balancer or Public
Gateway. Keeping instances off the public internet reduces attack surface.

### spec.publicIp.reverseDns

`string`

Reverse DNS hostname for the public IP.
Example: "web-01.example.com"

A matching DNS A record pointing to the public IP must already exist
before setting this field. Useful for email servers (SPF/DKIM compliance)
and professional appearance in network logs.

Optional. If omitted, reverse DNS defaults to Scaleway's generated
hostname.

### spec.securityGroupId

`string | valueFrom`

Security group to attach to the instance.

Controls inbound and outbound firewall rules for the instance. Can be
a literal security group UUID or a reference to a
ScalewayInstanceSecurityGroup resource's output.

If omitted, Scaleway assigns its default security group (which allows
all inbound and outbound traffic). For production, always specify an
explicit security group with a restrictive policy.

In infra charts, this is typically wired via valueFrom:

  securityGroupId:
    valueFrom:
      kind: ScalewayInstanceSecurityGroup
      name: web-sg
      fieldPath: status.outputs.security_group_id

- references: ScalewayInstanceSecurityGroup (`status.outputs.security_group_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: ScalewayInstanceSecurityGroup, name: <that resource's name>, fieldPath: status.outputs.security_group_id}} -- a bare string does not parse

### spec.privateNetworkId

`string | valueFrom`

Private Network to attach the instance to.

When set, the instance receives a private NIC on this network, enabling
communication with other resources on the same network (databases,
load balancers, other instances) using private IPs.

Can be a literal Private Network UUID or a reference to a
ScalewayPrivateNetwork resource's output.

In infra charts, this is typically wired via valueFrom:

  privateNetworkId:
    valueFrom:
      kind: ScalewayPrivateNetwork
      name: app-network
      fieldPath: status.outputs.private_network_id

Optional. Omit only for isolated instances that don't need internal
network connectivity (rare in production).

- references: ScalewayPrivateNetwork (`status.outputs.private_network_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: ScalewayPrivateNetwork, name: <that resource's name>, fieldPath: status.outputs.private_network_id}} -- a bare string does not parse

### spec.rootVolume

`ScalewayInstanceRootVolume`

Root volume configuration.

Controls the boot disk's size, type, and lifecycle behavior. If omitted,
the image's default root volume settings are used (typically 10-20 GB
of local SSD depending on the instance type).

Set this when you need a larger root disk, want to use SBS (network-
attached) storage for the root volume, or need specific IOPS guarantees.

### spec.rootVolume.sizeInGb

`int32`

Root volume size in GB.

If omitted, uses the image's default root volume size. The minimum
depends on the image (typically 10 GB). Can only be increased, not
decreased, after creation.

Set this when the default size is insufficient for your workload.

### spec.rootVolume.volumeType

`string`

Volume type for the root disk.

Options:
  - "l_ssd"      -- Local SSD. High performance, data stored on the
                     physical server. Cannot be resized. This is the
                     default for most instance types.
  - "sbs_volume"  -- SBS (Scaleway Block Storage). Network-attached,
                     resizable, with independent lifecycle. Choose this
                     for production workloads that need volume snapshots,
                     resizing, or persistence across instance replacement.

Changing volume_type after creation **recreates the instance** (data loss
unless you snapshot first). Choose carefully.

### spec.rootVolume.deleteOnTermination

`bool`

Delete the root volume when the instance is terminated.

Default: true. Set to false to preserve the root volume after instance
deletion (useful for forensics or data recovery). Only meaningful for
SBS volumes -- local SSD volumes are always destroyed with the instance.

### spec.rootVolume.sbsIops

`int32`

SBS IOPS allocation (only for sbs_volume type).

Specifies the guaranteed I/O operations per second for the root volume.
Only relevant when volume_type is "sbs_volume". If omitted, Scaleway
uses the default IOPS tier for the volume size.

Higher IOPS = higher cost. Set this only when the default IOPS is
insufficient for your workload (e.g., database servers).

### spec.additionalVolumes

`[]ScalewayInstanceAdditionalVolume`

Additional local volumes to create and attach to the instance.

These volumes are created as part of this composite resource and have
the same lifecycle as the instance -- they are destroyed when the
instance is terminated. Use these for:
  - High-performance local caching (l_ssd)
  - Temporary processing storage (scratch)
  - Data volumes that don't need to survive instance replacement

For persistent storage that survives instance termination, use
ScalewayBlockVolume (a separate resource kind) instead.

Optional. Most instances only need the root volume.

### spec.additionalVolumes[].name

`string`

Volume name.

A descriptive name for the volume. Used in the Scaleway console and
API for identification. If omitted, a name is auto-generated from
the resource name and volume index.

### spec.additionalVolumes[].volumeType

`string` · required

Volume type (required).

Options:
  - "l_ssd"   -- Local SSD. High performance, tied to the physical
                  server. Cannot be resized after creation. Destroyed
                  when the instance is terminated.
  - "scratch" -- Temporary local storage. Same performance as l_ssd
                  but explicitly ephemeral. Data is lost on instance
                  stop or restart. Use for temporary processing data.

- default: `l_ssd`
- rule: {"required":true}

### spec.additionalVolumes[].sizeInGb

`int32` · required

Volume size in GB (required).

The available size depends on the instance type's local storage
capacity. The total of all local volumes (root + additional) cannot
exceed the instance type's maximum local storage.

- rule: {"required":true}

### spec.cloudInit

`string`

Cloud-init script for instance bootstrapping.

A cloud-init script (typically starting with `#!/bin/bash` or
`#cloud-config`) that runs when the instance first boots. Use this
to install packages, configure services, set up users, or join
the instance to configuration management systems.

Example:
  cloudInit: |
    #!/bin/bash
    apt-get update && apt-get install -y docker.io
    systemctl enable docker

Maximum size: ~127 KB. For larger bootstrap scripts, use cloud-init
to download and execute a script from object storage.

Optional. If omitted, the instance boots with the image's defaults.

### spec.state

`string`

Desired instance state after creation.

Options:
  - "started" (default) -- Instance is running. Normal operational state.
  - "stopped"  -- Instance is shut down. Not billed for compute but
                  still billed for attached volumes and IPs.
  - "standby"  -- Instance is suspended to RAM. Faster resume than
                  "stopped" but still billed for RAM reservation.

Use "stopped" to pre-provision instances without starting them, or
to park instances during maintenance windows.

- default: `started`

### spec.protected

`bool`

Protect the instance against accidental deletion via the Scaleway API.

When true, the instance cannot be deleted through the API (including
Terraform/Pulumi destroy) without first disabling protection. This is
a safety net for critical production instances.

Default: false. Enable for production workloads.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: ScalewayInstance, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.server_id` | `string` | The unique identifier of the created instance server. Format: zoned ID (e.g., "fr-par-1/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"). This is the primary output for referencing this instance in Scaleway APIs, monitoring dashboards, and automation scripts. |
| `status.outputs.public_ip_address` | `string` | The public IPv4 address assigned to the instance. Populated only when `spec.public_ip` is set. Empty string if the instance has no public IP. Use this output for: - DNS A records (via ScalewayDnsRecord) - SSH access configuration - External service whitelisting - Monitoring endpoint registration |
| `status.outputs.public_ip_id` | `string` | The unique identifier of the Flexible IP resource. Populated only when `spec.public_ip` is set. Empty string otherwise. Format: zoned ID. The Flexible IP has independent lifecycle -- it survives instance replacement, preserving DNS records and firewall rules. |
| `status.outputs.private_ip_address` | `string` | The private IP address on the attached Private Network. Populated only when `spec.private_network_id` is set. Empty string if no Private Network is configured. Use this output for: - Load Balancer backend server IPs - Internal service discovery - Database connection allowlists - Inter-instance communication |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.securityGroupId` | ScalewayInstanceSecurityGroup | `status.outputs.security_group_id` |
| `spec.privateNetworkId` | ScalewayPrivateNetwork | `status.outputs.private_network_id` |

## See Also

- [Overview](./README.md)
