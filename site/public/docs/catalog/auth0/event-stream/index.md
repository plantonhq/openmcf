---
title: "Event Stream"
description: "Event Stream deployment documentation"
icon: "package"
order: 100
componentName: "auth0eventstream"
---

# Auth0 Event Stream

Deploys an Auth0 Event Stream that delivers real-time Auth0 events to an external destination. Supports AWS EventBridge for serverless event processing and HTTPS webhooks for custom endpoint delivery, with configurable event type subscriptions. Integrates with Planton's Auth0 Provider Connection for credential management.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Auth0 Event Stream** -- a stream resource configured with the specified destination type, event subscriptions, and destination-specific settings

For EventBridge destinations, Auth0 creates a partner event source in the target AWS account that must be associated with an EventBridge event bus to begin receiving events. For webhook destinations, Auth0 delivers event payloads as HTTP POST requests to the configured HTTPS endpoint with the specified authorization credentials.

## Before You Deploy

### Planton Setup

- **Auth0 Provider Connection** -- an active connection in the Connect module with Auth0 domain, client ID, and client secret. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credential authentication.

### Auth0 Account

- **An Auth0 tenant** with Event Streams enabled.
- **An AWS account** with permissions to accept partner event sources, if using the EventBridge destination type. You will need the 12-digit AWS account ID and target region.
- **A publicly accessible HTTPS endpoint** if using the webhook destination type. The endpoint must respond to POST requests within 10 seconds and support either Basic or Bearer token authentication.

## Deploy

### Console

Open the deployment store, find **Auth0 Event Stream**, and click **Deploy**. The creation wizard walks you through environment and connection configuration and spec fields.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: auth0.planton.dev/v1
kind: Auth0EventStream
metadata:
  name: login-events
  org: acme-corp
  env: prod
spec:
  destinationType: webhook
  subscriptions:
    - authentication.success
    - authentication.failure
  webhookConfiguration:
    webhookEndpoint: "https://api.example.com/webhooks/auth0"
    webhookAuthorization:
      method: bearer
      token: "your-secret-token"
```

```shell
planton apply -f auth0-event-stream.yaml
```

This creates an event stream that delivers authentication success and failure events to the specified webhook endpoint using bearer token authorization. No EventBridge configuration is included. A Stack Job tracks the provisioning in real time.

## Key Configuration

These are the most important decisions when configuring an Auth0 Event Stream. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Destination type** -- The `destinationType` field determines where events are delivered. Use `eventbridge` for AWS-native event processing with Lambda, Step Functions, or SIEM integrations. Use `webhook` for custom HTTPS endpoints. EventBridge configurations cannot be updated after creation -- any change forces resource recreation.

**Event subscriptions** -- The `subscriptions` array specifies which event types the stream delivers. Common types include `authentication.success`, `authentication.failure`, `user.created`, `user.updated`, and `api.authorization.success`. At least one subscription is required. Subscribe only to the events you need to reduce noise in downstream systems.

**Webhook authorization** -- The `webhookConfiguration.webhookAuthorization.method` field sets how Auth0 authenticates with your endpoint. Use `bearer` with a securely generated token for most integrations, or `basic` with username and password for systems that require HTTP Basic authentication. Generate tokens with `openssl rand -base64 32`.

**EventBridge region** -- The `eventbridgeConfiguration.awsRegion` field determines which AWS region receives the partner event source. Choose the region where your event processing infrastructure runs to minimize latency. The `awsAccountId` must be a valid 12-digit AWS account number.

**Immutability constraints** -- EventBridge configurations are immutable after creation. Any change to `awsAccountId` or `awsRegion` forces the event stream to be destroyed and recreated, which generates a new partner event source name. Webhook configurations can be updated in place, including the endpoint URL and authorization settings.

## Outputs and Dependencies

### What This Component Consumes

This component has no foreign key dependencies.

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `id` | Unique Auth0 event stream identifier | Auth0 API operations |
| `name` | Event stream name derived from metadata | Monitoring dashboards |
| `status` | Current stream status (`active`, `suspended`, `disabled`) | Health checks, alerting |
| `destination_type` | Configured destination (`eventbridge` or `webhook`) | Downstream routing logic |
| `created_at` | ISO 8601 creation timestamp | Audit logs, lifecycle tracking |
| `updated_at` | ISO 8601 last-updated timestamp | Change tracking |
| `subscriptions` | List of subscribed event types | Downstream event filtering |
| `aws_partner_event_source` | AWS partner event source name (EventBridge only) | EventBridge event bus association |

## Common Patterns

No presets are available yet. Configure directly using the fields documented in the [API Explorer](#api-explorer) tab.

## Works With

This component operates independently and does not reference other deployment components.