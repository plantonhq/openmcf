# Queue on Cloudflare

Deploys a Cloudflare Queue -- a managed, guaranteed-delivery message queue for Cloudflare Workers. Producers (a Worker's `queues` binding or an R2 bucket's event notifications) write messages; a single consumer reads them, decoupling the two so a slow or offline consumer never drops a producer's message. The consumer is modeled inline: leave it unset to start with producers only, or attach a push (Worker) or HTTP-pull consumer. Queues are account-scoped and integrate with Planton's Provider Connections for Cloudflare credential management.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Queue** -- an account-scoped, named message queue with the chosen delivery settings
- **Consumer** (optional) -- a push (Worker) or HTTP-pull consumer provisioned as its own provider resource, so editing it never recreates the queue
- **Cloudflare Labels** -- resource metadata applied for organization and environment tracking

## Before You Deploy

### Planton Setup

- **Cloudflare Provider Connection** -- an active connection in the Connect module with a Cloudflare API token that has Queues edit access (and Workers Scripts access when attaching a Worker consumer). Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline API token authentication.

### Cloudflare Account

- **Queues enabled** -- Cloudflare Queues must be available on the account.
- **Account-level access** -- queues are created at the account level, so the API token must be scoped to the account.

## Deploy

### Console

Open the deployment store, find **Queue on Cloudflare**, and click **Deploy**. The creation wizard captures the owning account and queue name, optional delivery settings (delay, pause, retention), and an optional consumer. Leave the consumer as **None** to provision a producer-only queue you can wire a consumer to later.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: cloudflare.planton.dev/v1
kind: CloudflareQueue
metadata:
  name: orders-events
  org: acme-corp
  env: prod
spec:
  accountId: "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4"
  queueName: orders-events
  consumer:
    type: worker
    scriptName:
      valueFrom:
        kind: CloudflareWorker
        name: orders-worker
        fieldPath: status.outputs.script_name
```

```shell
planton apply -f cloudflare-queue.yaml
```

This creates a queue consumed by the `orders-worker` Worker. A Stack Job tracks the provisioning in real time.

## Key Configuration

These are the most important decisions when configuring a queue. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Queue Name (`queueName`)** -- The name producers and the consumer address the queue by. Editable in place, but every producer and consumer binding that references it must be updated to match, so choose a stable name.

**Consumer Type (`consumer.type`)** -- `worker` invokes a Worker automatically with batches (requires `scriptName`); `http_pull` is read by external clients over the REST API. Omit the consumer entirely to start with producers only.

**Delivery Settings (`settings`)** -- A global `deliveryDelay`, a `deliveryPaused` switch, and `messageRetentionPeriod` (how long an unconsumed message survives before it is dropped). Size retention to your consumer's worst-case downtime.

**Consumer Batching (`consumer.settings`)** -- `batchSize`, `maxRetries`, and `retryDelay` apply to both consumer types; `maxConcurrency` and `maxWaitTimeMs` are worker-only; `visibilityTimeoutMs` is HTTP-pull-only.

## Outputs and Dependencies

### What This Component Consumes

A queue with a worker consumer references a **CloudflareWorker** (via `consumer.scriptName`). The optional dead-letter queue references another **CloudflareQueue** (via `consumer.deadLetterQueue`). Both accept a literal name or a ValueFromRef.

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `queue_id` | The Cloudflare-assigned identifier of the queue | Verification, dashboards |
| `queue_name` | The queue name (echoed) | Referenced by a Worker's `queues` producer binding, an R2 bucket's event notifications, or another queue's dead-letter queue |
| `created_on` / `modified_on` | Provisioning timestamps | Auditing |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Event-driven Worker** -- A queue with a worker consumer that processes events a producer Worker or R2 bucket enqueues. Tune batch size and concurrency to the consumer's throughput.

**Dead-letter capture** -- Point `deadLetterQueue` at a second queue so messages that exhaust their retries are captured for inspection instead of dropped.

**Producer-only buffer** -- A queue with no consumer yet: producers enqueue immediately, and you attach a consumer once the downstream is ready.

## Works With

- [**Worker on Cloudflare**](/cloud-catalog/cloudflare-worker) -- produces to a queue (a `queues` binding) and/or consumes one (as the queue's worker consumer)
- [**R2 Bucket on Cloudflare**](/cloud-catalog/cloudflare-r2-bucket) -- emits event notifications to a queue
