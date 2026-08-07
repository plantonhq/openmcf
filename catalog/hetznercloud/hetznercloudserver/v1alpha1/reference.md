# HetznerCloudServer

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `hetzner-cloud.planton.dev/v1alpha1`

HetznerCloudServerSpec defines the specification for a Hetzner Cloud server.

A server is the core compute resource in Hetzner Cloud. It runs a chosen OS
image on a specific server type (vCPU/RAM/disk combination) in a given
location. Servers can be attached to private networks, assigned firewall
rules, placed into anti-affinity placement groups, and configured with SSH
keys for secure access -- all via references to other Planton components.

Public networking is enabled by default (auto-assigned IPv4 + IPv6). For
stable public IPs that survive server replacement, assign Primary IPs via
the public_net block.

Bundled provider resources:
  - hcloud_server:  The server itself.
  - hcloud_rdns:    Optional reverse DNS record for the server's auto-assigned
                    IPv4 address. Only created when dns_ptr is set. If you use
                    Primary IPs via public_net, manage rDNS on the
                    HetznerCloudPrimaryIp component instead.

Fields not exposed in this spec (hardcoded or derived in IaC modules):
  - name:   Derived from metadata.name.
  - labels: Derived from metadata (CG01 pattern). Standard labels take
            precedence over user-specified metadata.labels.
  - iso:    Mounting ISOs is an operational action, not declarative IaC.
  - rescue: Rescue mode is a runtime emergency action.
  - allow_deprecated_images: Safety valve best left to provider-level errors.
  - ignore_remote_firewall_ids: Terraform-specific drift suppression.
  - datacenter: Deprecated by Hetzner (removal after 2026-07-01).

## Example

```yaml
apiVersion: hetzner-cloud.planton.dev/v1alpha1
kind: HetznerCloudServer
metadata:
  name: hetznercloudserver-demo
  org: demo-org
  env: dev
  labels:
    team: platform
spec:
  serverType: cx22
  image: ubuntu-24.04
  location: fsn1
  sshKeys:
    - value: "my-ssh-key"
  userData: |
    #!/bin/bash
    apt-get update && apt-get install -y nginx
  firewallIds:
    - value: "12345"
  networks:
    - networkId:
        value: "67890"
      ip: "10.0.1.5"
  backups: true
  dnsPtr: demo.example.com
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.serverType` | `string` | yes |  |  |
| `spec.image` | `string` | yes |  |  |
| `spec.location` | `string` | yes |  |  |
| `spec.sshKeys` | `[]string \| valueFrom` |  |  | HetznerCloudSshKey (`status.outputs.ssh_key_id`) |
| `spec.userData` | `string` |  |  |  |
| `spec.placementGroupId` | `string \| valueFrom` |  |  | HetznerCloudPlacementGroup (`status.outputs.placement_group_id`) |
| `spec.firewallIds` | `[]string \| valueFrom` |  |  | HetznerCloudFirewall (`status.outputs.firewall_id`) |
| `spec.publicNet` | `PublicNet` |  |  |  |
| `spec.publicNet.ipv4Enabled` | `bool` |  | `true` |  |
| `spec.publicNet.ipv6Enabled` | `bool` |  | `true` |  |
| `spec.publicNet.ipv4` | `string \| valueFrom` |  |  | HetznerCloudPrimaryIp (`status.outputs.primary_ip_id`) |
| `spec.publicNet.ipv6` | `string \| valueFrom` |  |  | HetznerCloudPrimaryIp (`status.outputs.primary_ip_id`) |
| `spec.networks` | `[]NetworkAttachment` |  |  |  |
| `spec.networks[].networkId` | `string \| valueFrom` | yes |  | HetznerCloudNetwork (`status.outputs.network_id`) |
| `spec.networks[].ip` | `string` |  |  |  |
| `spec.networks[].aliasIps` | `[]string` |  |  |  |
| `spec.backups` | `bool` |  |  |  |
| `spec.keepDisk` | `bool` |  |  |  |
| `spec.deleteProtection` | `bool` |  |  |  |
| `spec.rebuildProtection` | `bool` |  |  |  |
| `spec.shutdownBeforeDeletion` | `bool` |  |  |  |
| `spec.dnsPtr` | `string` |  |  |  |

## Field Details

### spec.serverType

`string` · required

Server type that determines vCPU, RAM, and disk resources.

Examples: "cx22" (2 vCPU / 4 GB), "cpx11" (2 vCPU / 2 GB, AMD),
"cax11" (2 vCPU / 4 GB, ARM64).

Changing this value triggers a server resize. The server is temporarily
stopped during the resize. Use keep_disk to prevent irreversible disk
upgrades.

