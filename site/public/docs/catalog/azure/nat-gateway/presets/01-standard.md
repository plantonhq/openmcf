---
title: "Standard NAT Gateway"
description: "This preset creates a zonal Standard NAT Gateway that SNATs through one referenced `AzurePublicIp`, giving every subnet that attaches it (via the subnet's `natGatewayId`) stable outbound connectivity..."
type: "preset"
rank: "01"
presetSlug: "01-standard"
componentSlug: "nat-gateway"
componentTitle: "NAT Gateway"
provider: "azure"
icon: "package"
order: 1
---

# Standard NAT Gateway

This preset creates a zonal Standard NAT Gateway that SNATs through one referenced `AzurePublicIp`, giving every subnet that attaches it (via the subnet's `natGatewayId`) stable outbound connectivity with a known source address. The address is a first-class resource -- visible in the resource graph, allowlistable, and reusable -- never created invisibly inside the gateway.

## When to Use

- Subnets running AKS nodes, VMs, or containers that need outbound internet access
- Replacing Azure default outbound access (which Microsoft is retiring) with dedicated egress
- Workloads that require a known source IP for firewall allowlisting or compliance

## Key Configuration Choices

- **`publicIpIds` reference** -- resolves to the address's `public_ip_id` output. Each address provides 64,512 SNAT ports; add more addresses (or a prefix, see preset 02) to scale
- **Zonal** (`zones: ["1"]`) -- a Standard gateway lives in one availability zone; zone-resilient designs deploy one gateway per zone with per-zone subnets. The referenced address must be in the same zone
- **Idle timeout** (`idleTimeoutInMinutes: 10`) -- balances connection reuse for long-lived connections against SNAT port hold time. Azure default is 4 minutes; maximum is 120
- **Subnet attachment lives on the subnet** -- set `natGatewayId` on each `AzureSubnet` that should egress through this gateway

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<azure-region>` | Azure region (must match the subnets it will serve) | Your regional deployment strategy |
| `<your-resource-group-name>` | Name of the resource group | Azure portal or `AzureResourceGroup` status outputs |
| `<your-public-ip-resource-name>` | Planton metadata name of the `AzurePublicIp` (zone-matched) | Your public IP resource |
