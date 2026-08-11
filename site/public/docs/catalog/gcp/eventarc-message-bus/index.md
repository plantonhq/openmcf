---
title: "Eventarc Message Bus"
description: "Eventarc Message Bus deployment documentation"
icon: "package"
order: 100
componentName: "gcpeventarcmessagebus"
---

# GCP Eventarc Message Bus

Creates an Eventarc ADVANCED message bus with its satellites — the enterprise eventing hub: Google API sources publish events INTO the bus, enrollments select messages with CEL expressions, and pipelines deliver them OUT (to HTTP endpoints, Pub/Sub topics, Workflows, or other buses) with per-pipeline auth, payload conversion, transformation, and retries. Distinct from Eventarc Standard (GcpEventarcTrigger): Standard routes single event types point-to-point; Advanced is the many-sources-many-destinations hub with mediation.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Message bus** -- an `eventarc.MessageBus` (the central conduit)
- **Google API sources** -- one `eventarc.GoogleApiSource` per spec entry, auto-wired to THIS bus
- **Pipelines** -- one `eventarc.Pipeline` per spec entry (destination, auth, payload formats, mediation, retries)
- **Enrollments** -- one `eventarc.Enrollment` per spec entry, binding CEL-selected messages to a sibling pipeline
- **Eventarc API enablement** -- `eventarc.googleapis.com` enabled in the target project (never disabled on destroy)

## Before You Deploy

### Planton Setup

- **GCP Provider Connection** -- an active connection in the Connect module with credentials for the target GCP project.
- **Planton Runner** -- required when using Runner-based credential delivery.

### GCP Project

- **A GCP project** to host the family (directly or via a GcpProject reference).
- **IAM**: the deploying identity needs `roles/eventarc.admin`.
- **Region**: Eventarc Advanced serves a [subset of regions](https://cloud.google.com/eventarc/docs/locations) — the API rejects unsupported ones at create time.

## Deploy

### Console

Open the deployment store, find **GCP Eventarc Message Bus**, and click **Deploy**. Start from the **Bus with Topic Pipeline** preset in the [Presets](#presets) tab.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpEventarcMessageBus
metadata:
  name: central-bus
  org: acme-corp
  env: prod
spec:
  location: us-central1
  logSeverity: INFO
  pipelines:
    - pipelineId: deliver-to-topic
      destination:
        topic:
          value: projects/my-project/topics/downstream
  enrollments:
    - enrollmentId: route-everything
      celMatch: "true"
      pipeline: deliver-to-topic
```

```shell
planton apply -f message-bus.yaml
```

### InfraChart

The central eventing hub in one chart: this bus with an audit-log API source, a GcpPubSubTopic destination for the analytics stream, and a GcpWorkflow destination for the remediation stream — one hub, two consumers.

## Key Configuration

**googleApiSources** -- publish Google-service events into the bus. Each source is AUTO-WIRED to this kind's own bus — a source feeding another bus belongs to that bus's manifest.

**enrollments** -- the routing table: `celMatch` selects messages (evaluated against the CloudEvent; `"true"` routes everything), `pipeline` names a sibling pipeline by its `pipelineId` (validated at manifest time — dangling ids never reach the API).

**pipelines** -- one destination each (API truth): an HTTPS endpoint via a VPC network attachment, a Pub/Sub topic, a Workflow, or another bus. Per-pipeline `authentication` (OIDC for Cloud Run/IAP, OAuth for Google APIs), avro/json/protobuf payload conversion, a single CEL `mediationTransformationTemplate`, and 1–100 attempt retries with 1–600s backoff.

**logSeverity** -- platform-log verbosity per resource. `INFO` while onboarding sources and pipelines, then scale back.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **GcpProject** (optional) | `projectId` | `status.outputs.project_id` |
| **GcpPubSubTopic** (pipeline arm) | `pipelines[].destination.topic` | `status.outputs.topic_id` |
| **GcpWorkflow** (pipeline arm) | `pipelines[].destination.workflow` | `status.outputs.workflow_id` |
| **GcpServiceAccount** (auth) | `pipelines[].authentication.*.serviceAccount` | `status.outputs.email` |
| **GcpKmsKey** (optional) | `cryptoKey` (bus/source/pipeline) | `status.outputs.key_id` |

### What This Component Provides

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `message_bus_name` | Full bus resource name | Another bus's `pipelines[].destination.messageBus` (cross-bus chaining); external publishers |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Bus with topic pipeline** -- everything into the bus, everything out to Pub/Sub; the smallest useful hub. Start from the **Bus with Topic Pipeline** preset.

**Audit fan-out** -- an audit-log API source with CEL-split enrollments routing storage events and IAM events to different pipelines. Start from the **Audit Fan-Out** preset.

## Works With

- [**GCP Pub/Sub Topic**](/cloud-catalog/gcp-pub-sub-topic) -- pipeline topic destinations
- [**GCP Workflow**](/cloud-catalog/gcp-workflow) -- pipeline workflow destinations
- [**GCP Service Account**](/cloud-catalog/gcp-service-account) -- pipeline auth identities
- [**GCP KMS Key**](/cloud-catalog/gcp-kms-key) -- CMEK on bus and satellites
- [**GCP Eventarc Trigger**](/cloud-catalog/gcp-eventarc-trigger) -- Eventarc Standard, for single-type point-to-point routes
