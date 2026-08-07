# Basic Database

Provisions a Spanner database on an existing instance with Google's SQL dialect and a 24-hour point-in-time recovery window. The database is named after `metadata.name`; schema management belongs to your migration tooling.

## When to Use

- The standard starting point for a new application database
- Teams using Google Standard SQL (interleaved tables, STRUCT types)
- When schema DDL is owned by migration tools (Liquibase, Flyway), not IaC

## Key Configuration

- **Instance by reference** — composes against a `GcpSpannerInstance` resource's `instance_name` output
- **24h version retention** — a full day of point-in-time recovery (GCP default is 1h; maximum is 7d)
- **Deletion protection ON by default** — both IaC engines refuse to destroy the database until `deletionProtection: false` is set explicitly

## Customization Notes

- `metadata.name` doubles as the database name when `databaseName` is omitted (2-30 chars; letters, digits, underscores, hyphens)
- Add initial schema via `ddl` (append-only after creation — editing an existing statement recreates the database)
- `project_id` falls back to the provider's default project

## Related Presets

- **02-postgresql-database** — PostgreSQL-dialect database
- **03-cmek-encrypted** — customer-managed encryption + drop protection
