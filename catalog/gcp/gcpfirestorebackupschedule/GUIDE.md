# GcpFirestoreBackupSchedule Guide

The judgment this guide protects: a backup schedule is a promise about the
FUTURE — it protects nothing retroactively, its recurrence is frozen at
creation, and deleting it silently stops the promise while leaving old
backups to age out. Set the cadence deliberately and guard it.

## The daily-plus-weekly pattern is two resources, on purpose

The API accepts exactly one recurrence per schedule, and a database
supports one daily and one weekly. The classic production posture is
both, composed as two of these nodes on the same database:

- **Daily, short retention** (`retention: "604800s"` — 7 days): the
  operational restore point for "yesterday's data".
- **Weekly, long retention** (up to `"8467200s"` — 14 weeks): the
  compliance and disaster horizon.

Recurrence (daily vs weekly, and the weekly day) is immutable — changing
cadence means replacing the schedule, which is safe: existing backups
outlive it.

## Retention arithmetic worth doing once

`retention` applies to backups created AFTER a change; existing backups
keep what they were born with. Cost scales with database size × backups
retained: a 7-day daily schedule holds ~7 copies, a 14-week weekly holds
~14. Shortening retention on a large database is the cheapest cost lever;
it just takes effect one backup at a time.

## Backups are the disaster plan; PITR is the oops plan

The database's point-in-time recovery covers at most 7 days INSIDE the
live database and dies with it. Scheduled backups survive database
deletion — they are what stands between a deleted production database and
a very bad week. If a database carries `deleteProtectionState` and
`deletionPolicy: PREVENT`, its weekly schedule deserves
`deletionPolicy: PREVENT` too: unguarded, a destroyed schedule stops
future backups without any immediate symptom.

## Conventions and gotchas

- Backup timing within the day is Firestore's choice; only weekly
  schedules pin a day. Do not build pipelines that assume a backup exists
  by a wall-clock hour.
- Deleting the schedule never deletes backups already taken — they age
  out per their own retention. `ABANDON` goes one step further: the
  schedule itself keeps running, unmanaged.
- The `database` reference resolves a `GcpFirestoreDatabase`'s
  `database_name` output; the schedule must live in the database's
  project.

## On the diagram

Schedules render as protection nodes hanging off their database — a
database with no schedule node visibly carries no backup promise, which
is exactly the review this composition enables.

## Pairs well with

- `GcpFirestoreDatabase` — the database being protected; pair this kind's
  PREVENT posture with the database's delete protection.
