# AzureSubnet Terraform Module

## Overview

This Terraform module provisions an Azure subnet using the `azurerm`
provider. It creates an `azurerm_subnet` -- the workload segment that
partitions a virtual network's address space -- plus one association
resource per declared attachment (route table, network security group,
NAT gateway). The subnet is an ARM child of its virtual network: the
module derives the resource group and network name from the parent
network's ARM ID. Attachments are declared subnet-side (matching Azure's
model), so one table, group, or gateway serves many subnets.

Everything except name and the parent network updates in place; those two
are the subnet's ARM identity, so changing either replaces the subnet and
everything deployed into it. Subnets are not tracked ARM resources, so
they carry no tags.

## Resources Created

- `azurerm_subnet.main` -- the subnet
- `azurerm_subnet_route_table_association.main` / `..._network_security_group_association.main` / `..._nat_gateway_association.main` -- one per declared attachment

## Variables

| Variable | Type | Description |
|----------|------|-------------|
| `metadata` | object | Planton metadata (name, org, env) |
| `spec` | object | Subnet specification |

Key spec fields:

| Field | Required | Description |
|-------|----------|-------------|
| `virtual_network_id` | yes | ARM ID of the parent network; resource group and network name are parsed from it |
| `name` | yes | Subnet name, unique within the network |
| `address_prefixes` | one of | Self-managed CIDR blocks |
| `ip_address_pool` | one of | Network Manager IPAM allocation (`id`, `number_of_ip_addresses`); exactly one address source |
| `service_endpoints` / `service_endpoint_policy_ids` | no | Backbone routing to Azure services, optionally narrowed |
| `delegations` | no | Hand the subnet to a PaaS service (`name`, `service_name`, optional `actions`) |
| `private_endpoint_network_policies` / `private_link_service_network_policies_enabled` | no | Policy toggles; unset applies Azure's defaults |
| `default_outbound_access_enabled` | no | Azure default true; production subnets set false and route egress explicitly |
| `route_table_id` / `network_security_group_id` / `nat_gateway_id` | no | The three attach seams, realized as association resources |

## Outputs

| Output | Description |
|--------|-------------|
| `subnet_id` | Full ARM ID of the subnet |
| `subnet_name` | The subnet's name within its network |
| `address_prefixes` | The CIDR blocks actually assigned (IPAM allocations surface here) |
| `virtual_network_name` | Parent network name, derived from its ARM ID |
| `resource_group_name` | Parent network's resource group, derived from its ARM ID |

## Usage

```hcl
module "subnet" {
  source = "./iac/tf"

  metadata = { name = "app", org = "mycompany", env = "production" }

  spec = {
    virtual_network_id        = "/subscriptions/xxx/resourceGroups/network-rg/providers/Microsoft.Network/virtualNetworks/prod-vnet"
    name                      = "app"
    address_prefixes          = ["10.0.1.0/24"]
    nat_gateway_id            = "/subscriptions/xxx/.../natGateways/prod-egress"

    default_outbound_access_enabled = false
  }
}
```

## Required Permissions

The deploying credential needs
`Microsoft.Network/virtualNetworks/subnets/write` on the parent network,
plus join rights on any attached route table, NSG, or NAT gateway --
held via Network Contributor, Contributor, or Owner.
