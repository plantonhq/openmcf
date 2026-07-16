# Azure MySQL Flexible Server

Creates an Azure Database for MySQL Flexible Server -- Azure's managed MySQL with per-server compute and storage sizing, zone-redundant high availability, a Microsoft Entra administrator, customer-managed-key encryption, read replicas, and point-in-time restore. Databases, firewall rules, server parameters, and the Entra administrator are declared on the server and managed with it.

## What Gets Created

When you deploy an AzureMysqlFlexibleServer resource, Planton provisions:

- **MySQL Flexible Server** -- an `azurerm_mysql_flexible_server` in the specified region and resource group, with your chosen compute SKU, storage profile (size, IOPS or elastic scaling, auto-grow), version, availability posture, networking, identities, and encryption
- **Databases** -- an `azurerm_mysql_flexible_database` for each entry in `databases`, each with its own charset and collation
- **Firewall Rules** -- an `azurerm_mysql_flexible_server_firewall_rule` for each entry in `firewallRules`, allowlisting IPv4 ranges on the public endpoint
- **Server Parameters** -- an `azurerm_mysql_flexible_server_configuration` for each `serverParameters` entry, applied as user overrides on Azure's per-SKU defaults
- **Entra Administrator** -- an `azurerm_mysql_flexible_server_active_directory_administrator` for the single `aadAdministrator` (MySQL supports exactly one), backed by a user-assigned identity attached to the server

A read replica or a restored server is simply another AzureMysqlFlexibleServer whose `createMode` and `sourceServerId` reference the source.

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **A resource group** to create the server in (an `AzureResourceGroup` in composed environments)
- **For VNet injection**: a subnet delegated to `Microsoft.DBforMySQL/flexibleServers` with no other resources, plus a private DNS zone for the server's name
- **For CMK encryption or the Entra administrator**: a user-assigned identity attached via `userAssignedIdentityIds` (with wrap/unwrap access on the Key Vault key for CMK)

## Quick Start

Create a file `mysql.yaml`:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureMysqlFlexibleServer
metadata:
  name: my-mysql
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureMysqlFlexibleServer.my-mysql
spec:
  region: eastus
  resourceGroup:
    value: my-rg
  serverName: myorg-dev-mysql
  administratorLogin: mysqladmin
  administratorPassword:
    value: "Ch@ngeMe1234!"
  skuName: B_Standard_B1ms
  databases:
    - name: myapp
```

Deploy:

```shell
planton apply -f mysql.yaml
```

This creates a Burstable MySQL 8.0.21 server with Azure's default 20 GiB auto-growing storage and 7-day backups, public access, and one application database. Read `status.outputs.fqdn` to build the connection string.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `region` | `string` | Azure region. Changing it replaces the server. | Required |
| `resourceGroup` | `StringValueOrRef` | Resource group name. Defaults to referencing an `AzureResourceGroup`'s name output. | Required |
| `serverName` | `string` | GLOBALLY unique -- becomes `{name}.mysql.database.azure.com`. Changing it replaces the server. | Required, 3-63 lowercase letters/digits/hyphens |
| `skuName` | `string` | Compute SKU, `{TIER}_Standard_{SIZE}` (`B_`/`GP_`/`MO_`). Required for a fresh server; a replica left unset inherits the source's. | Pattern-validated |

`administratorLogin` + `administratorPassword` are required for a fresh server; replicas and restores inherit them from the source. MySQL always keeps password auth on -- it cannot be disabled.

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `createMode` | `enum` | `DEFAULT` | `DEFAULT`, `REPLICA`, `POINT_IN_TIME_RESTORE`, `GEO_RESTORE`. Fixed at creation. |
| `sourceServerId` | `StringValueOrRef` | -- | The source server for replica/restore modes; references another server's `server_id` output. |
| `pointInTimeRestoreTimeInUtc` | `string` | -- | RFC-3339 restore instant (POINT_IN_TIME_RESTORE only; GEO_RESTORE takes no timestamp). |
| `replicationRole` | `enum` | -- | `NONE` promotes a replica to a standalone primary (irreversible, day-2 only). |
| `version` | `string` | `"8.0.21"` | MySQL version: `"5.7"`, `"8.0.21"`, `"8.4"`. Upgrades go 5.7 to 8.0.21 in place; downgrades replace the server. |
| `storage` | `object` | Azure defaults | `sizeGb` (20-16384, grows only), `iops` (360-48000) XOR `ioScalingEnabled`, `autoGrowEnabled` (default true), `logOnDiskEnabled`. |
| `zone` | `string` | Azure picks | Primary availability zone (`"1"`/`"2"`/`"3"`). |
| `highAvailability` | `object` | -- | `mode` (`ZONE_REDUNDANT`/`SAME_ZONE`) + optional `standbyAvailabilityZone`. Not supported on Burstable SKUs or replicas. |
| `maintenanceWindow` | `object` | system-managed | Weekly patching window (`dayOfWeek`/`startHour`/`startMinute`). |
| `backupRetentionDays` | `int32` | `7` | 1-35 days -- the point-in-time restore horizon. |
| `geoRedundantBackupEnabled` | `bool` | `false` | Replicates backups to the paired region (enables `GEO_RESTORE`). Fixed at creation. |
| `publicNetworkAccess` | `enum` | Azure derives | `ENABLED`/`DISABLED`. Leave unset: Azure derives ENABLED publicly, DISABLED when VNet-injected. |
| `delegatedSubnetId` | `StringValueOrRef` | -- | VNet injection: a delegated subnet (references `AzureSubnet`). Requires `privateDnsZoneId`. Fixed at creation. |
| `privateDnsZoneId` | `StringValueOrRef` | -- | The private DNS zone the injected server registers in (references `AzurePrivateDnsZone`). Fixed at creation. |
| `userAssignedIdentityIds` | `list` | `[]` | User-assigned identities attached to the server (MySQL supports no system-assigned flavor). Required for CMK and the Entra administrator. |
| `customerManagedKey` | `object` | -- | CMK encryption: `keyVaultKeyId` (references `AzureKeyVaultKey.versionless_id`), the unwrap identity, and an optional geo-backup key + identity pair. |
| `aadAdministrator` | `object` | -- | The single Entra admin: `identityId` (an attached identity), `login`, `objectId` (a directory object ID; a managed identity's CLIENT ID), optional `tenantId`. |
| `databases` | `list` | `[]` | Databases (`name`, `charset` default `utf8mb4`, `collation` default `utf8mb4_0900_ai_ci`). |
| `firewallRules` | `list` | `[]` | Public-endpoint IPv4 allowlist (`0.0.0.0`-`0.0.0.0` admits Azure-internal services only). |
| `serverParameters` | `map(string)` | `{}` | MySQL parameter overrides by name (e.g. `require_secure_transport`, `max_connections`). Static parameters need a restart. |
| `tags` | `map(string)` | `{}` | User tags, merged over Planton-derived tags (user wins on collision). |

## Examples

### Production: VNet-Injected with Zone-Redundant HA

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureMysqlFlexibleServer
metadata:
  name: prod-mysql
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      name: data-rg
  serverName: myorg-prod-mysql
  administratorLogin: mysqladmin
  administratorPassword:
    value: "Ch@ngeMe1234!"
  skuName: GP_Standard_D4ds_v4
  storage:
    sizeGb: 256
    ioScalingEnabled: true
  zone: "1"
  highAvailability:
    mode: ZONE_REDUNDANT
    standbyAvailabilityZone: "2"
  backupRetentionDays: 35
  geoRedundantBackupEnabled: true
  delegatedSubnetId:
    valueFrom:
      name: database-subnet
  privateDnsZoneId:
    valueFrom:
      name: mysql-dns-zone
  databases:
    - name: orders
```

