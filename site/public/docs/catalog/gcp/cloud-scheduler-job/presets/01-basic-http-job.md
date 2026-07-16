---
title: "Basic HTTP Job"
description: "The minimal scheduled trigger: an unauthenticated HTTP GET fired on a cron cadence, with GCP-managed defaults for retries and deadlines."
type: "preset"
rank: "01"
presetSlug: "01-basic-http-job"
componentSlug: "cloud-scheduler-job"
componentTitle: "Cloud Scheduler Job"
provider: "gcp"
icon: "package"
order: 1
---

# Basic HTTP Job

The minimal scheduled trigger: an unauthenticated HTTP GET fired on a
cron cadence, with GCP-managed defaults for retries and deadlines.

## What this preset creates

A Cloud Scheduler job named `uptime-ping` (from `metadata.name`) that
issues an HTTP GET to the configured URL every 5 minutes, in the ambient
project, with the default 180s attempt deadline and default retry policy.

## When to use

- Triggering a public webhook on a schedule
- Periodic health checks and uptime pings
- Any recurring call to an endpoint that needs no authentication

## Placeholders to replace

- `httpTarget.uri` — the real endpoint to call.
- `schedule` — your cron cadence (`"0 9 * * 1-5"` for weekday mornings).

## Remix ideas

- Add `oidcToken` under `httpTarget` when the endpoint requires
  authenticated invocation (see the secure-cloud-run-trigger preset).
- Add `timeZone` (e.g. `America/New_York`) to schedule in local time.
- Add `retryConfig` for endpoints where a missed fire must be retried.
