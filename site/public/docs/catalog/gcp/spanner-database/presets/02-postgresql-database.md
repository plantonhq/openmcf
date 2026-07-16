---
title: "PostgreSQL-Dialect Database"
description: "Provisions a Spanner database with the PostgreSQL interface — PostgreSQL syntax and tooling on Spanner's globally distributed, strongly consistent storage. The dialect choice is permanent."
type: "preset"
rank: "02"
presetSlug: "02-postgresql-database"
componentSlug: "spanner-database"
componentTitle: "Spanner Database"
provider: "gcp"
icon: "package"
order: 2
---

# PostgreSQL-Dialect Database

Provisions a Spanner database with the PostgreSQL interface — PostgreSQL syntax and tooling on Spanner's globally distributed, strongly consistent storage. The dialect choice is permanent.

## When to Use

- Teams standardized on PostgreSQL syntax, drivers, and tooling
- Migrating a PostgreSQL application to Spanner with minimal query changes
- Organizations that mandate open, portable SQL interfaces

## Key Configuration

- **POSTGRESQL dialect** — immutable and permanent; converting later means creating a new database and moving data
- **7d version retention** — the maximum point-in-time recovery window
- **Deletion protection ON by default** — set `deletionProtection: false` explicitly before an intentional teardown

## Customization Notes

- Some Spanner-specific features (interleaved tables, STRUCT types) are unavailable through the PostgreSQL dialect
- DDL statements in `ddl` must use PostgreSQL syntax (the provider quotes identifiers per dialect)
- `project_id` falls back to the provider's default project

## Related Presets

- **01-basic-database** — Google Standard SQL database
- **03-cmek-encrypted** — customer-managed encryption + drop protection