### Read Replica of a Primary

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureMysqlFlexibleServer
metadata:
  name: prod-mysql-replica
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      name: data-rg
  serverName: myorg-prod-mysql-replica
  createMode: REPLICA
  sourceServerId:
    valueFrom:
      kind: AzureMysqlFlexibleServer
      name: prod-mysql
      fieldPath: status.outputs.server_id
```

The replica inherits the source's SKU, storage, and version. Promote it later by setting `replicationRole: NONE`. The source's `replica_capacity` output reports how many more replicas it can accept.

### Hardened: Customer-Managed Key and an Entra Administrator

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureMysqlFlexibleServer
metadata:
  name: hardened-mysql
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      name: data-rg
  serverName: myorg-hardened-mysql
  administratorLogin: mysqladmin
  administratorPassword:
    value: "Ch@ngeMe1234!"
  skuName: GP_Standard_D4ds_v4
  userAssignedIdentityIds:
    - valueFrom:
        name: mysql-cmk-identity
  customerManagedKey:
    keyVaultKeyId:
      valueFrom:
        kind: AzureKeyVaultKey
        name: mysql-cmk
        fieldPath: status.outputs.versionless_id
    primaryUserAssignedIdentityId:
      valueFrom:
        name: mysql-cmk-identity
  aadAdministrator:
    identityId:
      valueFrom:
        name: mysql-cmk-identity
    login: dba-team
    objectId:
      value: "11111111-2222-3333-4444-555555555555"
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `server_id` | `string` | The server's ARM ID -- referenced by `AzurePrivateEndpoint` and by replica/restore servers' `sourceServerId` |
| `server_name` | `string` | The server's name |
| `fqdn` | `string` | `{name}.mysql.database.azure.com` -- the connection-string host; resolves privately when VNet-injected |
| `administrator_login` | `string` | The admin login, echoed for connection strings |
| `database_ids` | `map(string)` | Each declared database's ARM ID, keyed by name |
| `replica_capacity` | `int32` | How many read replicas the server can still accept (burstable SKUs report 0) |

## Related Components

- [AzureResourceGroup](/docs/catalog/azure/azureresourcegroup) — provides the resource group for server placement
- [AzureSubnet](/docs/catalog/azure/subnet) — the delegated subnet for VNet injection
- [AzurePrivateDnsZone](/docs/catalog/azure/private-dns-zone) — private name resolution for the injected server
- [AzurePrivateEndpoint](/docs/catalog/azure/azureprivateendpoint) — Private Link connectivity as an alternative to VNet injection
- [AzureUserAssignedIdentity](/docs/catalog/azure/user-assigned-identity) — the CMK unwrap identity and the Entra administrator's backing identity
- [AzureKeyVaultKey](/docs/catalog/azure/key-vault-key) — the customer-managed encryption key
