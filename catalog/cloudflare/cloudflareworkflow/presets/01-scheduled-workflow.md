# Scheduled workflow

A nightly batch workflow: the `NightlyReconcile` class exported by the `reconcile-worker` script runs every day at 03:00 UTC, capped at 512 steps per instance. Errored instances stay queryable for a week (debugging window); successful ones for a day.

**When to use it:** any recurring multi-step job that must survive failures mid-run -- reconciliation, report generation, cleanup sweeps.

**What to change:** point `script_name` at your Worker (or a `CloudflareWorker` reference), name the exported class, and set the cron cadence. Remember the workflow name is an UPSERT key at Cloudflare -- pick one no other team uses.
