---
title: "Standard SQL-Auth Server"
description: "This preset creates the simplest useful Azure SQL logical server: SQL authentication, the public endpoint with the Azure-services firewall rule, and Azure's defaults everywhere else (TLS 1.2 floor,..."
type: "preset"
rank: "01"
presetSlug: "01-standard-sql-auth"
componentSlug: "sql-server"
componentTitle: "SQL Server"
provider: "azure"
icon: "package"
order: 1
---

# Standard SQL-Auth Server

This preset creates the simplest useful Azure SQL logical server: SQL
authentication, the public endpoint with the Azure-services firewall
rule, and Azure's defaults everywhere else (TLS 1.2 floor, Default
connection policy). The server itself is free -- compute and billing live
on the AzureMssqlDatabase and AzureMssqlElasticPool resources you create
on it.

## When to Use

- The starting point for any Azure SQL estate that authenticates with
  SQL logins
- Development and production alike -- the server is only the
  administrative container; harden or scale by what you attach to it

## Key Configuration Choices

- **SQL authentication** -- the classic login + password pair; layer an
  `azureadAdministrator` later without recreating the server
- **Public endpoint with a firewall allowlist** -- the 0.0.0.0 rule
  admits Azure-internal services only; each human or CI network gets its
  own rule
- **`server_name` is a global DNS name** -- it becomes
  `{name}.database.windows.net`, so pick something org-prefixed

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<resource-group-name>` | The resource group to create the server in | The resource group's `status.outputs.resource_group_name` |
| `<globally-unique-server-name>` | 1-63 lowercase chars, globally unique | Your naming convention (org prefix recommended) |
| `<admin-login>` | The SQL admin login (not `admin`/`sa`/`root`) | Your convention |
| `<admin-password>` | 8-128 chars from 3+ character classes | A secret manager; never commit literals |
| `<cost-center>` | Your org's cost-attribution tag value | Your tagging convention |

## Downstream Wiring

Create databases as AzureMssqlDatabase resources referencing this
server's `server_id` output, then connect with:

```text
Server={status.outputs.fqdn},1433;Database={database};User ID={admin};Password={password};Encrypt=True;
```
