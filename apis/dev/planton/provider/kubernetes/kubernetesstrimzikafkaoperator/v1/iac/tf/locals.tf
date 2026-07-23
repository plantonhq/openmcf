# Computed values for the KubernetesStrimziKafkaOperator module. Every
# resolution here has an exact twin in the Pulumi module's locals.go /
# values.go — keep them in lockstep.
#
# HCL DISCIPLINE (applies to every conditional object in this file):
# conditional entries are written as `key = cond ? value : null` inside ONE
# object literal, pruned with `{ for k, v in {...} : k => v if v != null }`.
# The two tempting alternatives are both broken: `cond ? {...} : {}`
# ternaries fail plan-time type unification when branches carry different
# attributes, and `merge(concat(cond ? [{...}] : [], ...)...)` silently
# UNIFIES primitive-only sibling objects into map(string) — numbers and
# booleans arrive in the chart values as strings. The null-prune form
# preserves every value's type.
#
# Optional nested blocks are read with try(): HCL's && does NOT
# short-circuit, so chained null checks still dereference the null — and
# var.spec is typed 'any', so an absent attribute is an error, not a null.

locals {
  # Chart identity — must stay byte-identical with the Pulumi module's
  # vars: cross-engine chart drift deploys two different products from one
  # manifest. The SERVED index (https://strimzi.io/charts/) governs the
  # version — the Chart.yaml inside the Strimzi source tree carries a
  # build-time placeholder and never reflects the served version.
  helm_chart_repo = "https://strimzi.io/charts/"
  helm_chart_name = "strimzi-kafka-operator"

  # Release name = metadata.name. The chart's operand-facing resources use
  # fixed Strimzi names (strimzi-cluster-operator), so a SECOND install in
  # one cluster additionally needs create_global_resources false — see the
  # spec comment.
  release_name = var.metadata.name

  # Chart version resolved to the pinned default when unset, so both
  # engines install the same chart whether or not the platform's
  # defaulting middleware ran — mirror of the Pulumi module's
  # DefaultChartVersion. Chart and operator versions move TOGETHER for
  # this chart (chart 1.1.0 ships operator 1.1.0).
  chart_version = coalesce(try(var.spec.chart_version, null), "1.1.0")

  namespace = var.spec.namespace

  # Resource-identity labels stamped on the namespace this module creates
  # (never injected into the chart's own resources — Helm owns those).
  labels = merge(
    {
      "planton.ai/resource"      = "true"
      "planton.ai/resource-name" = var.metadata.name
      "planton.ai/resource-kind" = "KubernetesStrimziKafkaOperator"
    },
    var.metadata.id != null && var.metadata.id != "" ? { "planton.ai/resource-id" = var.metadata.id } : {},
    var.metadata.org != null && var.metadata.org != "" ? { "planton.ai/organization" = var.metadata.org } : {},
    var.metadata.env != null && var.metadata.env != "" ? { "planton.ai/environment" = var.metadata.env } : {}
  )

  # ---- operator container resources (shared ContainerResources shape) -------
  # Twin of the Pulumi module's resourcesMap. The chart SHIPS default
  # requests/limits (requests 200m/384Mi, limits 1000m/384Mi) — the
  # resources key renders only when the spec sets them, so the chart
  # defaults survive an empty spec. Helm deep-merges per key, so a
  # partial spec block overrides only the halves it carries.
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

  # ---- typed chart values (twin of the Pulumi module's buildHelmValues) ------
  # Chart-default-matching values render only on divergence (the watch
  # scope, the true-defaulted toggles, the image source), so the rendered
  # values stay minimal on both engines.
  typed_values = {
    for k, v in {
      replicas  = try(var.spec.replicas, null)
      resources = local.operator_resources != null && length(local.operator_resources) > 0 ? local.operator_resources : null

      # Watch scope maps to two independent chart values: cluster-wide
      # RBAC (watchAnyNamespace) or a namespace LIST (watchNamespaces —
      # the installation namespace is always watched in addition). Spec
      # CEL rules make the two arms mutually exclusive, so at most one
      # renders; both unset means the operator watches its own namespace
      # (the chart default).
      watchAnyNamespace = try(var.spec.watch.any_namespace, false) ? true : null
      watchNamespaces   = length(try(var.spec.watch.namespaces, [])) > 0 ? var.spec.watch.namespaces : null

      # Both integers in the chart's values file; rendered on presence so
      # the chart defaults (120000 / 300000) survive an empty spec.
      fullReconciliationIntervalMs = try(var.spec.full_reconciliation_interval_ms, null)
      operationTimeoutMs           = try(var.spec.operation_timeout_ms, null)

      # logLevel renders whenever the spec carries a value (the chart
      # default is INFO via an env-substituted expression; a literal
      # level replaces it cleanly).
      logLevel     = try(var.spec.log_level, "") != "" ? var.spec.log_level : null
      featureGates = try(var.spec.feature_gates, "") != "" ? var.spec.feature_gates : null

      kubernetesServiceDnsDomain = try(var.spec.kubernetes_service_dns_domain, "") != "" ? var.spec.kubernetes_service_dns_domain : null

      # The chart nests a single enable flag. Rendered on presence — an
      # explicit true re-states the chart default harmlessly, an explicit
      # false is the actual opt-out.
      leaderElection = try(var.spec.leader_election_enabled, null) != null ? { enable = var.spec.leader_election_enabled } : null

      # Operand policy generation toggles — both default true in the
      # chart; rendered on presence so an explicit opt-out reaches the
      # chart and an empty spec changes nothing.
      generateNetworkPolicy       = try(var.spec.generate_network_policy, null)
      generatePodDisruptionBudget = try(var.spec.generate_pod_disruption_budget, null)

      # createGlobalResources defaults true; false is the second-install
      # posture (the fixed-name ClusterRoles are owned by the first
      # release).
      createGlobalResources = try(var.spec.create_global_resources, null)

      nodeSelector = length(try(var.spec.node_selector, {})) > 0 ? var.spec.node_selector : null
      tolerations = length(try(var.spec.tolerations, [])) > 0 ? [
        for t in var.spec.tolerations : {
          for tk, tv in {
            key               = try(t.key, "") != "" ? t.key : null
            operator          = try(t.operator, "") != "" ? t.operator : null
            value             = try(t.value, "") != "" ? t.value : null
            effect            = try(t.effect, "") != "" ? t.effect : null
            tolerationSeconds = try(t.toleration_seconds, null)
          } : tk => tv if tv != null
        }
      ] : null

      # defaultImageRegistry/Repository/Tag steer EVERY Strimzi image
      # (the operator and all operand images it deploys) — the air-gap
      # path. Pull secrets ride the chart's image.imagePullSecrets list
      # (raw Kubernetes object list, piped into the pod spec).
      defaultImageRegistry   = try(var.spec.image.registry, "") != "" ? var.spec.image.registry : null
      defaultImageRepository = try(var.spec.image.repository, "") != "" ? var.spec.image.repository : null
      defaultImageTag        = try(var.spec.image.tag, "") != "" ? var.spec.image.tag : null
      image = length(try(var.spec.image_pull_secrets, [])) > 0 ? {
        imagePullSecrets = [for s in var.spec.image_pull_secrets : { name = s }]
      } : null
    } : k => v if v != null
  }
}
