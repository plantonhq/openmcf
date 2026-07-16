---
title: "Private Hardened Server with CMK and Auditing"
description: "This preset creates a compliance-oriented Azure SQL logical server: no public endpoint (private endpoints only), a customer-managed transparent-data-encryption key unwrapped through a user-assigned..."
type: "preset"
rank: "03"
presetSlug: "03-private-hardened"
componentSlug: "sql-server"
componentTitle: "SQL Server"
provider: "azure"
icon: "package"
order: 3
---

# Private Hardened Server with CMK and Auditing

This preset creates a compliance-oriented Azure SQL logical server: no
public endpoint (private endpoints only), a customer-managed
transparent-data-encryption key unwrapped through a user-assigned
identity, outbound network restriction, server-wide auditing to Azure
Monitor, and Microsoft Defender threat detection.

## When to Use

- Regulated environments where the encryption key must live in your own
  Key Vault with your own rotation policy
- Estates whose databases must never be reachable from the public
  internet, in either direction

## Key Configuration Choices

- **`publicNetworkAccessEnabled: false`** -- only an
  AzurePrivateEndpoint (subresource `sqlServer`) reaches the server;
  firewall rules become irrelevant
- **TDE CMK against a VERSIONED key** (`key_id`) -- ARM pins the exact
  key version at the server level; rotate by updating the reference
  (unlike the versionless seams elsewhere in the catalog, this is ARM's
  contract for SQL). The identity needs wrap/unwrap on the key's vault
  BEFORE the server is created
- **Outbound restriction** -- elastic queries and external tables can
  only reach the allowlisted FQDNs, closing the exfiltration path
- **Auditing to Azure Monitor** -- no storage account to manage; consume
  `SQLSecurityAuditEvents` through a diagnostic setting

## Prerequisites

- An `AzureUserAssignedIdentity` with wrap/unwrap access on the key's
  vault (a "Key Vault Crypto Service Encryption User" role assignment)
- An `AzureKeyVaultKey` in a vault with purge protection enabled
- An `AzurePrivateEndpoint` composed against this server's `server_id`
  output for connectivity

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<resource-group-name>` | The resource group to create the server in | The resource group's `status.outputs.resource_group_name` |
| `myorg-prod-sql` | 1-63 lowercase chars, globally unique | Your naming convention |
| `<admin-login>` / `<admin-password>` | SQL admin credentials | A secret manager; never commit literals |
| `<tde-identity-resource-name>` | The user-assigned identity's Planton resource name | Your identity composition |
| `<tde-key-resource-name>` | The Key Vault key's Planton resource name | Your Key Vault composition |
| `<allowed-peer-fqdn>` | An outbound-allowed FQDN (e.g. a peer server) | Your architecture |
| `<security-team-email>` | Where Defender alerts go | Your security team |
| `<cost-center>` | Your org's cost-attribution tag value | Your tagging convention |
