# HetznerCloudNetwork

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `hetzner-cloud.planton.dev/v1alpha1`

HetznerCloudNetworkSpec defines the specification for a Hetzner Cloud network.

A network provides private IPv4 connectivity between Hetzner Cloud resources.
It is defined by a top-level CIDR block (ip_range) that must be one of the
RFC 1918 private ranges (10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16). Resources
such as servers and load balancers attach to subnets within the network, not
directly to the network itself.

Subnets carve the network's ip_range into smaller blocks assigned to specific
network zones. At least one subnet is required because a Hetzner Cloud network
is unusable without subnets. Three subnet types exist:

  - cloud:   Standard subnet for cloud servers (most common).
  - server:  Subnet for connecting dedicated (Robot) servers via vSwitch.
  - vswitch: Subnet linked to a Hetzner Robot vSwitch for hybrid connectivity.

Routes define custom static routing within the network. They are optional --
default routing handles most use cases. Custom routes are useful for VPN
gateways, NAT instances, or inter-network routing.

The network is referenced by other components (HetznerCloudServer,
HetznerCloudLoadBalancer) via its network_id output through StringValueOrRef.

## Example

```yaml
apiVersion: hetzner-cloud.planton.dev/v1alpha1
kind: HetznerCloudNetwork
metadata:
  name: hetznercloudnetwork-demo
spec:
  ipRange: "10.0.0.0/16"
  subnets:
    - type: cloud
      networkZone: eu-central
      ipRange: "10.0.1.0/24"
    - type: server
      networkZone: eu-central
      ipRange: "10.0.2.0/24"
  routes:
    - destination: "172.16.0.0/12"
      gateway: "10.0.0.1"
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.ipRange` | `string` | yes |  |  |
| `spec.subnets` | `[]Subnet` | yes |  |  |
| `spec.subnets[].type` | `enum` | yes |  |  |
| `spec.subnets[].networkZone` | `string` | yes |  |  |
| `spec.subnets[].ipRange` | `string` | yes |  |  |
| `spec.subnets[].vswitchId` | `int64` |  |  |  |
| `spec.routes` | `[]Route` |  |  |  |
| `spec.routes[].destination` | `string` | yes |  |  |
| `spec.routes[].gateway` | `string` | yes |  |  |
| `spec.deleteProtection` | `bool` |  |  |  |
| `spec.exposeRoutesToVswitch` | `bool` |  |  |  |

## Field Details

### spec.ipRange

`string` · required

CIDR block for the network. Must be one of the private IPv4 ranges
defined in RFC 1918 (10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16).
All subnet ip_ranges must fall within this range.

Changing this value forces replacement of the network resource.

- rule: {"string":{"minLen":"1"}}

### spec.subnets

`[]Subnet` · required

Subnets within the network. At least one subnet is required because
a Hetzner Cloud network is unusable without subnets -- servers and
other resources attach to subnets, not directly to the network.

- rule: {"repeated":{"minItems":"1"}}
- rule: vswitch_id is required when subnet type is vswitch

### spec.subnets[].type

`enum` · required

Type of subnet. Determines how the subnet connects to infrastructure.

  - cloud:   Standard subnet for cloud servers (most common).
  - server:  Subnet for connecting dedicated (Robot) servers.
  - vswitch: Subnet linked to a Hetzner Robot vSwitch.

- rule: {"required":true,"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `subnet_type_unspecified`
- `cloud`
- `server`
- `vswitch`

### spec.subnets[].networkZone

`string` · required

Hetzner Cloud network zone for this subnet.

Known zones: "eu-central", "us-east", "us-west", "ap-southeast".
All subnets in a network can be in different zones, enabling
multi-region private connectivity.

- rule: {"string":{"minLen":"1"}}

### spec.subnets[].ipRange

`string` · required

IP range for this subnet in CIDR notation. Must be a subset of the
network's ip_range and must not overlap with other subnets or route
destinations in the same network.

Changing this value forces replacement of the subnet.

- rule: {"string":{"minLen":"1"}}

### spec.subnets[].vswitchId

`int64`

Hetzner Robot vSwitch ID. Required when type is vswitch, ignored
otherwise. Links this subnet to a dedicated server vSwitch for
hybrid cloud/dedicated connectivity.

### spec.routes

`[]Route`

Static routes for the network. Optional -- default routing handles
most use cases. Custom routes are needed for VPN gateways, NAT
instances, or routing between networks.

### spec.routes[].destination

`string` · required

Destination CIDR for this route. Must not overlap with any subnet
ip_range or other route destinations in the same network, and must
not be the first IP of the network's ip_range.

- rule: {"string":{"minLen":"1"}}

### spec.routes[].gateway

`string` · required

Gateway IP address within one of the network's subnets. Cannot be
the first IP of the network's ip_range or 172.31.1.1 (reserved for
the public network interface gateway).

- rule: {"string":{"minLen":"1"}}

### spec.deleteProtection

`bool`

Prevent accidental deletion of the network via the Hetzner Cloud API.
When enabled, the network cannot be deleted until protection is removed.

### spec.exposeRoutesToVswitch

`bool`

Expose the network's routes to vSwitch connections. Only takes effect
when a vSwitch connection is active (via a vswitch-type subnet).

## Outputs

Reference an output from another manifest as `valueFrom: {kind: HetznerCloudNetwork, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.network_id` | `string` | The Hetzner Cloud numeric ID of the created network (as a string). Referenced by HetznerCloudServer, HetznerCloudLoadBalancer, and other networking components via StringValueOrRef. |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| HetznerCloudLoadBalancer | `spec.network.networkId` | `status.outputs.network_id` |
| HetznerCloudServer | `spec.networks[].networkId` | `status.outputs.network_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
