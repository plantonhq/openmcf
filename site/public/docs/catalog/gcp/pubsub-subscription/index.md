---
title: "Pub/Sub Subscription"
description: "Pub/Sub Subscription deployment documentation"
icon: "package"
order: 100
componentName: "gcppubsubsubscription"
---

# GCP Pub/Sub Subscription

Deploys a Pub/Sub subscription attached to a topic with configurable delivery method (pull, push, BigQuery, or Cloud Storage), acknowledgement deadlines, message ordering, exactly-once delivery, dead-letter handling, and retry policies. The component integrates with Planton's Provider Connections for GCP credential management and supports ValueFromRef wiring to GCP projects, topics, and GCS buckets.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Pub/Sub Subscription** -- a named subscription resource in the specified GCP project, attached to the referenced topic, with the configured delivery method and message handling policies
- **Push Configuration** -- created only when `pushConfig` is set; configures HTTP POST delivery to an HTTPS endpoint with optional OIDC authentication and unwrapped payload mode
- **BigQuery Configuration** -- created only when `bigqueryConfig` is set; streams messages directly into a BigQuery table with optional schema mapping and metadata columns
- **Cloud Storage Configuration** -- created only when `cloudStorageConfig` is set; batches messages into Cloud Storage objects based on size, duration, or message count thresholds
- **Dead Letter Policy** -- created only when `deadLetterPolicy` is set; forwards messages that exceed the maximum delivery attempts to a dead-letter topic
- **Retry Policy** -- created only when `retryPolicy` is set; configures exponential backoff between delivery attempts after a NACK or ack deadline exceeded event
- **GCP Labels** -- resource metadata labels (resource name, kind, organization, environment) applied automatically for tracking and governance

## Before You Deploy

### Planton Setup

- **GCP Provider Connection** -- an active connection in the Connect module with credentials for the target GCP project. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials or browser OAuth authentication modes.

### GCP Project

- **A GCP project** where the subscription will be created. Provide the project ID directly or reference a GcpProject Cloud Resource via ValueFromRef.
- **An existing Pub/Sub topic** to subscribe to. Provide the fully qualified topic ID or reference a GcpPubSubTopic Cloud Resource via ValueFromRef.
- **Cloud Pub/Sub API** enabled in the target project.
- **BigQuery Data Editor role** on the target table (if using BigQuery delivery).
- **Storage Object Creator role** on the target bucket (if using Cloud Storage delivery).

## Deploy

### Console

Open the deployment store, find **GCP Pub/Sub Subscription**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Basic Pull** preset in the [Presets](#presets) tab to pre-populate a minimal pull subscription.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpPubSubSubscription
metadata:
  name: order-processor
  org: acme-corp
  env: prod
spec:
  projectId:
    value: "acme-prod-12345"
  subscriptionName: order-processor
  topic:
    value: "projects/acme-prod-12345/topics/order-events"
```

```shell
planton apply -f pubsub-subscription.yaml
```

This creates a pull subscription with a 10-second ack deadline, 7-day message retention, no dead-letter handling, and no message ordering. A Stack Job tracks the provisioning in real time.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the subscription to a project and topic deployed in the same InfraPipeline:

```yaml
spec:
  projectId:
    valueFrom:
      kind: GcpProject
      name: production-project
      fieldPath: status.outputs.project_id
  topic:
    valueFrom:
      kind: GcpPubSubTopic
      name: order-events
      fieldPath: status.outputs.topic_id
```

The InfraPipeline resolves the dependency graph, deploys the project and topic first, then provisions the subscription with the resolved topic reference.

## Key Configuration

These are the most important decisions when configuring a Pub/Sub subscription. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Delivery method** -- Choose between pull (default, no config needed), push (`pushConfig`), BigQuery (`bigqueryConfig`), or Cloud Storage (`cloudStorageConfig`). Only one delivery method can be active. Pull requires application-side polling; push delivers to an HTTPS endpoint; BigQuery and Cloud Storage write directly without consumer code.

**Acknowledgement deadline** -- Set `ackDeadlineSeconds` (10-600, default 10) to control how long Pub/Sub waits for an acknowledgement before redelivering a message. Set higher for long-running processing tasks. Too low causes unnecessary redeliveries; too high delays retry on actual failures.

**Dead-letter handling** -- Configure `deadLetterPolicy` with a dead-letter topic and `maxDeliveryAttempts` (5-100) to route unprocessable messages instead of retrying indefinitely. The Pub/Sub service account needs Subscriber permissions on this subscription and Publisher permissions on the dead-letter topic.

**Exactly-once delivery** -- Set `enableExactlyOnceDelivery: true` to guarantee each message is delivered exactly once within the ack deadline window. Adds latency to acknowledgement processing. Does not prevent duplicates from publishers sending the same message multiple times.

**Message ordering** -- Set `enableMessageOrdering: true` to deliver messages with the same ordering key in publish order. Immutable after creation. Requires publishers to set an ordering key on messages.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **GcpProject** | `projectId` | `status.outputs.project_id` |
| **GcpPubSubTopic** | `topic` | `status.outputs.topic_id` |
| **GcpPubSubTopic** (optional) | `deadLetterPolicy.deadLetterTopic` | `status.outputs.topic_id` |
| **GcpGcsBucket** (optional) | `cloudStorageConfig.bucket` | `status.outputs.bucket_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `subscription_id` | Fully qualified subscription ID (`projects/{project}/subscriptions/{name}`) | Application configuration, monitoring filters |
| `subscription_name` | Short subscription name | Display, logging, consumer identification |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Basic pull** -- Minimal pull subscription with default ack deadline and retention. Suitable for application-driven message consumption where the subscriber controls the pull rate. Start from the **Basic Pull** preset.

**Push with OIDC** -- Push delivery to an HTTPS endpoint with OIDC-based authentication. Suitable for webhook integrations, Cloud Run services, and serverless backends that receive messages as HTTP POST requests. Start from the **Push with OIDC** preset.

**BigQuery delivery** -- Direct streaming into a BigQuery table with topic schema mapping and metadata columns. Suitable for analytics pipelines, event archival, and real-time dashboards without writing consumer code. Start from the **BigQuery Delivery** preset.

**Dead letter** -- Subscription with exactly-once delivery, 14-day retention, retained ack'd messages, dead-letter routing after 10 attempts, and exponential retry backoff. Suitable for payment processing, order fulfillment, and other workflows that require reliable delivery with poison message isolation. Start from the **Dead Letter** preset.

## Works With

- [**GCP Project**](/cloud-catalog/gcp-project) -- provides the GCP project where the subscription is created
- [**GCP Pub/Sub Topic**](/cloud-catalog/gcp-pub-sub-topic) -- provides the topic that the subscription receives messages from, and optionally the dead-letter topic
- [**GCP GCS Bucket**](/cloud-catalog/gcp-gcs-bucket) -- provides the destination bucket for Cloud Storage delivery