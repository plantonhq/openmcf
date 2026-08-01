---
title: "Default preset"
description: "The standard operator install: the pinned chart (1.6.2 = operator v1.6.2) into its own `ray-system` namespace, watching every namespace on the cluster, leader election on (the chart default — safe..."
type: "preset"
rank: "01"
presetSlug: "01-default"
componentSlug: "kuberay-operator"
componentTitle: "KubeRay Operator"
provider: "kubernetes"
icon: "package"
order: 1
---

# Default preset

The standard operator install: the pinned chart (1.6.2 = operator
v1.6.2) into its own `ray-system` namespace, watching every namespace
on the cluster, leader election on (the chart default — safe for one
replica, required for standbys), control-plane metrics on. No
admission webhook exists — no certificate machinery, no cert-manager
dependency.

What the install owns: the three `ray.io` CRDs (`rayclusters`,
`rayjobs`, `rayservices`) ride the chart's `crds/` directory —
installed once, never upgraded by chart bumps, and KEPT on uninstall
so removing the operator never deletes Ray declarations. One operator
per cluster is the grain.

Declare Ray clusters with `KubernetesRayCluster` against this install
— that kind carries the head/worker topology, autoscaling, GCS fault
tolerance, and the token-auth credential handle.

Keep `metadata.name` at 47 characters or fewer — the longest derived
child name suffixes `-leader-election` and Kubernetes caps names at
63; both engines fail loudly over the budget.

Change first: nothing, usually. Reach for `watchNamespaces` to fence
Ray to specific team namespaces, and `serviceMonitorEnabled` once
Prometheus is on the cluster (it needs the monitoring.coreos.com
CRDs — the install fails without them).

See [01-default.yaml](./01-default.yaml) for the manifest.
