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
edit (both directions live-verified).

Know the limit of declarative here: an out-of-band
`gcloud tasks queues pause` is NOT corrected by the next apply. The
provider tracks `desiredState` as a config-only value and never reads
the queue's live dispatch state back, so an apply whose spec value did
not change plans zero changes and leaves the queue paused
(live-verified). Two consequences: a console/CLI pause during an
incident safely survives routine deploys — and it also means nobody's
apply will resume the queue for you. Resume the way you paused, or
flip the spec `PAUSED` → apply → `RUNNING` → apply so the value change
triggers the provider's resume call.

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
