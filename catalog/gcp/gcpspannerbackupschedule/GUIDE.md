# GcpSpannerBackupSchedule Guide

The judgment this guide protects: schedules are cheap, first-class, and
many-per-database — the right mental model is one schedule per recovery
objective (a daily incremental AND a weekly full, side by side), not
one schedule stretched to cover everything.

## Spanner accepts four cadences — design within them

The cron surface looks general but the API accepts a bounded set:
every 12 hours, daily, weekly, or monthly (evaluated in UTC). Put the
backup window where your write load is lowest; the cadence updates in
place, so tightening it later is a field edit, not a migration.

## FULL vs INCREMENTAL is an edition decision as much as a cost one

`INCREMENTAL` chains store only changes since the previous backup —
dramatically cheaper storage at identical restore semantics — but
require the parent instance to be ENTERPRISE or ENTERPRISE_PLUS, and
`backupType` is immutable. The robust pattern for important data:
a frequent INCREMENTAL schedule for cheap short-range recovery plus an
independent FULL schedule as the chain-independent anchor.

## Retention applies forward, never backward

`retentionDuration` (max 366 days) governs backups created AFTER a
change; existing backups keep the retention they were born with.
Shortening retention does not reclaim old backups early — plan storage
cost from the retention you set at creation.

## Backups outlive everything here — including the schedule

Deleting the schedule (or the whole database) stops FUTURE backups;
existing backups live until their retention expires, and they are what
makes destroying the parent instance fail unless its `forceDestroy` is
armed. `deletionPolicy` on the schedule follows the same truth: DELETE
removes only the cadence, PREVENT protects a recovery objective from
riding along with a stack teardown, ABANDON keeps it producing backups
outside management.

## Encryption: inherit unless compliance says otherwise

Omitting `encryptionConfig` gives USE_DATABASE_ENCRYPTION — CMEK
databases produce CMEK backups with zero configuration. Explicit
CUSTOMER_MANAGED_ENCRYPTION exists for the case where backups must use
a DIFFERENT key (one per region on multi-region configs); explicit
GOOGLE_DEFAULT_ENCRYPTION for deliberately de-CMEKed backups.

## What is deliberately absent

The provider's `full_backup_spec`/`incremental_backup_spec` marker
blocks carry no arguments — the spec's `backupType` enum IS that
choice. One-off (unscheduled) backups are an operational action, not
IaC surface.
