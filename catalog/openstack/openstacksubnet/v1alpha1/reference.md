# OpenStackSubnet

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `openstack.planton.dev/v1alpha1`

OpenStackSubnetSpec defines the configuration for an OpenStack Neutron subnet.

A subnet provides IP address allocation for a network. It defines a CIDR block,
gateway, DNS servers, and DHCP settings. Every OpenStack workload that needs IP
connectivity requires at least one subnet attached to a network.

The subnet name is derived from metadata.name.

Terraform resource: openstack_networking_subnet_v2
Pulumi resource: openstack.networking.Subnet

## Example

```yaml
apiVersion: openstack.planton.dev/v1alpha1
kind: OpenStackSubnet
metadata:
  name: test-subnet
spec:
  network_id:
    value: "e0a1f622-9aab-4a48-8c8c-3b0c7e2a9b1d"
  cidr: "192.168.1.0/24"
  dns_nameservers:
    - "8.8.8.8"
    - "8.8.4.4"
  description: "Test subnet for local development"
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.networkId` | `string \| valueFrom` | yes |  | OpenStackNetwork (`status.outputs.network_id`) |
| `spec.cidr` | `string` | yes |  |  |
| `spec.ipVersion` | `int32` |  | `4` |  |
| `spec.gatewayIp` | `string` |  |  |  |
| `spec.noGateway` | `bool` |  |  |  |
| `spec.enableDhcp` | `bool` |  | `true` |  |
| `spec.dnsNameservers` | `[]string` |  |  |  |
| `spec.allocationPools` | `[]AllocationPool` |  |  |  |
| `spec.allocationPools[].start` | `string` | yes |  |  |
| `spec.allocationPools[].end` | `string` | yes |  |  |
| `spec.description` | `string` |  |  |  |
| `spec.tags` | `[]string` |  |  |  |
| `spec.region` | `string` |  |  |  |

## Field Details

### spec.networkId

`string | valueFrom` · required

network_id is the ID of the network to which this subnet belongs.
This is the defining relationship -- every subnet must belong to exactly one network.
Can reference an OpenStackNetwork resource's output or be a literal network UUID.

- references: OpenStackNetwork (`status.outputs.network_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: OpenStackNetwork, name: <that resource's name>, fieldPath: status.outputs.network_id}} -- a bare string does not parse

### spec.cidr

`string` · required

cidr is the IP address range for this subnet in CIDR notation.
Example IPv4: "192.168.1.0/24". Example IPv6: "2001:db8::/64".
The CIDR must be valid for the selected ip_version.

- rule: cidr must be valid CIDR notation (e.g., '192.168.1.0/24' or '2001:db8::/64')
- rule: {"required":true}

### spec.ipVersion

`int32` · optional (explicit presence)

ip_version is the IP protocol version for this subnet.
Must be 4 (IPv4) or 6 (IPv6).
Default: 4

- default: `4`
- rule: ip_version must be 4 (IPv4) or 6 (IPv6)

### spec.gatewayIp

`string`

gateway_ip is the IP address of the subnet gateway.
If omitted (and no_gateway is false), OpenStack automatically assigns the first
usable IP in the CIDR as the gateway.
Mutually exclusive with no_gateway.

### spec.noGateway

`bool`

no_gateway disables the gateway on this subnet when set to true.
Use this for isolated subnets that do not need routing (e.g., storage networks).
Mutually exclusive with gateway_ip.

### spec.enableDhcp

`bool` · optional (explicit presence)

enable_dhcp controls whether DHCP is enabled on this subnet.
When enabled, OpenStack's DHCP agent assigns IP addresses to ports on this subnet.
Default: true

- default: `true`

### spec.dnsNameservers

`[]string`

dns_nameservers is a list of DNS server IP addresses for this subnet.
These are pushed to instances via DHCP.
Example: ["8.8.8.8", "8.8.4.4"]

### spec.allocationPools

`[]AllocationPool`

allocation_pools defines the sub-ranges of the CIDR from which IPs are allocated.
If omitted, the entire CIDR (minus gateway and network/broadcast addresses) is used.
Each pool specifies a start and end IP address.

### spec.allocationPools[].start

`string` · required

start is the first IP address in the allocation range.
Example: "192.168.1.100"

- rule: {"required":true}

### spec.allocationPools[].end

`string` · required

end is the last IP address in the allocation range.
Example: "192.168.1.200"

- rule: {"required":true}

### spec.description

`string`

description is a human-readable description of the subnet.
This is stored on the OpenStack resource and visible in Horizon and API responses.

### spec.tags

`[]string`

tags are string tags to associate with the subnet in OpenStack.
Tags are stored on the OpenStack resource and can be used for filtering
and organization in the OpenStack API and Horizon dashboard.

- rule: {"repeated":{"unique":true}}

### spec.region

`string`

region overrides the region from the provider config for this subnet.
If omitted, the region from the OpenStack provider config is used.
Example: "RegionOne"

## Validation Rules

- `gateway.mutual_exclusion`: gateway_ip and no_gateway are mutually exclusive -- set one or neither, not both

## Outputs

Reference an output from another manifest as `valueFrom: {kind: OpenStackSubnet, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.subnet_id` | `string` | subnet_id is the unique identifier (UUID) of the subnet in OpenStack. This is the primary output used as a foreign key by downstream components. |
| `status.outputs.name` | `string` | name is the name of the subnet (derived from metadata.name). |
| `status.outputs.cidr` | `string` | cidr is the CIDR block of the subnet. |
| `status.outputs.gateway_ip` | `string` | gateway_ip is the gateway IP address of the subnet. This may be empty if no_gateway was set to true. |
| `status.outputs.network_id` | `string` | network_id is the ID of the parent network. Included for convenience -- downstream components can reference this without needing a separate lookup to the OpenStackNetwork resource. |
| `status.outputs.region` | `string` | region is the OpenStack region where the subnet was created. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.networkId` | OpenStackNetwork | `status.outputs.network_id` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| OpenStackLoadBalancer | `spec.vipSubnetId` | `status.outputs.subnet_id` |
| OpenStackLoadBalancerMember | `spec.subnetId` | `status.outputs.subnet_id` |
| OpenStackNetworkPort | `spec.fixedIps[].subnetId` | `status.outputs.subnet_id` |
| OpenStackRouterInterface | `spec.subnetId` | `status.outputs.subnet_id` |

## See Also

- [Overview](../README.md)
