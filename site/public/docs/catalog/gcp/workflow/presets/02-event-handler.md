---
title: "Event Handler"
description: "An event-driven step function: a GcpEventarcTrigger (destination: workflow) starts one execution per event; the workflow decodes the CloudEvent and fans out the work."
type: "preset"
rank: "02"
presetSlug: "02-event-handler"
componentSlug: "workflow"
componentTitle: "Workflow"
provider: "gcp"
icon: "package"
order: 2
---

# Event Handler

An event-driven step function: a GcpEventarcTrigger (destination:
workflow) starts one execution per event; the workflow decodes the
CloudEvent and fans out the work.

## What it configures

- A Pub/Sub-shaped event decoder (base64 + JSON) feeding an
  authenticated handler call.
- Detailed execution history — step inputs/outputs visible in the
  console while the event shapes are still being learned.

## Adjust before deploying

- **decode step** — the shape shown fits Pub/Sub messagePublished
  events; Storage and audit-log events carry their payload directly in
  `event.data` (drop the base64 step).
- **serviceAccount** — the trigger's service account ALSO needs
  `roles/workflows.invoker` on this workflow; grant both in the chart.

## After deploying

Create the GcpEventarcTrigger with `destination.workflow` referencing
this workflow's `workflow_id` output — events then start executions
automatically.

## When to choose something else

For request/response orchestration invoked directly (not by events),
start from the **HTTP Orchestrator** preset.
