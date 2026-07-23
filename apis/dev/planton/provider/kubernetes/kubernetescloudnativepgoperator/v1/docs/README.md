# KubernetesCloudNativePgOperator: Research and Design

## Introduction

CloudNativePG is the CNCF PostgreSQL operator: it reconciles `Cluster`
custom resources into highly available PostgreSQL clusters — streaming
replication, automated failover with a safe primary election, rolling
updates, declarative roles and storage, and plugin-based backups. This
component installs it from the official Helm chart (`cloudnative-pg` at
`https://cloudnative-pg.github.io/charts`; pinned default chart 0.29.0,
which ships operator 1.30.0 — chart and app versions move separately,
and the chart pin governs; both pins verified against the chart-repo
index). The optional Barman Cloud plugin installs from the same
repository (`plugin-barman-cloud`, pinned default chart 0.7.0, which
ships plugin v0.13.0 — the plugin chart versions independently of the
operator chart, and each release carries its own pin).

## Upstream Architecture

An installation is:

1. **The operator Deployment** — a controller reconciling the
   `postgresql.cnpg.io` custom resources (Cluster, ScheduledBackup,
   Backup, Pooler, Database, ...). Unlike most database operators,
   CloudNativePG manages instances DIRECTLY — no StatefulSet in
   between: each PostgreSQL instance is a Pod (plus its PVCs) the
   operator creates, watches, and replaces itself, which is what makes
   its failover and rolling-update choreography possible.
2. **The instance manager** — PID 1 of every instance pod, injected by
   the operator. It supervises the `postgres` process, reports
   readiness, applies configuration, and executes the operator's
   promotion/demotion decisions. The operator holds the control plane;
   the instance manager is its agent inside every pod.
3. **Cluster-scoped webhooks and CRDs** — a mutating/validating webhook
   pair defaulting and rejecting Cluster changes at apply time. The
   webhook's service name (`cnpg-webhook-service`) is fixed by the
   chart and baked into the webhook certificate.

One installation per cluster: the CRDs and webhooks are cluster-scoped
and their names are fixed, so a second installation would fight over
both. The release name is therefore fixed to `cnpg` and never derives
from `metadata.name`.

## Engine vs Database Declarations

This component installs and configures the ENGINE. The databases are
declared with KubernetesPostgres resources — one per PostgreSQL cluster
— which the operator reconciles. The split matters for lifecycle and
ownership: the operator is platform infrastructure upgraded on the
platform team's schedule; databases are application infrastructure
declared next to the applications they serve. Day-2 actions against a
running database (promote, hibernate, one-off backup) belong to the
`cnpg` kubectl plugin and the CRDs, not to either spec.

## Plugin-Based Backups: Why the Barman Plugin Rides Along

CloudNativePG's in-tree `barmanObjectStore` backup support is deprecated
upstream and scheduled for removal; object-store backups now travel
through CNPG-I (the operator's plugin interface) with the Barman Cloud
plugin as the reference implementation. The spec models the plugin as an
arm of THIS component because the plugin is cluster infrastructure
exactly like the operator: one plugin installation serves every
database's backup blocks.

Two upstream constraints shape the module:

- **A separate Helm release, same namespace** — upstream forbids folding
  the plugin into the operator's release (Helm ownership of shared
  resources would conflict). The module installs `plugin-barman-cloud`
  as its own release AFTER the operator, so the plugin's CNPG-I
  registration always lands on a running operator; uninstall unwinds in
  reverse for free. The plugin release name is fixed for the same
  singleton reason as the operator's: its gRPC service name
  (`barman-cloud`) is baked into its TLS certificate.
- **cert-manager is REQUIRED by the plugin arm** — the plugin chart
  renders cert-manager Issuer/Certificate resources unconditionally; its
  operator↔sidecar TLS is cert-manager-issued. Without cert-manager on
  the cluster (KubernetesCertManager) the plugin release fails to
  install — atomic rolls it back with a clear readiness timeout.

