# Kubernetes SigNoz — design notes

## Grain

One resource = one SigNoz Helm release (chart `signoz`,
https://charts.signoz.io; the chart version tracks the application in
lockstep). The release is named after `metadata.name`;
`fullnameOverride` pins the server Service and `<name>-otel-collector`
— several SigNoz instances coexist in one cluster and the exported
outputs are deterministic. Keep the resource name at 32 characters or
fewer: the collector Deployment (`<name>-otel-collector`) plus its pod
suffixes must fit Kubernetes' 63-character cap, and both modules fail
loudly over the budget.

## The database seam

The `clickhouse` connection is the load-bearing design choice: the
telemetry store is COMPOSED, never bundled. Every field
default-references a `KubernetesClickHouse` resource's outputs, nothing
database-related installs, and the release carries no operator and no
CRDs — uninstall is ordinary object deletion on both sides.

Verified live, and the reason bundling is rejected rather than merely
deprioritized: a chart that packs the ClickHouse operator and its
installation into ONE release cannot uninstall cleanly — Helm deletes
the operator within seconds while the installation's deletion finalizer
waits for that same operator, and the release deadlocks until its
timeout with the database orphaned. Upstream documents "delete the
namespace" as the workaround. Separate components destroy in the
correct order by construction.

The composed ClickHouse's contract (taught on the spec fields):
SigNoz's tested ClickHouse version (25.12.5 at chart 0.133.0 —
verified live: older servers fail the schema migrations with
UNKNOWN_SETTING), a keeper-backed coordination service (SigNoz
migrates `ON CLUSTER`; explicit `managed_keeper` on single-replica
topologies), a user with
either no grants (unrestricted config-user semantics, verified live) or
a grant set including `GRANT CLUSTER ON *.*`, explicit `networks` on
that user (verified live: the operator fences a networks-less user to
the ClickHouse pods and localhost, and the rejection reads as a
password failure), and the password Secret in SigNoz's OWN namespace
(secretKeyRef cannot cross namespaces — co-locate or replicate).

## Secret discipline

Fully secret-native: the ClickHouse password reaches the server as a
secretKeyRef (existingSecret / existingSecretPasswordKey — never a
rendered value), and the SMTP password rides a valueFrom secretKeyRef
env entry. No credential is generated, rendered, or owned by the
modules.

## The collector's receiver contract

The receiver toggles drive the Service ports AND the collector pipeline
receiver lists from one derivation, so a receiver is never exposed
without being wired or wired without being exposed. Pipeline lists
REPLACE under Helm merge — both modules always render them.

## Cross-engine parity

The Terraform and Pulumi modules render byte-identical chart values from
the same typed spec. Every chart image is the split registry+repository
form deferring to `global.imageRegistry`, so one `image_registry` value
re-points the server, collector and schema-migrator images.

## Deliberate exclusions

The bundled `clickhouse` subchart (permanently disabled — see the
database seam), the `postgresql` and `signoz-otel-gateway` subcharts
(enterprise/licensed surfaces, default off), the `redpanda` subchart (a
banned license family — never modeled, never enabled), the server
replica knob (community = single-writer SQLite), per-image tag
overrides (the chart's tested pairing governs), and the `k8s-infra`
telemetry-shipper chart (the shipping role belongs to
`KubernetesOtelCollector`). Telemetry cold-storage tiering is a
database concern and belongs on the component that owns the database.
All chart surfaces beyond the typed fields remain reachable through
`helm_values` — never the primary interface.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
