# AzurePrivateEndpoint: Research & Deployment Guide

## What is Azure Private Endpoint?

Azure Private Endpoint is a network interface that connects you privately and securely to a service powered by Azure Private Link. It uses a private IP address from your Virtual Network (VNet), effectively bringing the service into your VNet. The service could be an Azure PaaS service (Azure SQL, PostgreSQL, Storage, Key Vault, etc.) or a custom Private Link Service.

Private Endpoints enable three critical capabilities:

1. **Private connectivity** -- Access Azure PaaS services over a private IP instead of the public internet. Traffic stays entirely on the Microsoft backbone network, reducing latency and improving security.

2. **Data exfiltration protection** -- Each private endpoint maps to a specific sub-resource (e.g., "postgresqlServer", "vault", "blob"), not the entire service. Clients can only connect to the specific resource, preventing lateral data access to other resources in the same service account.

3. **Simplified network architecture** -- No need for service endpoints, NAT devices, or public IP addresses to reach Azure services from the VNet. Private endpoints eliminate the complexity of managing public endpoints and firewall rules.

## How Private Link Works

The Private Link flow involves several Azure resources working together:

```
VNet → Subnet → Private Endpoint → Network Interface → Private IP → Target Service
```

1. **VNet** -- The virtual network where your workloads reside
2. **Subnet** -- A dedicated subnet (typically named "pe-subnet" or "private-endpoints") where private endpoints are deployed
3. **Private Endpoint** -- The Azure resource that creates the connection
4. **Network Interface** -- Azure automatically creates a network interface (NIC) for each private endpoint
5. **Private IP** -- A private IP address allocated from the subnet's address space
6. **Target Service** -- The Azure PaaS service (PostgreSQL, Key Vault, Storage, etc.) being accessed privately

When a client in the VNet connects to the service's FQDN (e.g., `myserver.postgres.database.azure.com`), DNS resolution should point to the private IP address. This is where Private DNS Zones come in -- they ensure the FQDN resolves to the private IP instead of the public one.

## Deployment Landscape

### Manual (Azure Portal / CLI)

```bash
# Create private endpoint
az network private-endpoint create \
  --resource-group myRG \
  --name myPE \
  --vnet-name myVnet \
  --subnet pe-subnet \
  --private-connection-resource-id /subscriptions/.../Microsoft.DBforPostgreSQL/flexibleServers/myserver \
  --connection-name myConnection \
  --group-id postgresqlServer

# Create DNS zone group (optional)
az network private-endpoint dns-zone-group create \
  --resource-group myRG \
  --endpoint-name myPE \
  --name myZoneGroup \
  --private-dns-zone /subscriptions/.../privateDnsZones/privatelink.postgres.database.azure.com \
  --zone-name privatelink.postgres.database.azure.com
```

### Terraform

```hcl
resource "azurerm_private_endpoint" "postgres" {
  name                = "pg-private-endpoint"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  subnet_id           = azurerm_subnet.pe.id

  private_service_connection {
    name                           = "pg-connection"
    private_connection_resource_id = azurerm_postgresql_flexible_server.example.id
    subresource_names              = ["postgresqlServer"]
    is_manual_connection           = false
  }

  private_dns_zone_group {
    name                 = "pg-zone-group"
    private_dns_zone_ids = [azurerm_private_dns_zone.postgres.id]
  }
}
```

### Pulumi (Go)

```go
endpoint, _ := network.NewPrivateEndpoint(ctx, "pg-pe", &network.PrivateEndpointArgs{
    Name:                pulumi.String("pg-private-endpoint"),
    ResourceGroupName:   rg.Name,
    Location:            rg.Location,
    SubnetId:            subnet.ID(),
    PrivateServiceConnection: &network.PrivateEndpointPrivateServiceConnectionArgs{
        Name:                           pulumi.String("pg-connection"),
        PrivateConnectionResourceId:     postgresql.ID(),
        SubresourceNames:                pulumi.StringArray{pulumi.String("postgresqlServer")},
        IsManualConnection:              pulumi.Bool(false),
    },
    PrivateDnsZoneGroup: &network.PrivateEndpointPrivateDnsZoneGroupArgs{
        Name:                 pulumi.String("pg-zone-group"),
        PrivateDnsZoneIds:     pulumi.StringArray{zone.ID()},
    },
})
```

## Why the DNS Zone Group Is Part of the Endpoint

A private endpoint without DNS zone group registration does not resolve
correctly inside the VNet: the service FQDN resolves to the PUBLIC IP
instead of the private one, and clients silently bypass the private
endpoint entirely. Because that guarantee must be atomic, the DNS zone
group is part of this resource rather than a separate node -- the endpoint
and its A-record registration deploy together or not at all. The private
DNS zones themselves are first-class (`AzurePrivateDnsZone`) and referenced
here by id; `private_dns_zone_ids` accepts one or more so a single endpoint
can register into every zone a client population needs.

The zone group is optional only for the case where DNS is managed
externally (custom DNS servers, a hub-spoke forwarder). When present, every
referenced zone gets an A record for the endpoint's private IP.

## Field Mapping (azurerm → spec)

