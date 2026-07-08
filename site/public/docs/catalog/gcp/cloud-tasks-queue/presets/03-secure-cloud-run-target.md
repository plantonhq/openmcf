---
title: "Secure Cloud Run Target"
description: "The recommended modern Cloud Tasks pattern: the queue owns authentication and routing, so producers enqueue bare payloads and every task is dispatched to one Cloud Run service with an automatically..."
type: "preset"
rank: "03"
presetSlug: "03-secure-cloud-run-target"
componentSlug: "cloud-tasks-queue"
componentTitle: "Cloud Tasks Queue"
provider: "gcp"
icon: "package"
order: 3
---

# Secure Cloud Run Target

The recommended modern Cloud Tasks pattern: the queue owns authentication
and routing, so producers enqueue bare payloads and every task is
dispatched to one Cloud Run service with an automatically minted OIDC
token.

## What this preset creates

A Cloud Tasks queue named `cloud-run-dispatch` whose queue-level HTTP
target rewrites every task's URI to the configured Cloud Run host
(`enforceMode: ALWAYS`), attaches an OIDC token generated as the
referenced `GcpServiceAccount`, sets a JSON content type, and dispatches
at most 100 tasks/second.

## When to use

- Cloud Run microservice task dispatch (private services requiring
  authenticated invocation)
- Cloud Functions background processing
- Serverless webhook delivery where producers should not carry credentials

## Key configuration choices

- `oidcToken.serviceAccountEmail` — a `GcpServiceAccount` reference; the
  account needs `roles/run.invoker` on the target service, and the deploy
  identity needs `iam.serviceAccounts.actAs` on it.
- `uriOverride.enforceMode: ALWAYS` — the queue's routing wins over
  whatever URI the producer enqueued; use `IF_NOT_EXISTS` to only fill
  gaps.
- `audience` — for Cloud Run, the service URL. Cloud Run rejects tokens
  minted for a different audience.

## Placeholders to replace

- `audience` and `uriOverride.host` — the real URL/host of your Cloud Run
  service.
- The `task-invoker` service account reference — the name of your
  `GcpServiceAccount` resource.

## Related presets

- `02-rate-limited-processing` — the same rate/retry discipline without
  queue-level auth.
