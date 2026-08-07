# Azure PostgreSQL Flexible Server

Creates an Azure Database for PostgreSQL Flexible Server -- Azure's managed PostgreSQL with per-server compute and storage sizing, zone-redundant high availability, Microsoft Entra authentication, customer-managed-key encryption, read replicas, and point-in-time restore. Databases, firewall rules, server parameters, and Entra administrators are declared on the server and managed with it.

## What Gets Created

When you deploy an AzurePostgresqlFlexibleServer resource, Planton provisions:

- **PostgreSQL Flexible Server** -- an `azurerm_postgresql_flexible_server` in the specified region and resource group, with your chosen compute SKU, storage size and performance tier, version, availability posture, networking, authentication, and encryption
- **Databases** -- an `azurerm_postgresql_flexible_server_database` for each entry in `databases`, each with its own charset and collation
- **Firewall Rules** -- an `azurerm_postgresql_flexible_server_firewall_rule` for each entry in `firewallRules`, allowlisting IPv4 ranges on the public endpoint
- **Server Parameters** -- an `azurerm_postgresql_flexible_server_configuration` for each `serverParameters` entry, applied as user overrides on Azure's per-SKU defaults
- **Entra Administrators** -- an `azurerm_postgresql_flexible_server_active_directory_administrator` for each `aadAdministrators` entry, granting the principal the server's Entra admin role

A read replica or a restored server is simply another AzurePostgresqlFlexibleServer whose `createMode` and `sourceServerId` reference the source.

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **A resource group** to create the server in (an `AzureResourceGroup` in composed environments)
- **For VNet injection**: a subnet delegated to `Microsoft.DBforPostgreSQL/flexibleServers` with no other resources, plus a private DNS zone for the server's name
- **For CMK encryption**: a user-assigned identity with wrap/unwrap access on the Key Vault key, granted before the server is created

## Quick Start

Create a file `postgresql.yaml`:

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzurePostgresqlFlexibleServer
metadata:
  name: my-pg
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzurePostgresqlFlexibleServer.my-pg
spec:
  region: eastus
  resourceGroup:
    value: my-rg
  serverName: myorg-dev-pg
  administratorLogin: pgadmin
  administratorPassword:
    value: "Ch@ngeMe1234!"
  skuName: B_Standard_B1ms
  databases:
    - name: myapp
```

Deploy:

```shell
planton apply -f postgresql.yaml
```

This creates a Burstable PostgreSQL 16 server with Azure's default 32 GiB storage and 7-day backups, public access, and one application database. Read `status.outputs.fqdn` to build the connection string.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `region` | `string` | Azure region. Changing it replaces the server. | Required |
| `resourceGroup` | `StringValueOrRef` | Resource group name. Defaults to referencing an `AzureResourceGroup`'s name output. | Required |
| `serverName` | `string` | GLOBALLY unique -- becomes `{name}.postgres.database.azure.com`. Changing it replaces the server. | Required, 3-63 lowercase letters/digits/hyphens |
| `skuName` | `string` | Compute SKU, `{TIER}_Standard_{SIZE}` (`B_`/`GP_`/`MO_`). Required for a fresh server; a replica left unset inherits the source's. | Pattern-validated |

`administratorLogin` + `administratorPassword` are required for a fresh server with password auth enabled, and must be OMITTED when password auth is disabled (Entra-only).

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `createMode` | `enum` | `DEFAULT` | `DEFAULT`, `REPLICA`, `POINT_IN_TIME_RESTORE`, `GEO_RESTORE`, `REVIVE_DROPPED`. Fixed at creation. |
| `sourceServerId` | `StringValueOrRef` | -- | The source server for replica/restore modes; references another server's `server_id` output. |
| `pointInTimeRestoreTimeInUtc` | `string` | -- | RFC-3339 restore instant (restore modes only). |
| `replicationRole` | `enum` | -- | `NONE` promotes a replica to a standalone primary (irreversible, day-2 only). |
| `version` | `string` | `"16"` | PostgreSQL major version, `"11"`-`"18"`. Upgrades go up only; elastic clusters need 17+. |
| `storageMb` | `int32` | `32768` | Fixed ladder from 32768 (32 GiB) to 33553408 (32 TiB). Grows only. |
| `storageTier` | `enum` | size default | `P4`-`P80` IOPS class, validated against the size's supported range. |
| `autoGrowEnabled` | `bool` | `false` | Automatic storage growth when free space runs low. |
| `zone` | `string` | Azure picks | Primary availability zone (`"1"`/`"2"`/`"3"`). |
| `highAvailability` | `object` | -- | `mode` (`ZONE_REDUNDANT`/`SAME_ZONE`) + optional `standbyAvailabilityZone`. Not supported on Burstable SKUs. |
| `maintenanceWindow` | `object` | system-managed | Weekly patching window (`dayOfWeek`/`startHour`/`startMinute`). |
| `backupRetentionDays` | `int32` | `7` | 7-35 days -- the point-in-time restore horizon. |
| `geoRedundantBackupEnabled` | `bool` | `false` | Replicates backups to the paired region (enables `GEO_RESTORE`). Fixed at creation. |
| `publicNetworkAccessEnabled` | `bool` | `true` | The public endpoint. Must be explicitly `false` when VNet-injected. |
| `delegatedSubnetId` | `StringValueOrRef` | -- | VNet injection: a delegated subnet (references `AzureSubnet`). Requires `privateDnsZoneId`. Fixed at creation. |
| `privateDnsZoneId` | `StringValueOrRef` | -- | The private DNS zone the injected server registers in (references `AzurePrivateDnsZone`). |
| `authentication` | `object` | password on | `activeDirectoryAuthEnabled`, `passwordAuthEnabled`, optional `tenantId` (defaults to the deploying credential's tenant). |
| `aadAdministrators` | `list` | `[]` | Entra admin grants: `objectId` (references a user-assigned identity's `principal_id`), `principalName`, `principalType` (`USER`/`GROUP`/`SERVICE_PRINCIPAL`). |
| `identity` | `object` | -- | `SYSTEM_ASSIGNED`, `USER_ASSIGNED`, or `SYSTEM_AND_USER_ASSIGNED` + identity references. |
| `customerManagedKey` | `object` | -- | CMK encryption: `keyVaultKeyId` (references `AzureKeyVaultKey.versionless_id`), the unwrap identity, and an optional geo-backup key pair. Fixed at creation. |
| `cluster` | `object` | -- | Elastic cluster (PG 17+): `size` 1-20 nodes + optional `defaultDatabaseName`. |
| `databases` | `list` | `[]` | Databases (`name`, `charset` default `UTF8`, `collation` default `en_US.utf8`). |
| `firewallRules` | `list` | `[]` | Public-endpoint IPv4 allowlist (`0.0.0.0`-`0.0.0.0` admits Azure-internal services only). |
| `serverParameters` | `map(string)` | `{}` | PostgreSQL parameter overrides by name (e.g. `azure.extensions`, `max_connections`). Static parameters need a restart. |
| `tags` | `map(string)` | `{}` | User tags, merged over Planton-derived tags (user wins on collision). |

## Examples

### Production: VNet-Injected with Zone-Redundant HA

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzurePostgresqlFlexibleServer
metadata:
  name: prod-pg
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      name: data-rg
  serverName: myorg-prod-pg
  administratorLogin: pgadmin
  administratorPassword:
    value: "Ch@ngeMe1234!"
  skuName: GP_Standard_D4ds_v5
  storageMb: 262144
  autoGrowEnabled: true
  zone: "1"
  highAvailability:
    mode: ZONE_REDUNDANT
    standbyAvailabilityZone: "2"
  backupRetentionDays: 35
  geoRedundantBackupEnabled: true
  publicNetworkAccessEnabled: false
  delegatedSubnetId:
    valueFrom:
      name: database-subnet
  privateDnsZoneId:
    valueFrom:
      name: postgres-dns-zone
  databases:
    - name: orders
```

