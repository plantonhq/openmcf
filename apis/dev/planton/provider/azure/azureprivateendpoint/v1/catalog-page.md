# Azure Private Endpoint

Deploys an Azure Private Endpoint that gives a Private Link-enabled service a private IP inside your virtual network, so traffic reaches the service over the Microsoft backbone instead of the public internet. The private DNS zone group is part of the resource: without it the service FQDN resolves to the public IP and clients silently bypass the private link.

## What Gets Created

When you deploy an AzurePrivateEndpoint resource, Planton provisions:

- **Private Endpoint** — an `azurerm_private_endpoint` with a network interface in the specified subnet and a single private service connection to the target
- **Private DNS Zone Group** (when `privateDnsZoneIds` is set) — registers the endpoint's private IP as an A record in each referenced private DNS zone, so in-VNet DNS resolves the service FQDN to the private IP
- **Application Security Group associations** (when `applicationSecurityGroupIds` is set) — one association per group, member-side

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **A subnet** to place the endpoint in (an `AzureSubnet` in composed environments)
- **A Private Link-enabled target** (PostgreSQL, Key Vault, Storage, Cosmos DB, ...) or a Private Link Service alias
- **Network write rights**: `Microsoft.Network/privateEndpoints/write` (Network Contributor, Contributor, or Owner)

## Quick Start

```yaml
apiVersion: azure.planton.dev/v1
kind: AzurePrivateEndpoint
metadata:
  name: pg-private-endpoint
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzurePrivateEndpoint.pg-private-endpoint
spec:
  region: eastus
  resourceGroup:
    value: network-rg
  name: pg-private-endpoint
  subnetId:
    value: /subscriptions/.../subnets/pe-subnet
  privateServiceConnection:
    privateConnectionResourceId:
      value: /subscriptions/.../Microsoft.DBforPostgreSQL/flexibleServers/prod-postgres
    subresourceNames:
      - postgresqlServer
  privateDnsZoneIds:
    - value: /subscriptions/.../privateDnsZones/privatelink.postgres.database.azure.com
```

Deploy:

```shell
planton apply -f private-endpoint.yaml
```

After deployment, read `status.outputs.private_ip_address` for the IP the service FQDN now resolves to inside the VNet.

## Configuration Reference

### Required Fields

| Field | Type | Description |
|-------|------|-------------|
| `region` | `string` | Azure region; must match the subnet. Fixed at creation. |
| `resourceGroup` | `StringValueOrRef` | Resource group. Defaults to an `AzureResourceGroup` reference. |
| `name` | `string` | Endpoint name, unique in the resource group. Fixed at creation. |
| `subnetId` | `StringValueOrRef` | Subnet the private IP is drawn from. Defaults to an `AzureSubnet` reference. Fixed at creation. |
| `privateServiceConnection` | `object` | The private link connection (target + sub-resource). |

### `privateServiceConnection`

| Field | Type | Description |
|-------|------|-------------|
| `privateConnectionResourceId` | `StringValueOrRef` | Target service by ARM ID (polymorphic). Set this OR `connectionAlias`. |
| `connectionAlias` | `string` | Private Link Service alias (partner cross-tenant). Set this OR `privateConnectionResourceId`. |
| `subresourceNames` | `string[]` | Sub-resource / group ID (e.g. `postgresqlServer`, `blob`, `vault`). |
| `isManualConnection` | `bool` | Requires owner approval (default `false`). |
| `requestMessage` | `string` | Approval message; required iff `isManualConnection` is true (1-140 chars). |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `privateDnsZoneIds` | `StringValueOrRef[]` | `[]` | Private DNS zones to register the IP into. Defaults to `AzurePrivateDnsZone` references. When empty, no DNS zone group is created. |
| `ipConfigurations` | `object[]` | `[]` | Static IP assignments (`name`, `privateIpAddress`, `subresourceName`, `memberName`); empty means dynamic allocation. |
| `applicationSecurityGroupIds` | `StringValueOrRef[]` | `[]` | ASGs the endpoint's NIC joins. Defaults to `AzureApplicationSecurityGroup` references. |
| `customNetworkInterfaceName` | `string` | `""` | Custom name for the auto-created NIC. Fixed at creation. |
| `tags` | `map(string)` | `{}` | User tags, merged over Planton-derived tags (user wins). |

## Examples

### Storage Account (Blob)

```yaml
apiVersion: azure.planton.dev/v1
kind: AzurePrivateEndpoint
metadata:
  name: blob-private-endpoint
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      name: data-rg
  name: blob-private-endpoint
  subnetId:
    valueFrom:
      name: pe-subnet
  privateServiceConnection:
    privateConnectionResourceId:
      valueFrom:
        kind: AzureStorageAccount
        name: prod-storage
        fieldPath: status.outputs.storage_account_id
    subresourceNames:
      - blob
  privateDnsZoneIds:
    - valueFrom:
        name: blob-privatelink-zone
```

### Cross-Tenant via Private Link Service Alias (Manual Approval)

```yaml
spec:
  privateServiceConnection:
    connectionAlias: partner-pls.d20286c8-4ea5-11eb-9584.centralus.azure.privatelinkservice
    isManualConnection: true
    requestMessage: "please approve access for the analytics workload"
```

## Stack Outputs

| Output | Type | Description |
|--------|------|-------------|
| `private_endpoint_id` | `string` | Full ARM ID of the endpoint |
| `private_endpoint_name` | `string` | The endpoint's name as deployed |
| `private_ip_address` | `string` | The private IP allocated from the subnet |
| `network_interface_id` | `string` | ARM ID of the auto-created network interface |

## Related Components

- [AzureSubnet](/docs/catalog/azure/subnet) — provides the subnet the endpoint draws its IP from
- [AzurePrivateDnsZone](/docs/catalog/azure/private-dns-zone) — the privatelink zones the endpoint registers into
- [AzureApplicationSecurityGroup](/docs/catalog/azure/application-security-group) — workload groups the endpoint's NIC joins
- [AzureResourceGroup](/docs/catalog/azure/azureresourcegroup) — provides the resource group for placement
