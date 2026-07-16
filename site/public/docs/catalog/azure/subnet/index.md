---
title: "Subnet"
description: "Subnet deployment documentation"
icon: "package"
order: 100
componentName: "azuresubnet"
---

# Azure Subnet

Deploys an Azure Subnet within an existing Virtual Network, with configurable address prefixes (or IPAM-delegated allocation), service endpoints, service delegations, private endpoint network policies, and optional route table, network security group, and NAT gateway attachments. Subnets partition a VNet's address space into segments for different workloads, tiers, or service delegations.

## What Gets Created

When you deploy an AzureSubnet resource, Planton provisions:

- **Subnet** — a `network.Subnet` resource inside the referenced Virtual Network, configured with the given address prefixes (or an IPAM pool allocation), service endpoints, delegations, and network policies
- **Service Endpoints** — optimized routes over the Azure backbone to specified Azure services, bypassing the public internet (when `serviceEndpoints` is provided)
- **Service Delegations** — grants an Azure PaaS service permission to inject service-specific resources and network rules into the subnet (when `delegations` is provided)
- **Attachments** — associations to a route table, network security group, and/or NAT gateway (when `routeTableId`, `networkSecurityGroupId`, or `natGatewayId` is provided)

Subnets are not tracked ARM resources, so they carry no tags of their own.

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **An Azure Virtual Network** with an address space that contains the desired subnet CIDR blocks (can reference an AzureVirtualNetwork resource). The resource group and VNet name are derived from the VNet's ARM ID — there is no separate resource group field.
- **Network planning** — every subnet address prefix must be a subset of the parent VNet's address space and must not overlap with other subnets in the same VNet. Azure reserves 5 IPs per subnet (first 4 + last) for internal use.
- **Optional attachments** — a route table, network security group, or NAT gateway to attach (can reference AzureRouteTable, AzureNetworkSecurityGroup, or AzureNatGateway resources)

## Quick Start

Create a file `subnet.yaml`:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureSubnet
metadata:
  name: my-subnet
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureSubnet.my-subnet
spec:
  virtualNetworkId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-rg/providers/Microsoft.Network/virtualNetworks/my-vnet
  name: my-subnet
  addressPrefixes:
    - "10.0.1.0/24"
