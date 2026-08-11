# GcpCloudSchedulerJob Guide

The judgment this guide protects: a scheduler job is a promise that
something happens on time, unattended. The failure mode is silent — a
job that stops firing pages nobody until the downstream gap is noticed.

## Schedules read in a timezone; write it down

The cron expression is interpreted in `time_zone` (default `Etc/UTC`).
A "9am report" job left in UTC fires at the wrong local hour twice a
year even in UTC-pinned regions' heads — set the zone the humans mean
and say so in the job name or description. There is no catch-up: a
missed window (job paused, target down past retries) is skipped, not
replayed. If every tick must be processed, the target needs its own
reconciliation.

## Retries protect the tick, deadlines protect the queue

`retry_config` re-attempts a failed run with backoff; `attempt_deadline`
caps how long one attempt may hang (15s–30min for HTTP). Size the
deadline to the target's real p99 — a deadline shorter than the work
guarantees spurious retries, and with a non-idempotent target that
means double execution. Same law as every event system: make the
handler idempotent first.

## Auth the target properly

For Cloud Run/Functions targets, use the OIDC token arm with a
dedicated invoker service account; for `*.googleapis.com` targets, the
OAuth arm. Unauthenticated HTTP targets belong only where the endpoint
itself is public by design.

## Destroy semantics

`deletion_policy: DELETE` (default) removes the job — the schedule
simply stops, which is usually the intent. `PREVENT` suits jobs whose
missed runs break downstream systems (billing exports, cache warmers
feeding SLAs). `ABANDON` keeps the job FIRING unmanaged — an abandoned
scheduler job is a background process nobody owns, so prefer pausing
(`paused: true`) before abandoning.
