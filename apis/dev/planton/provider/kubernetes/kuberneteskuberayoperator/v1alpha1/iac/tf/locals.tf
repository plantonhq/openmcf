# Computed values for the KubernetesKubeRayOperator module. Every
# resolution here has an exact twin in the Pulumi module's locals.go /
# values.go — keep them in lockstep.
#
# HCL DISCIPLINE (applies to every conditional object in this file):
# conditional entries are written as `key = cond ? value : null` inside ONE
# object literal, pruned with `{ for k, v in {...} : k => v if v != null }`.
# The two tempting alternatives are both broken: `cond ? {...} : {}`
# ternaries fail plan-time type unification when branches carry different
# attributes, and merge() over conditional lists silently UNIFIES
# primitive-only sibling objects into map(string). The null-prune form
# preserves every value's type.
#
# Optional nested blocks are read with try(): HCL's && does NOT
# short-circuit, so chained null checks still dereference the null.

locals {
  # Chart identity — must stay byte-identical with the Pulumi module's
  # vars: cross-engine chart drift deploys two different products from one
  # manifest.
  helm_chart_repo = "https://ray-project.github.io/kuberay-helm"
  helm_chart_name = "kuberay-operator"

  # Release name = metadata.name.
  release_name = var.metadata.name

  # Chart version resolved to the pinned default when unset, so both
  # engines install the same chart whether or not the platform's
  # defaulting middleware ran — mirror of the Pulumi module's
  # DefaultChartVersion. Chart 1.6.2 pairs with operator image
  # quay.io/kuberay/operator:v1.6.2. The three ray.io CRDs ship from the
  # chart's crds/ directory: Helm installs them once and NEVER upgrades
  # them — bumping this version does not touch the CRDs (apply the new
  # release's CRD files manually when a bump changes them).
  chart_version = coalesce(try(var.spec.chart_version, null), "1.6.2")

  namespace = var.spec.namespace

  # Resource-identity labels stamped on the namespace this module creates
  # (never injected into the chart's own resources — Helm owns those).
  labels = merge(
    {
      "planton.ai/resource"      = "true"
      "planton.ai/resource-name" = var.metadata.name
      "planton.ai/resource-kind" = "KubernetesKubeRayOperator"
    },
    var.metadata.id != null && var.metadata.id != "" ? { "planton.ai/resource-id" = var.metadata.id } : {},
    var.metadata.org != null && var.metadata.org != "" ? { "planton.ai/organization" = var.metadata.org } : {},
    var.metadata.env != null && var.metadata.env != "" ? { "planton.ai/environment" = var.metadata.env } : {}
  )

  # ---- watch scope (chart: watchNamespace — singular key, list value) --------
  # Empty = cluster-wide (the chart default and the normal
  # one-operator-per-cluster posture) — the key stays unrendered so the
  # chart neither passes --watch-namespace nor scopes its RBAC. Non-empty
  # feeds the operator's --watch-namespace flag verbatim.
  watch_namespaces = try(var.spec.watch_namespaces, [])

  # ---- feature gates (chart: featureGates — a LIST) --------------------------
  # Helm LISTS REPLACE, never merge: rendering only the spec's entries
  # would silently DROP every chart-default gate. So when the spec flips
  # any gate, the FULL list renders: the chart's defaults (below, verified
  # against the pinned 1.6.2 values.yaml) overridden by name from the
  # spec, then spec gates the defaults don't know appended. Keep this
  # default list in lockstep with the Pulumi module's
  # chartDefaultFeatureGates AND re-verify it on every chart bump.
  feature_gate_defaults = [
    { name = "RayClusterStatusConditions", enabled = true },
    { name = "RayJobDeletionPolicy", enabled = true },
    { name = "RayMultiHostIndexing", enabled = true },
    { name = "RayServiceIncrementalUpgrade", enabled = false },
    { name = "RayCronJob", enabled = false },
  ]
  spec_feature_gates = try(var.spec.feature_gates, [])
  spec_gate_by_name  = { for g in local.spec_feature_gates : g.name => g.enabled }
  default_gate_names = [for g in local.feature_gate_defaults : g.name]
  feature_gates = concat(
    [for g in local.feature_gate_defaults : {
      name    = g.name
      enabled = lookup(local.spec_gate_by_name, g.name, g.enabled)
    }],
    [for g in local.spec_feature_gates : {
      name    = g.name
      enabled = g.enabled
    } if !contains(local.default_gate_names, g.name)]
  )

  # ---- metrics (chart: metrics.enabled / metrics.serviceMonitor) -------------
  # enabled renders only on an EXPLICIT false (proto optional bool; chart
  # default true). The ServiceMonitor requires the monitoring.coreos.com
  # CRDs on the cluster (KubernetesKubePrometheusStack) — the install
  # FAILS without them.
  metrics_block = {
    for k, v in {
      enabled        = try(var.spec.metrics_enabled, null) == false ? false : null
      serviceMonitor = try(var.spec.service_monitor_enabled, false) ? { enabled = true } : null
    } : k => v if v != null
  }

  # ---- operator container resources (shared ContainerResources) --------------
  # Twin of the Pulumi module's resourcesMap. The chart ships REAL
  # defaults here (100m CPU / 512Mi limits — upstream sizes ~500MB per
  # 500 managed Ray pods) — the resources key renders only when the spec
  # sets them, so the upstream-tested sizing stands otherwise. Helm
  # deep-merges per key: a partial block overrides only the halves it
  # carries.
  operator_resources = try(var.spec.resources, null) == null ? null : {
    for k, v in {
      limits = try(var.spec.resources.limits, null) == null ? null : {
        for lk, lv in {
          cpu    = try(var.spec.resources.limits.cpu, "") != "" ? var.spec.resources.limits.cpu : null
          memory = try(var.spec.resources.limits.memory, "") != "" ? var.spec.resources.limits.memory : null
        } : lk => lv if lv != null
      }
      requests = try(var.spec.resources.requests, null) == null ? null : {
        for rk, rv in {
          cpu    = try(var.spec.resources.requests.cpu, "") != "" ? var.spec.resources.requests.cpu : null
          memory = try(var.spec.resources.requests.memory, "") != "" ? var.spec.resources.requests.memory : null
        } : rk => rv if rv != null
      }
    } : k => v if v != null && length(v) > 0
  }

  # ---- operator image (air-gap/private-mirror registry replacement) ----------
  # image_registry replaces ONLY the registry part of the operator image
  # (chart default quay.io/kuberay/operator — the swap drops quay.io);
  # the tag stays the chart's appVersion-locked default. Ray CLUSTER
  # images ride each KubernetesRayCluster's own image field — this never
  # rewrites those. Twin of the Pulumi module.
  operator_image = {
    for k, v in {
      repository = try(var.spec.image_registry, "") != "" ? "${var.spec.image_registry}/kuberay/operator" : null
    } : k => v if v != null
  }

  # ---- typed chart values (twin of the Pulumi module's buildHelmValues) ------
  # Chart-default-matching values render only on divergence — with ONE
  # deliberate always-rendered set: the NAME PINS. The chart hardcodes
  # nameOverride, fullnameOverride, AND serviceAccount.name to
  # "kuberay-operator" in its values, so every install collapses onto
  # the same child names (and the same ServiceAccount) by construction.
  # Pinning all three to metadata.name keeps instances distinguishable
  # (the catalog's Helm-kind identity convention).
  typed_values = {
    for k, v in {
      nameOverride     = local.release_name
      fullnameOverride = local.release_name
      serviceAccount   = { name = local.release_name }

      watchNamespace = length(local.watch_namespaces) > 0 ? local.watch_namespaces : null

      # Explicit false only — an unset optional keeps the chart default
      # (true: safe for single replicas, required for standbys).
      leaderElectionEnabled = try(var.spec.leader_election_enabled, null) == false ? false : null

      # batchScheduler.name is the standard knob; batchScheduler.enabled
      # is the deprecated legacy flag and MUTUALLY EXCLUSIVE with name
      # (the chart errors when both are set) — never render it.
      batchScheduler = try(var.spec.batch_scheduler, "") != "" ? { name = var.spec.batch_scheduler } : null

      featureGates = length(local.spec_feature_gates) > 0 ? local.feature_gates : null

      metrics = length(local.metrics_block) > 0 ? local.metrics_block : null

      resources = local.operator_resources != null && length(local.operator_resources) > 0 ? local.operator_resources : null

      image = length(local.operator_image) > 0 ? local.operator_image : null
    } : k => v if v != null && (!can(length(v)) || length(v) > 0)
  }
}
