---
title: "Preset: Shared Database Link"
description: "A resource link: a local pointer to a Glue database that lives in another AWS account (or region) and was shared via AWS RAM / Lake Formation. Queries against the link — from Athena, Redshift..."
type: "preset"
rank: "03"
presetSlug: "03-shared-database-link"
componentSlug: "glue-catalog-database"
componentTitle: "Glue Catalog Database"
provider: "aws"
icon: "package"
order: 3
---

# Preset: Shared Database Link

A resource link: a local pointer to a Glue database that lives in another AWS
account (or region) and was shared via AWS RAM / Lake Formation. Queries
against the link — from Athena, Redshift Spectrum, or EMR — resolve to the
target database's tables, subject to the permissions granted on the share.

## What This Configures

- A resource-link database pointing at `sales_curated` in the sharing
  account's catalog.
- No storage, schema, or parameters of its own — a link carries none.

## When to Use

- Consuming a data-platform team's curated datasets from your own account.
- Multi-account data mesh architectures where producers share databases and
  consumers link them locally.

## Customization Points

- Set `targetDatabase.region` for cross-region links.
- The sharing account must have granted your account access via AWS RAM /
  Lake Formation before the link resolves.
