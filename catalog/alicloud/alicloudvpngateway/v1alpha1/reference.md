# AliCloudVpnGateway

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `alicloud.planton.dev/v1alpha1`

AliCloudVpnGatewaySpec defines the configuration for an Alibaba Cloud VPN
Gateway with bundled customer gateways and IPsec VPN connections.

A VPN Gateway provides encrypted site-to-site connectivity between an
Alibaba Cloud VPC and remote networks (on-premises data centers, branch
offices, or other cloud environments) over the public internet using IPsec.

This component bundles the VPN Gateway, customer gateways, and VPN
connections into a single deployable unit (per DD07 composite bundling)
because a VPN Gateway without at least one connection is non-functional.

The bundled flow:
  1. Create the VPN Gateway in the specified VPC/VSwitch.
  2. For each connection entry, create a Customer Gateway from the
     remote device's public IP.
  3. For each connection entry, create a VPN Connection linking the
     VPN Gateway to the Customer Gateway with IKE/IPsec negotiation
     parameters.

The gateway optionally supports SSL VPN for remote client access when
enable_ssl is true.

Provider resources:
  Terraform: alicloud_vpn_gateway + alicloud_vpn_customer_gateway + alicloud_vpn_connection
  Pulumi:    vpn.Gateway + vpn.CustomerGateway + vpn.Connection

## Example

```yaml
apiVersion: alicloud.planton.dev/v1alpha1
kind: AliCloudVpnGateway
metadata:
  name: alicloudvpngateway-demo
spec:
  region: cn-hangzhou
  vpcId:
    value: vpc-demo123
  vswitchId:
    value: vsw-demo123
  vpnGatewayName: demo-vpn
  description: Demo VPN Gateway for local testing
  bandwidth: 10
  connections:
    - name: demo-site
      customerGatewayIp: "203.0.113.1"
      localSubnets:
        - "10.0.0.0/8"
      remoteSubnets:
        - "192.168.0.0/16"
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.vpcId` | `string \| valueFrom` | yes |  | AliCloudVpc (`status.outputs.vpc_id`) |
| `spec.vswitchId` | `string \| valueFrom` | yes |  | AliCloudVswitch (`status.outputs.vswitch_id`) |
| `spec.vpnGatewayName` | `string` | yes |  |  |
| `spec.description` | `string` |  |  |  |
| `spec.bandwidth` | `int32` | yes |  |  |
| `spec.paymentType` | `string` |  | `PayAsYouGo` |  |
| `spec.enableSsl` | `bool` |  | `false` |  |
| `spec.sslConnections` | `int32` |  |  |  |
| `spec.tags` | `map<string, string>` |  |  |  |
| `spec.resourceGroupId` | `string` |  |  |  |
| `spec.connections` | `[]AliCloudVpnConnection` |  |  |  |
| `spec.connections[].name` | `string` | yes |  |  |
| `spec.connections[].customerGatewayIp` | `string` | yes |  |  |
| `spec.connections[].customerGatewayAsn` | `string` |  |  |  |
| `spec.connections[].localSubnets` | `[]string` | yes |  |  |
| `spec.connections[].remoteSubnets` | `[]string` | yes |  |  |
| `spec.connections[].enableDpd` | `bool` |  | `true` |  |
| `spec.connections[].enableNatTraversal` | `bool` |  | `true` |  |
| `spec.connections[].effectImmediately` | `bool` |  | `true` |  |
| `spec.connections[].ikeConfig` | `AliCloudIkeConfig` |  |  |  |
| `spec.connections[].ikeConfig.psk` | `string` (sensitive) |  |  |  |
| `spec.connections[].ikeConfig.ikeVersion` | `string` |  | `ikev2` |  |
| `spec.connections[].ikeConfig.ikeMode` | `string` |  | `main` |  |
| `spec.connections[].ikeConfig.ikeEncAlg` | `string` |  | `aes` |  |
| `spec.connections[].ikeConfig.ikeAuthAlg` | `string` |  | `sha1` |  |
| `spec.connections[].ikeConfig.ikePfs` | `string` |  | `group2` |  |
| `spec.connections[].ikeConfig.ikeLifetime` | `int32` |  | `86400` |  |
| `spec.connections[].ipsecConfig` | `AliCloudIpsecConfig` |  |  |  |
| `spec.connections[].ipsecConfig.ipsecEncAlg` | `string` |  | `aes` |  |
| `spec.connections[].ipsecConfig.ipsecAuthAlg` | `string` |  | `md5` |  |
| `spec.connections[].ipsecConfig.ipsecPfs` | `string` |  | `group2` |  |
| `spec.connections[].ipsecConfig.ipsecLifetime` | `int32` |  | `86400` |  |
| `spec.connections[].healthCheckConfig` | `AliCloudVpnHealthCheckConfig` |  |  |  |
| `spec.connections[].healthCheckConfig.enable` | `bool` |  | `false` |  |
| `spec.connections[].healthCheckConfig.sip` | `string` |  |  |  |
| `spec.connections[].healthCheckConfig.dip` | `string` |  |  |  |
| `spec.connections[].healthCheckConfig.interval` | `int32` |  | `3` |  |
| `spec.connections[].healthCheckConfig.retry` | `int32` |  | `3` |  |

