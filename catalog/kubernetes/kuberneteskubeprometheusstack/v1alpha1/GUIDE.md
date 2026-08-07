# KubernetesKubePrometheusStack Guide

The judgment this guide carries: this stack is closer to cluster
infrastructure than to an app — one per cluster, and HALF THE CATALOG's
monitoring toggles silently depend on the CRDs it installs. Whether to
run it at all is the assembled-vs-all-in-one decision, which lives in the
[observability-stack pattern](../../../_patterns/observability-stack.md).

## The serviceMonitor seam — the trap that spans the catalog

Components across the catalog carry a `serviceMonitor` (or metrics)
toggle — ingress-nginx, databases, operators. Every one of those toggles
creates a ServiceMonitor, and ServiceMonitor is a CRD THIS stack installs:
enable any component's monitor before the stack exists and THAT
component's release fails to install. When an architecture turns on any
monitoring toggle anywhere, this stack (or its CRDs) must already be on
the cluster — deploy it early, in the shared-cluster chart.

## One per cluster, by CRD physics

The monitoring CRDs are cluster-scoped singletons; a second stack must
skip CRDs and scope its discovery — an advanced posture the reference
page describes, never the default. Treat the stack like the operators:
one install, shared-cluster layer, its own namespace (`createNamespace: true` is
the [namespace-ownership pattern](../../../_patterns/namespace-ownership.md)'s
sole-tenant case). Note the reference page's 26-character name budget —
the chart silently truncates longer fullnames.

## The bundled Grafana boundary

The stack ships a Grafana pre-loaded with its dashboards — enough when
this stack is the only datasource. The moment dashboards must read Loki,
Tempo, or anything else, the standalone
[KubernetesGrafana](../../kubernetesgrafana/v1alpha1/GUIDE.md) is the
composition hub; run it INSTEAD of the bundled one, not beside it.

## On the diagram

The stack renders as a shared-layer node; a standalone Grafana draws a
datasource edge into its `prometheus_endpoint`, and Loki's alert routing
can draw an edge into its Alertmanager. The serviceMonitor dependency,
like every prerequisite-shaped coupling
([operator-prerequisite pattern](../../../_patterns/operator-prerequisite.md)),
draws NOTHING — reviewers check
for the stack node whenever any component's monitoring toggle is on.

## Pairs well with

- KubernetesGrafana — the hub, when datasources go beyond this stack.
- KubernetesLoki — log alerts routed through this Alertmanager
  (`ruler.alertmanagerUrl`).
- Every kind with a serviceMonitor toggle — they all assume this stack.
