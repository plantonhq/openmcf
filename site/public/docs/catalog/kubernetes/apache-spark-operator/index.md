---
title: "Apache Spark Operator"
description: "Apache Spark Operator deployment documentation"
icon: "package"
order: 100
componentName: "kubernetessparkoperator"
---

# Apache Spark Operator

The official ASF controller that runs Spark workloads declared as
`SparkApplication` (one job, run to completion) and `SparkCluster`
(a long-lived standalone cluster) custom resources, from the official
`spark-kubernetes-operator` Helm chart (1.8.0 = operator 1.0.0). One
operator per cluster: it watches cluster-wide by default, and Spark
jobs are submitted per pipeline run — typically by an orchestrator —
not declared as catalog resources.

## Highlights

- **Removing the operator never deletes the workloads' declarations**
  — the two spark.apache.org CRDs ride the chart's `crds/` directory:
  installed once, never upgraded by chart bumps, kept on uninstall by
  upstream design.
- **No webhook, no cert-manager** — the operator validates in its
  reconcile loop; there is no certificate machinery and nothing that
  can fail-close the cluster's write path.
- **The fence is the watch scope** — list workload namespaces and the
  chart creates them, plants the `spark` service account and a
  namespace-scoped Role in each, and the operator watches exactly that
  list; empty means cluster-wide with a workload ClusterRole.
- **Instances coexist** — the chart hardcodes its cluster-scoped RBAC
  names, so a second install would collide by construction; the
  modules derive every RBAC name from the release identity instead.
- **Fail-loud, not fail-later** — names over the 40-character budget
  are rejected at apply time, and the atomic install waits for the
  operator to become Available instead of surfacing as jobs that never
  reconcile.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
