# Kafka on Kubernetes: Why the Operator Is the Product

## Introduction

Apache Kafka is the archetypal hard-to-run stateful system: a
distributed commit log whose brokers carry identity, whose partitions
carry placement, and whose safety guarantees depend on a metadata
quorum staying healthy while individual machines fail. Running it on
Kubernetes without database-aware software is the same trap as running
databases on raw StatefulSets — Day-1 deployment is easy to script,
and Day-2 (rolling a broker without violating in-sync replicas,
rotating certificates, growing storage, upgrading Kafka versions in
the mandated order) is where naive approaches quietly lose data.

**Strimzi** is the CNCF project that encodes that Day-2 expertise into
software. It extends Kubernetes with custom resources — `Kafka`,
`KafkaNodePool`, `KafkaTopic`, `KafkaUser`, and the rest of the family
— and runs a cluster operator that continuously reconciles them into
brokers, controllers, listeners, certificates, and per-cluster entity
operators. This document records why Planton's Kafka catalog is built
on Strimzi, what the operator component actually installs, and the
design decisions behind the typed spec.

## The Deployment Landscape

### Raw StatefulSets: the anti-pattern

A Kafka broker is not a fungible pod. It has a node ID baked into the
cluster metadata, a set of partitions only it holds, and a position in
the controller quorum. A bare StatefulSet can restart a crashed broker
pod, but it cannot know that rolling two brokers at once drops a
partition below its minimum in-sync replica count, cannot sequence a
Kafka version upgrade, and cannot re-issue the TLS certificates
brokers authenticate each other with. Everything that makes Kafka safe
is invisible to it.

### Helm charts without an operator: Day-1 only

A chart that templates StatefulSets automates the anti-pattern.
Helm's involvement ends when `helm install` returns; nothing watches
the cluster afterwards. The modern pattern — and the one this
component implements — is to use Helm for what it is good at
(installing the OPERATOR, a stateless controller) and let the operator
own the lifecycle of the stateful thing.

### The operator: continuous reconciliation

The Strimzi cluster operator runs the standard reconcile loop — read
the declared `Kafka`/`KafkaNodePool` resources, observe the actual
StrimziPodSets, roll exactly the pods that need rolling, in the order
that keeps quorums intact — plus a periodic full pass (the chart's
`fullReconciliationIntervalMs`, default two minutes) that repairs
drift event-driven reconciliation missed. Certificates are its job
too: it generates and renews a cluster CA and a clients CA, and rolls
pods when renewals require it (optionally only inside declared
maintenance windows).

## Why Strimzi

- **CNCF-governed, Apache-2.0 licensed** — the operator, the CRDs, and
  the container images (quay.io/strimzi) are open source with no
  open-core split in the operator itself.
- **KRaft-native** — from Strimzi 0.46 onward, support for
  ZooKeeper-based clusters (and ZooKeeper-to-KRaft migration) is
  REMOVED. The 1.x line is KRaft-only: the metadata quorum runs inside
  Kafka itself, on nodes carrying the `controller` role. There is no
  ZooKeeper anywhere in this architecture — one less distributed
  system to operate.
- **The `v1` API generation** — Strimzi 1.0.0 moved the CRDs to the
  served-and-stored `v1` API and removed the old `v1beta2`/`v1alpha1`
  versions. This catalog is built against `kafka.strimzi.io/v1`
  exclusively.
- **Node pools as the only topology primitive** — brokers and
  controllers are declared as `KafkaNodePool` resources bound to their
  cluster by the `strimzi.io/cluster` label, each pool independently
  scalable with its own storage shape and sizing.
- **A complete satellite ecosystem** — per-cluster entity operators
  (topics and users as resources), Cruise Control integration for
  partition rebalancing, Kafka Exporter for consumer-lag metrics — all
  reconciled by this one operator.

## What This Component Installs

One Helm release of the official `strimzi-kafka-operator` chart from
`https://strimzi.io/charts/`, named after `metadata.name`: the
operator Deployment, its ServiceAccount, watch-scoped RBAC, and the
Strimzi CRDs.

### The chart version is the served version

`chart_version` (default `1.1.0`) must exist as a SERVED chart in the
repository index at https://strimzi.io/charts/. The `Chart.yaml`
inside the Strimzi source tree carries a build-time placeholder
(`version: 0.1.1`) and never reflects the served version — the index
governs. Chart and operator versions move together for this chart:
chart 1.1.0 ships operator 1.1.0.

### The `crds/` directory posture (read this before uninstalling)

The chart ships its CRDs in the Helm-native `crds/` directory, which
gives them a deliberate lifecycle Helm applies to nothing else:

- **Installed on first install** of the release.
- **Never deleted on uninstall.** Removing the operator release leaves
  every CRD — and therefore every Kafka cluster, topic and user —
  untouched. Clusters keep serving traffic; they simply stop being
  reconciled until an operator returns. This is the upstream safety
  posture, and it is why uninstalling this component can never
  cascade-delete a Kafka cluster.
