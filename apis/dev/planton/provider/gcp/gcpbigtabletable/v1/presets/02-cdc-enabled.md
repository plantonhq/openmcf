# CDC-Enabled Table with Backups

An operational table wired for downstream processing: a change-stream
feed for Dataflow CDC pipelines, daily automated backups, and a combined
age-plus-versions retention policy.

## When to use

Tables other systems react to — order flows, state machines, audit
trails — where changes must be streamed out and the data is valuable
enough to back up automatically.

## What to customize

- `changeStreamRetention` — 1-7 days; size to your pipeline's worst-case
  catch-up window.
- `automatedBackupPolicy` — backup cadence and retention.
- The `UNION` GC policy — a cell is collected when either the age bound
  or the version cap is hit; switch to `INTERSECTION` to require both.

## Composes with

`GcpBigtableInstance` upstream; Dataflow pipelines consume the change
stream downstream.
