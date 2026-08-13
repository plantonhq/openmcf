# GcpCloudTasksQueue Guide

The judgment this guide protects: a queue is a shock absorber between
producers and a target that must not be overrun. Its settings are the
backpressure contract — and its backlog is real data.

## The name is a seven-day commitment

Queue names are immutable and a deleted queue's ID stays reserved for
up to 7 days. A rename is therefore a replace PLUS a week-long burn on
the old identifier. The spec requires the name explicitly for exactly
this reason — choose one that outlives the current service name.

## Pause is declarative — and that changes incident response

`desiredState: PAUSED` holds every task in the queue while producers
keep enqueueing — the right first move when the target is degraded:
nothing is lost, nothing dispatches, and resuming is a one-line spec
edit. The flip side of declarative: an out-of-band
`gcloud tasks queues pause` is REVERTED on the next apply. Pause in the
manifest, not in the console, or the next deploy resumes dispatch into
the incident.

## Rate limits are the target's contract, not the queue's

`maxDispatchesPerSecond` and `maxConcurrentDispatches` should encode
what the TARGET survives, with headroom for its other callers. Cloud
Tasks also self-throttles on 429/503 — but that is damage control, not
capacity planning.

## Queue-level HTTP targets: enqueue payloads, own routing here

`httpTarget` with an OIDC token and URI override is the modern shape:
producers enqueue bare payloads, the queue owns auth and destination.
That makes endpoint migrations a queue edit instead of a producer
redeploy — and it means this spec, not application code, is where the
security review looks.

## Destroy semantics: the backlog dies with the queue

`deletionPolicy: DELETE` (default) removes the queue AND every task
still in it — undelivered work, gone. Drain first: pause producers, let
dispatch empty the queue, then destroy. `PREVENT` is the posture for
queues whose backlog is money (orders, webhooks). `ABANDON` keeps the
queue running and dispatching unmanaged.
