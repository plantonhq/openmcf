# AwsS3TableBucket — Operational Guide

Live-earned judgment lands here as proof runs and adopter operations teach it; the notes below are the forge-time seed.

## The spec's schema is a birth certificate, not a contract

The Iceberg schema in the spec creates the table; it is never read back, never drifts, and never updates — real schema evolution happens through engines (ALTER TABLE ADD COLUMN). Treat a spec schema edit as what it is: a table REPLACEMENT that drops the data. Adding tables is safe; editing their schemas is destructive.

## KMS encryption has a silent failure mode

With `aws:kms`, the S3 Tables maintenance principal (`maintenance.s3tables.amazonaws.com`) must be able to use the key. Miss the grant and nothing errors — compaction and snapshot expiry just stop, and query performance degrades over weeks as small files pile up. Grant the key BEFORE the first write.

## Maintenance dials are cost-performance trades

`unreferenced_days`/`non_current_days` decide how long dead files bill before cleanup; `min_snapshots_to_keep`/`max_snapshot_age_hours` bound time travel. Shrinking retention saves storage and shrinks your undo window in the same edit — decide per table, not globally.

## Namespaces and tables are create-only names

Renaming either replaces it (and a table replacement drops data). Naming conventions matter more here than in most kinds: settle `namespace.table` naming before the first deploy.

## force_destroy is the difference between teardown and a support ticket

A non-empty table bucket refuses deletion resource-by-resource (tables, then namespaces, then the bucket). `force_destroy: true` drains it declaratively — leave it false everywhere data matters.
