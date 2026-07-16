---
title: "Development Burstable Server"
description: "This preset creates the smallest practical PostgreSQL Flexible Server: a single Burstable instance on the public endpoint, one application database, and the Azure-services firewall rule. Version,..."
type: "preset"
rank: "01"
presetSlug: "01-dev-burstable"
componentSlug: "postgresql-flexible-server"
componentTitle: "PostgreSQL Flexible Server"
provider: "azure"
icon: "package"
order: 1
---

# Development Burstable Server

This preset creates the smallest practical PostgreSQL Flexible Server: a
single Burstable instance on the public endpoint, one application database,
and the Azure-services firewall rule. Version, storage (32 GiB), backup
retention (7 days), and authentication (password on, Entra off) all ride
Azure's defaults.

## When to Use

- Development and test environments where cost matters more than
  availability
- Short-lived preview environments a CI pipeline creates and destroys

## Key Configuration Choices

- **`B_Standard_B1ms`** -- 1 vCPU / 2 GiB, the cheapest real server;
  burstable SKUs cannot use high availability or read replicas, so promote
  to `GP_Standard_*` before production
- **Public endpoint with a firewall allowlist** -- the 0.0.0.0 rule admits
  Azure-internal services only; each human or CI network gets its own rule
- **`server_name` is a global DNS name** -- it becomes
  `{name}.postgres.database.azure.com`, so pick something org-prefixed

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<resource-group-name>` | The resource group to create the server in | The resource group's `status.outputs.resource_group_name` |
| `myorg-dev-postgres` | 3-63 lowercase chars, globally unique | Your naming convention (org prefix recommended) |
| `<admin-login>` | The PostgreSQL admin login (not `admin`/`root`/`pg_*`) | Your convention |
| `<admin-password>` | 8-128 chars from 3+ character classes | A secret manager; never commit literals |
| `appdb` | The application's database | Your application |
| `<cost-center>` | Your org's cost-attribution tag value | Your tagging convention |

## Downstream Wiring

Applications construct their connection string from the outputs:

```text
postgresql://{administrator_login}:{password}@{status.outputs.fqdn}:5432/{database}?sslmode=require
```
