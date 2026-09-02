# CloudflareQueue guide

Operational judgment for Cloudflare Queues. The README covers what each field is; this covers how the pieces interact.

## One consumer per queue is the platform's shape

A queue has at most one consumer; fan-out means multiple queues (a producer Worker can write to several). The consumer is modeled inside this kind because it has no life without its queue — but it is a separate provider resource underneath, so consumer edits never recreate the queue or lose messages.

## Push and pull are different operational worlds

A `worker` (push) consumer is invoked by Cloudflare automatically — batching, concurrency, and retries all run inside the platform. An `http_pull` consumer moves all of that to your clients: they poll, they ack, and `visibility_timeout_ms` is what stops two clients from processing one message. Choose push unless the consumer genuinely lives outside Cloudflare.

## The dead-letter queue must already exist

The DLQ is named by reference, not created on demand: point it at another CloudflareQueue that is already deployed, or the consumer's create fails. Ordering matters in charts — DLQ first, then the queue that names it.

## Queues ride the Workers Paid plan

Cloudflare has historically gated Queues behind the Workers Paid plan. On a free account the create call fails at the API, not at validation — budget the entitlement before wiring queues into a design, and remember R2 event notifications deliver INTO a queue, so they inherit this gate too.

## Adopting a queue that already has a consumer

The queue itself imports cleanly (`{account_id}/{queue_id}`), but its consumer cannot come along: the provider ships no consumer importer, and re-creating the consumer from your declaration is refused while the live one exists — Cloudflare enforces one consumer per queue and answers the duplicate create with 400 code 11004 "already has a consumer" (measured 2026-08-26). To bring an existing queue fully under management, delete its consumer first (dashboard or API) and let the next apply create the declared one — deleting a consumer never touches the queue or its messages. Alternatively, import the queue without declaring a `consumer` and leave the existing consumer managed where it is.
