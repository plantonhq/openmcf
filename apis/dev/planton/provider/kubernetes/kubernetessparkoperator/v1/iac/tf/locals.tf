# Computed values for the KubernetesSparkOperator module. Every resolution
# here has an exact twin in the Pulumi module's locals.go / values.go —
# keep them in lockstep.
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
  helm_chart_repo = "https://apache.github.io/spark-kubernetes-operator"
  helm_chart_name = "spark-kubernetes-operator"

  # Release name = metadata.name.
  release_name = var.metadata.name

  # Chart version resolved to the pinned default when unset, so both
  # engines install the same chart whether or not the platform's
  # defaulting middleware ran — mirror of the Pulumi module's
  # DefaultChartVersion. 1.8.0 is the newest SERVED chart (= operator
  # appVersion 1.0.0, verified against the repository index). The
  # spark.apache.org CRDs ship from the chart's crds/ directory: Helm
  # installs them once and NEVER upgrades them — bumping this version
  # does not touch the CRDs (apply the new release's CRD files manually
  # when a bump changes them).
  chart_version = coalesce(try(var.spec.chart_version, null), "1.8.0")

  namespace = var.spec.namespace

  # Resource-identity labels stamped on the namespace this module creates
  # (never injected into the chart's own resources — Helm owns those).
  labels = merge(
    {
      "planton.ai/resource"      = "true"
      "planton.ai/resource-name" = var.metadata.name
      "planton.ai/resource-kind" = "KubernetesSparkOperator"
    },
    var.metadata.id != null && var.metadata.id != "" ? { "planton.ai/resource-id" = var.metadata.id } : {},
    var.metadata.org != null && var.metadata.org != "" ? { "planton.ai/organization" = var.metadata.org } : {},
    var.metadata.env != null && var.metadata.env != "" ? { "planton.ai/environment" = var.metadata.env } : {}
  )

  # ---- workload surface (chart: workloadResources) ---------------------------
  # The watch scope and the workload RBAC are ONE chart surface, decided
  # together. Empty namespaces = cluster-wide: the chart's defaults stand
  # (workload ClusterRole, operator watches everywhere, no namespaces
  # created). Non-empty = the fenced posture: the chart CREATES each
  # listed namespace, plants the service account and a namespace-scoped
  # Role/RoleBinding in each, drops the workload ClusterRole, and
  # overrideWatchedNamespaces (chart default true) wires the operator's
  # spark.kubernetes.operator.watched.namespaces property from the same
  # list — one value, one truth.
  workload_namespaces      = try(var.spec.workload.namespaces, [])
  workload_fenced          = length(local.workload_namespaces) > 0
  workload_service_account = coalesce(try(var.spec.workload.service_account, null), "spark")

  # ---- operator properties (chart: operatorConfiguration) --------------------
  # The operator is properties-file configured. The chart APPENDS this
  # document over its built-in defaults (operatorConfiguration.append,
  # chart default true — kept). Leader election is module-owned: any
  # replica count beyond 1 REQUIRES it (the chart's own contract), so the
  # property renders exactly when replicas > 1 — never a spec knob that
  # could drift from the replica count.
  leader_election_needed = try(var.spec.replicas, 1) > 1
  operator_properties = merge(
    try(var.spec.operator_properties, {}),
    local.leader_election_needed ? { "spark.kubernetes.operator.leaderElection.enabled" = "true" } : {}
  )
  operator_properties_file = join("\n", [
    for k in sort(keys(local.operator_properties)) : "${k}=${local.operator_properties[k]}"
  ])

  # ---- dynamic (hot-reload) properties ---------------------------------------
  # Rendered only when enabled: the chart creates the ConfigMap
  # (dynamicConfig.create) and the RBAC that lets the operator watch it
  # (operatorRbac.configManagement, chart default true — kept).
  dynamic_config_enabled = try(var.spec.dynamic_config.enabled, false)
  dynamic_config = local.dynamic_config_enabled ? {
    enable = true
    create = true
    data   = try(var.spec.dynamic_config.properties, {})
  } : null

  # ---- operator container resources (shared ContainerResources) --------------
  # Twin of the Pulumi module's resourcesMap. The chart ships REAL
  # defaults here (1 CPU / 2Gi, requests = limits) — the resources key
  # renders only when the spec sets them, so the upstream-tested sizing
  # stands otherwise. Helm deep-merges per key: a partial block overrides
  # only the halves it carries.
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

  # ---- operator container block ----------------------------------------------
  operator_container = {
    for k, v in {
      jvmArgs   = try(var.spec.jvm_args, "") != "" ? var.spec.jvm_args : null
      resources = local.operator_resources != null && length(local.operator_resources) > 0 ? local.operator_resources : null
    } : k => v if v != null
  }

  # ---- operator pod block (scheduling + image pulls) -------------------------
  operator_pod = {
    for k, v in {
      operatorContainer = length(local.operator_container) > 0 ? local.operator_container : null
      nodeSelector      = length(try(var.spec.scheduling.node_selector, {})) > 0 ? var.spec.scheduling.node_selector : null
      tolerations = length(try(var.spec.scheduling.tolerations, [])) > 0 ? [
        for t in var.spec.scheduling.tolerations : {
          for tk, tv in {
            key               = try(t.key, "") != "" ? t.key : null
            operator          = try(t.operator, "") != "" ? t.operator : null
            value             = try(t.value, "") != "" ? t.value : null
            effect            = try(t.effect, "") != "" ? t.effect : null
            tolerationSeconds = try(t.toleration_seconds, null)
          } : tk => tv if tv != null
        }
      ] : null
      priorityClassName = try(var.spec.scheduling.priority_class_name, "") != "" ? var.spec.scheduling.priority_class_name : null
    } : k => v if v != null
  }

  # ---- operator image (air-gap/private-mirror registry replacement) ----------
  # image_registry replaces ONLY the registry part of the operator image
  # (chart default `apache/spark-kubernetes-operator`, Docker Hub
  # implied); the tag stays the chart's appVersion-locked default. Spark
  # WORKLOAD images ride each SparkApplication's own image field — this
  # never rewrites those. Twin of the Pulumi module.
  operator_image = {
    for k, v in {
      repository = try(var.spec.image_registry, "") != "" ? "${var.spec.image_registry}/apache/spark-kubernetes-operator" : null
    } : k => v if v != null
  }

  # ---- typed chart values (twin of the Pulumi module's buildHelmValues) ------
  # Chart-default-matching values render only on divergence — with ONE
  # deliberate always-rendered block: the RBAC NAME RE-PINS. The chart
  # hardcodes every cluster-scoped RBAC name as a plain value
  # ("spark-operator-clusterrole", …), which makes a second install
  # anywhere on the cluster collide by construction. Deriving the names
  # from the release identity makes instances coexist — the same defense
  # as the fullname pin, applied to the chart's values-borne names.
  # The workload SERVICE ACCOUNT name deliberately stays the upstream
  # contract ("spark" unless overridden): SparkApplications reference it
  # by that conventional name.
  typed_values = {
    for k, v in {
      # nameOverride is THIS chart's identity pin: every named object
      # (the operator Deployment, PDB selector, NetworkPolicy) renders
      # from the `spark-operator.name` helper (default .Chart.Name |
      # nameOverride) — the chart defines a fullname helper but NO
      # template consumes it, so a fullnameOverride pin is a no-op and
      # the Deployment would keep the chart's constant name
      # `spark-kubernetes-operator` (verified live: the pinned name was
      # NotFound while the chart-named Deployment served).
      nameOverride = local.release_name

      operatorDeployment = {
        for dk, dv in {
          # Rendered on presence — an explicit 1 re-states the chart
          # default harmlessly; >1 pairs with the leader-election
          # property rendered into operator_properties above (the chart
          # REFUSES multi-replica installs without it, by design).
          replicas    = try(var.spec.replicas, null)
          operatorPod = length(local.operator_pod) > 0 ? local.operator_pod : null
        } : dk => dv if dv != null
      }

      operatorRbac = {
        serviceAccount = { name = local.release_name }
        clusterRole = {
          name = "${local.release_name}-clusterrole"
        }
        clusterRoleBinding = {
          name = "${local.release_name}-clusterrolebinding"
        }
        configManagement = {
          roleName        = "${local.release_name}-config-monitor"
          roleBindingName = "${local.release_name}-config-monitor-binding"
        }
      }

      workloadResources = {
        for wk, wv in {
          # The fenced posture: chart creates the listed namespaces and
          # watches exactly them (overrideWatchedNamespaces default true
          # wires the property from this list).
          namespaces = local.workload_fenced ? { create = true, data = local.workload_namespaces } : null

          serviceAccount = { name = local.workload_service_account }

          # Cluster-wide: the workload ClusterRole (chart default) with a
          # release-derived name. Fenced: per-namespace Roles replace it.
          clusterRole = {
            create = !local.workload_fenced
            name   = "${local.release_name}-workload-clusterrole"
          }
          role = {
            create = local.workload_fenced
            name   = "${local.release_name}-workload-role"
          }
          # The chart derives this binding's roleRef ITSELF from
          # clusterRole.create (ClusterRole when true, Role when false) —
          # only the name is ours to pin (template-verified).
          roleBinding = {
            name = "${local.release_name}-workload-rolebinding"
          }
        } : wk => wv if wv != null
      }

      operatorConfiguration = {
        for ck, cv in {
          "spark-operator.properties" = length(local.operator_properties) > 0 ? local.operator_properties_file : null
          dynamicConfig               = local.dynamic_config
        } : ck => cv if cv != null
      }

      image = length(local.operator_image) > 0 ? local.operator_image : null

      imagePullSecrets = length(try(var.spec.image_pull_secrets, [])) > 0 ? [for s in var.spec.image_pull_secrets : { name = s }] : null
    } : k => v if v != null && (!can(length(v)) || length(v) > 0)
  }
}
