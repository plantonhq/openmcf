# Weekly Long-Retention Archive

A weekly backup on Sunday with the maximum 14-week retention — the
long-lived archive leg of the daily-plus-weekly pattern.

## When to use

Compliance or disaster-recovery archives where backups must survive
database deletion and extend well beyond PITR's version window.

## What to customize

- `weeklyRecurrence.day` — pin the weekly run (immutable after creation).
- `retention` — up to `8467200s` (14 weeks).

## Composes with

`GcpFirestoreDatabase` upstream and `01-daily-short-retention` on the
same database (Firestore allows one daily and one weekly schedule each).
