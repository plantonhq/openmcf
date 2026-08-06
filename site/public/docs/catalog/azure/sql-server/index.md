---
title: "SQL Server"
description: "SQL Server deployment documentation"
icon: "package"
order: 100
componentName: "azuremssqlserver"
---

# Azure SQL Server

Creates an Azure SQL Database logical server -- the administrative container carrying the login endpoint, SQL and Microsoft Entra authentication, networking posture, transparent-data-encryption key, auditing, and Microsoft Defender settings. Databases and elastic pools are first-class resources (AzureMssqlDatabase, AzureMssqlElasticPool) created on the server, each with its own compute and billing; the server itself is free.

## What Gets Created

When you deploy an AzureMssqlServer resource, Planton provisions:

- **SQL logical server** -- an `azurerm_mssql_server` with your authentication posture (SQL, Entra, or mixed), managed identity, TDE customer-managed key, connection policy, TLS floor, and network dials
- **Firewall Rules** -- an `azurerm_mssql_firewall_rule` per `firewallRules` entry (IPv4 allowlist on the public endpoint)
- **Virtual Network Rules** -- an `azurerm_mssql_virtual_network_rule` per `virtualNetworkRules` entry (subnet service-endpoint allowlist)
- **Outbound Firewall Rules** -- an `azurerm_mssql_outbound_firewall_rule` per `outboundFirewallRules` FQDN while outbound restriction is on
- **Extended Auditing** -- an `azurerm_mssql_server_extended_auditing_policy` when `extendedAuditing` is present (blob storage and/or Azure Monitor)
- **Defender Alert Policy** -- an `azurerm_mssql_server_security_alert_policy` when `securityAlertPolicy` is present

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **A resource group** to create the server in (an `AzureResourceGroup` in composed environments)
- **For the TDE customer-managed key**: an identity with wrap/unwrap access on the Key Vault key, granted before the server is created

## Quick Start

Create a file `sqlserver.yaml`:

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureMssqlServer
metadata:
  name: my-sql
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureMssqlServer.my-sql
spec:
  region: eastus
  resourceGroup:
    value: my-rg
  serverName: myorg-dev-sql
  administratorLogin: sqladmin
  administratorPassword:
    value: "Ch@ngeMe1234!"
  firewallRules:
    - name: allow-azure-services
      startIpAddress: "0.0.0.0"
      endIpAddress: "0.0.0.0"