## CRD Lifecycle: Databases Survive Uninstall

The chart stamps `helm.sh/resource-policy: keep` on every CRD
UNCONDITIONALLY, so uninstalling the release never cascade-deletes the
Cluster resources — and the databases behind them. This is the upstream
safety posture, kept as-is: `crds.install` (chart default true) exists
only for clusters where something else manages the CRDs, and there is no
destructive cleanup switch to misconfigure.

## Watch Scope and Operator Configuration

`watch.cluster_wide` (chart default true) is the normal posture: RBAC is
created as ClusterRoles and the operator reconciles every namespace.
Setting it false with `watch.namespaces` fences the operator into
specific namespaces with namespace-scoped RBAC — CEL rules enforce that
namespaces are present exactly when the fence is on.

`operator_config` passes through the chart's `config.data` map (the
operator's configuration ConfigMap: `INHERITED_ANNOTATIONS`,
`INHERITED_LABELS`, `PULL_SECRET_NAME`, ... — see the
operator-configuration page of the upstream documentation for the full
vocabulary). One precedence rule keeps the two from fighting: the typed
watch field OWNS the `WATCH_NAMESPACE` key — a user entry under that key
is always stripped, and the key renders only from `watch.namespaces`
(comma-joined) when the fence is on. Both engines implement the same
stripping.

`max_concurrent_reconciles` (chart default 10) is the throughput knob —
`replicas` is not: extra operator replicas are leader-elected warm
standbys that shorten failover of the operator itself.

## Typed Surface vs Escape Hatch

The typed spec covers namespace and lifecycle, chart version, CRD
lifecycle, operator sizing and scheduling, watch scope, the operator
configuration map, reconcile concurrency, the Barman Cloud plugin arm,
telemetry, and image overrides.

Deliberately unmodeled as typed fields (all reachable via
`helm_values`, which scopes to the OPERATOR chart only):

- **Webhook tuning** (probes, port, TLS knobs) — the chart's defaults
  are the tested posture; tuning them is an expert move
- **Host network mode** — a niche posture for clusters whose CNI blocks
  webhook traffic from the control plane
- **Topology spread constraints, update strategy, security contexts,
  RBAC/service-account tuning** — the chart's operational long tail

The plugin's typed surface is deliberately minimal — container resources
only; everything else rides the plugin chart's defaults. `helm_values`
does NOT flow to the plugin release: the two charts share value keys
(`resources`, `image`, ...), so forwarding one document to both would
misconfigure the plugin.

## Install Semantics

Both engines install real Helm releases, atomically, with cleanup on
fail and a 600s timeout, waiting for the operator to become Available. A
release that can never become ready — a PodMonitor rendered without the
Prometheus operator CRDs (THE classic install failure), the plugin
without cert-manager — fails the deploy with a rollback instead of
surfacing later as Cluster resources that never reconcile. Typed values
render first; `helm_values` merges last with Helm `-f` semantics
(identical on both engines). The module (not Helm) owns namespace
creation via `create_namespace`.

## Outputs

`namespace`, `release_name` (fixed `cnpg`), and
`barman_plugin_release_name` (the fixed `plugin-barman-cloud` when the
plugin arm is enabled; empty otherwise — the handle KubernetesPostgres
backup blocks key off to know whether object-store backups can work).

## E2E

The behavioral facts are properties of the platform, not of any one test
run:

- The readiness proof is the release becoming Available under the
  atomic/wait posture — a webhook that never serves or a plugin without
  cert-manager fails the install, not the first database.
- The plugin scenario chains a cert-manager prerequisite; ordinary
  operator installs need none.
- The PodMonitor arm fails the release on clusters without the
  Prometheus operator CRDs, by design.
- Uninstall keeps the CRDs (and every Cluster resource) by the chart's
  unconditional keep policy.
