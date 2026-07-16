# Production Private Server with Zone-Redundant HA

This preset creates the production baseline: a General Purpose server
injected into a delegated subnet (no public endpoint), a zone-redundant
standby with automatic failover, geo-redundant 35-day backups, storage
auto-grow, and a Sunday 03:00 maintenance window.

## When to Use

- Production workloads that must survive an availability-zone failure
- Any environment whose security posture forbids a public database
  endpoint

## Key Configuration Choices

- **VNet injection over private endpoints** -- the server itself lives on
  the delegated subnet; requires public access explicitly OFF and a
  private DNS zone (conventionally `privatelink.postgres.database.azure.com`
  or a `*.postgres.database.azure.com` zone) so the fqdn resolves privately
- **`ZONE_REDUNDANT` HA** -- the standby replicates synchronously in zone
  "2"; failover is automatic; the standby zone is fixed at creation
- **`geo_redundant_backup_enabled`** -- unlocks cross-region `GEO_RESTORE`
  disaster recovery; fixed at creation, so decide before the server exists
- **Storage auto-grow** -- steps up the storage ladder before the disk
  fills; storage never shrinks

## Prerequisites

The delegated subnet must carry the `Microsoft.DBforPostgreSQL/flexibleServers`
delegation and hold no other resources:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureSubnet
metadata:
  name: database-subnet
spec:
  virtualNetworkId:
    valueFrom:
      kind: AzureVirtualNetwork
      name: prod-network
      fieldPath: status.outputs.virtual_network_id
  addressPrefixes:
    - 10.0.20.0/24
  delegations:
    - name: postgres-delegation
      serviceName: Microsoft.DBforPostgreSQL/flexibleServers
      actions:
        - Microsoft.Network/virtualNetworks/subnets/join/action
```

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<resource-group-name>` | The resource group to create the server in | The resource group's `status.outputs.resource_group_name` |
| `myorg-prod-postgres` | 3-63 lowercase chars, globally unique | Your naming convention |
| `<admin-login>` / `<admin-password>` | PostgreSQL admin credentials | A secret manager; never commit literals |
| `<database-subnet-resource-name>` | The delegated AzureSubnet resource | Your network manifests |
| `<postgres-dns-zone-resource-name>` | The AzurePrivateDnsZone resource | Your network manifests |
| `appdb` | The application's database | Your application |
| `<cost-center>` | Your org's cost-attribution tag value | Your tagging convention |
