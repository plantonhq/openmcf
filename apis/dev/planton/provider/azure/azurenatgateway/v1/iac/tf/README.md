# AzureNatGateway Terraform Module

## Overview

This Terraform module provisions an Azure NAT Gateway using the `azurerm`
provider. It creates an `azurerm_nat_gateway` -- the managed SNAT service
that gives every workload in its attached subnets stable, scalable
outbound connectivity -- plus one association resource per referenced
public IP and public IP prefix.

The gateway is just the gateway. The addresses it SNATs through are
referenced first-class AzurePublicIp / AzurePublicIpPrefix resources (each
address adds 64,512 SNAT ports), and the subnets it serves attach
themselves via AzureSubnet's `nat_gateway_id` -- matching Azure's model,
so one gateway serves many subnets without listing them. A gateway with
no associated addresses deploys but cannot translate anything.

Idle timeout, tags, and the IP/prefix associations update in place. Name,
SKU, and zone are the gateway's identity -- changing any of them replaces
it, briefly interrupting egress for every attached subnet.

## Resources Created

- `azurerm_nat_gateway.main` -- the gateway itself
- `azurerm_nat_gateway_public_ip_association.main` -- one per referenced public IP
- `azurerm_nat_gateway_public_ip_prefix_association.main` -- one per referenced prefix

## Variables

| Variable | Type | Description |
|----------|------|-------------|
| `metadata` | object | Planton metadata (name, org, env) |
| `spec` | object | NAT gateway specification |

Key spec fields:

| Field | Required | Description |
|-------|----------|-------------|
| `region` | yes | Azure region; a gateway only serves subnets in its own region |
| `resource_group` | yes | Resource group name |
| `name` | yes | Gateway name, unique within the resource group |
| `sku_name` | no | `STANDARD` (Azure's default, zonal) or `STANDARD_V2` (zone-redundant, needs StandardV2 addresses, empty zones) |
| `idle_timeout_in_minutes` | no | SNAT port hold time for idle TCP, 4-120; Azure default 4 |
| `zones` | no | Single zone to pin a STANDARD gateway to; empty for STANDARD_V2 |
| `public_ip_ids` / `public_ip_prefix_ids` | no | ARM IDs of the addresses and ranges to SNAT through |
| `tags` | no | User tags, merged over metadata-derived tags |

## Outputs

| Output | Description |
|--------|-------------|
| `nat_gateway_id` | Full ARM ID of the gateway -- the join key subnets attach through |
| `nat_gateway_name` | The gateway's name as deployed |
| `resource_guid` | The immutable GUID ARM assigns the gateway (billing/support correlation) |

## Usage

```hcl
module "nat_gateway" {
  source = "./iac/tf"

  metadata = { name = "prod-egress", org = "mycompany", env = "production" }

  spec = {
    region               = "eastus"
    resource_group       = "network-rg"
    name                 = "prod-egress"
    public_ip_prefix_ids = ["/subscriptions/xxx/.../publicIPPrefixes/prod-egress"]
  }
}
```

## Required Permissions

The deploying credential needs `Microsoft.Network/natGateways/write` on
the resource group, plus join rights on the associated public IPs and
prefixes -- held via Network Contributor, Contributor, or Owner.
