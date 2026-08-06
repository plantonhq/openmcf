---
title: "Strimzi Kafka Operator"
description: "Strimzi Kafka Operator deployment documentation"
icon: "package"
order: 100
componentName: "kubernetesstrimzikafkaoperator"
---

# Kubernetes Strimzi Kafka Operator

Installs the Strimzi cluster operator — the CNCF project that runs
Apache Kafka on Kubernetes — from the official `strimzi-kafka-operator`
Helm chart, with a typed spec over the chart's meaningful configuration
surface. The operator reconciles `Kafka` and `KafkaNodePool` custom
resources into KRaft-mode Kafka clusters (Kafka's built-in Raft
metadata quorum — ZooKeeper was removed from Strimzi in 0.46), and its
per-cluster entity operators reconcile `KafkaTopic` / `KafkaUser`
resources into real topics and authenticated users. This component
installs the ENGINE; the Kafka clusters themselves are declared with
KubernetesKafka resources — one per cluster.

## What Gets Created

- **Namespace** (optional) — the installation namespace, created and
  owned when `create_namespace` is set
- **Helm Release** (named `metadata.name`) — the operator Deployment,
  its RBAC, and the Strimzi CRDs (Kafka, KafkaNodePool, KafkaTopic,
  KafkaUser, and the rest of the family). The CRDs ship in the chart's
  Helm-native `crds/` directory: installed on first install, never
  upgraded or deleted by Helm — uninstalling the release NEVER
  cascade-deletes the Kafka clusters. The same posture means a
  `chart_version` upgrade runs new operator code against the EXISTING
  CRDs — apply the new release's CRDs yourself when an upgrade's
  release notes call for it.

## Watch Scope

By default the operator watches ITS OWN namespace only (the chart
default — Kafka clusters live beside their operator). Two knobs widen
it, mutually exclusive:

- **`watch.any_namespace: true`** — one operator manages Kafka
  clusters in every namespace (cluster-wide RBAC)
- **`watch.namespaces: [...]`** — an explicitly fenced set of
  namespaces (the installation namespace is always watched in
  addition)

A KubernetesKafka cluster declared in a namespace the operator does
not watch is silently never reconciled — install the operator in the
cluster's namespace, or widen the watch.

## Prerequisites

- A cluster not already running a Strimzi operator with overlapping
  watch scope — the CRDs and the chart's ClusterRoles are
  cluster-scoped singletons (a second fenced install needs
  `create_global_resources: false`)
- Nothing else: Kafka clusters, topics, and users are concerns of the
  KubernetesKafka / KubernetesKafkaTopic / KubernetesKafkaUser
  resources composed on top

## Quick Start

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesStrimziKafkaOperator
metadata:
  name: strimzi-operator
spec:
  namespace:
    value: kafka
  create_namespace: true
```

The operator becomes Available and starts reconciling; KubernetesKafka
resources declared in the `kafka` namespace turn into running KRaft
clusters.

## Configuration Surface

- **`chart_version`** — the chart pin (default `1.1.0`; chart and
  operator versions move together). Versions must exist as SERVED
  charts in the repository index (https://strimzi.io/charts/) — the
  Chart.yaml inside the Strimzi source tree is a build-time
  placeholder. Remember the CRD posture: upgrades never upgrade CRDs
- **`replicas`** — operator pods (default 1); extras are
  leader-elected warm standbys for the OPERATOR itself, not
  reconciliation throughput
- **`watch`** — the scope block above
- **`full_reconciliation_interval_ms`** — the periodic full-repair
  pass over every watched resource (default 120000 — 2 minutes)
- **`operation_timeout_ms`** — timeout for internal operations like
  rolling a pod (default 300000 — 5 minutes); raise on slow storage
  where broker restarts legitimately exceed it
- **`log_level`** — ERROR / WARN / INFO (default) / DEBUG / TRACE
- **`feature_gates`** — Strimzi gates in the operator's own syntax
  (`+Gate,-OtherGate`); an unknown gate fails operator startup
- **`kubernetes_service_dns_domain`** — set only on clusters with a
  non-default DNS domain; a mismatch produces TLS certificates whose
  SANs don't match the advertised service DNS names
- **`leader_election_enabled`** — meaningful only with `replicas`
  above 1
- **`generate_network_policy` / `generate_pod_disruption_budget`** —
  the operand policy toggles (both default true)
- **`create_global_resources`** — set false only for a SECOND operator
  install; the fixed-name ClusterRoles belong to the first release
- **`resources`, `node_selector`, `tolerations`,
  `image_pull_secrets`, `image`** — operator pod placement, sizing,
  and the image source override — `image` steers EVERY Strimzi image
  (operator and all operands: Kafka, entity operators, Cruise Control,
  exporter), the air-gap path
- **`helm_values`** — escape hatch: extra chart values merged LAST
  (Helm `-f` semantics) for the surface beyond the typed fields —
  never the primary interface, never for secrets

## Stack Outputs

| Output | Description |
|---|---|
| `namespace` | Namespace the operator runs in |
| `release_name` | Helm release name of the operator (equals `metadata.name`) |

## Next Steps

Declare Kafka clusters as KubernetesKafka resources — one per cluster,
in a namespace the operator watches; the operator reconciles them into
brokers, controllers, listeners, and certificates, and its per-cluster
entity operators serve KubernetesKafkaTopic / KubernetesKafkaUser
declarations. Pin `chart_version` deliberately and upgrade the operator
on the platform's schedule — removing this component never deletes the
CRDs or the Kafka clusters behind them.