### Read Replica of a Primary

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzurePostgresqlFlexibleServer
metadata:
  name: prod-pg-replica
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      name: data-rg
  serverName: myorg-prod-pg-replica
  createMode: REPLICA
  sourceServerId:
    valueFrom:
      kind: AzurePostgresqlFlexibleServer
      name: prod-pg
      fieldPath: status.outputs.server_id
```

The replica inherits the source's SKU, storage, and version. Promote it later by setting `replicationRole: NONE`.

### Entra-Only with Customer-Managed Key

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzurePostgresqlFlexibleServer
metadata:
  name: hardened-pg
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      name: data-rg
  serverName: myorg-hardened-pg
  skuName: GP_Standard_D4ds_v5
  authentication:
    activeDirectoryAuthEnabled: true
    passwordAuthEnabled: false
  aadAdministrators:
    - objectId:
        value: "11111111-2222-3333-4444-555555555555"
      principalName: dba-team
      principalType: GROUP
  identity:
    type: USER_ASSIGNED
    identityIds:
      - valueFrom:
          name: pg-cmk-identity
  customerManagedKey:
    keyVaultKeyId:
      valueFrom:
        kind: AzureKeyVaultKey
        name: pg-cmk
        fieldPath: status.outputs.versionless_id
    primaryUserAssignedIdentityId:
      valueFrom:
        name: pg-cmk-identity
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `server_id` | `string` | The server's ARM ID -- referenced by `AzurePrivateEndpoint` and by replica/restore servers' `sourceServerId` |
| `server_name` | `string` | The server's name |
| `fqdn` | `string` | `{name}.postgres.database.azure.com` -- the connection-string host; resolves privately when VNet-injected |
| `administrator_login` | `string` | The admin login (empty on Entra-only servers) |
| `database_ids` | `map(string)` | Each declared database's ARM ID, keyed by name |
| `identity_principal_id` | `string` | The system-assigned identity's principal ID -- the `AzureRoleAssignment` seam |

## Related Components

- [AzureResourceGroup](/docs/catalog/azure/azureresourcegroup) — provides the resource group for server placement
- [AzureSubnet](/docs/catalog/azure/subnet) — the delegated subnet for VNet injection
- [AzurePrivateDnsZone](/docs/catalog/azure/private-dns-zone) — private name resolution for the injected server
- [AzurePrivateEndpoint](/docs/catalog/azure/azureprivateendpoint) — Private Link connectivity as an alternative to VNet injection
- [AzureUserAssignedIdentity](/docs/catalog/azure/user-assigned-identity) — the CMK unwrap identity and Entra administrator principals
- [AzureKeyVaultKey](/docs/catalog/azure/key-vault-key) — the customer-managed encryption key