- **Never upgraded on chart upgrade.** A `chart_version` bump re-runs
  the release with new operator code against the EXISTING CRDs.
  Strimzi minor releases regularly add CRD fields — when an upgrade's
  release notes call for it, apply the new release's CRDs yourself.
  Skipping this leaves new spec surface silently unusable until the
  CRDs catch up.

The operator registers no admission webhooks and creates no
cluster-scoped singletons at runtime (verified in the operator
source), so uninstall leaves nothing stranded beyond the
deliberately-kept CRDs.

### Watch scope

The chart models scope with two independent values, and the spec
mirrors them as a `watch` block with mutually exclusive arms
(validated by CEL at the spec):

| Spec | Chart value | Meaning |
|---|---|---|
| (no `watch` block) | — | Own namespace only — the chart default; clusters live beside their operator |
| `watch.any_namespace: true` | `watchAnyNamespace: true` | Every namespace, cluster-wide RBAC |
| `watch.namespaces: [...]` | `watchNamespaces: [...]` | An explicit fence; the installation namespace is always watched in addition |

### Multiple installs and `create_global_resources`

The chart's ClusterRoles carry fixed Strimzi names. The first release
owns them; a second install in the same cluster (for a disjoint
namespace fence) must set `create_global_resources: false` or Helm
fails the install on the ownership conflict. One cluster-wide operator
is the simpler topology whenever tenancy allows it.

### The image override is the air-gap path

`image.registry` / `image.repository` / `image.tag` render the chart's
`defaultImageRegistry` / `defaultImageRepository` / `defaultImageTag`
— which steer EVERY Strimzi image: the operator itself and every
operand image it deploys (Kafka, entity operators, Cruise Control,
Kafka Exporter). Empty means the chart defaults (quay.io/strimzi at
the chart version). `image_pull_secrets` rides the chart's
`image.imagePullSecrets` list.

## Design Decisions

- **Typed fields render only on divergence.** The chart ships real
  defaults (watch scope off, reconciliation at 120000 ms, operation
  timeout 300000 ms, resources requests 200m/384Mi and limits
  1000m/384Mi, `createGlobalResources: true`); the modules render a
  value only when the spec diverges, so an empty spec installs the
  chart exactly as upstream ships it — on both IaC engines.
- **Atomic, waited installs.** The release installs with
  wait + atomic + cleanup-on-fail and a 600 s timeout. An operator
  that never becomes ready (an unpullable image from a private mirror
  is the classic case) fails THIS deploy with a readiness timeout
  instead of surfacing later as Kafka resources that mysteriously
  never reconcile.
- **The module owns namespace creation.** `create_namespace` drives a
  labeled namespace resource; Helm's own create-namespace is always
  false so the governance labels are never lost.
- **`replicas` is honest about what it buys.** Extra operator replicas
  are leader-elected warm standbys — exactly one replica reconciles at
  a time. The knob exists for failover, not throughput.
- **`helm_values` merges LAST.** The escape hatch is a YAML document
  merged over the typed rendering with Helm `-f` semantics (maps
  deep-merge, later document wins, lists replace) — identical on both
  engines. It is a safety valve, never the primary interface.

## Deliberately Unmodeled

Kept off the typed spec on purpose; every item remains reachable
through `helm_values`:

- **Per-component image overrides** — the chart can override each
  operand image individually; the typed `image` block models only the
  global registry/repository/tag switch (the actual air-gap need).
- **Extra operator env vars and container security contexts** —
  operational edges with no cross-environment story.
- **The operator's own PodDisruptionBudget / NetworkPolicy values** —
  distinct from the operand toggles the spec does model.
- **Grafana dashboards, aggregate ClusterRoles, Connect build
  timeout** — chart conveniences outside this component's job of
  running the reconciler.
- **Kafka-cluster concerns** — everything about a specific cluster
  (versions, listeners, storage, authorization) belongs to
  KubernetesKafka, never here.

## Planton's Approach

Separation of concerns, one layer per lifecycle:

1. **KubernetesStrimziKafkaOperator** (this component) manages the
   operator: chart pin, watch scope, sizing, image source.
2. **KubernetesKafka** declares each Kafka cluster — node pools,
   listeners, broker config, authorization — rendering
   `Kafka` + `KafkaNodePool` resources the operator reconciles.
3. **KubernetesKafkaTopic / KubernetesKafkaUser** declare topics and
   users, reconciled by each cluster's entity operators, placed in the
   cluster's namespace with the `strimzi.io/cluster` binding label.

Operator upgrades never touch cluster declarations; one operator
serves many clusters; and both layers keep an escape hatch
(`helm_values` here, KubernetesManifest for raw CR surface).

## Further Reading

- **[Strimzi documentation](https://strimzi.io/documentation/)** — the
  operator's own deployment and configuration guides
- **[Strimzi Helm chart repository index](https://strimzi.io/charts/)**
  — the served chart versions `chart_version` must come from
- **[Strimzi GitHub repository](https://github.com/strimzi/strimzi-kafka-operator)**
  — source, CHANGELOG (read before every `chart_version` bump — CRD
  updates ride release notes), and the chart source
