# GcpCloudTasksQueue - Design Documentation

## GCP Cloud Tasks Overview

Cloud Tasks is a fully managed service that manages the execution, dispatch, and delivery of a large number of distributed tasks. It provides at-least-once delivery with configurable rate limiting and retry policies.

### Core Concepts

- **Queue**: A named entity that manages the dispatch of tasks. Controls rate, retries, and routing.
- **Task**: A unit of work defined as an HTTP request (or App Engine request). Tasks are added to queues and dispatched according to queue configuration.
- **Dispatch**: The act of sending the task's HTTP request to its target handler.

### Task Target Types

1. **HTTP tasks** (modern, recommended): Tasks dispatched to any HTTP endpoint. The target URL, method, headers, and body are defined per-task or overridden at the queue level via `http_target`.
2. **App Engine tasks** (legacy): Tasks dispatched to App Engine request handlers. Routing is based on App Engine service/version/instance; the queue can pin routing for all its App Engine tasks via `app_engine_routing_override`.

## Terraform Resource: google_cloud_tasks_queue

**Provider**: `hashicorp/google ~> 6.0`

### Key Schema Fields

| Field | Type | Required | ForceNew | Description |
|-------|------|----------|----------|-------------|
| `name` | string | yes | yes | Queue name |
| `location` | string | yes | yes | GCP region |
| `project` | string | no | yes | GCP project |

### Nested Blocks

- `rate_limits` -- max_dispatches_per_second, max_concurrent_dispatches, max_burst_size (computed)
- `retry_config` -- max_attempts, max_retry_duration, min_backoff, max_backoff, max_doublings
- `stackdriver_logging_config` -- sampling_ratio
- `http_target` -- http_method, header_overrides, oauth_token, oidc_token, uri_override
- `app_engine_routing_override` -- service, version, instance

### Notable Behaviors

- `max_burst_size` in rate_limits is **computed-only** -- GCP calculates it from `max_dispatches_per_second`; the effective value is exported as a stack output
- `name`, `location`, and `project` are all `ForceNew` -- changing them destroys and recreates the queue
- A deleted queue's ID is reserved by the Cloud Tasks API for up to 7 days; it cannot be reused within that window
- Duration fields (`min_backoff`, `max_backoff`, `max_retry_duration`) use GCP's duration format (e.g., "300s", "0.100s")

## Pulumi Resource: cloudtasks.Queue

**SDK**: `pulumi-gcp/sdk/v9/go/gcp/cloudtasks`

The Pulumi SDK mirrors the provider schema closely. The bridge tracks a newer provider major than the released Terraform 6.x line, so it carries surfaces the Terraform module cannot reach (queue `state`/`desired_state`, client-side `deletion_policy`). The Pulumi module deliberately does not use them -- both engines implement the identical contract:

- `deletion_policy` is pinned to `DELETE` so destroy really deletes the queue on both engines
- `desired_state` is not set; pause/resume is a runtime operation (see below)

## Design Decisions

### Included Features

1. **http_target** -- Queue-level HTTP configuration is the modern Cloud Tasks pattern. It enables:
   - Central authentication config (OIDC/OAuth) for all tasks in a queue
   - URI overrides for consistent routing
   - Header defaults
   This is critical for the Cloud Run + Cloud Tasks integration pattern.

2. **app_engine_routing_override** -- Completes the provider surface for queues dispatching App Engine tasks: the queue pins its tasks to one service/version/instance instead of relying on each task's own routing.

3. **rate_limits** and **retry_config** -- Core queue behavior that most users will configure.

4. **stackdriver_logging_config** -- Operational observability for task dispatch.

### Excluded Features (with reasons)

1. **Pause/resume (`desired_state`/`state`)** -- Present only on provider majors newer than the released 6.x line this catalog builds against, so modeling it would make one engine silently ignore user intent. Pause/resume is also a runtime operation rather than infrastructure definition: use `gcloud tasks queues pause|resume` on a provisioned queue. The surface returns when the catalog moves to the newer provider major.

