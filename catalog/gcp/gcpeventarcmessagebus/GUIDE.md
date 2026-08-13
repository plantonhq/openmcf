# GcpEventarcMessageBus Guide

Operational judgment for running Eventarc Advanced as code — the things
the spec reference cannot tell you.

## Advanced is not Standard with more knobs

Eventarc Standard (GcpEventarcTrigger) binds ONE event type to ONE
destination — simple, per-route, mature. Advanced is a HUB: sources
publish everything, enrollments select, pipelines deliver with
conversion and mediation. Choose Advanced when multiple consumers read
overlapping slices of the same event stream, or when messages need
transformation in flight. A single Pub/Sub-to-Cloud-Run route on the
bus is over-engineering — use a trigger.

## Region availability is decided by the API, not the provider

The provider schema accepts any location string; Eventarc Advanced
serves a SUBSET of regions and the create call rejects the rest. Check
the locations page before choosing, and expect the availability set to
grow — a region that failed last quarter may work today.

## The bus is a same-project hub

API sources must live in the bus's project; pipelines and the bus must
be same-project too (provider-documented). Cross-project flows enter
through the pipeline's cross-bus destination (`destination.messageBus`
pointing at another project's bus) — chain buses rather than reaching
across projects from satellites.

## Enrollment CEL runs against the CloudEvent, not the payload

`celMatch` sees the event ENVELOPE (`message.type`, `message.source`,
attributes) — not the decoded body. Route on type/source in the
enrollment; inspect payloads in the destination (or with a mediation
template). `"true"` routes everything and is the honest starting point;
tighten as consumers specialize.

## Payload conversion needs schemas on BOTH sides

Converting between Avro and Protobuf requires input AND output formats
with schema definitions; JSON is the schemaless middle. Mismatched
input format vs actual message bytes surfaces as per-message delivery
failures in platform logs (set logSeverity INFO while onboarding), not
at apply time.

## Retry tuning is per-pipeline, and the delays interact

max_attempts 1–100 with 1–600s min/max delays; the API's own doc notes
surfaces where min and max are required to be EQUAL — if an apply
rejects your asymmetric pair, that is the constraint you hit, not a
provider bug. Defaults (5 attempts, 5s–60s) are right for most HTTP
destinations.

## One mediation per pipeline is an API ceiling

The transformation template is a single CEL expression; the API allows
exactly one mediation. Complex reshaping belongs in the consumer, or in
a Workflow destination that transforms and re-publishes — pipelines
mediate, they do not compute.

## Teardown discipline

`DELETE` removes the family; undelivered messages in the bus are lost
with it. `ABANDON` leaves the hub running unmanaged — the escape when
the bus must survive an IaC migration; note the satellites are
abandoned WITH it (one lever, every resource). `PREVENT` is the posture
once production consumers depend on the hub.
