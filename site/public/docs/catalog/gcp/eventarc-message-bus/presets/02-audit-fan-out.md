---
title: "Audit Fan-Out"
description: "One hub, two consumers: a Google API source feeds the bus, and CEL enrollments split the stream — storage events to an analytics topic, everything else to an ops workflow."
type: "preset"
rank: "02"
presetSlug: "02-audit-fan-out"
componentSlug: "eventarc-message-bus"
componentTitle: "Eventarc Message Bus"
provider: "gcp"
icon: "package"
order: 2
---

# Audit Fan-Out

One hub, two consumers: a Google API source feeds the bus, and CEL
enrollments split the stream — storage events to an analytics topic,
everything else to an ops workflow.

## What it configures

- A `googleApiSources` entry publishing Google-service events into the
  bus (auto-wired — never someone else's bus).
- Two pipelines (topic + authenticated workflow) and two complementary
  CEL enrollments covering the whole stream.

## Adjust before deploying

- **celMatch pair** — the two expressions shown are complementary;
  keep them exhaustive or accept that unmatched messages go nowhere.
- **destination.topic / destination.workflow** — reference your
  GcpPubSubTopic `topic_id` and GcpWorkflow `workflow_id` outputs.
- **authentication.googleOidc.serviceAccount** — the Eventarc service
  agent needs `roles/iam.serviceAccountTokenCreator` on it; the account
  needs `roles/workflows.invoker` on the workflow.

## After deploying

Watch platform logs at INFO while the slices stabilize — per-message
delivery failures (auth, format) surface there, not at apply time.

## When to choose something else

For a single consumer, start from the **Bus with Topic Pipeline**
preset — fan-out earns its complexity only with multiple slices.
