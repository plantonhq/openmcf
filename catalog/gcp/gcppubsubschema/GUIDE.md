# GcpPubSubSchema Guide

The judgment this guide protects: a schema is a CONTRACT shared by every
topic that attaches it and every publisher that writes to those topics.
Its own resource is tiny — the expensive mistakes are contract evolution
done casually and teardown done in the wrong order.

## One schema, many topics

Resist the reflex to create one schema per topic. The schema is the
shareable statement of what an event IS; topics are delivery channels.
When five topics carry `OrderEvent`, one schema validates all five and
the contract evolves in exactly one place. Name schemas after the event
type (`order-events`), never after the topic consuming them.

## Evolve additively, and mind the 20-revision ceiling

Changing `definition` commits a new revision in place — no replacement,
no downtime. Two sharp edges:

- Attached topics accept messages conforming to ANY available revision
  (unless the topic pins a revision range), so a widening change is
  immediately live for every publisher. Evolve additively — add optional
  fields, never retype or remove — or old consumers break silently.
- A schema holds at most 20 revisions. Beyond that, commits FAIL until
  old revisions are deleted manually (`gcloud pubsub schemas
  delete-revision`). A CI pipeline that commits a revision per deploy
  hits this ceiling in twenty deploys; commit revisions when the
  contract changes, not when the code does.

## Revision pinning is the consumer's safety valve

The `revision_id` output exists for exactly one composition: a topic's
`schemaSettings.firstRevisionId`/`lastRevisionId` referencing this
schema pins validation to the revision a deploy produced. Pin both
bounds to freeze the contract exactly during a coordinated migration;
leave them open for the normal any-revision posture.

## AVRO unless you have a reason

`AVRO` reads in pull requests, carries logical types, and is what
BigQuery delivery and Cloud Storage Avro export understand natively.
Choose `PROTOCOL_BUFFER` only when publishers already serialize protobuf
and binary encoding is the point. Keep the type stable for the schema's
life — a type flip invalidates every existing publisher.

## Teardown order matters

Deleting a schema while topics still reference it leaves those topics
validating against the `_deleted-schema_` sentinel — every publish then
FAILS. Detach or destroy the topics first, schema last. For a schema
many topics depend on, set `deletionPolicy: PREVENT` and make the
mistake impossible; `ABANDON` (leave it serving, unmanaged) is the
handover lever, not a cleanup.
