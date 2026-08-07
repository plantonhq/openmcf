# AzureRouteTable Terraform Module

## Overview

This Terraform module provisions an Azure route table using the `azurerm`
provider. It creates a single `azurerm_route_table` with its user-defined
routes managed inline (a route has no life of its own in Azure) and BGP
propagation control.

Routes, BGP propagation, and tags update in place -- and take effect
immediately for every subnet attached to the table. Name, region, and
resource group are the table's ARM identity; changing any of them replaces
the table, detaching it from every subnet until re-attached.

The subnet-side attachment is deliberately not modeled here: a subnet
declares which route table it uses (matching Azure's model), so one table
serves many subnets without listing them.

## Resources Created

- `azurerm_route_table.main` -- the route table with inline routes

## Variables

| Variable | Type | Description |
|----------|------|-------------|
| `metadata` | object | Planton metadata (name, org, env) |
| `spec` | object | Route table specification |

Key spec fields:

| Field | Required | Description |
|-------|----------|-------------|
| `region` | yes | Azure region (must match the attaching subnets' networks) |
| `resource_group` | yes | Resource group name |
| `name` | yes | Table name, unique within the resource group |
| `routes` | no | User-defined routes: `name`, `address_prefix` (CIDR or service tag), `next_hop_type` (enum name string), `next_hop_in_ip_address` (VIRTUAL_APPLIANCE only) |
| `bgp_route_propagation_enabled` | no | Azure default true; false is the forced-tunneling hardening |
| `tags` | no | User tags, merged over metadata-derived tags |

## Outputs

| Output | Description |
|--------|-------------|
| `route_table_id` | Full ARM ID of the table |
| `route_table_name` | The table's name as deployed |

## Usage

```hcl
module "route_table" {
  source = "./iac/tf"

  metadata = {
    name = "egress-via-firewall"
    org  = "mycompany"
    env  = "production"
  }

  spec = {
    region         = "eastus"
    resource_group = "network-rg"
    name           = "egress-via-firewall"
    routes = [
      {
        name                   = "default-via-firewall"
        address_prefix         = "0.0.0.0/0"
        next_hop_type          = "VIRTUAL_APPLIANCE"
        next_hop_in_ip_address = "10.0.1.4"
      }
    ]
    bgp_route_propagation_enabled = false
  }
}
```

## Required Permissions

The deploying credential needs `Microsoft.Network/routeTables/write` on
the resource group -- held via Network Contributor, Contributor, or Owner.
