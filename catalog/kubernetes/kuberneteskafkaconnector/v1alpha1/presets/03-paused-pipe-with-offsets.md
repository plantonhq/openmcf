# Paused pipe with offsets preset

The operational-lifecycle preset: a connector declared `paused` with
both offset ConfigMap targets wired. Use this shape when a pipe must
exist but not move data yet (a cutover waiting on a downstream
system), or as the template for offset surgery on any existing
connector.

The mechanics worth internalizing:

- **`paused` vs `stopped`**: paused keeps tasks allocated and resumes
  cheaply; stopped deallocates tasks — and is the REQUIRED state
  before offsets can be altered.
- **Offsets are declared targets, annotated verbs.** The spec names
  the ConfigMaps; nothing happens until the resource carries the
  `strimzi.io/connector-offsets` annotation. `list` writes the
  connector's current offsets into `audit-mirror-offsets`
  (non-disruptive, run it anytime); `alter` applies offsets from
  `audit-mirror-offsets-override` while the connector is stopped —
  the replay, skip-poison-record, and migration-cutover mechanism.
  The operator removes the annotation when the operation completes.
- **Sinks track position as a consumer group** (`connect-<name>`), so
  the listed offsets for a sink are consumer offsets; a source's are
  its own source-position records. The ConfigMap format differs
  accordingly — list first, edit, then alter.

A stock-image MirrorSource pipe stands in for real sink classes (S3,
Elasticsearch, warehouse connectors — those arrive through the
Connect cluster's `image`, `plugins`, or `build` arms, since the
stock image carries only the MirrorMaker 2 connector set); the
lifecycle and offset mechanics taught here are identical for any
class.

See [03-paused-pipe-with-offsets.yaml](./03-paused-pipe-with-offsets.yaml)
for the manifest.
