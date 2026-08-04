# ScalewayPublicGateway

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `scaleway.planton.dev/v1`

ScalewayPublicGatewaySpec defines the specification for a Scaleway Public Gateway.

A Scaleway Public Gateway is a managed network appliance that sits at the edge
of a Private Network and provides:
  - NAT (masquerade): Allows resources in the Private Network to reach the
    internet through a single public IP.
  - SSH bastion: A secure access point for SSH connections to resources that
    have no public IP.
  - Port forwarding (PAT): Maps public ports on the gateway's IP to private
    IP:port pairs inside the attached Private Network.

This resource bundles three Scaleway resources into a single declarative unit:
  1. A dedicated Flexible IP (public IPv4 address).
  2. The Public Gateway appliance.
  3. A GatewayNetwork attachment (connecting the gateway to a Private Network).

Public Gateways are zonal resources (e.g., "fr-par-1"), unlike VPCs and Private
Networks which are regional. The zone must be within the same region as the
target Private Network.

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.privateNetworkId` | `string \| valueFrom` | yes |  | ScalewayPrivateNetwork (`status.outputs.private_network_id`) |
| `spec.zone` | `string` | yes |  |  |
| `spec.type` | `string` | yes | `VPC-GW-S` |  |
| `spec.enableMasquerade` | `bool` |  | `true` |  |
| `spec.bastion` | `ScalewayPublicGatewayBastion` |  |  |  |
| `spec.bastion.enabled` | `bool` |  |  |  |
| `spec.bastion.port` | `int32` |  | `22` |  |
| `spec.bastion.allowedIpRanges` | `[]string` |  |  |  |
| `spec.enableSmtp` | `bool` |  |  |  |
| `spec.reverseDns` | `string` |  |  |  |
| `spec.patRules` | `[]ScalewayPublicGatewayPatRule` |  |  |  |
| `spec.patRules[].privateIp` | `string` | yes |  |  |
| `spec.patRules[].privatePort` | `int32` | yes |  |  |
| `spec.patRules[].publicPort` | `int32` | yes |  |  |
| `spec.patRules[].protocol` | `string` |  | `both` |  |

## Field Details

### spec.privateNetworkId

`string | valueFrom` · required

The Private Network to attach this gateway to.

The gateway provides NAT, SSH bastion, and port forwarding for resources
in this network. Can be a literal Private Network UUID or a reference to
a ScalewayPrivateNetwork resource's output.

In infra charts, this is typically wired via valueFrom:

  privateNetworkId:
    valueFrom:
      kind: ScalewayPrivateNetwork
      name: app-network
      fieldPath: status.outputs.private_network_id

- references: ScalewayPrivateNetwork (`status.outputs.private_network_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: ScalewayPrivateNetwork, name: <that resource's name>, fieldPath: status.outputs.private_network_id}} -- a bare string does not parse

### spec.zone

`string` · required

The Scaleway zone where the gateway will be created.
Examples: "fr-par-1", "nl-ams-1", "pl-waw-1"

Public Gateways are zonal resources (unlike VPCs and Private Networks
which are regional). The zone must be within the same region as the
target Private Network. For example, if the Private Network is in
region "fr-par", the gateway zone must be "fr-par-1", "fr-par-2", etc.

This field is required and cannot be changed after creation.

- rule: {"required":true}

### spec.type

`string` · required

Gateway type determines bandwidth and performance tier.

Available types:
  - "VPC-GW-S"  -- Standard gateway. Sufficient for most workloads.
  - "VPC-GW-XL" -- High-bandwidth gateway (up to 10 Gbps). Available
                    only in Paris region zones (fr-par-1, fr-par-2, etc.).

Choose "VPC-GW-S" unless you have specific high-throughput requirements.

- default: `VPC-GW-S`
- rule: {"required":true}

### spec.enableMasquerade

`bool`

Whether to enable NAT masquerade on the Private Network attachment.

When enabled, resources in the attached Private Network can reach the
internet through the gateway's public IP. This is the primary use case
for a Public Gateway and is almost always enabled.

Disable only if using the gateway solely as an SSH bastion without NAT
functionality, or if outbound internet access is handled by other means.

Default: true

- default: `true`

### spec.bastion

