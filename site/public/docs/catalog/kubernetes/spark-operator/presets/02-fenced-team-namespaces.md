---
title: "Fenced team-namespaces preset"
description: "The multi-tenant posture: Spark runs ONLY in the listed namespaces. The chart CREATES each namespace (know this before pointing at names you manage elsewhere), plants the `spark` service account and..."
type: "preset"
rank: "02"
presetSlug: "02-fenced-team-namespaces"
componentSlug: "spark-operator"
componentTitle: "Spark Operator"
provider: "kubernetes"
icon: "package"
order: 2
---

# Fenced team-namespaces preset

The multi-tenant posture: Spark runs ONLY in the listed namespaces.
The chart CREATES each namespace (know this before pointing at names
you manage elsewhere), plants the `spark` service account and a
namespace-scoped Role in each — replacing the cluster-wide workload
ClusterRole — and the operator watches exactly this list.
SparkApplications anywhere else are ignored without an error, so a
missing namespace in this list looks like a job that never starts.

Know the keep truth: the chart marks workload namespaces and their
RBAC `helm.sh/resource-policy: keep` — they SURVIVE uninstall by
upstream design, so running jobs never lose their identity mid-flight.

`replicas: 2` runs a warm standby; the module renders the operator's
leader-election property automatically (the chart refuses
multi-replica installs without it). The reconciler-interval property
is the operator's own key space — the full catalog ships with the
operator's docs at the pinned version.

Change first: the namespace list — it is the fence, the watch scope,
and the RBAC placement in one value.

See [02-fenced-team-namespaces.yaml](./02-fenced-team-namespaces.yaml)
for the manifest.
