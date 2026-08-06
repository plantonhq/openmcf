# KubernetesLoki Guide

The judgment this guide carries: Loki STORES logs — it ships none. An
architecture with Loki but no shipper collects nothing, silently; and the
mode choice is really a storage commitment.

## Nothing ships logs by itself — compose the shipper

Deploy a [KubernetesOtelCollector](../../kubernetesotelcollector/v1alpha1/GUIDE.md)
in daemonset mode with the cluster-logs pipeline (its presets carry it),
pointed at this component's `gateway_endpoint` output — that pair is the
cluster's log pipeline. Grafana reads the logs back through a `loki`
datasource on the same endpoint. The full wired composition:
[observability-stack pattern](../../../patterns/observability-stack.md).
Loki alone looks healthy and receives nothing — the gap appears only when
someone searches for logs that were never shipped.

## The mode choice is a storage commitment

`monolithic` (default) runs everything in one StatefulSet — right for
dev and small production volumes. `simpleScalable` splits write/read/
backend tiers AND REQUIRES object storage — on-cluster, that means
composing a KubernetesSeaweedFs (the same S3 move the Flink and Tempo
guides make). Choosing the scalable mode without the storage is a
deployment that cannot come up; the microservices mode is deliberately
not modeled (the reference page says why).

## Alerts route through the one Alertmanager

`ruler.alertmanagerUrl` is a typed reference to the
kube-prometheus-stack's Alertmanager endpoint — wire it so log-driven
alerts join the same routing, silencing, and paging as everything else,
instead of growing a second alert system.

## Namespace ownership

Observability components conventionally share one namespace (the
pattern's example uses `observability`) — wire `spec.namespace` through
that KubernetesNamespace
([namespace-ownership pattern](../../../patterns/namespace-ownership.md)).

## Pairs well with

- KubernetesOtelCollector (daemonset) — the shipper; without it, nothing
  arrives.
- KubernetesGrafana — the reader (`loki` datasource on the gateway).
- KubernetesKubePrometheusStack — Alertmanager for log-driven alerts.
- KubernetesSeaweedFs — object storage when the mode demands it.
