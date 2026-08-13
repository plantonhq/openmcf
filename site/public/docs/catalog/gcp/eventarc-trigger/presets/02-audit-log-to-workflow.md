---
title: "Audit Log to Workflow"
description: "React to GCP control-plane changes with an orchestrated response: every `storage.buckets.create` API call starts a Workflow execution that can notify, tag, or remediate."
type: "preset"
rank: "02"
presetSlug: "02-audit-log-to-workflow"
componentSlug: "eventarc-trigger"
componentTitle: "Eventarc Trigger"
provider: "gcp"
icon: "package"
order: 2
---

# Audit Log to Workflow

React to GCP control-plane changes with an orchestrated response: every
`storage.buckets.create` API call starts a Workflow execution that can
notify, tag, or remediate.

## What it configures

- An audit-log trigger (`google.cloud.audit.log.v1.written`) filtered
  to one service + method — the broadest trigger family, able to watch
  almost any GCP API call.
- A GcpWorkflow destination (referenced by its `workflow_id` output).

## Adjust before deploying

- **serviceName / methodName** — pick from the audit-log method names
  of the service you watch (`gcloud logging read` shows them live).
- **PREREQUISITE**: enable Data Access audit logs for the watched
  service (IAM & Admin → Audit Logs) — without them nothing fires,
  silently.
- **serviceAccount** — REQUIRED for audit-log triggers; also needs
  `roles/workflows.invoker` on the destination workflow.

## After deploying

Perform the watched action (create a bucket); an execution appears in
the workflow within a minute.

## When to choose something else

For your own application events, start from the **Pub/Sub to Cloud
Run** preset — audit-log triggers are for GCP's control plane, not your
data plane.
