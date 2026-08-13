# GcpSpannerDatabase Guide

The judgment this guide protects: four of this kind's fields are
permanent (`databaseName`, `databaseDialect`, `encryptionConfig`,
`instance`) and one is append-only (`ddl`). A Spanner database is cheap
to create and expensive to recreate — decide dialect and encryption
before the first byte of production data, not after.

## Dialect is a one-time identity decision

`GOOGLE_STANDARD_SQL` is the full feature set (interleaved tables,
STRUCTs); `POSTGRESQL` buys familiar syntax and tooling at the cost of
some Spanner-native features. There is no conversion — a dialect change
is a new database and a data migration. Teams rarely regret standard
SQL; they regret discovering a missing PostgreSQL feature mid-project.

## Treat IaC-owned DDL as the initial schema only

The `ddl` list is append-only: new statements apply via UpdateDDL, but
editing or removing an existing entry forces database recreation. Use
it to bootstrap the schema a fresh environment needs, then hand ongoing
migrations to a migration tool (Liquibase, Flyway) — IaC diffs are the
wrong review surface for schema evolution.

## Point-in-time recovery is a dial, not a default

`versionRetentionPeriod` defaults to 1h; raising it (up to 7d) widens
the PITR window at the cost of extra storage on every write-heavy
table. Set it deliberately for databases where "restore to just before
the bad deploy" is the recovery story.

## CMEK has two shapes — match the instance topology

`kmsKeyName` (one key) for regional instance configs;
`kmsKeyNames` (one key per region) for multi-region configs. The keys
must live in the instance's regions, and the choice is immutable. If
the org requires CMEK, wire it through `GcpKmsKey` references at
creation — retrofitting means a new database.

## Three deletion levers, three enforcement points

`deletionProtection` (default TRUE) is the IaC-side guard — a destroy
plan fails before touching GCP. `enableDropProtection` is the API-side
compliance lock — NO interface can delete the database, and the parent
instance becomes undeletable too (remember to unset it before an
intentional instance teardown). `deletionPolicy` decides what a
PERMITTED destroy does: `PREVENT` as a third wall, `ABANDON` to leave
the database (and its bill) running outside management. Production
posture: leave deletionProtection true and add enableDropProtection
for regulated data.

## What is deliberately absent

Backup schedules are `GcpSpannerBackupSchedule` (many per database).
Labels do not exist at the database level in GCP — organize at the
instance.
