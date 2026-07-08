---
title: "Secure Cloud Run Trigger"
description: "The authenticated cron pattern: an HTTP POST to a private Cloud Run service with an OIDC token minted per attempt as a referenced service account — no credentials stored anywhere."
type: "preset"
rank: "03"
presetSlug: "03-secure-cloud-run-trigger"
componentSlug: "cloud-scheduler-job"
componentTitle: "Cloud Scheduler Job"
provider: "gcp"
icon: "package"
order: 3
---

# Secure Cloud Run Trigger

The authenticated cron pattern: an HTTP POST to a private Cloud Run
service with an OIDC token minted per attempt as a referenced service
account — no credentials stored anywhere.

## What this preset creates

A Cloud Scheduler job named `daily-report-trigger` that POSTs a JSON
payload to a Cloud Run endpoint at 09:00 America/New_York on weekdays,
authenticating with an OIDC token generated as the referenced
`GcpServiceAccount`, allowing the handler 10 minutes per attempt and
retrying up to 3 times with 5s–600s exponential backoff.

## When to use

- Scheduled Cloud Run work (daily reports, ETL kickoffs, cleanup)
- Invoking authenticated Cloud Functions on a cron schedule
- Periodic calls to any private endpoint that validates Google OIDC

## Key configuration choices

- `oidcToken.serviceAccountEmail` — a `GcpServiceAccount` reference; the
  account needs `roles/run.invoker` on the target service, and the deploy
  identity needs `iam.serviceAccounts.actAs` on it.
- `audience` — for Cloud Run, the service URL; a token minted for another
  audience is rejected.
- `attemptDeadline: "600s"` — size to the handler's real worst case; HTTP
  targets accept 15 seconds to 30 minutes.

## Placeholders to replace

- `httpTarget.uri` and `oidcToken.audience` — your Cloud Run service URL.
- The `report-invoker` service account reference — the name of your
  `GcpServiceAccount` resource.
- `body` — base64 of your real payload.

## Related presets

- `01-basic-http-job` — unauthenticated endpoints.
- `02-pubsub-publisher` — fan out through an event instead of calling one
  service.
