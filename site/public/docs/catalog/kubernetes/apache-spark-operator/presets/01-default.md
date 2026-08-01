---
title: "Default preset"
description: "The standard operator install: the pinned chart (1.8.0 = operator 1.0.0, the official ASF distribution) into its own `spark-operator` namespace, watching every namespace on the cluster, with the..."
type: "preset"
rank: "01"
presetSlug: "01-default"
componentSlug: "apache-spark-operator"
componentTitle: "Apache Spark Operator"
provider: "kubernetes"
icon: "package"
order: 1
---

# Default preset

The standard operator install: the pinned chart (1.8.0 = operator
1.0.0, the official ASF distribution) into its own `spark-operator`
namespace, watching every namespace on the cluster, with the workload
ClusterRole and the `spark` service account created in the release
namespace. No admission webhook exists in this operator — there is no
certificate machinery and no cert-manager dependency.

What the install owns: the two `spark.apache.org` CRDs
(`sparkapplications`, `sparkclusters`) ride the chart's `crds/`
directory — installed once, never upgraded by chart bumps, and KEPT on
uninstall so removing the operator never deletes workload
declarations. One operator per cluster is the grain.

Run Spark jobs by applying `SparkApplication` resources (per pipeline
run — typically submitted by an orchestrator such as
`KubernetesAirflow`, or declared through `KubernetesManifest`),
referencing the `spark` service account this install creates.

Keep `metadata.name` at 40 characters or fewer — the module derives
RBAC names like `<name>-config-monitor-binding` and Kubernetes caps
names at 63; both engines fail loudly over the budget.

Change first: nothing, usually. Reach for `workload.namespaces` when
Spark should be fenced to specific team namespaces, and `replicas: 2`
for a leader-elected warm standby.

See [01-default.yaml](./01-default.yaml) for the manifest.
