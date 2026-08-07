# AliCloud VPN Gateway

Deploys an Alibaba Cloud VPN Gateway with bundled customer gateways and IPsec VPN connections. The component provisions encrypted site-to-site connectivity between an Alibaba Cloud VPC and remote networks (on-premises data centers, branch offices, or other cloud environments) over the public internet using IPsec. The VPN Gateway, customer gateways, and VPN connections are deployed as a single atomic unit because a VPN Gateway without at least one connection is non-functional.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **VPN Gateway** -- a `alicloud_vpn_gateway` resource placed in the specified VPC and VSwitch with configurable bandwidth, billing, and optional SSL VPN support
- **Customer Gateways** -- one `alicloud_vpn_customer_gateway` per connection entry, representing each remote VPN device's public IP and optional BGP ASN
- **VPN Connections** -- one `alicloud_vpn_connection` per connection entry, linking the VPN Gateway to a customer gateway with IKE/IPsec negotiation parameters, local/remote subnet routing, and optional health checks

## Before You Deploy

### Planton Setup

- **AliCloud Provider Connection** -- an active connection in the Connect module with credentials for the target Alibaba Cloud account. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery.

### Alibaba Cloud Account

- **An existing VPC and VSwitch** -- the VPN Gateway requires placement in a VSwitch within a VPC.
- **Remote device details** -- public IP address, local/remote CIDR blocks, and IKE/IPsec parameters of the on-premises or remote VPN device.
- **Bandwidth selection** -- the `bandwidth` field is immutable after creation. Choose based on expected traffic volume (5, 10, 20, 50, 100, 200, 500, or 1000 Mbps).

## Deploy

### Console

Open the deployment store, find **AliCloud VPN Gateway**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields including VPC, VSwitch, bandwidth, and VPN connections.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: alicloud.planton.dev/v1
kind: AliCloudVpnGateway
metadata:
  name: dc-vpn
  org: acme-corp
  env: prod
spec:
  region: cn-hangzhou
  vpcId:
    value: vpc-bp1234567890
  vswitchId:
    value: vsw-nat-zone-a
  vpnGatewayName: datacenter-vpn
  bandwidth: 100
  connections:
    - name: primary-dc
      customerGatewayIp: "203.0.113.1"
      localSubnets:
        - "10.0.0.0/8"
      remoteSubnets:
        - "192.168.0.0/16"
```

```shell
planton apply -f alicloud-vpn.yaml
```

This creates a 100 Mbps VPN Gateway with one IPsec connection to a remote data center. A Stack Job tracks the provisioning in real time.

### InfraChart

When deploying as part of a hybrid network stack, use ValueFromRef to wire VPC and VSwitch dependencies:

```yaml
spec:
  region: cn-hangzhou
  vpcId:
    valueFrom:
      kind: AliCloudVpc
      name: platform-vpc
      fieldPath: status.outputs.vpc_id
  vswitchId:
    valueFrom:
      kind: AliCloudVswitch
      name: vpn-vswitch
      fieldPath: status.outputs.vswitch_id
  vpnGatewayName: datacenter-vpn
  bandwidth: 100
  connections:
    - name: primary-dc
      customerGatewayIp: "203.0.113.1"
      localSubnets:
        - "10.0.0.0/8"
      remoteSubnets:
        - "192.168.0.0/16"
      ikeConfig:
        ikeVersion: ikev2
        ikeEncAlg: aes256
        ikeAuthAlg: sha256
        ikePfs: group14
```

The InfraPipeline resolves the dependency graph and provisions VPC and VSwitch before the VPN Gateway.

## Key Configuration

These are the most important decisions when configuring a VPN Gateway. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Bandwidth** -- The `bandwidth` field selects the gateway's maximum throughput in Mbps (5 to 1000). This is immutable after creation -- plan for peak traffic.

**IKE configuration** -- Each connection's `ikeConfig` controls Phase 1 negotiation. IKEv2 is recommended. Configure encryption (aes256), authentication (sha256), and DH group (group14) to match the remote device.

**IPsec configuration** -- Each connection's `ipsecConfig` controls Phase 2 data encryption. Align encryption algorithm, authentication, and PFS settings with the remote device.

**SSL VPN** -- The `enableSsl` field adds remote client access capability. When enabled, configure `sslConnections` to set the concurrent client limit.

**Health checks** -- Each connection supports `healthCheckConfig` for monitoring tunnel liveness. Enable for production deployments with active-standby failover.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **AliCloudVpc** | `vpcId` | `status.outputs.vpc_id` |
| **AliCloudVswitch** | `vswitchId` | `status.outputs.vswitch_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `vpn_gateway_id` | VPN Gateway ID (e.g., vpn-xxxxx) | Monitoring, route management |
| `internet_ip` | Gateway's public IP for IPsec tunnels | Remote device configuration |
| `ssl_vpn_internet_ip` | SSL VPN endpoint IP (when SSL is enabled) | Client VPN configuration |
| `connection_ids` | Map of connection names to VPN connection IDs | Advanced management |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Basic site-to-site** -- A single IPsec connection to one remote network with default IKE/IPsec settings. Start from the **Basic Site to Site** preset.

**Production multi-site** -- High-bandwidth gateway with multiple connections, AES-256 encryption, health checks, and organizational tags. Start from the **Production Multi Site** preset.

**SSL-enabled** -- A VPN Gateway with both IPsec site-to-site and SSL VPN remote client access. Start from the **SSL Enabled** preset.

## Works With

- [**AliCloud VPC**](/cloud-catalog/ali-cloud-vpc) -- the VPC this VPN Gateway connects
- [**AliCloud VSwitch**](/cloud-catalog/ali-cloud-vswitch) -- placement VSwitch for the VPN Gateway
- [**AliCloud CEN Instance**](/cloud-catalog/ali-cloud-cen-instance) -- alternative connectivity for multi-VPC/multi-region topologies
