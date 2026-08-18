# AzureNetworkInterface Terraform Module

## Overview

This Terraform module provisions an Azure Network Interface using the
`azurerm` provider. It creates an `azurerm_network_interface` -- the
attachment point that gives a virtual machine its presence in a subnet --
plus association resources for the NIC-level network security group,
application security group memberships, load-balancer backend-pool and
inbound-NAT-rule memberships, and Application Gateway backend-pool
memberships. Memberships are expressed here, from the member side --
Azure's own model: the load balancer declares its pools and rules, and
each of this NIC's ip_configurations declares which it joins.

The NIC is the attachment point, not the machine. An AzureVirtualMachine
consumes this module's `network_interface_id` output via its
`network_interface_ids` (a VM references its NICs, never contains them),
each IP configuration deploys into a referenced AzureSubnet and may front
a referenced AzurePublicIp, and the NIC-level NSG is the per-workload
complement to subnet-level filtering (inbound traffic must pass the
subnet NSG then the NIC NSG when both are attached).

Name, region, and edge zone are the NIC's identity -- changing any of
them replaces the NIC, detaching it from its VM. Everything else
(configurations, DNS, acceleration, forwarding, associations, tags)
updates in place. The MAC address is assigned when the NIC attaches to a
running VM, not at creation.

## Resources Created

- `azurerm_network_interface.main` -- the NIC itself
- `azurerm_network_interface_security_group_association.main` -- when a NIC-level NSG is referenced
- `azurerm_network_interface_application_security_group_association.main` -- one per application security group membership
- `azurerm_network_interface_backend_address_pool_association.main` -- one per load-balancer backend-pool membership declared on an ip_configuration
- `azurerm_network_interface_nat_rule_association.main` -- one per single-target inbound NAT rule an ip_configuration completes
- `azurerm_network_interface_application_gateway_backend_address_pool_association.main` -- one per Application Gateway backend-pool membership

## Variables

| Variable | Type | Description |
|----------|------|-------------|
| `metadata` | object | Planton metadata (name, org, env) |
| `spec` | object | Network interface specification |

Key spec fields:

| Field | Required | Description |
|-------|----------|-------------|
| `region` | yes | Azure region; must match the virtual network and the VM that attaches the NIC |
| `resource_group` | yes | Resource group name |
| `name` | yes | NIC name, unique within the resource group; renaming replaces the NIC and detaches it from any VM |
| `ip_configurations` | yes | At least one; each places a private address in a subnet (DYNAMIC or STATIC, IPv4 or IPv6) and may front a public IP. With multiple, the first must be primary (ARM's contract, spec-enforced) |
| `ip_configurations[].load_balancer_backend_address_pool_ids` | no | Load-balancer backend pools this configuration joins, as resolved ARM IDs (the platform resolves references to the LB's name-keyed `backend_pool_ids` output); one association resource each |
| `ip_configurations[].load_balancer_inbound_nat_rule_ids` | no | Single-target inbound NAT rules this configuration completes, as resolved ARM IDs (resolved from the LB's name-keyed `nat_rule_ids` output); one association resource each |
| `ip_configurations[].application_gateway_backend_address_pool_ids` | no | Application Gateway backend pools this configuration joins, as plain ARM IDs; one association resource each |
| `dns_servers` | no | DNS servers overriding the virtual network's DNS for this NIC only |
| `internal_dns_name_label` | no | VNet-internal DNS label other VMs can resolve this NIC by |
| `accelerated_networking_enabled` | no | SR-IOV; Azure defaults to false -- enable on every supported VM size |
| `ip_forwarding_enabled` | no | Forwarding of traffic not addressed to the NIC; network virtual appliances only |
| `auxiliary_mode` / `auxiliary_sku` | no | Preview NVA acceleration (subscription must be enrolled); both or neither, unset sends nothing |
| `edge_zone` | no | Edge Zone pinning; fixed at creation |
| `network_security_group_id` | no | The NIC-level NSG, as a resolved ARM ID; realized as an association resource |
| `application_security_group_ids` | no | ASG memberships, as ARM IDs; one association resource each |
| `tags` | no | User tags, merged over metadata-derived tags (user wins) |

## Outputs

| Output | Description |
|--------|-------------|
| `network_interface_id` | Full ARM ID of the NIC -- the join key AzureVirtualMachine attaches through |
| `network_interface_name` | The NIC's name as deployed |
| `private_ip_address` | The primary configuration's private IP address |
| `private_ip_addresses` | All configurations' private addresses, in configuration order |
| `mac_address` | The NIC's MAC address (populated once attached to a running VM) |
| `internal_domain_name_suffix` | The DNS suffix completing `internal_dns_name_label` into a resolvable FQDN |

## Usage

```hcl
module "network_interface" {
  source = "./iac/tf"

  metadata = { name = "app-nic", org = "mycompany", env = "production" }

  spec = {
    region         = "eastus"
    resource_group = "app-rg"
    name           = "app-nic"

    ip_configurations = [{
      name      = "primary"
      subnet_id = "/subscriptions/xxx/.../subnets/app-subnet"
    }]

    accelerated_networking_enabled = true
  }
}
```

## Required Permissions

See [`../permissions.yaml`](../permissions.yaml) for the least-privilege actions the deploying principal needs.