- rule: {"string":{"minLen":"1"}}

### spec.image

`string` · required

OS image name or ID to provision the server with.

Examples: "ubuntu-24.04", "debian-12", "rocky-9", "45346857".

Changing this value forces replacement of the server.

- rule: {"string":{"minLen":"1"}}

### spec.location

`string` · required

Hetzner Cloud location for the server (e.g., "fsn1", "nbg1", "hel1",
"ash", "hil", "sin"). Determines the physical datacenter.

Primary IPs and Floating IPs assigned to the server must be in the
same location.

Changing this value forces replacement of the server.

- rule: {"string":{"minLen":"1"}}

### spec.sshKeys

`[]string | valueFrom`

SSH keys to inject into the server at creation time.

Each entry accepts a literal SSH key name or numeric ID (as a string),
or a reference to a HetznerCloudSshKey resource's output via valueFrom.

SSH keys are only injected during initial server creation. Changing this
list after creation forces replacement of the server.

Example (literal name):
  sshKeys:
    - value: "my-ssh-key"

Example (reference):
  sshKeys:
    - valueFrom:
        kind: HetznerCloudSshKey
        name: prod-key
        fieldPath: status.outputs.ssh_key_id

- references: HetznerCloudSshKey (`status.outputs.ssh_key_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: HetznerCloudSshKey, name: <that resource's name>, fieldPath: status.outputs.ssh_key_id}} -- a bare string does not parse

### spec.userData

`string`

Cloud-init user data script or configuration to execute on first boot.

Accepts raw shell scripts (starting with #!/bin/bash) or cloud-config
YAML (starting with #cloud-config). Maximum size is 32 KB.

Changing this value forces replacement of the server.

### spec.placementGroupId

`string | valueFrom`

Placement group to assign the server to for anti-affinity scheduling.

Accepts a literal Hetzner Cloud placement group ID (as a string) or a
reference to a HetznerCloudPlacementGroup resource's output via
valueFrom. Servers in a "spread" placement group are guaranteed to run
on different physical hosts.

Example (reference):
  placementGroupId:
    valueFrom:
      kind: HetznerCloudPlacementGroup
      name: ha-group
      fieldPath: status.outputs.placement_group_id

- references: HetznerCloudPlacementGroup (`status.outputs.placement_group_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: HetznerCloudPlacementGroup, name: <that resource's name>, fieldPath: status.outputs.placement_group_id}} -- a bare string does not parse

### spec.firewallIds

`[]string | valueFrom`

Firewalls to apply to the server at creation time.

Each entry accepts a literal Hetzner Cloud firewall ID (as a string) or
a reference to a HetznerCloudFirewall resource's output via valueFrom.

Example (reference):
  firewallIds:
    - valueFrom:
        kind: HetznerCloudFirewall
        name: web-firewall
        fieldPath: status.outputs.firewall_id

- references: HetznerCloudFirewall (`status.outputs.firewall_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: HetznerCloudFirewall, name: <that resource's name>, fieldPath: status.outputs.firewall_id}} -- a bare string does not parse

### spec.publicNet

`PublicNet`

Public network configuration for the server.

If omitted, the server receives auto-assigned public IPv4 and IPv6
addresses (provider default). Set this block to disable public
networking or to attach existing Primary IPs.

### spec.publicNet.ipv4Enabled

`bool` · optional (explicit presence)

Enable public IPv4 for the server.

Default: true

- default: `true`

### spec.publicNet.ipv6Enabled

`bool` · optional (explicit presence)

Enable public IPv6 for the server.

Default: true

- default: `true`

### spec.publicNet.ipv4

`string | valueFrom`

Existing Primary IP (IPv4) to attach to the server instead of
auto-assigning. Accepts a literal Hetzner Cloud Primary IP ID (as a
string) or a reference to a HetznerCloudPrimaryIp output via
valueFrom.

The Primary IP must be IPv4 type and in the same location as the
server.

Example (reference):
  ipv4:
    valueFrom:
      kind: HetznerCloudPrimaryIp
      name: web-ipv4
      fieldPath: status.outputs.primary_ip_id

- references: HetznerCloudPrimaryIp (`status.outputs.primary_ip_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: HetznerCloudPrimaryIp, name: <that resource's name>, fieldPath: status.outputs.primary_ip_id}} -- a bare string does not parse

### spec.publicNet.ipv6

`string | valueFrom`

Existing Primary IP (IPv6) to attach to the server instead of
auto-assigning. Same semantics as ipv4 above, but for IPv6.

- references: HetznerCloudPrimaryIp (`status.outputs.primary_ip_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: HetznerCloudPrimaryIp, name: <that resource's name>, fieldPath: status.outputs.primary_ip_id}} -- a bare string does not parse

### spec.networks

`[]NetworkAttachment`

Private networks to attach the server to.

Each entry creates an inline network attachment on the server resource.
The server receives an IP from the network's subnet range (either
auto-assigned or explicitly specified via the ip field).

### spec.networks[].networkId

`string | valueFrom` · required

Network to attach the server to.

Accepts a literal Hetzner Cloud network ID (as a string) or a
reference to a HetznerCloudNetwork resource's output via valueFrom.

Example (reference):
  networkId:
    valueFrom:
      kind: HetznerCloudNetwork
      name: main-vpc
      fieldPath: status.outputs.network_id

- references: HetznerCloudNetwork (`status.outputs.network_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: HetznerCloudNetwork, name: <that resource's name>, fieldPath: status.outputs.network_id}} -- a bare string does not parse

### spec.networks[].ip

`string`

Specific IP address to assign to the server within the network's
subnet range. If omitted, Hetzner Cloud auto-assigns an IP.

### spec.networks[].aliasIps

`[]string`

Additional IP addresses for the server within this network.

Alias IPs allow a single server to listen on multiple private IPs
within the same network, useful for hosting multiple services or
virtual hosts.

### spec.backups

`bool`

Enable automatic daily backups for the server.

Backups are stored for 14 days and can be used to restore the server.
Incurs an additional 20% of the server price.

### spec.keepDisk

`bool`

Preserve the existing disk size when changing server_type.

When false (default), a server_type upgrade also upgrades the disk,
which is irreversible -- you cannot later downgrade to a smaller
server_type. Set to true to only change vCPU and RAM, keeping the
disk at its current size.

### spec.deleteProtection

`bool`

Prevent accidental deletion of the server via the Hetzner Cloud API.
When enabled, the server cannot be deleted until protection is removed.

### spec.rebuildProtection

`bool`

Prevent accidental rebuild (re-image) of the server via the Hetzner
Cloud API. When enabled, the server cannot be rebuilt until protection
is removed.

### spec.shutdownBeforeDeletion

`bool`

Attempt a graceful shutdown of the server before Terraform destroys it.

When true, Terraform sends an ACPI shutdown signal and waits for the
server to power off before deletion. When false (default), the server
is deleted immediately.

### spec.dnsPtr

`string`

Reverse DNS pointer record for the server's auto-assigned public IPv4
address. If set, an hcloud_rdns resource is created mapping the
server's IPv4 to this hostname.

Only use this when the server has auto-assigned public IPv4 (default
behavior or public_net with ipv4_enabled=true and no ipv4 Primary IP
reference). If you assign Primary IPs via public_net.ipv4, manage
rDNS on the HetznerCloudPrimaryIp component instead to avoid
conflicting management of the same IP's rDNS record.

Example: "web.example.com"

## Outputs

Reference an output from another manifest as `valueFrom: {kind: HetznerCloudServer, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.server_id` | `string` | The Hetzner Cloud numeric ID of the created server (as a string). Referenced by HetznerCloudVolume, HetznerCloudSnapshot, HetznerCloudFloatingIp (assignment), and HetznerCloudLoadBalancer (targets) via StringValueOrRef. |
| `status.outputs.ipv4_address` | `string` | The public IPv4 address assigned to the server. Empty if public IPv4 is disabled via public_net.ipv4_enabled = false. |
| `status.outputs.ipv6_address` | `string` | The first IPv6 address of the server's assigned /64 network. Empty if public IPv6 is disabled via public_net.ipv6_enabled = false. |
| `status.outputs.status` | `string` | The current status of the server (e.g., "running", "off", "rebuilding", "migrating"). |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.sshKeys` | HetznerCloudSshKey | `status.outputs.ssh_key_id` |
| `spec.placementGroupId` | HetznerCloudPlacementGroup | `status.outputs.placement_group_id` |
| `spec.firewallIds` | HetznerCloudFirewall | `status.outputs.firewall_id` |
| `spec.publicNet.ipv4` | HetznerCloudPrimaryIp | `status.outputs.primary_ip_id` |
| `spec.publicNet.ipv6` | HetznerCloudPrimaryIp | `status.outputs.primary_ip_id` |
| `spec.networks[].networkId` | HetznerCloudNetwork | `status.outputs.network_id` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| HetznerCloudFloatingIp | `spec.serverId` | `status.outputs.server_id` |
| HetznerCloudLoadBalancer | `spec.serverTargets[].serverId` | `status.outputs.server_id` |
| HetznerCloudSnapshot | `spec.serverId` | `status.outputs.server_id` |
| HetznerCloudVolume | `spec.serverId` | `status.outputs.server_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