`ScalewayPublicGatewayBastion`

SSH bastion configuration.

When enabled, the gateway acts as a jump host for SSH connections to
resources in the Private Network. Users SSH to the gateway's public IP,
and the gateway proxies the connection to the target private IP.

Optional. If omitted, the SSH bastion is disabled.

### spec.bastion.enabled

`bool`

Whether to enable SSH bastion.

When true, the gateway listens for SSH connections and proxies them
to resources in the attached Private Network.

### spec.bastion.port

`int32`

Port for the SSH bastion to listen on.

Default: 22. Change only if port 22 is blocked by corporate firewalls
or if you need a non-standard port for security-through-obscurity.

- default: `22`

### spec.bastion.allowedIpRanges

`[]string`

CIDR ranges allowed to connect to the bastion.
Examples: ["203.0.113.0/24", "198.51.100.10/32"]

Restricts SSH bastion access to specific source IP ranges. If empty,
all IP ranges are allowed -- not recommended for production.

Best practice: Restrict to your office IP range, VPN exit IPs, or
CI/CD runner IPs.

### spec.enableSmtp

`bool`

Whether to enable outbound SMTP through the gateway.

By default, outbound SMTP (port 25) is blocked to prevent spam abuse.
Enable only if resources in the Private Network need to send email
directly (e.g., an email relay or notification service).

Default: false

### spec.reverseDns

`string`

Reverse DNS hostname for the gateway's public IP.
Example: "gateway.example.com"

A matching DNS A record pointing to the public IP must already exist
before setting this field. Useful for email servers (SPF/DKIM compliance),
compliance requirements, and professional appearance in network logs.

Optional. If omitted, reverse DNS is not configured.

### spec.patRules

`[]ScalewayPublicGatewayPatRule`

Port forwarding (PAT) rules.

Each rule maps a public port on the gateway's IP to a private IP and
port inside the attached Private Network. This enables inbound access
to specific services without giving them public IPs.

Example: Forward public port 8080 to a web server at 10.0.1.5:80.

Optional. If omitted, no port forwarding rules are created.

### spec.patRules[].privateIp

`string` · required

Private IP address to forward traffic to.

Must be an IP address within the attached Private Network's subnet.
Example: "10.0.1.5"

- rule: {"required":true}

### spec.patRules[].privatePort

`int32` · required

Private port to forward traffic to.
Example: 80 (for a web server), 5432 (for PostgreSQL)

- rule: {"required":true}

### spec.patRules[].publicPort

`int32` · required

Public port to listen on the gateway's public IP.

External traffic hitting this port on the gateway's public IP will
be forwarded to private_ip:private_port.
Example: 8080

- rule: {"required":true}

### spec.patRules[].protocol

`string`

Protocol for the rule: "tcp", "udp", or "both".

Default: "both" (forwards both TCP and UDP traffic).
Use "tcp" for web servers, SSH, databases. Use "udp" for DNS, game
servers, media streaming. Use "both" when unsure.

- default: `both`

## Outputs

Reference an output from another manifest as `valueFrom: {kind: ScalewayPublicGateway, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.gateway_id` | `string` | The unique identifier of the created Public Gateway. Format: zoned ID (e.g., "fr-par-1/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"). This ID can be used to manage the gateway or reference it in PAT rules. |
| `status.outputs.public_ip_address` | `string` | The public IPv4 address assigned to the gateway. This is the IP that external traffic sees and that NAT masquerade uses as the source address for outbound connections. Useful for: - DNS A records pointing to your gateway - Firewall allowlists on external services - Connectivity diagnostics and monitoring |
| `status.outputs.public_ip_id` | `string` | The unique identifier of the Flexible IP resource. Format: zoned ID. The Flexible IP is the public IP resource that is attached to the gateway. This ID is useful if you need to manage the IP independently (e.g., reassign to a replacement gateway without changing the public IP). |
| `status.outputs.gateway_network_id` | `string` | The unique identifier of the gateway-to-network attachment. Format: zoned ID. Represents the GatewayNetwork resource that binds the gateway to the Private Network with NAT/masquerade settings. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.privateNetworkId` | ScalewayPrivateNetwork | `status.outputs.private_network_id` |

## See Also

- [Overview](./README.md)
