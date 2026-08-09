# GcpBigQueryDataset Guide

The judgment this guide protects: a dataset is a permission and billing
boundary, not a folder. Its location decides where every table it will
ever hold can live, and its access list is authoritative — the two
decisions worth slowing down for.

## Location is forever, and it is gravity

`location` is immutable and every table inherits it. Cross-location
queries don't join; moving a dataset means copying every table. Pick the
location where the data's consumers run — and if an organization
standard exists (US, EU), follow it even for "temporary" datasets,
because temporary datasets aren't.

## The access list is authoritative — own it or don't touch it

`access` REPLACES the dataset's ACL wholesale: an entry removed from the
spec is a grant revoked in GCP, including the defaults BigQuery seeds
(project owners/writers/readers), so carry the baseline grants
explicitly once you take ownership. Authorized views/routines/datasets
belong here — they are resource authorizations, not principal grants,
and they are how a locked-down dataset exposes curated slices.

## Two destroy levers, checked in order

`deletionPolicy` decides whether destroy is attempted at all: `PREVENT`
fails it, `ABANDON` walks away leaving the dataset serving. When it is
`DELETE` (the default), `deleteContentsOnDestroy` gates the second door:
GCP refuses to remove a non-empty dataset until it is true. The safe
production posture is `PREVENT` + `deleteContentsOnDestroy: false`; the
CI-fixture posture is `DELETE` + `true`.

## CMEK is a default, not an enforcement

`kmsKeyName` sets the default key for NEW tables; existing tables keep
their encryption, and a table can still override with its own key. If
the requirement is "everything under this key", set it before the first
table lands and audit per-table keys separately.

## Expiration defaults are fleet policy

`defaultTableExpirationMs` / `defaultPartitionExpirationMs` apply to
tables created AFTER the setting; they silently delete data on
schedule. They are the right tool for scratch/staging datasets and the
wrong one anywhere an analyst might park something important.