```

Deploy:

```shell
planton apply -f sqlserver.yaml
```

This creates a SQL-auth logical server with the TLS 1.2 floor and public access allowlisted to Azure-internal services. Create databases on it as AzureMssqlDatabase resources referencing `status.outputs.server_id`.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `region` | `string` | Azure region. Changing it replaces the server. | Required |
| `resourceGroup` | `StringValueOrRef` | Resource group name. Defaults to referencing an `AzureResourceGroup`'s name output. | Required |
| `serverName` | `string` | GLOBALLY unique -- becomes `{name}.database.windows.net`. Changing it replaces the server. | Required, 1-63 lowercase letters/digits/hyphens |

At least one authentication mechanism is required: `administratorLogin` + `administratorPassword` (SQL auth), and/or an `azureadAdministrator` with `azureadAuthenticationOnly: true` (Entra-only -- SQL credentials must then be omitted).

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `version` | `string` | `"12.0"` | `"12.0"` (current Azure SQL) or `"2.0"` (legacy identifier). Fixed at creation. |
| `azureadAdministrator` | `object` | -- | `loginUsername`, `objectId` (a directory object ID; defaults to referencing a UAI's `principal_id`), optional `tenantId`, `azureadAuthenticationOnly`. |
| `identity` | `object` | -- | `SYSTEM_ASSIGNED`, `USER_ASSIGNED`, or `SYSTEM_AND_USER_ASSIGNED` + identity references. |
| `primaryUserAssignedIdentityId` | `StringValueOrRef` | -- | Which attached identity ARM uses for Key Vault access. Required with USER_ASSIGNED. |
| `transparentDataEncryptionKeyVaultKeyId` | `StringValueOrRef` | -- | Server-level TDE CMK, a VERSIONED key (references `AzureKeyVaultKey.key_id`). Requires an identity. |
| `connectionPolicy` | `enum` | Azure's Default | `DEFAULT`, `PROXY`, `REDIRECT` (lowest latency; needs ports 11000-11999). |
| `minimumTlsVersion` | `string` | `"1.2"` | The TLS floor (only `"1.2"` on current API versions). Cannot be removed once set. |
| `publicNetworkAccessEnabled` | `bool` | `true` | When false, only private endpoints reach the server. |
| `outboundNetworkRestrictionEnabled` | `bool` | `false` | Restrict OUTBOUND connections to `outboundFirewallRules`. |
| `outboundFirewallRules` | `list(string)` | `[]` | Allowed outbound FQDNs (requires the restriction toggle). |
| `expressVulnerabilityAssessmentEnabled` | `bool` | `false` | Microsoft Defender's agentless SQL scanning. |
| `firewallRules` | `list` | `[]` | Public-endpoint IPv4 allowlist (`0.0.0.0`-`0.0.0.0` admits Azure-internal services only). |
| `virtualNetworkRules` | `list` | `[]` | Subnets admitted via Microsoft.Sql service endpoints (references `AzureSubnet`). |
| `extendedAuditing` | `object` | -- | Server-wide audit trail to blob storage (`storageEndpoint` + sensitive key) and/or Azure Monitor (`logMonitoringEnabled`, default true). |
| `securityAlertPolicy` | `object` | -- | Defender threat detection: `state`, `disabledAlerts`, emails, retention, optional storage export. |
| `tags` | `map(string)` | `{}` | User tags, merged over Planton-derived tags (user wins on collision). |

## Examples

### Entra-Only (Passwordless) Server

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureMssqlServer
metadata:
  name: entra-sql
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      name: data-rg
  serverName: myorg-entra-sql
  azureadAdministrator:
    loginUsername: dba-team
    objectId:
      value: "11111111-2222-3333-4444-555555555555"
    azureadAuthenticationOnly: true
  securityAlertPolicy:
    state: ENABLED
    emailAccountAdmins: true
```

### Private Server with a Customer-Managed TDE Key

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureMssqlServer
metadata:
  name: hardened-sql
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      name: data-rg
  serverName: myorg-hardened-sql
  administratorLogin: sqladmin
  administratorPassword:
    value: "Ch@ngeMe1234!"
  publicNetworkAccessEnabled: false
  identity:
    type: USER_ASSIGNED
    identityIds:
      - valueFrom:
          name: sql-tde-identity
  primaryUserAssignedIdentityId:
    valueFrom:
      name: sql-tde-identity
  transparentDataEncryptionKeyVaultKeyId:
    valueFrom:
      kind: AzureKeyVaultKey
      name: sql-tde-key
      fieldPath: status.outputs.key_id
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `server_id` | `string` | The server's ARM ID -- referenced by AzureMssqlDatabase/AzureMssqlElasticPool `serverId` and AzurePrivateEndpoint |
| `server_name` | `string` | The server's name |
| `fqdn` | `string` | `{name}.database.windows.net` -- the connection-string host |
| `administrator_login` | `string` | The admin login (empty on Entra-only servers) |
| `identity_principal_id` | `string` | The system-assigned identity's principal ID -- the `AzureRoleAssignment` seam |

## Related Components

- [AzureResourceGroup](/docs/catalog/azure/resource-group) — provides the resource group for server placement
- [AzureMssqlDatabase](/docs/catalog/azure/sql-database) — databases on the server (the unit of compute and billing)
- [AzureMssqlElasticPool](/docs/catalog/azure/sql-elastic-pool) — shared compute databases can join
- [AzurePrivateEndpoint](/docs/catalog/azure/private-endpoint) — private connectivity (subresource `sqlServer`)
- [AzureSubnet](/docs/catalog/azure/subnet) — subnets admitted via service-endpoint rules
- [AzureUserAssignedIdentity](/docs/catalog/azure/user-assigned-identity) — the TDE unwrap identity and Entra administrator principals
- [AzureKeyVaultKey](/docs/catalog/azure/key-vault-key) — the customer-managed TDE key
