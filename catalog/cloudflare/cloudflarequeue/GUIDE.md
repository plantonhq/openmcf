# CloudflareQueue guide

Operational judgment for Cloudflare Queues. The README covers what each field is; this covers how the pieces interact.

## One consumer per queue is the platform's shape

A queue has at most one consumer; fan-out means multiple queues (a producer Worker can write to several). The consumer is modeled inside this kind because it has no life without its queue — but it is a separate provider resource underneath, so consumer edits never recreate the queue or lose messages.

## Push and pull are different operational worlds

A `worker` (push) consumer is invoked by Cloudflare automatically — batching, concurrency, and retries all run inside the platform. An `http_pull` consumer moves all of that to your clients: they poll, they ack, and `visibility_timeout_ms` is what stops two clients from processing one message. Choose push unless the consumer genuinely lives outside Cloudflare.

## The dead-letter queue must already exist

The DLQ is named by reference, not created on demand: point it at another CloudflareQueue that is already deployed, or the consumer's create fails. Ordering matters in charts — DLQ first, then the queue that names it.

## Queues ride the Workers Paid plan

Cloudflare has historically gated Queues behind the Workers Paid plan (~$5/mo). On a free account the create call fails at the API, not at validation — budget the entitlement before wiring queues into a design, and remember R2 event notifications deliver INTO a queue, so they inherit this gate too.