| azurerm | spec | Notes |
|---|---|---|
| `name` | `name` | Required, ForceNew |
| `location` | `region` | Required, ForceNew |
| `resource_group_name` | `resource_group` | FK → AzureResourceGroup |
| `subnet_id` | `subnet_id` | FK → AzureSubnet, ForceNew |
| `private_service_connection` (MaxItems 1) | `private_service_connection` message | Folded singular block |
| `.private_connection_resource_id` XOR `.private_connection_resource_alias` | `.private_connection_resource_id` / `.connection_alias` | ExactlyOneOf → message CEL |
| `.subresource_names` | `.subresource_names` | ForceNew |
| `.is_manual_connection` | `.is_manual_connection` | Required in azurerm; optional bool here (defaults false) |
| `.request_message` | `.request_message` | Paired with manual (CustomizeDiff → message CEL) |
| `private_dns_zone_group` (MaxItems 1) | `private_dns_zone_ids` | Group name auto-derived; the id list is the substance |
| `ip_configuration` | `ip_configurations` | Static IP assignments |
| `custom_network_interface_name` | `custom_network_interface_name` | ForceNew |
| ASG association (member-side resource) | `application_security_group_ids` | Realized as association resources both engines |
| `tags` | `tags` | User tags merged over Planton-derived tags |
| `private_service_connection.0.private_ip_address` (computed) | `private_ip_address` output | |
| `network_interface.0.id` (computed) | `network_interface_id` output | |

## Design Decisions

**Polymorphic connection target.** `private_connection_resource_id` carries
no `default_kind` because the target can be any Private Link-enabled service
(PostgreSQL, MySQL, Key Vault, Storage, Cosmos DB, Redis, SQL, ...).
Reference the service's own output id with an explicit `kind`/`fieldPath` in
composed environments. `connection_alias` is the alternative when the target
is exposed through a Private Link Service alias (a partner's cross-tenant
service); the two are mutually exclusive (message CEL).

**Manual-approval flow modeled.** `is_manual_connection` + `request_message`
support cross-tenant and cross-subscription connections that need the target
owner's approval. The provider's CustomizeDiff pairing -- request message
required iff manual -- is front-loaded as a message CEL so the error is
caught at validation, not apply.

**Static IP assignments modeled.** `ip_configurations` pins sub-resources to
fixed addresses when firewall allowlists or hard-coded DNS require it; leave
empty for dynamic allocation (the common case).

**ASG membership member-side.** Azure models application security group
membership from the member side, as its own association resource. The spec
exposes `application_security_group_ids`; both engines realize each as an
association resource.

**Auto-derived internal names.** The private service connection name and DNS
zone group name are internal handles Azure requires but nothing references;
both are derived from the endpoint name to keep the spec free of noise.

## Best Practices

1. **Dedicated subnet for private endpoints** -- Create a dedicated subnet (e.g., "pe-subnet") for all private endpoints. This simplifies network policy management and IP address planning.

2. **One endpoint per service instance** -- Each database, Key Vault, or Storage Account instance should have its own private endpoint. Don't try to share endpoints across instances.

3. **Always use DNS zone groups** -- Unless you have a specific reason to manage DNS externally, always provide `private_dns_zone_ids` to enable seamless DNS resolution.

4. **Match zone names exactly** -- For Private Link, the DNS zone name must exactly match Azure's predefined name for the service (e.g., `privatelink.postgres.database.azure.com`). A typo means DNS resolution fails silently.

5. **Sub-resource names are service-specific** -- Each Azure service defines its own sub-resource names. Use the correct name for your service type (see Common Sub-Resource Names table in README.md).

6. **Region must match subnet** -- The private endpoint's region must match the subnet's region. Subnets inherit their region from the parent VNet.

## Composition Seams

The endpoint consumes references on four sides: `subnet_id` →
`AzureSubnet`, `private_service_connection.private_connection_resource_id` →
any Private Link-enabled service (polymorphic), `private_dns_zone_ids` →
`AzurePrivateDnsZone`, and `application_security_group_ids` →
`AzureApplicationSecurityGroup`. Its outputs (`private_endpoint_id`,
`private_endpoint_name`, `private_ip_address`, `network_interface_id`) are
available for operational visibility and downstream references.

## Infra Chart Integration

### Database Stack Pattern

The database-stack infra chart creates one private endpoint per database instance:

```
VPC → Subnet → PrivateDnsZone → Database Server → PrivateEndpoint → DNS Zone Group
```

Each database server (PostgreSQL, MySQL, MSSQL, Redis) gets its own private endpoint. The endpoint is wired to:
1. The subnet (for private IP allocation)
2. The database server (via `private_connection_resource_id`)
3. The corresponding private DNS zone (via `private_dns_zone_ids`)

The DNS zone group automatically registers the private IP as an A-record in the zone, ensuring the database FQDN resolves to the private IP within the VNet.

### Enterprise Network Foundation

Optional component for organizations that pre-create private endpoints as part of their networking foundation, before any databases or services are deployed. However, private endpoints are typically created alongside their target services, not as standalone networking infrastructure.

---

**Status**: Production Ready
**Azure Provider Version**: ~> 4.0
**Pulumi Provider Version**: v6
