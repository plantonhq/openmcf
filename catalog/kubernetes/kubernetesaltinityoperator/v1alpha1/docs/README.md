# KubernetesAltinityOperator: Research and Design

## Introduction

KubernetesAltinityOperator installs the Altinity ClickHouse operator
from the official `altinity-clickhouse-operator` Helm chart
(https://docs.altinity.com/clickhouse-operator/, pinned 0.27.2) as a
single Helm release named after `metadata.name`. The operator is the
ENGINE of the ClickHouse story in this catalog: KubernetesClickHouse
declares `ClickHouseInstallation` (CHI) custom resources — and, for
its managed-Keeper arm, `ClickHouseKeeperInstallation` (CHK) — and
this operator reconciles them into per-host StatefulSets, generated
server configuration, Services, and safe rolling restarts.

## The Deployment Landscape

ClickHouse without an operator is the classic stateful anti-pattern:
shard and replica topology, ZooKeeper/Keeper coordination, schema
propagation across replicas, and restart ordering during configuration
changes are Day-2 concerns no plain StatefulSet or Helm chart encodes.
The Altinity operator carries that expertise as a reconciler, which is
why the catalog splits the concern in two: this kind installs the
engine once, KubernetesClickHouse declares each cluster.

### The chart/image lockstep truth

The served repository index
(https://docs.altinity.com/clickhouse-operator/index.yaml) is the
source of truth for what actually installs:

- **Chart versions track operator releases one-to-one** — chart
  0.27.2 runs operator image 0.27.2 (`altinity/clickhouse-operator`)
  and the same-versioned metrics-exporter sidecar
  (`altinity/metrics-exporter`). A `chart_version` bump is an operator
  upgrade; when overriding the operator image for a mirror, keep the
  tag matched to the chart.
- **All default images at the pinned chart verified pullable** on
  Docker Hub — with one aging exception, the CRD hook's kubectl image
  (below).

## The CRD Lifecycle: Chart-Owned, Keep-on-Uninstall

The chart ships its four CRDs (`clickhouseinstallations`,
`clickhouseinstallationtemplates`, `clickhouseoperatorconfigurations`
in the `clickhouse.altinity.com` group, and
`clickhousekeeperinstallations` in `clickhouse-keeper.altinity.com`)
in its `crds/` directory. That single upstream choice makes this
component's CRD story the opposite of sibling operator kinds whose
charts template CRDs as release-owned resources and force the modules
to own the lifecycle:

- **Helm's native `crds/` handling installs the CRDs on FIRST install
  and never deletes them on uninstall.** Keep-on-uninstall is
  inherent — destroying this resource removes the operator but NEVER
  the CRDs, so ClickHouse clusters and their data are never
  cascade-deleted. The modules deliberately manage no CRDs at all.
- **Native `crds/` handling also never UPGRADES them** — which is
  where the chart's `crd_hook` comes in: a pre-install/pre-upgrade
  hook job (enabled by default) that server-side-applies the four
  CRDs on every install and upgrade. Leave it enabled; disabling it
  means chart upgrades silently run new operators against old schemas
  — only sound when CRD lifecycle is managed elsewhere.
- **The hook's default image will silently age.** It runs
  `bitnami/kubectl:latest` — verified pullable today, but FROZEN since
  Bitnami's 2025 public-catalog retirement, and `latest` is exactly
  the tag that resolves differently across nodes. `crd_hook.image` is
  the override; the private-mirror preset shows the posture.
- **Server-side apply** keeps the large CRD schemas free of
  client-side last-applied-configuration annotation bloat and lets
  repeated hook runs never conflict with themselves.

## Watch Scope and RBAC Posture

By default the operator watches ONLY its own namespace — the inverse
of many operators' cluster-wide default, and the single most common
misconfiguration: a fenced operator silently ignores CHI/CHK resources
in every other namespace, which presents as KubernetesClickHouse
resources that never reconcile.

- **`watch_namespaces`** entries are regular expressions (exact names
  work as-is). Every namespace that will hold KubernetesClickHouse
  resources must be covered; `[".*"]` watches the whole cluster — the
  standard preset's choice, since one engine per cluster is the usual
  platform shape.
- **`namespace_scoped_rbac`** (chart `rbac.namespaceScoped`) swaps
  cluster-wide RBAC for namespace-scoped Roles/RoleBindings — the
  tenancy posture, paired with a watch list of exactly the install
  namespace. A wide watch with namespace-scoped RBAC would let the
  operator SEE clusters it cannot manage.
- **The CRDs stay cluster-scoped regardless** (all Kubernetes CRDs
  are): the first install on a cluster still needs permission to
  create them, and CRD upgrades ride whichever install's hook runs
  first. The per-team isolation is about workloads and RBAC, not the
  API types.

## The Operator's Own Credentials

The operator logs into every ClickHouse cluster it manages — host
management, schema propagation, and metrics scraping all run as the
`operator_credentials` user. The chart provisions the pair as a
chart-managed Secret (exported as `credentials_secret_name`) and the
operator auto-injects the matching user into every managed cluster,
network-restricted to the operator pod's address.

The chart's defaults are publicly documented
(`clickhouse_operator` / `clickhouse_operator_password`), so an unset
`operator_credentials` is unsafe outside throwaway environments. The
modules render the chart's `secret` block only when the message is
present — and the presets ship a placeholder password precisely to
force the conversation.

## Metrics

The metrics-exporter sidecar (enabled by default upstream) serves
Prometheus metrics for EVERY managed cluster on port 8888, exposed
through the chart's `<fullname>-metrics` Service and exported as
`metrics_endpoint` (empty when disabled). Its memory grows with the
number and size of managed clusters — `metrics.resources` exists for
that. `service_monitor_enabled` adds a ServiceMonitor covering both
the per-cluster endpoint (8888) and the operator's own (9999); it
requires the Prometheus Operator CRDs on the cluster, and enabling it
without them fails the install.

## Design Decisions

- **The install is blocking.** The Helm release waits for the operator
  Deployment to become Available (atomic, 600s timeout, cleanup on
  fail): an operator that never becomes ready — an unpullable image
  from a private mirror is the classic case — fails THIS apply with a
  readiness timeout instead of surfacing later as ClickHouse clusters
  that mysteriously never reconcile. The timeout also covers the
  pre-install CRD hook job.
- **The module owns namespace creation** (`create_namespace`), never
  the Helm release — pre-existing-namespace installs leave the flag
  false.
- **Chart-default-matching values render only on divergence** (watch
  scope, RBAC posture, credentials, the true-defaulted toggles, the
  image sources), so the rendered values stay minimal on both engines.
- **`fullnameOverride` pins to the resource name.** Every generated
  child name and every exported output (the Deployment, the
  credentials Secret, the `<fullname>-metrics` Service) hangs off the
  chart's fullname, so the modules pin it to `metadata.name` — and
  re-pin it AFTER the `helm_values` escape-hatch merge, the one
  deliberate exception to the escape hatch's last-word contract. This
  is also where the naming budget comes from: the longest generated
  child suffix is `-keeper-templatesd-files` (24 characters), so the
  resource name must stay at 39 characters or fewer against the
  Kubernetes 63-character cap.
- **`metrics_endpoint` replays the chart's own naming.** The exported
  URL derives from the fullname-pinned `<name>-metrics` Service and
  the exporter's 8888 port, and clears only on an explicit
  `metrics.enabled: false` — the presence-tracked field mirrors the
  upstream true default.

## Version Pins and Naming Contracts

| What | Value | Notes |
|---|---|---|
| Chart | `altinity-clickhouse-operator` at https://docs.altinity.com/clickhouse-operator/ | Pinned 0.27.2 (spec default) |
| Operator image | `altinity/clickhouse-operator` at the chart's appVersion | Chart and image versions move in lockstep |
| Metrics exporter | `altinity/metrics-exporter` at the chart's appVersion | Sidecar; per-cluster metrics on 8888 |
| CRD hook image | `bitnami/kubectl:latest` | Frozen upstream since Bitnami's 2025 catalog retirement — override for durable installs |
| CRD API groups | `clickhouse.altinity.com`, `clickhouse-keeper.altinity.com` | Four CRDs, chart-owned (`crds/` + hook) |
| Fullname | pinned to the resource name | Resource name ≤ 39 chars (longest child suffix is 24 chars) |
| Watch scope | the operator's own namespace (chart default) | `watch_namespaces` regexps widen; `[".*"]` = cluster-wide |

## IaC Twins

Pulumi (`module/values.go`) and Terraform (`locals.tf` + `main.tf`)
render identical chart values, the same fullname pinning (including
the post-merge re-pin), and the same derived outputs. Neither engine
manages CRDs — the chart owns that lifecycle on both. Keep the
typed-value rendering and the output derivations in lockstep.

## Validation Status

The rebuilt component is offline-validated: the spec's validation
tests pass, and both engines' modules carry offline rendering proofs
against the pinned chart. Live end-to-end verification on a running
cluster is pending.
