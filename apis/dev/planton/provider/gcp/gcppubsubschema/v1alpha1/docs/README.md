# GcpPubSubSchema — Deep Dive

## The problem this resource solves

Without schemas, a Pub/Sub topic is a byte pipe: any producer can publish anything, and contract violations surface as processing failures deep inside consumers — often hours later, in a different team's service. Pub/Sub schemas invert that: the topic validates every published message against a typed contract and rejects non-conforming payloads at the publisher, where the producing team sees the error immediately and owns the fix.

The schema is deliberately a **separate, shareable resource** rather than inline topic configuration. One contract governs many topics (per-environment copies of the same stream, fan-out topics carrying the same event), and the contract evolves in one place through revisions.

## Where it sits in the composition

- **GcpPubSubSchema** — this resource: the typed definition (`AVRO` or `PROTOCOL_BUFFER`).
- **GcpPubSubTopic** — attaches the schema through `schemaSettings.schema`, a `StringValueOrRef` whose default reference resolves this schema's `status.outputs.schema_id` (`projects/{p}/schemas/{name}`), together with the message encoding (`JSON` or `BINARY`).
- **GcpPubSubSubscription** — downstream, benefits indirectly: BigQuery delivery's `useTopicSchema` maps message fields to table columns, and Cloud Storage delivery's `avroConfig.useTopicSchema` writes Avro files in the schema's layout.

## Lifecycle contract

| Property | Behavior |
|---|---|
| `schemaName`, `projectId` | Immutable (ForceNew) — renaming replaces the resource |
| `definition` | Mutable — each change commits a new **revision** in place |
| `type` | Accepted at create; keep it stable across revisions (a mid-life type flip invalidates existing publishers) |
| Revision retention | At most 20 revisions per schema; beyond that, commits fail until old revisions are deleted manually (`gcloud pubsub schemas delete-revision`) |
| Deletion | Deleting a schema topics still reference leaves them validating against the `_deleted-schema_` sentinel — publishes fail |

Teardown ordering follows: topics (or their schema attachments) first, schema last. In dependency-graph terms the schema behaves like a KMS key or a network — long-lived shared infrastructure beneath ephemeral consumers.

## Validation semantics

A topic with `schemaSettings` validates every publish:

- **`encoding: JSON`** — the message payload must be the JSON rendering of the schema (Avro JSON or protojson for protobuf schemas). Human-debuggable; slightly larger on the wire.
- **`encoding: BINARY`** — the payload must be Avro binary or serialized protobuf. Compact; suits high-volume streams whose producers already serialize natively.

Messages that fail validation are rejected with `INVALID_ARGUMENT` at publish time — they never enter the topic, never reach subscribers, and never poison a delivery pipeline.

## Revision evolution rules of thumb

Attached topics accept messages conforming to **any available revision**, so a new revision widens the accepted set. Evolve additively:

- **Avro**: add fields with defaults; never remove or retype existing fields consumers rely on.
- **Protobuf**: add new tagged fields; never renumber or repurpose existing tags.

An incompatible revision does not break the topic — it breaks the *consumers*, silently, which is exactly the failure mode schemas exist to prevent. Review definition changes like API changes.

## Deliberately not modeled (recorded reasons)

- **`deletion_policy`** — a client-side provider lever (`DELETE`/`PREVENT`/`ABANDON`) that conflicts with Planton-managed destroy semantics; catalog-wide skip.
- **`revision_id` output** — the released `google ~> 6.x` line does not expose the committed revision ID as a resource attribute (the provider's unreleased line adds it). Exporting it from one engine only would break output parity; revisit when the catalog's provider line carries it.
- **Per-schema IAM (`google_pubsub_schema_iam_*`)** — resource-scoped IAM stays out of the catalog pending concrete pull (the additive project-level grant, `GcpProjectIamMember`, covers the real cases).

## Provider mapping

Maps to `google_pubsub_schema` (`google/services/pubsub/resource_pubsub_schema.go`): `schemaName` → `name` (ForceNew), `type` → `type`, `definition` → `definition` (update = commit revision), `projectId` → `project` (ForceNew, ambient fallback when empty).
