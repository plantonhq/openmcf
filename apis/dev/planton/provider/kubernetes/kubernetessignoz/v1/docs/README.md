# Kubernetes SigNoz — design notes

## Grain

One resource = one SigNoz Helm release (chart `signoz`,
https://charts.signoz.io; the chart version tracks the application in
lockstep). The release is named after `metadata.name`;
`fullnameOverride` pins the server Service and `<name>-otel-collector`,
and the bundled installation's fullname is pinned to
`<name>-clickhouse` — several SigNoz instances coexist in one cluster
and the exported outputs are deterministic. Keep the resource name at
30 characters or fewer: the bundled ClickHouse wraps it in ~27
characters of operator scaffolding (`chi-<name>-clickhouse-cluster-0-0`)
inside Kubernetes' 63-character cap, and both modules fail loudly over
the budget.

## The database seam

The `database` oneof is the load-bearing design choice. Empty or
`managed_clickhouse` = the appliance: the chart's own ClickHouse stack
(a namespace-fenced Altinity operator, the installation, ZooKeeper)
with capacity/topology knobs only — deep ClickHouse control is
deliberately NOT re-modeled here. `external_clickhouse` = composition:
every field default-references a `KubernetesClickHouse` resource's
outputs, nothing database-related installs, and the referenced user
needs `access_management` plus cluster-wide DDL grants (SigNoz runs
`ON CLUSTER` DDL and owns its schema migrations).

## Secret discipline

The bundled arm's admin password is module-generated (random, 24 chars),
reaches the chart outside the values documents (set_sensitive /
a Pulumi secret Output), and is exported through the module-owned
`<name>-clickhouse-auth` Secret. KNOW THE UPSTREAM GRAIN: the chart
renders that password into the ClickHouseInstallation object and literal
container env — namespace readers can see it; what the module guarantees
is a unique per-install credential and that the chart's
publicly-documented default never ships. The external arm and SMTP are
fully secret-native (secretKeyRef wiring). Cold-storage keys, when
declared instead of the IRSA role arm, render into ClickHouse's storage
configuration — the spec field says so and recommends the keyless arm.

## The collector's receiver contract

The receiver toggles drive the Service ports AND the collector pipeline
receiver lists from one derivation, so a receiver is never exposed
without being wired or wired without being exposed. Pipeline lists
REPLACE under Helm merge — both modules always render them.

## Cross-engine parity

The Terraform and Pulumi modules render byte-identical chart values from
the same typed spec. Every chart image is the split registry+repository
form deferring to `global.imageRegistry`, so one `image_registry` value
re-points the server, collector, ClickHouse, operator, metrics-exporter
and ZooKeeper images; the bundled arm's UDF init-container image follows
its own chart key (`helm_values` territory).

## Deliberate exclusions

The `postgresql` and `signoz-otel-gateway` subcharts (enterprise/licensed
surfaces, default off), the `redpanda` subchart (a banned license family
— never modeled, never enabled), the server replica knob (community =
single-writer SQLite), per-image tag overrides (the chart's tested
pairing governs), and the `k8s-infra` telemetry-shipper chart (the
shipping role belongs to `KubernetesOtelCollector`). All chart surfaces
beyond the typed fields remain reachable through `helm_values` — never
the primary interface.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
