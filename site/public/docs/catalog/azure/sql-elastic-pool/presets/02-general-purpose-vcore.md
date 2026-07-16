---
title: "General Purpose vCore Pool with Hybrid Benefit"
description: "This preset creates a 4-vCore General Purpose elastic pool with a 512 GB shared storage cap, fractional per-database bounds, and Azure Hybrid Benefit licensing (bring your own SQL Server license)."
type: "preset"
rank: "02"
presetSlug: "02-general-purpose-vcore"
componentSlug: "sql-elastic-pool"
componentTitle: "SQL Elastic Pool"
provider: "azure"
icon: "package"
order: 2
---

# General Purpose vCore Pool with Hybrid Benefit

This preset creates a 4-vCore General Purpose elastic pool with a 512 GB
shared storage cap, fractional per-database bounds, and Azure Hybrid
Benefit licensing (bring your own SQL Server license).

## When to Use

- Database fleets sized in vCores that want independent compute/storage
  dials
- Organizations holding SQL Server licenses with Software Assurance
  (Hybrid Benefit cuts the compute rate by up to 55%)

## Key Configuration Choices

- **`GP_Gen5` + `capacity: 4`** -- the tier and hardware family are
  derived from the SKU name, so a mismatched combination is
  unrepresentable
- **Fractional per-database bounds** -- vCore pools accept quarters
  (0.25/0.5/...): each database keeps a quarter vCore warm and caps at 2
- **`licenseType: BASE_PRICE`** -- Hybrid Benefit; drop the field or use
  `LICENSE_INCLUDED` for pay-as-you-go

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<server-resource-name>` | The AzureMssqlServer's Planton resource name | Your server composition |
| `<server-region>` | The parent server's region (must match) | The server's spec |
| `<pool-name>` | The pool's name on the server | Your convention |
| `<cost-center>` | Your org's cost-attribution tag value | Your tagging convention |
