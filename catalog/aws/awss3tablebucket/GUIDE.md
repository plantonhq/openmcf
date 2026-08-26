# AwsS3TableBucket — Operational Guide

Live-earned judgment lands here as proof runs and adopter operations teach it; the notes below are the forge-time seed.

## The spec's schema is a birth certificate, not a contract

The Iceberg schema in the spec creates the table; it is never read back, never drifts, and never updates — real schema evolution happens through engines (ALTER TABLE ADD COLUMN). Both modules therefore IGNORE post-create changes to the schema field (live-proven 2026-08-26): editing it after the table exists does nothing — no replacement, no data loss, no drift noise — and importing an existing table plans zero-diff against whatever schema the manifest carries. Adding tables is safe; to change a live table's columns, use a query engine, then update the manifest to match for the record.

## KMS encryption has a silent failure mode

With `aws:kms`, the S3 Tables maintenance principal (`maintenance.s3tables.amazonaws.com`) must be able to use the key. Miss the grant and nothing errors — compaction and snapshot expiry just stop, and query performance degrades over weeks as small files pile up. Grant the key BEFORE the first write.

## Maintenance dials are cost-performance trades

`unreferenced_days`/`non_current_days` decide how long dead files bill before cleanup; `min_snapshots_to_keep`/`max_snapshot_age_hours` bound time travel. Shrinking retention saves storage and shrinks your undo window in the same edit — decide per table, not globally.

## Namespaces and tables are create-only names

Renaming either replaces it (and a table replacement drops data). Naming conventions matter more here than in most kinds: settle `namespace.table` naming before the first deploy.

## force_destroy is the difference between teardown and a support ticket

A non-empty table bucket refuses deletion resource-by-resource (tables, then namespaces, then the bucket). `force_destroy: true` drains it declaratively — leave it false everywhere data matters.

## A deleted bucket's name stays reserved for a short window

Recreating a table bucket under a just-deleted name fails with 409 ConflictException — "The bucket is in a transitional state because of a previous deletion attempt. Try again later" (observed live: the window covers at least the seconds-to-minutes range). This is not an error in your manifest; wait the window out, or pick a fresh name when the rebuild is part of an automated flow that cannot wait.
