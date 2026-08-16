# AzureNetworkSecurityGroup Terraform Module

## Overview

This Terraform module provisions an Azure Network Security Group using the
`azurerm` provider. It creates one `azurerm_network_security_group` plus a
standalone `azurerm_network_security_rule` per rule -- the standalone form
enforces the singular/plural conflicts at plan time and gives each rule
its own plan line. (The Pulumi module manages the same rules inline on the
group; both engines put the identical rule set into ARM.)

Rules and tags update in place and take effect immediately for every
subnet and NIC the group guards. Name, region, and resource group are the
group's ARM identity; changing any of them replaces it, detaching it from
every subnet until re-attached.

A group with no rules is meaningful: Azure's implicit defaults then govern
(allow VNet-internal and load-balancer traffic, deny other inbound, allow
all outbound). The subnet-side attachment is deliberately not modeled
here: a subnet declares which NSG guards it, so one group serves many
subnets.

## Resources Created

- `azurerm_network_security_group.main` -- the group
- `azurerm_network_security_rule.rules` -- one per security rule

## Variables

| Variable | Type | Description |
|----------|------|-------------|
| `metadata` | object | Planton metadata (name, org, env) |
| `spec` | object | NSG specification |

Key spec fields:

| Field | Required | Description |
|-------|----------|-------------|
| `region` | yes | Azure region (must match the guarded subnets/NICs) |
| `resource_group` | yes | Resource group name |
| `name` | yes | Group name, unique within the resource group |
| `security_rules` | no | 5-tuple filters: `priority` (100-4096, lowest first), `direction` (`INBOUND`/`OUTBOUND`), `access` (`ALLOW`/`DENY`), `protocol` (`ANY`/`TCP`/`UDP`/`ICMP`/`AH`/`ESP`; `ANY` is ARM's `*`); ports and each address side take exactly one form, unset means any |
| `tags` | no | User tags, merged over metadata-derived tags |

## Outputs

| Output | Description |
|--------|-------------|
| `network_security_group_id` | Full ARM ID of the group |
| `network_security_group_name` | The group's name as deployed |

## Usage

```hcl
module "nsg" {
  source = "./iac/tf"

  metadata = { name = "web-tier", org = "mycompany", env = "production" }

  spec = {
    region         = "eastus"
    resource_group = "network-rg"
    name           = "web-tier"
    security_rules = [
      {
        name                   = "allow-https-inbound"
        priority               = 100
        direction              = "INBOUND"
        access                 = "ALLOW"
        protocol               = "TCP"
        source_address_prefix  = "Internet"
        destination_port_range = "443"
      }
    ]
  }
}
```

## Required Permissions

See [`../permissions.yaml`](../permissions.yaml) for the least-privilege action set the deploying principal needs.
