# AwsSagemakerFeatureGroup — Component Guide

Authored operational judgment for the feature group component: the
design decisions behind the spec's shape, and what to know before
running a feature store in production.

## Design decisions

- **The group's AWS name derives from `metadata.name`.** SageMaker's
  charset (letters, digits, hyphens) accepts the manifest name as-is —
  no separate name field to drift from the resource's identity.
- **`vector_dimension` pairs exactly with the Vector collection type.**
  A CEL rule enforces AWS's own pairing in both directions: a Vector
  without a dimension and a dimension without a Vector both fail at
  manifest time, not at apply.
- **At least one store, and `enabled: false` does not count.** A
  feature group without any store stores nothing; the CEL rule requires
  an enabled online store or an offline store before AWS ever sees the
  request.
- **The special features must be members of the schema.** AWS requires
  `record_identifier_feature_name` and `event_time_feature_name` to be
  defined in `feature_definitions`; two CEL rules catch the mismatch at
  manifest time. The reserved names (`is_deleted`, `write_time`,
  `api_invocation_time`) are excluded the same way.
- **Provisioned capacity units pair with Provisioned mode.** The spec
  rejects capacity units under on-demand billing, mirroring the
  provider's create behavior.

## Running a feature store in production

- **Treat the schema as forever.** Everything except the online TTL and
  throughput is create-time only — changing a feature, a store, or even
  the description replaces the group. Design the schema before data
  lands.
- **TTL is your one online lever.** Records hard-delete at
  `EventTime + ttl`; size it to serving freshness, and adjust it freely
  — it is the only online-store setting that updates in place.
- **Offline data outlives the group.** Deleting the group leaves its S3
  objects in place by AWS design — budget for the bucket's lifecycle
  separately, and clean it yourself when a group's data must go.
- **Collections need InMemory storage.** List / Set / Vector features
  require the online store's InMemory tier (server-enforced), which
  bills an at-rest floor — reach for vectors only when embeddings truly
  serve online.
- **Start on-demand.** Throughput updates in place, so move to
  provisioned capacity once traffic is predictable rather than guessing
  up front.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
