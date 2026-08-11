---
title: "HTTP Orchestrator"
description: "Call internal services in order with per-step OIDC auth and retries — the canonical Workflows shape: the workflow is the coordination skeleton, the services do the work."
type: "preset"
rank: "01"
presetSlug: "01-http-orchestrator"
componentSlug: "workflow"
componentTitle: "Workflow"
provider: "gcp"
icon: "package"
order: 1
---

# HTTP Orchestrator

Call internal services in order with per-step OIDC auth and retries —
the canonical Workflows shape: the workflow is the coordination
skeleton, the services do the work.

## What it configures

- A two-step orchestration (validate, then process) with exponential
  backoff retries on the processing step.
- A dedicated service account identity for every authenticated call.
- `LOG_ERRORS_ONLY` call logging and `PREVENT` teardown posture.

## Adjust before deploying

- **sourceContents** — point the step URLs at your services; add steps
  as the flow grows (keep per-item logic in the services, not inline).
- **serviceAccount** — reference the GcpServiceAccount granted invoker
  roles on exactly the services the steps call.
- **region** — co-locate with the services being called.

## After deploying

Run it: `gcloud workflows run http-orchestrator --data='{...}'`. Every
source change deploys a new revision — the `revision_id` output confirms
a deploy rolled.

## When to choose something else

For executions started by EVENTS (a Pub/Sub message, a Storage object,
an audit log entry), start from the **Event Handler** preset and pair
with a GcpEventarcTrigger.
