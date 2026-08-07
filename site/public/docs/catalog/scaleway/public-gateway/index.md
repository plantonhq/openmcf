---
title: "Public Gateway"
description: "Public Gateway deployment documentation"
icon: "package"
order: 100
componentName: "scalewaypublicgateway"
---

# Scaleway Public Gateway

Deploys a managed Public Gateway on a Scaleway Private Network with NAT masquerade for outbound internet access, optional SSH bastion for secure jump-host connectivity, and configurable port forwarding (PAT) rules for inbound traffic routing. The IaC module bundles a dedicated Flexible IP, the gateway appliance, and the network attachment into a single declarative resource. Supports ValueFromRef for Private Network dependency wiring in InfraCharts.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Flexible IP** -- a dedicated public IPv4 address (`scaleway_vpc_public_gateway_ip`) assigned to the gateway, with optional reverse DNS configuration. The IP persists independently of the gateway appliance.
- **Public Gateway** -- a `scaleway_vpc_public_gateway` appliance in the specified zone with the configured gateway type, SSH bastion settings, and SMTP policy
- **Gateway Network Attachment** -- a `scaleway_vpc_gateway_network` binding the gateway to the referenced Private Network with NAT masquerade settings
- **PAT Rules** -- created only when `patRules` entries are provided; `scaleway_vpc_public_gateway_pat_rule` resources mapping public ports to private IP:port pairs inside the attached network
- **Scaleway Tags** -- resource metadata tags (resource name, kind, organization, environment) applied to the gateway and Flexible IP for tracking and governance

## Before You Deploy

### Scaleway Account

- **A Scaleway Private Network** in the target region. Provide the Private Network UUID directly or reference a ScalewayPrivateNetwork Cloud Resource via ValueFromRef.
- **Choose an Availability Zone** -- Public Gateways are zonal resources (`fr-par-1`, `nl-ams-1`, `pl-waw-1`). The zone must be within the same region as the target Private Network.
- **Gateway type** -- `VPC-GW-S` (standard) is sufficient for most workloads. `VPC-GW-XL` (high-bandwidth, up to 10 Gbps) is available only in Paris region zones.

## Deploy

### Console

Open the deployment store, find **Scaleway Public Gateway**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **NAT Gateway** preset in the [Presets](#presets) tab for a standard gateway with NAT masquerade enabled.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: scaleway.planton.dev/v1
kind: ScalewayPublicGateway
metadata:
  name: app-gateway
  org: acme-corp
  env: prod
spec:
  privateNetworkId:
    value: "abc12345-6789-def0-1234-567890abcdef"
  zone: fr-par-1
  type: VPC-GW-S
  enableMasquerade: true
```

```shell
planton apply -f scaleway-public-gateway.yaml
```

This creates a standard Public Gateway in Paris-1 with NAT masquerade enabled, allowing all resources in the attached Private Network to reach the internet through the gateway's public IP. No SSH bastion or port forwarding is configured. A Stack Job tracks the provisioning in real time.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the gateway to a Private Network deployed in the same InfraPipeline:

```yaml
spec:
  privateNetworkId:
    valueFrom:
      kind: ScalewayPrivateNetwork
      name: app-network
      fieldPath: status.outputs.private_network_id
```

The InfraPipeline resolves the dependency graph, deploys the Private Network first, then provisions the Public Gateway attached to it.

## Key Configuration

These are the most important decisions when configuring a Public Gateway. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Gateway type** -- The `type` field selects the performance tier. `VPC-GW-S` provides standard bandwidth and is sufficient for most workloads. `VPC-GW-XL` offers up to 10 Gbps but is available only in Paris zones (`fr-par-1`, `fr-par-2`).

**NAT masquerade** -- The `enableMasquerade` field controls whether resources in the attached Private Network can reach the internet through the gateway's public IP. Almost always enabled. Disable only when using the gateway solely as an SSH bastion.

**SSH bastion** -- The `bastion` block enables the gateway as a secure SSH jump host. Configure `bastion.port` (default 22) and restrict access with `bastion.allowedIpRanges` to specific CIDRs (office IP, VPN exit). Leaving the allowlist empty permits all IPs -- not recommended for production.

**Port forwarding** -- The `patRules` array maps public ports on the gateway's IP to private IP:port pairs inside the network. Each rule specifies `publicPort`, `privateIp`, `privatePort`, and `protocol` (`tcp`, `udp`, or `both`). Use for exposing specific services without assigning individual public IPs.

**SMTP** -- Set `enableSmtp` to true only when resources in the Private Network need to send email directly on port 25. Blocked by default to prevent spam abuse.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **ScalewayPrivateNetwork** | `privateNetworkId` | `status.outputs.private_network_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `gateway_id` | Zoned ID of the Public Gateway | Scaleway API operations, monitoring dashboards |
| `public_ip_address` | Public IPv4 address of the gateway | DNS A records, firewall allowlists, connectivity diagnostics |
| `public_ip_id` | Zoned ID of the Flexible IP resource | IP lifecycle management, reassignment to replacement gateways |
| `gateway_network_id` | Zoned ID of the gateway-to-network attachment | Network attachment management, NAT diagnostics |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**NAT gateway** -- A standard gateway with NAT masquerade for outbound internet access. The most common configuration for Private Networks hosting Kapsule nodes, databases, and application servers that need to reach external services without individual public IPs. Start from the **NAT Gateway** preset.

**Bastion-enabled gateway** -- A gateway with both NAT masquerade and SSH bastion for secure operator access. Provides a single auditable entry point for SSH connections to private instances, eliminating the need for public IPs on individual machines. Start from the **Bastion-Enabled** preset.

## Works With

- [**Scaleway Private Network**](/cloud-catalog/scaleway-private-network) -- provides the Private Network that the gateway attaches to for NAT, bastion, and port forwarding