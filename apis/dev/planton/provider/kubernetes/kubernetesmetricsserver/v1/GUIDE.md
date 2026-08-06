# KubernetesMetricsServer Guide

The judgment this guide carries: this is the smallest component with the
largest silent blast radius — without it, every HPA on the cluster
deploys cleanly and simply never scales, and `kubectl top` goes dark.
When an architecture enables ANY resource-based autoscaling, this node
belongs in the shared-cluster chart.

## The silent-HPA trap

HorizontalPodAutoscalers (the standalone kind AND the inline autoscaling
blocks on workload kinds and controllers like ingress-nginx) read
CPU/memory utilization from the `metrics.k8s.io` APIService that ONLY
metrics-server provides. Absent, nothing errors: HPAs deploy, report no
metrics, and never scale (the spec doc states it verbatim). The failure
is invisible until the traffic spike that autoscaling was supposed to
absorb. Any proposal that turns on utilization-based scaling includes
this component.

## One per cluster

It registers the cluster-wide `metrics.k8s.io` APIService — a singleton;
one install in the shared-cluster chart serves everything. Do not confuse
it with the kube-prometheus-stack: metrics-server feeds AUTOSCALING and
`kubectl top` (live, unstored); the stack feeds MONITORING (scraped,
stored, alerted). Clusters routinely need both — the
[observability-stack pattern](../../../patterns/observability-stack.md)
covers the monitoring side.

## Namespace ownership — the infra exception

A dedicated single-tenant namespace with `createNamespace: true` — the
[namespace-ownership pattern](../../../patterns/namespace-ownership.md)'s
sole-tenant case.

## Pairs well with

- KubernetesHorizontalPodAutoscaler and every inline autoscaling block —
  the consumers that silently need this.
- KubernetesKubePrometheusStack — the monitoring complement, not a
  substitute in either direction.
