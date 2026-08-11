# GcpPubSubSubscription Guide

The judgment this guide protects: the subscription is where delivery
semantics live — ordering, exactly-once, dead-lettering, and the choice
of who does the work (your client, an HTTPS endpoint, BigQuery, or
Cloud Storage). Most production incidents around Pub/Sub are really
subscription-configuration incidents.

## Choose the delivery mode by who owns the consumer

- **Pull** (no delivery config): you run the consumer; maximum control,
  client-library batching. The default for services.
- **Push**: Pub/Sub calls your HTTPS endpoint — the serverless pattern.
  Reference a `GcpCloudRun` service and pair with `oidcToken` (the
  service account needs `run.invoker`); unauthenticated push is for
  public webhooks only.
- **BigQuery / Cloud Storage**: no consumer at all — Pub/Sub writes
  straight to the table or bucket. Analytics and archival, not
  application logic.

Exactly one may be set; switching modes later is an in-place update but
a semantic migration for whatever consumed before.

## The immutables

`subscriptionName`, `filter`, and `enableMessageOrdering` are ForceNew —
and replacing a subscription DROPS its backlog. Decide the filter up
front (an over-broad filter can be narrowed only by replacement), and
decide ordering before traffic exists: ordered delivery serializes per
ordering key, which is a throughput trade you cannot toggle.

## Dead-letter before you need it

Without `deadLetterPolicy`, a poison message redelivers forever and can
head-of-line-block an ordered subscription. Wire the dead-letter topic
on day one (5–10 `maxDeliveryAttempts` is the practical band) and put a
monitoring subscription on the dead-letter topic itself — an unwatched
dead-letter topic is where failures go to be forgotten. The Pub/Sub
service agent needs Subscriber here and Publisher on the dead-letter
topic; grant both with the additive `GcpProjectIamMember` pattern.

## Retention, replay, and the expiration trap

`expirationPolicy` defaults to 31-day auto-delete on INACTIVITY — a
subscription for a quarterly batch job silently vanishes. Set
`ttl: ""` for anything that must outlive quiet periods.
`retainAckedMessages` + `messageRetentionDuration` is the replay lever
(seek back through acknowledged history); it is also the lever that
makes storage grow, so scope it to subscriptions that genuinely replay.

## Transforms reshape for THIS consumer

A subscription transform changes what this consumer sees without
touching the topic's canonical form. Each step is exactly one arm:
`javascriptUdf` for mechanical reshapes, `aiInference` (a Vertex AI
endpoint, referenced or literal) for per-consumer enrichment. If every
subscription needs the same change, it belongs on the topic instead.

## Tags bind at create time only

`resourceManagerTags` changes REPLACE the subscription — backlog
included. Decide them with the name; labels stay the mutable lever.

## Teardown

`deletionPolicy: DELETE` (the default) drops the unacknowledged backlog
immediately. For a consumer whose backlog is the system of record until
processed, set `PREVENT` — a destroy that would silently discard queued
work should fail loudly instead. `ABANDON` keeps the subscription
accumulating in GCP, unmanaged: a handover lever, and otherwise a
storage-bill leak.