2. **max_burst_size (as input)** -- Computed-only in both engines. GCP calculates it automatically from `max_dispatches_per_second`; the effective value is exported as the `max_burst_size` stack output.

3. **GCP labels** -- Cloud Tasks queues do NOT support GCP labels (confirmed in both the Terraform schema and the Pulumi SDK). This is a GCP API limitation, not a design choice.

4. **Per-queue IAM (member/binding/policy trio)** -- Resource-scoped IAM glue is deliberately not modeled as kinds; grant `roles/cloudtasks.enqueuer` on the project via `GcpProjectIamMember`.

### Flattening Decisions

The provider resource has several unnecessarily nested structures:

- `http_target.uri_override.path_override.path` -- Flattened to `http_target.uri_override.path`
- `http_target.uri_override.query_override.query_params` -- Flattened to `http_target.uri_override.query_params`
- `http_target.header_overrides[].header.{key,value}` -- Flattened to `http_target.header_overrides[].{key,value}`

These wrapper messages add zero semantic value. The spec provides a cleaner, more intuitive API while the IaC modules handle the mapping back to the provider's nested structure. The modules also only send the nested blocks when their value is set, which avoids the provider's known perpetual-diff on an empty `query_override` block.

### Naming and Project Contracts

- `queue_name` is **deliberately required** (never derived from `metadata.name`): a deleted queue's ID is reserved for up to 7 days, so the name deserves an explicit, stable choice.
- `project_id` is optional: when omitted, the queue is created in the provider's default project (the ambient-project contract shared across the GCP catalog).

### StringValueOrRef Fields (Infra-Chart Composability)

Three StringValueOrRef fields enable infra-chart composition:

1. **`project_id`** -> GcpProject (`status.outputs.project_id`)
2. **`oauth_token.service_account_email`** -> GcpServiceAccount (`status.outputs.email`)
3. **`oidc_token.service_account_email`** -> GcpServiceAccount (`status.outputs.email`)

This enables a common infra-chart pattern: create a dedicated service account for task dispatch, then wire its email into the queue's authentication config.

## Cloud Tasks Pricing

- **First 1 million tasks/month**: Free
- **Beyond that**: per-million-task pricing tiers

No charges for queue creation or management. Costs are driven entirely by task creation volume.

## Common Integration Patterns

### Cloud Run + Cloud Tasks

The most common modern pattern:

```
Producer Service -> Cloud Tasks Queue -> Cloud Run Handler Service
                    (rate limited)      (OIDC authenticated)
```

The queue manages dispatch rate and retries. OIDC tokens authenticate requests to the Cloud Run service. This is configured entirely at the queue level via `http_target`.

### Cloud Functions + Cloud Tasks

Similar to Cloud Run, but targeting Cloud Functions:

```
Event Source -> Cloud Tasks Queue -> Cloud Function
               (retry config)     (OIDC authenticated)
```

### Cloud Scheduler + Cloud Tasks

For scheduled task batches:

```
Cloud Scheduler Job -> Creates Tasks -> Cloud Tasks Queue -> HTTP Handler
(cron trigger)        (via API)        (rate limited)       (authenticated)
```

## Comparison: Cloud Tasks vs Pub/Sub

| Feature | Cloud Tasks | Pub/Sub |
|---------|------------|---------|
| Pattern | Task dispatch (1:1) | Message fan-out (1:N) |
| Delivery | At-least-once | At-least-once |
| Rate limiting | Built-in | Not built-in |
| Retry | Configurable per-queue | Per-subscription |
| Target | HTTP endpoint | Push/Pull/BigQuery/GCS |
| Task dedup | Task naming | Message ID |
| Use case | Background jobs, API calls | Event streaming, pub/sub |

**Rule of thumb**: Use Cloud Tasks when you want to control *when* and *how fast* work is done. Use Pub/Sub when you want to broadcast events to multiple consumers.