```

Deploy:

```shell
planton apply -f subnet.yaml
```

This creates a /24 subnet (251 usable IPs) with private endpoint network policies disabled, private link service network policies enabled, default outbound access enabled, and no service endpoints, delegations, or attachments.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `virtualNetworkId` | `StringValueOrRef` | Azure Resource Manager ID of the parent Virtual Network. The subnet is created inside this VNet; the resource group and VNet name are derived from the ARM ID. Can reference an AzureVirtualNetwork resource via `valueFrom`. Changing the parent replaces the subnet. | Required |
| `name` | `string` | Name of the subnet. Must be unique within the VNet. Allowed characters: alphanumerics, underscores, periods, and hyphens; must start with a letter or number and end with a letter, number, or underscore. Changing the name replaces the subnet. | Required, 1–80 characters |
| `addressPrefixes` | `string[]` | CIDR blocks for the subnet (e.g., `["10.0.1.0/24"]`). Every block must be a subset of the parent VNet's address space and must not overlap with other subnets. Multiple blocks support dual-stack (IPv4 + IPv6) subnets. | Exactly one of `addressPrefixes` or `ipAddressPool` must be set |
| `ipAddressPool` | `object` | Delegated address allocation from an Azure Network Manager IPAM pool, as an alternative to hand-planned prefixes. `id` is the pool's ARM ID; `numberOfIpAddresses` is a positive number in string form (e.g., `"256"`). The provisioned CIDR surfaces in the `address_prefixes` output. | Exactly one of `addressPrefixes` or `ipAddressPool` must be set |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `serviceEndpoints` | `string[]` | `[]` | Azure service endpoints to enable. Creates optimized routes over the Azure backbone. Common values: `Microsoft.Storage`, `Microsoft.Sql`, `Microsoft.KeyVault`, `Microsoft.AzureCosmosDB`, `Microsoft.ServiceBus`, `Microsoft.EventHub`, `Microsoft.Web`, `Microsoft.ContainerRegistry`. |
| `serviceEndpointPolicyIds` | `string[]` | `[]` | ARM IDs of Service Endpoint Policies to associate with the subnet. Endpoint policies narrow a service endpoint's reach to specific resources (e.g., only your storage accounts instead of all of Azure Storage). |
| `delegations` | `object[]` | `[]` | Service delegations granting an Azure PaaS service permission to inject resources into the subnet. Nearly every real subnet has zero or one; the list form mirrors ARM's shape. See delegation fields below. |
| `delegations[].name` | `string` | — | A user-chosen label for the delegation (e.g., `postgresql`, `container-apps`). Required per delegation. |
| `delegations[].serviceName` | `string` | — | The Azure service to delegate to. Required per delegation. Common values: `Microsoft.DBforPostgreSQL/flexibleServers`, `Microsoft.DBforMySQL/flexibleServers`, `Microsoft.App/environments`, `Microsoft.Web/serverFarms`, `Microsoft.ContainerInstance/containerGroups`, `Microsoft.Netapp/volumes`. |
| `delegations[].actions` | `string[]` | `[]` | Actions the delegated service is permitted to perform. If omitted, Azure uses the service's default action set. Common action: `Microsoft.Network/virtualNetworks/subnets/join/action`. |
| `privateEndpointNetworkPolicies` | `enum` | unset (Azure's `Disabled`) | Controls whether network policies apply to private endpoints in the subnet. Leave unset for Azure's default (no policy evaluation — what most private-endpoint subnets want). Values: `ENABLED` (both NSG and route table), `NETWORK_SECURITY_GROUP_ENABLED` (NSG only), `ROUTE_TABLE_ENABLED` (route table only). |
| `privateLinkServiceNetworkPoliciesEnabled` | `bool` | `true` | Controls whether network policies apply to Private Link Service resources. Set to `false` only on subnets hosting a Private Link Service, which requires policies off. |
| `defaultOutboundAccessEnabled` | `bool` | `true` | Whether workloads get Azure's implicit default outbound internet access. Microsoft is retiring implicit outbound; production subnets should set this `false` and route egress explicitly through a NAT gateway, load balancer outbound rules, or a firewall via a route table. |
| `sharingScope` | `enum` | unset (not shared) | Opt the subnet into cross-tenant sharing through Azure Virtual Network Manager. `TENANT` is the only accepted value, and it requires `defaultOutboundAccessEnabled` to be explicitly `false`. |
| `routeTableId` | `StringValueOrRef` | none | ARM ID of a route table to attach, replacing Azure's default system routing with the table's user-defined routes. Can reference an AzureRouteTable resource via `valueFrom`. Detach by removing the field. |
| `networkSecurityGroupId` | `StringValueOrRef` | none | ARM ID of a network security group to attach; its rules filter all traffic in the subnet. Can reference an AzureNetworkSecurityGroup resource via `valueFrom`. |
| `natGatewayId` | `StringValueOrRef` | none | ARM ID of a NAT gateway to attach; it owns the subnet's outbound connectivity through the gateway's public IPs. Can reference an AzureNatGateway resource via `valueFrom`. |

## Examples

### General-Purpose Workload Subnet

A /24 subnet for general workloads with no special endpoints or delegations:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureSubnet
metadata:
  name: workload-subnet
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureSubnet.workload-subnet
spec:
  virtualNetworkId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/dev-rg/providers/Microsoft.Network/virtualNetworks/dev-vnet
  name: workload-subnet
  addressPrefixes:
    - "10.0.1.0/24"
```

### Database Subnet with Delegation and Service Endpoints

A subnet delegated to PostgreSQL Flexible Server with service endpoints for secure access to Storage and Key Vault:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureSubnet
metadata:
  name: postgres-subnet
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureSubnet.postgres-subnet
spec:
  virtualNetworkId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/prod-rg/providers/Microsoft.Network/virtualNetworks/prod-vnet
  name: postgres-subnet
  addressPrefixes:
    - "10.0.10.0/24"
  serviceEndpoints:
    - Microsoft.Storage
    - Microsoft.KeyVault
  delegations:
    - name: postgresql
      serviceName: Microsoft.DBforPostgreSQL/flexibleServers
      actions:
        - Microsoft.Network/virtualNetworks/subnets/join/action
```

### Private Endpoint Subnet with Network Policies

A subnet for private endpoints with NSG policies enabled for zero-trust architectures:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureSubnet
metadata:
  name: pe-subnet
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureSubnet.pe-subnet
spec:
  virtualNetworkId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/prod-rg/providers/Microsoft.Network/virtualNetworks/prod-vnet
  name: pe-subnet
  addressPrefixes:
    - "10.0.20.0/28"
  privateEndpointNetworkPolicies: NETWORK_SECURITY_GROUP_ENABLED
  privateLinkServiceNetworkPoliciesEnabled: true
```

