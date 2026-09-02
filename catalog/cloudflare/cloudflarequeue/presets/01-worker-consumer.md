# Queue with a Worker (push) consumer

A queue consumed automatically by a Worker — the most common Queues setup. Cloudflare
invokes the Worker with batches of messages as they arrive and autoscales the
number of concurrent invocations.

## When to use

- You have a Worker that should process messages off a queue (emails, webhooks,
  uploads, fan-out work) without blocking the request path.

## Key choices

- `consumer.scriptName`: reference the consuming `CloudflareWorker` so the graph
  deploys it before wiring the consumer.
- `consumer.deadLetterQueue`: messages that exhaust `maxRetries` are sent here
  instead of being dropped — point it at another `CloudflareQueue`.
- `consumer.settings.maxConcurrency`: leave unset to autoscale (recommended); set
  a number to cap concurrent invocations (e.g. to protect a rate-limited
  downstream).

## Placeholders

| Placeholder | Description |
|---|---|
| `0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d` | 32-character Cloudflare account ID |
| `<consumer-worker-name>` | Name of the Worker that consumes this queue |

## Producing to the queue

Reference this queue from a `CloudflareWorker` `queues` producer binding:

```yaml
queues:
  - name: ORDERS
    queueName:
      valueFrom:
        kind: CloudflareQueue
        name: orders-queue
        fieldPath: status.outputs.queue_name
```
