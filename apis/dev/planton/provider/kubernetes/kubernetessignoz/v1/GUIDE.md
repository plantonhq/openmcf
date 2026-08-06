# KubernetesSignoz Guide

The judgment this guide carries: Signoz is the all-in-one answer to
observability — but "all-in-one" does NOT include its database. The
deepest composition chain in the catalog sits underneath it: the Altinity
operator, then a ClickHouse, then Signoz — and skipping a link leaves a
UI with nothing behind it.

## The choice this kind anchors

Signoz vs the assembled stack (kube-prometheus-stack + Loki + Tempo +
Grafana) is the comparison the
[observability-stack pattern](../../../patterns/observability-stack.md)
owns — one tool for all three signals versus independently scalable,
swappable pieces. Propose ONE shape per signal set; running both without
a stated split is double-tooling.

## The three-link chain underneath

Nothing ClickHouse-related is bundled — the telemetry store is composed
(the spec doc is explicit, and explains why: your telemetry outlives
Signoz reinstalls, and upgrades roll independently). A complete Signoz
proposal is three kinds deep. The first link is a prerequisite that
draws no diagram edge (the
[operator-prerequisite pattern](../../../patterns/operator-prerequisite.md));
the Signoz-to-ClickHouse link, by contrast, IS a drawn `valueFrom` edge:

1. [KubernetesAltinityOperator](../../kubernetesaltinityoperator/v1/GUIDE.md)
   — watching the ClickHouse's namespace (its guide's own trap).
2. [KubernetesClickHouse](../../kubernetesclickhouse/v1/GUIDE.md) — with
   coordination if replicated.
3. This kind — its `clickhouse` block default-references that kind's
   outputs, so the wiring is one `valueFrom` per field.

## Ingestion: its collector and yours are collaborators

The bundled SigNoz OpenTelemetry Collector is the app-facing OTLP
gateway — applications send there. For cluster/node telemetry (logs,
kubelet metrics, events), deploy a
[KubernetesOtelCollector](../../kubernetesotelcollector/v1/GUIDE.md) in
daemonset mode pointed at that same gateway (the spec doc prescribes
exactly this). The two are complementary layers, not alternatives.

## On the diagram

Signoz, its ClickHouse, and the operator render as three nodes with the
`clickhouse` reference edges drawn — a reviewer can see the database the
telemetry actually lives in, which is the point of composing it. All
services stay ClusterIP; exposing the UI composes from route kinds.

## Pairs well with

- KubernetesClickHouse + KubernetesAltinityOperator — the required chain.
- KubernetesOtelCollector (daemonset) — cluster telemetry into the same
  gateway.
- KubernetesIngress / route kinds — UI exposure, composed as always.