### Using Foreign Key References and Attachments

Reference Planton-managed resources instead of hardcoding ARM IDs — for the parent network and for the route table, NSG, and NAT gateway attachments:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureSubnet
metadata:
  name: app-subnet
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureSubnet.app-subnet
spec:
  virtualNetworkId:
    valueFrom:
      name: my-vnet
  name: app-subnet
  addressPrefixes:
    - "10.0.2.0/24"
  serviceEndpoints:
    - Microsoft.Sql
    - Microsoft.Storage
    - Microsoft.KeyVault
    - Microsoft.Web
  defaultOutboundAccessEnabled: false
  routeTableId:
    valueFrom:
      name: prod-egress-rt
  networkSecurityGroupId:
    valueFrom:
      name: app-tier-nsg
  natGatewayId:
    valueFrom:
      name: prod-egress-nat
```

Each reference resolves through the field's default kind (AzureVirtualNetwork, AzureRouteTable, AzureNetworkSecurityGroup, AzureNatGateway) against that resource's outputs.

### Container App Environment Subnet

A subnet delegated to Azure Container App Environments with the minimum /23 sizing recommended by Azure:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureSubnet
metadata:
  name: cae-subnet
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureSubnet.cae-subnet
spec:
  virtualNetworkId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/prod-rg/providers/Microsoft.Network/virtualNetworks/prod-vnet
  name: cae-subnet
  addressPrefixes:
    - "10.0.32.0/23"
  delegations:
    - name: container-apps
      serviceName: Microsoft.App/environments
      actions:
        - Microsoft.Network/virtualNetworks/subnets/join/action
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `subnet_id` | `string` | Azure Resource Manager ID of the subnet. This is the most referenced Azure output in Planton, consumed by AzureAksCluster, AzureContainerAppEnvironment, AzurePostgresqlFlexibleServer, AzureMysqlFlexibleServer, AzureRedisCache, AzurePrivateEndpoint, AzureApplicationGateway, AzureLoadBalancer, AzureVirtualMachine, AzureFunctionApp, and AzureLinuxWebApp. |
| `subnet_name` | `string` | Name of the subnet within the VNet |
| `address_prefixes` | `string[]` | CIDR blocks actually assigned to the subnet. For self-managed subnets this echoes `addressPrefixes`; for IPAM-allocated subnets it carries the ranges the Network Manager pool provisioned. Useful in NSG rules, firewall rules, and network planning downstream. |
| `virtual_network_name` | `string` | Name of the parent Virtual Network, derived from the referenced network's ARM ID |
| `resource_group_name` | `string` | Name of the resource group the subnet (and its parent network) lives in, derived from the referenced network's ARM ID |

## Related Components

- [AzureVirtualNetwork](/docs/catalog/azure/virtual-network) — provides the parent Virtual Network that contains this subnet
- [AzureRouteTable](/docs/catalog/azure/route-table) — attaches via `routeTableId` to steer the subnet's egress
- [AzureNetworkSecurityGroup](/docs/catalog/azure/network-security-group) — attaches via `networkSecurityGroupId` to filter the subnet's traffic
- [AzureNatGateway](/docs/catalog/azure/nat-gateway) — attaches via `natGatewayId` to own the subnet's outbound connectivity
- [AzureAksCluster](/docs/catalog/azure/aks-cluster) — references `subnet_id` for node pool placement
- [AzureContainerAppEnvironment](/docs/catalog/azure/container-app-environment) — requires a delegated subnet for VNet integration
- [AzurePostgresqlFlexibleServer](/docs/catalog/azure/postgresql-flexible-server) — requires a delegated subnet for VNet integration
- [AzureMysqlFlexibleServer](/docs/catalog/azure/mysql-flexible-server) — requires a delegated subnet for VNet integration
- [AzurePrivateEndpoint](/docs/catalog/azure/private-endpoint) — deployed into a subnet for private connectivity to Azure PaaS services
- [AzureApplicationGateway](/docs/catalog/azure/application-gateway) — requires a dedicated subnet (minimum /27)
- [AzureKeyVault](/docs/catalog/azure/key-vault) — can restrict access to specific subnet IDs via network ACLs
