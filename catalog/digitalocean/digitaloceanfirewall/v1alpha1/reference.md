# DigitalOceanFirewall

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `digital-ocean.planton.dev/v1alpha1`

DigitalOceanFirewallSpec defines the user configuration for a DigitalOcean Cloud Firewall.

## Example

```yaml
apiVersion: digital-ocean.planton.dev/v1alpha1
kind: DigitalOceanFirewall
metadata:
  name: first-firewall                 # Kubernetes object name
spec:
  name: first-firewall                 # DigitalOcean firewall name
  inboundRules:
    - protocol: tcp
      portRange: "80"
      sourceAddresses:
        - "0.0.0.0/0"                # allow HTTP from anywhere
    - protocol: tcp
      portRange: "443"
      sourceAddresses:
        - "0.0.0.0/0"                # allow HTTPS from anywhere
  outboundRules:
    - protocol: tcp
      portRange: "1-65535"
      destinationAddresses:
        - "0.0.0.0/0"                # allow all outbound traffic
  dropletIds: []
  tags:
    - planton
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.name` | `string` | yes |  |  |
| `spec.inboundRules` | `[]DigitalOceanFirewallInboundRule` |  |  |  |
| `spec.inboundRules[].protocol` | `string` | yes |  |  |
| `spec.inboundRules[].portRange` | `string` |  |  |  |
| `spec.inboundRules[].sourceAddresses` | `[]string` |  |  |  |
| `spec.inboundRules[].sourceDropletIds` | `[]int64` |  |  |  |
| `spec.inboundRules[].sourceTags` | `[]string` |  |  |  |
| `spec.inboundRules[].sourceKubernetesIds` | `[]string` |  |  |  |
| `spec.inboundRules[].sourceLoadBalancerUids` | `[]string` |  |  |  |
| `spec.outboundRules` | `[]DigitalOceanFirewallOutboundRule` |  |  |  |
| `spec.outboundRules[].protocol` | `string` | yes |  |  |
| `spec.outboundRules[].portRange` | `string` |  |  |  |
| `spec.outboundRules[].destinationAddresses` | `[]string` |  |  |  |
| `spec.outboundRules[].destinationDropletIds` | `[]int64` |  |  |  |
| `spec.outboundRules[].destinationTags` | `[]string` |  |  |  |
| `spec.outboundRules[].destinationKubernetesIds` | `[]string` |  |  |  |
| `spec.outboundRules[].destinationLoadBalancerUids` | `[]string` |  |  |  |
| `spec.dropletIds` | `[]int64` |  |  |  |
| `spec.tags` | `[]string` |  |  |  |

## Field Details

### spec.name

`string` · required

Name of the firewall for identification (must be unique per account/project).

- rule: {"string":{"minLen":"1","maxLen":"255"}}

### spec.inboundRules

`[]DigitalOceanFirewallInboundRule`

Inbound rules: traffic allowed *to* Droplets on specific ports from specified sources.

### spec.inboundRules[].protocol

`string` · required

"tcp", "udp", or "icmp". Required.

- rule: {"string":{"minLen":"1"}}

### spec.inboundRules[].portRange

`string`

Ports to allow (e.g., "80", "8000-9000", or "1-65535"; empty or "1-65535" means all ports for tcp/udp).

### spec.inboundRules[].sourceAddresses

`[]string`

IPv4 or IPv6 addresses or CIDR ranges (e.g., "192.0.2.0/24", "0.0.0.0/0").

### spec.inboundRules[].sourceDropletIds

`[]int64`

IDs of Droplets from which traffic is allowed.

### spec.inboundRules[].sourceTags

`[]string`

Names of Droplet tags; any Droplet with these tags is allowed.

### spec.inboundRules[].sourceKubernetesIds

`[]string`

IDs of Kubernetes clusters from which traffic is allowed.

### spec.inboundRules[].sourceLoadBalancerUids

`[]string`

IDs of Load Balancers from which traffic is allowed.

### spec.outboundRules

`[]DigitalOceanFirewallOutboundRule`

Outbound rules: traffic allowed *from* Droplets on specific ports to specified destinations.

### spec.outboundRules[].protocol

`string` · required

"tcp", "udp", or "icmp". Required.

- rule: {"string":{"minLen":"1"}}

### spec.outboundRules[].portRange

`string`

Ports to allow (format as in inbound rules; required for tcp/udp).

### spec.outboundRules[].destinationAddresses

`[]string`

IPv4/IPv6 addresses or CIDRs to which traffic is allowed.

### spec.outboundRules[].destinationDropletIds

`[]int64`

IDs of Droplets to which traffic is allowed.

### spec.outboundRules[].destinationTags

`[]string`

Names of Droplet tags whose members are allowed destinations.

### spec.outboundRules[].destinationKubernetesIds

`[]string`

IDs of Kubernetes clusters to which traffic is allowed.

### spec.outboundRules[].destinationLoadBalancerUids

`[]string`

IDs of Load Balancers which are allowed as destinations.

### spec.dropletIds

`[]int64`

The Droplet IDs to which this firewall is applied (max 10).
These Droplets will have the firewall's rules enforced.

### spec.tags

`[]string`

The names of Droplet tags to which this firewall is applied (max 5).
Any Droplet with these tags will be protected by this firewall.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: DigitalOceanFirewall, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.firewall_id` | `string` |  |

## See Also

- [Overview](../README.md)
