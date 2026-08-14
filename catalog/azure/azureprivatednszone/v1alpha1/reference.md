# AzurePrivateDnsZone

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**AzurePrivateDnsZoneSpec** defines the configuration for creating an
Azure Private DNS zone: name resolution inside virtual networks without
running a DNS server. Private DNS zones serve two primary scenarios:

1. **Private Link DNS resolution** -- zones like
   "privatelink.postgres.database.azure.com" enable automatic DNS
   resolution for Azure Private Endpoints. When a private endpoint is
   created for an Azure PaaS service, its private IP is registered in the
   corresponding privatelink zone, so VNet-connected clients resolve the
   service's FQDN to its private IP instead of its public one.

2. **Custom internal DNS** -- zones like "corp.internal" provide internal
   name resolution for VMs, containers, and other resources.

The zone is deliberately just the zone -- a global record container with
no reach of its own. Which networks can resolve it is declared through
AzurePrivateDnsZoneVirtualNetworkLink resources referencing this zone's
zone_id output: one link per network, added and removed without touching
the zone. A zone with no links answers nobody, so every deployment pairs
the zone with at least one link.

**Note:** Private DNS zones are global Azure resources -- they have no
region. This is why this spec omits the `region` field that other Azure
resources include.

## Example

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzurePrivateDnsZone
metadata:
  name: pg-private-dns
  org: test-org
  env: dev
spec:
  resourceGroup:
    value: test-rg
  name: privatelink.postgres.database.azure.com
  tags:
    cost-center: platform
    owner: network-team
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.resourceGroup` | `string \| valueFrom` | yes |  | AzureResourceGroup (`status.outputs.resource_group_name`) |
| `spec.name` | `string` | yes |  |  |
| `spec.soaRecord` | `AzurePrivateDnsZoneSoaRecord` |  |  |  |
| `spec.soaRecord.email` | `string` | yes |  |  |
| `spec.soaRecord.expireTime` | `int64` |  | `2419200` |  |
| `spec.soaRecord.minimumTtl` | `int64` |  | `10` |  |
| `spec.soaRecord.refreshTime` | `int64` |  | `3600` |  |
| `spec.soaRecord.retryTime` | `int64` |  | `300` |  |
| `spec.soaRecord.ttl` | `int64` |  | `3600` |  |
| `spec.soaRecord.tags` | `map<string, string>` |  |  |  |
| `spec.tags` | `map<string, string>` |  |  |  |

## Field Details

### spec.resourceGroup

`string | valueFrom` · required

The Azure resource group the private DNS zone will be created in.
Can be a literal resource-group name or a reference to an
AzureResourceGroup's name output.

- references: AzureResourceGroup (`status.outputs.resource_group_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureResourceGroup, name: <that resource's name>, fieldPath: status.outputs.resource_group_name}} -- a bare string does not parse

### spec.name

`string` · required

The name of the private DNS zone. Must be a valid DNS domain name;
changing it replaces the zone and every record in it. For Private Link
scenarios, this must match the Azure-defined privatelink zone name for
the target service.

Common Private Link zone names:
- "privatelink.postgres.database.azure.com" -- PostgreSQL Flexible Server
- "privatelink.mysql.database.azure.com" -- MySQL Flexible Server
- "privatelink.database.windows.net" -- Azure SQL Database
- "privatelink.documents.azure.com" -- Cosmos DB
- "privatelink.redis.cache.windows.net" -- Azure Cache for Redis
- "privatelink.blob.core.windows.net" -- Azure Blob Storage
- "privatelink.vaultcore.azure.net" -- Azure Key Vault

For custom internal DNS, use any valid domain (e.g., "corp.internal").

- rule: Zone name must be a valid DNS domain (e.g., privatelink.postgres.database.azure.com or corp.internal)
- rule: {"required":true}

### spec.soaRecord

`AzurePrivateDnsZoneSoaRecord`

Customize the zone's Start of Authority record. Nearly every
deployment leaves this unset and takes Azure's defaults; set it only
when operational tooling requires a specific contact email or
negative-caching behavior. The SOA record is created with the zone and
cannot be customized afterwards (changing this replaces the zone).

### spec.soaRecord.email

`string` · required

The email contact for the zone, in SOA host format (dots instead of
"@", e.g. "azureprivatedns-host.microsoft.com" or
"dnsadmin.contoso.com" for dnsadmin@contoso.com).

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.soaRecord.expireTime

`int64` · optional (explicit presence)

Seconds a secondary considers the zone valid without a successful
refresh. Azure's default is 2419200 (28 days).

- default: `2419200`
- rule: {"int64":{"gte":"0"}}

### spec.soaRecord.minimumTtl

`int64` · optional (explicit presence)

Negative-caching TTL: how long resolvers cache a "name does not
exist" answer, in seconds. Azure's default is 10. The one field with
real operational impact -- raising it slows how quickly newly-created
records become visible to clients that asked before the record
existed.

- default: `10`
- rule: {"int64":{"gte":"0"}}

### spec.soaRecord.refreshTime

`int64` · optional (explicit presence)

Seconds between secondary refresh attempts. Azure's default is 3600.

- default: `3600`
- rule: {"int64":{"gte":"0"}}

### spec.soaRecord.retryTime

`int64` · optional (explicit presence)

Seconds a secondary waits before retrying a failed refresh. Azure's
default is 300.

- default: `300`
- rule: {"int64":{"gte":"0"}}

### spec.soaRecord.ttl

`int64` · optional (explicit presence)

TTL of the SOA record itself, in seconds. Azure's default is 3600.

- default: `3600`
- rule: {"int64":{"lte":"2147483647","gte":"0"}}

### spec.soaRecord.tags

`map<string, string>`

Tags applied to the SOA record set (distinct from the zone's own
tags). Rarely used; exists because ARM models record-set tags.

### spec.tags

`map<string, string>`

Free-form tags applied to the zone, merged over the Planton-derived
resource tags (organization, environment, resource id); a user tag
with the same key wins. Tags are Azure's governance surface -- Azure
Policy enforces them and Microsoft Cost Management groups by them --
so carry your org's ownership/cost-center conventions here. Updatable
in place.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzurePrivateDnsZone, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.zone_id` | `string` | The Azure Resource Manager ID of the private DNS zone. Format: /subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.Network/privateDnsZones/{name} This is the primary output referenced by downstream resources via StringValueOrRef. |
| `status.outputs.zone_name` | `string` | The name of the private DNS zone (e.g., "privatelink.postgres.database.azure.com"). Echoed from the spec for convenience -- useful in IaC modules that need the zone name for creating DNS records. |
| `status.outputs.resource_group_name` | `string` | The resource group the zone lives in. Echoed for downstream tooling that addresses records or links by zone name + resource group rather than parsing the ARM ID. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.resourceGroup` | AzureResourceGroup | `status.outputs.resource_group_name` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AzureAksCluster | `spec.privateDnsZoneId` | `status.outputs.zone_id` |
| AzureMysqlFlexibleServer | `spec.privateDnsZoneId` | `status.outputs.zone_id` |
| AzurePostgresqlFlexibleServer | `spec.privateDnsZoneId` | `status.outputs.zone_id` |
| AzurePrivateDnsRecord | `spec.privateDnsZoneId` | `status.outputs.zone_id` |
| AzurePrivateDnsZoneVirtualNetworkLink | `spec.privateDnsZoneId` | `status.outputs.zone_id` |
| AzurePrivateEndpoint | `spec.privateDnsZoneIds` | `status.outputs.zone_id` |

## See Also

- [Overview](../README.md)
