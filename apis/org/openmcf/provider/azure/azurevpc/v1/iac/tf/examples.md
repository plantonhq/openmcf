# Azure VPC - Terraform Examples

## Minimal Configuration

```hcl
module "vpc" {
  source = "./path/to/module"

  metadata = {
    name = "dev-vpc"
  }

  spec = {
    address_space_cidr = "10.10.0.0/16"
    nodes_subnet_cidr  = "10.10.0.0/18"
  }

  location = "eastus"
}
```

## Development Environment

```hcl
module "dev_vpc" {
  source = "./path/to/module"

  metadata = {
    name = "dev-vpc"
    env  = "development"
  }

  spec = {
    address_space_cidr = "10.10.0.0/16"
    nodes_subnet_cidr  = "10.10.0.0/18"
    tags = {
      purpose = "development"
      team    = "platform"
    }
  }

  location = "eastus"
}

output "dev_vnet_id" {
  value = module.dev_vpc.vnet_id
}

output "dev_subnet_id" {
  value = module.dev_vpc.nodes_subnet_id
}
```

## Production with NAT Gateway

```hcl
module "prod_vpc" {
  source = "./path/to/module"

  metadata = {
    name = "prod-vpc"
    org  = "mycompany"
    env  = "production"
  }

  spec = {
    address_space_cidr     = "10.0.0.0/16"
    nodes_subnet_cidr      = "10.0.0.0/18"
    is_nat_gateway_enabled = true
    tags = {
      cost-center = "infrastructure"
      compliance  = "soc2"
    }
  }

  location = "eastus"
}

# Use outputs for AKS cluster
output "prod_vnet_id" {
  description = "VNet ID for AKS cluster deployment"
  value       = module.prod_vpc.vnet_id
}

output "prod_subnet_id" {
  description = "Subnet ID for AKS node pool"
  value       = module.prod_vpc.nodes_subnet_id
}

output "prod_nat_gateway_id" {
  description = "NAT Gateway ID for outbound connectivity"
  value       = module.prod_vpc.nat_gateway_id
}
```

## With Private DNS Zone Links

```hcl
# First, create Private DNS zones (or reference existing ones)
resource "azurerm_private_dns_zone" "sql" {
  name                = "privatelink.database.windows.net"
  resource_group_name = "shared-dns-rg"
}

resource "azurerm_private_dns_zone" "storage" {
  name                = "privatelink.blob.core.windows.net"
  resource_group_name = "shared-dns-rg"
}

module "vpc_with_dns" {
  source = "./path/to/module"

  metadata = {
    name = "prod-vpc-dns"
    org  = "mycompany"
    env  = "production"
  }

  spec = {
    address_space_cidr     = "10.0.0.0/16"
    nodes_subnet_cidr      = "10.0.0.0/18"
    is_nat_gateway_enabled = true
    dns_private_zone_links = [
      azurerm_private_dns_zone.sql.id,
      azurerm_private_dns_zone.storage.id
    ]
    tags = {
      environment = "production"
    }
  }

  location = "eastus"

  depends_on = [
    azurerm_private_dns_zone.sql,
    azurerm_private_dns_zone.storage
  ]
}
```

## Multi-Environment Setup

```hcl
locals {
  environments = {
    dev = {
      cidr              = "10.10.0.0/16"
      subnet_cidr       = "10.10.0.0/18"
      nat_gateway       = false
      location          = "eastus"
    }
    staging = {
      cidr              = "10.20.0.0/16"
      subnet_cidr       = "10.20.0.0/18"
      nat_gateway       = true
      location          = "eastus"
    }
    prod = {
      cidr              = "10.0.0.0/16"
      subnet_cidr       = "10.0.0.0/18"
      nat_gateway       = true
      location          = "eastus"
    }
  }
}

module "vpcs" {
  source   = "./path/to/module"
  for_each = local.environments

  metadata = {
    name = "${each.key}-vpc"
    env  = each.key
    org  = "mycompany"
  }

  spec = {
    address_space_cidr     = each.value.cidr
    nodes_subnet_cidr      = each.value.subnet_cidr
    is_nat_gateway_enabled = each.value.nat_gateway
    tags = {
      environment = each.key
      managed-by  = "terraform"
    }
  }

  location = each.value.location
}

output "vnet_ids" {
  value = { for k, v in module.vpcs : k => v.vnet_id }
}
```

## Use Cases

- **Basic VNet**: Minimal configuration for development/testing environments
- **Production VNet**: With NAT Gateway for reliable outbound connectivity and predictable egress IPs
- **DNS Integration**: VNet with Private DNS zones for Azure PaaS services (SQL Database, Storage, etc.)
- **Multi-Environment**: Consistent network deployment across dev, staging, and production
- **AKS Foundation**: Network infrastructure for Azure Kubernetes Service clusters

## NAT Gateway Benefits

- **Eliminates SNAT Port Exhaustion**: Prevents connection failures from large-scale workloads
- **Provides Predictable Egress IPs**: Essential for firewall allow-listing on external services
- **Scales to 1M+ Ports**: Uses public IP prefixes for massive scale
- **Automatic Precedence**: Takes priority over Load Balancer outbound rules and instance-level public IPs

## Best Practices

1. **Address Space Planning**: Use non-overlapping CIDR blocks across environments
2. **Subnet Sizing**: For AKS, ensure subnet has enough IPs for max node count + pods (if using Azure CNI)
3. **NAT Gateway**: Enable for production clusters to avoid SNAT issues
4. **Tags**: Use consistent tagging for cost allocation and resource management
5. **Location**: Choose regions close to your users and other Azure services

## Variables Reference

See [variables.tf](./variables.tf) for the complete variable definitions.

## Outputs Reference

See [outputs.tf](./outputs.tf) for all available outputs.

