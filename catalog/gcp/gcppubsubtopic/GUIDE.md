# GcpPubSubTopic Guide

The judgment this guide protects: the topic is the permanent name in an
event architecture — publishers hardcode it, subscriptions attach to it,
and almost everything sharp about it (name, tags, region policy) is
decided at create time.

## Names are forever

`topicName` is ForceNew, and replacing a topic detaches every
subscription (they survive but drain to empty, receiving nothing).
Treat topic names like API endpoints: versioned when the contract
breaks (`orders-v2`), never renamed in place.

## Topic retention vs subscription retention

`messageRetentionDuration` on the topic is the replay lever: it lets ANY
subscription — including one created tomorrow — seek back through the
window. Subscription-level retention only covers that subscription's own
backlog. Pay for topic retention when replay-for-new-consumers is part
of the design (event sourcing, backfill); skip it for fire-and-forget
fan-out where subscriptions manage their own tails.

## Schema validation: attach early, pin deliberately

Attaching `schemaSettings` after a topic is live is a contract change on
publishers — messages that used to pass start bouncing. Attach from day
one when a schema exists. The revision bounds
(`firstRevisionId`/`lastRevisionId`, referencing the schema's
`revision_id` output) are the migration tool: pin both to one revision
to freeze the contract during a coordinated rollout, then widen.

## Ingestion sources replace publisher fleets

One `ingestionDataSourceSettings` source (Kinesis, MSK, Event Hubs,
Cloud Storage, Confluent) turns the topic into a managed importer — no
bridge service to run. The federation setup (AWS role ARN + GCP service
account) is the real work; the reference fields wire the GCP side. Run
ONE source per topic: multiple sources make provenance and ordering
unreasonable for consumers.

## Transforms belong at the boundary that owns the change

A topic-level transform rewrites what EVERY subscription sees — use it
for redaction and normalization that define the canonical event. A
per-consumer reshape belongs on the subscription instead. Each step is
exactly one arm: `javascriptUdf` for mechanical rewrites, `aiInference`
(a Vertex AI endpoint, referenced or literal) for classification and
enrichment a function cannot express. The `disabled` flag is the staging
lever — park a transform in position, flip it live later.

## Tags bind at create time only

`resourceManagerTags` exist for org policies and IAM conditions that key
on tags. Changing them later REPLACES the topic — and takes every
subscription attachment with it. Decide them with the name. Labels stay
the mutable organizational lever.

## Teardown

Deleting a topic strands its subscriptions (alive, empty). For the
event-bus topics everything publishes into, set
`deletionPolicy: PREVENT`; the default DELETE is right for
per-environment topics that come and go with their stack.