## Field Details

### spec.region

`string` · required

Alibaba Cloud region where the VPN Gateway will be created.
Must match the region of the VPC and VSwitch.
Examples: "cn-hangzhou", "cn-shanghai", "us-west-1", "ap-southeast-1".

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.vpcId

`string | valueFrom` · required

VPC ID that the VPN Gateway belongs to.

- references: AliCloudVpc (`status.outputs.vpc_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AliCloudVpc, name: <that resource's name>, fieldPath: status.outputs.vpc_id}} -- a bare string does not parse

### spec.vswitchId

`string | valueFrom` · required

VSwitch ID where the VPN Gateway is placed. Must be in the same VPC as
vpc_id. The VPN Gateway consumes a private IP from this VSwitch.

- references: AliCloudVswitch (`status.outputs.vswitch_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AliCloudVswitch, name: <that resource's name>, fieldPath: status.outputs.vswitch_id}} -- a bare string does not parse

### spec.vpnGatewayName

`string` · required

VPN Gateway name. 2-128 characters; must start with a letter or Chinese
character; can contain digits, underscores, periods, and hyphens.

- rule: {"required":true,"string":{"minLen":"2","maxLen":"128"}}

### spec.description

`string`

Human-readable description of the VPN Gateway's purpose.

### spec.bandwidth

`int32` · required

Maximum public network bandwidth for the VPN Gateway in Mbps.
This is immutable after creation (ForceNew in the provider).

- rule: bandwidth must be one of: 5, 10, 20, 50, 100, 200, 500, 1000 (Mbps)
- rule: {"required":true}

### spec.paymentType

`string` · optional (explicit presence)

Billing method for the VPN Gateway.
"PayAsYouGo" for on-demand billing, "Subscription" for reserved pricing.
Default: "PayAsYouGo"

- default: `PayAsYouGo`
- rule: payment_type must be one of: PayAsYouGo, Subscription

### spec.enableSsl

`bool` · optional (explicit presence)

Enable SSL VPN on the gateway for remote client access.
When true, the gateway allocates a dedicated SSL VPN IP and supports
SSL VPN server/client configuration. Default: false

- default: `false`

### spec.sslConnections

`int32` · optional (explicit presence)

Maximum number of concurrent SSL VPN client connections.
Only relevant when enable_ssl is true. The value determines
the SSL VPN license tier on the gateway.

### spec.tags

`map<string, string>`

Tags to apply to the VPN Gateway resource.

### spec.resourceGroupId

`string`

Resource group ID for access control. Leave empty to use the default
resource group.

### spec.connections

`[]AliCloudVpnConnection`

IPsec VPN connections bundled with this gateway. Each entry creates a
Customer Gateway (from the remote device's public IP) and a VPN
Connection with the specified IKE/IPsec negotiation parameters.

### spec.connections[].name

`string` · required

Connection name. Used as both the customer_gateway_name and
vpn_connection_name in the provider. 2-128 characters.

- rule: {"required":true,"string":{"minLen":"2","maxLen":"128"}}

### spec.connections[].customerGatewayIp

`string` · required

Public IP address of the remote VPN device (on-premises router, firewall,
or other cloud gateway). This is immutable after the customer gateway is
created.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.connections[].customerGatewayAsn

`string`

BGP Autonomous System Number of the remote device. Used when the remote
gateway participates in BGP routing. Format: plain 32-bit ASN (e.g.
"65001") or dot notation (e.g. "4200000001").

### spec.connections[].localSubnets

`[]string` · required

VPC-side CIDR blocks that should be reachable through this tunnel.
The VPN connection routes traffic for these CIDRs through the tunnel.
At least 1, at most 10 CIDR blocks.
Example: ["10.0.0.0/8", "172.16.0.0/12"]

- rule: {"repeated":{"minItems":"1","maxItems":"10"}}

### spec.connections[].remoteSubnets

`[]string` · required

Remote-site CIDR blocks reachable on the other end of the tunnel.
At least 1, at most 10 CIDR blocks.
Example: ["192.168.0.0/16"]

- rule: {"repeated":{"minItems":"1","maxItems":"10"}}

### spec.connections[].enableDpd

`bool` · optional (explicit presence)

Enable Dead Peer Detection (DPD). When enabled, the VPN gateway sends
DPD packets to verify the remote peer is alive and re-negotiates if the
peer becomes unreachable. Default: true

- default: `true`

### spec.connections[].enableNatTraversal

`bool` · optional (explicit presence)

Enable NAT traversal (NAT-T). Required when the remote VPN device is
behind a NAT device. Encapsulates IPsec ESP packets in UDP to traverse
NAT. Default: true

- default: `true`

### spec.connections[].effectImmediately

`bool` · optional (explicit presence)

Initiate IPsec negotiation immediately after creation rather than waiting
for traffic to trigger the tunnel. Default: true

- default: `true`

### spec.connections[].ikeConfig

`AliCloudIkeConfig`

IKE (Phase 1) negotiation parameters. When omitted, the provider uses
sensible defaults (IKEv2, main mode, AES encryption, SHA1 auth, DH group2,
auto-generated pre-shared key, 86400s lifetime).

### spec.connections[].ikeConfig.psk

`string` · sensitive

Pre-shared key (PSK) for IKE authentication. 1-100 characters.
If omitted, the provider auto-generates a random PSK.

- rule: {"string":{"maxLen":"100"}}

### spec.connections[].ikeConfig.ikeVersion

`string` · optional (explicit presence)

IKE protocol version. IKEv2 is recommended for better performance and
security. Default: "ikev2"

- default: `ikev2`
- rule: ike_version must be one of: ikev1, ikev2

### spec.connections[].ikeConfig.ikeMode

`string` · optional (explicit presence)

IKE negotiation mode. "main" provides identity protection (6 messages);
"aggressive" is faster (3 messages) but exposes identities.
Default: "main"

- default: `main`
- rule: ike_mode must be one of: main, aggressive

### spec.connections[].ikeConfig.ikeEncAlg

`string` · optional (explicit presence)

IKE encryption algorithm for Phase 1.
Default: "aes"

- default: `aes`
- rule: ike_enc_alg must be one of: aes, aes192, aes256, des, 3des

### spec.connections[].ikeConfig.ikeAuthAlg

`string` · optional (explicit presence)

IKE authentication algorithm for Phase 1.
Default: "sha1"

- default: `sha1`
- rule: ike_auth_alg must be one of: md5, sha1, sha256, sha384, sha512

### spec.connections[].ikeConfig.ikePfs

`string` · optional (explicit presence)

Diffie-Hellman key exchange group for Perfect Forward Secrecy in Phase 1.
Higher groups provide stronger security at the cost of performance.
Default: "group2" (1024-bit MODP)

- default: `group2`
- rule: ike_pfs must be one of: group1, group2, group5, group14

### spec.connections[].ikeConfig.ikeLifetime

`int32` · optional (explicit presence)

IKE SA (Security Association) lifetime in seconds. When this expires, the
SA is re-negotiated. Range: 0-86400. Default: 86400 (24 hours)

- default: `86400`
- rule: ike_lifetime must be between 0 and 86400 seconds

### spec.connections[].ipsecConfig

`AliCloudIpsecConfig`

IPsec (Phase 2) negotiation parameters. When omitted, the provider uses
sensible defaults (AES encryption, MD5 auth, DH group2, 86400s lifetime).

### spec.connections[].ipsecConfig.ipsecEncAlg

`string` · optional (explicit presence)

IPsec encryption algorithm for Phase 2 data traffic.
Default: "aes"

- default: `aes`
- rule: ipsec_enc_alg must be one of: aes, aes192, aes256, des, 3des

### spec.connections[].ipsecConfig.ipsecAuthAlg

`string` · optional (explicit presence)

IPsec authentication algorithm for Phase 2 data traffic.
Default: "md5"

- default: `md5`
- rule: ipsec_auth_alg must be one of: md5, sha1, sha256, sha384, sha512

### spec.connections[].ipsecConfig.ipsecPfs

`string` · optional (explicit presence)

Diffie-Hellman key exchange group for Perfect Forward Secrecy in Phase 2.
"disabled" turns off PFS (the Phase 1 key is reused).
Default: "group2"

- default: `group2`
- rule: ipsec_pfs must be one of: disabled, group1, group2, group5, group14

### spec.connections[].ipsecConfig.ipsecLifetime

`int32` · optional (explicit presence)

IPsec SA lifetime in seconds. Range: 0-86400. Default: 86400 (24 hours)

- default: `86400`
- rule: ipsec_lifetime must be between 0 and 86400 seconds

### spec.connections[].healthCheckConfig

`AliCloudVpnHealthCheckConfig`

Health check configuration for monitoring tunnel connectivity. When
enabled, the VPN gateway periodically probes the remote endpoint and
can trigger failover if the tunnel becomes unhealthy.

### spec.connections[].healthCheckConfig.enable

`bool` · optional (explicit presence)

Enable health checks for this connection. Default: false

- default: `false`

### spec.connections[].healthCheckConfig.sip

`string`

Source IP address for health check probes. Should be a private IP
within the local VPC that is routable through the tunnel.

### spec.connections[].healthCheckConfig.dip

`string`

Destination IP address for health check probes. Should be an IP on the
remote network reachable through the tunnel.

### spec.connections[].healthCheckConfig.interval

`int32` · optional (explicit presence)

Interval between health check probes in seconds. Default: 3

- default: `3`

### spec.connections[].healthCheckConfig.retry

`int32` · optional (explicit presence)

Maximum number of consecutive probe failures before the tunnel is
declared unhealthy. Default: 3

- default: `3`

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AliCloudVpnGateway, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.vpn_gateway_id` | `string` | The VPN Gateway ID assigned by Alibaba Cloud (e.g., "vpn-xxxxx"). |
| `status.outputs.internet_ip` | `string` | The VPN Gateway's public internet IP address, used as the local endpoint for IPsec tunnels. Remote VPN devices connect to this IP. |
| `status.outputs.ssl_vpn_internet_ip` | `string` | The SSL VPN internet IP address, populated only when enable_ssl is true. SSL VPN clients connect to this IP for remote access. |
| `status.outputs.connection_ids` | `map<string, string>` | Map of connection name to VPN connection ID. Keys are the names specified in spec.connections[].name. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.vpcId` | AliCloudVpc | `status.outputs.vpc_id` |
| `spec.vswitchId` | AliCloudVswitch | `status.outputs.vswitch_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
