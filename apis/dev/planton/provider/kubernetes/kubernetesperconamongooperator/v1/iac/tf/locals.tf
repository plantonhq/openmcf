# Computed values for the KubernetesPerconaMongoOperator module. Every
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
  # manifest. Percona publishes the operator charts for every database
  # product from this one repository.
  helm_chart_repo = "https://percona.github.io/percona-helm-charts"
  helm_chart_name = "psmdb-operator"

  # Release name = metadata.name. The chart derives every resource name
  # from the release (Deployment, ServiceAccount, RBAC through its
  # fullname helper), so distinct release names keep multiple
  # namespace-scoped installations from colliding.
  release_name = var.metadata.name

  # Chart version resolved to the pinned default when unset, so both
  # engines install the same chart whether or not the platform's
  # defaulting middleware ran — mirror of the Pulumi module's
  # DefaultChartVersion. Chart and operator versions move TOGETHER for
  # this chart (chart 1.22.0 ships operator 1.22.0); the chart pin
  # governs.
  chart_version = coalesce(try(var.spec.chart_version, null), "1.22.0")

  namespace = var.spec.namespace

  # Resource-identity labels stamped on the namespace this module creates
  # (never injected into the chart's own resources — Helm owns those).
  labels = merge(
    {
      "planton.ai/resource"      = "true"
      "planton.ai/resource-name" = var.metadata.name
      "planton.ai/resource-kind" = "KubernetesPerconaMongoOperator"
    },
    var.metadata.id != null && var.metadata.id != "" ? { "planton.ai/resource-id" = var.metadata.id } : {},
    var.metadata.org != null && var.metadata.org != "" ? { "planton.ai/organization" = var.metadata.org } : {},
    var.metadata.env != null && var.metadata.env != "" ? { "planton.ai/environment" = var.metadata.env } : {}
  )

  # ---- operator container resources (shared ContainerResources shape) -------
  # Twin of the Pulumi module's resourcesMap. The chart ships no default
  # requests/limits — rendering is purely additive.
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

  # ---- image ------------------------------------------------------------------
  # The image override renders only the halves that are set — the chart
  # composes "<repository>:<tag>" itself, so an unset half keeps the
  # chart's default for it (repository
  # percona/percona-server-mongodb-operator; tag = the chart version).
  image_values = {
    for k, v in {
      repository = try(var.spec.image.repository, "") != "" ? var.spec.image.repository : null
      tag        = try(var.spec.image.tag, "") != "" ? var.spec.image.tag : null
    } : k => v if v != null
  }

  # ---- typed chart values (twin of the Pulumi module's buildHelmValues) ------
  # Chart-default-matching values render only on divergence
  # (watchAllNamespaces, logStructured, disableTelemetry all default false
  # upstream), so the rendered values stay minimal on both engines.
  typed_values = {
    for k, v in {
      replicaCount = try(var.spec.replicas, null)
      resources    = local.operator_resources != null && length(local.operator_resources) > 0 ? local.operator_resources : null

      # Watch scope maps to two independent chart values: cluster-wide
      # RBAC (watchAllNamespaces) or a comma-joined namespace fence
      # (watchNamespace). Spec CEL rules make the two arms mutually
      # exclusive, so at most one renders; both unset means the operator
      # watches its own namespace (the upstream default). The chart's own
      # createNamespace value (create the WATCHED namespaces) is never
      # rendered — the module owns only the installation namespace, and
      # watched namespaces must already exist.
      watchAllNamespaces = try(var.spec.watch.cluster_wide, false) ? true : null
      watchNamespace     = length(try(var.spec.watch.namespaces, [])) > 0 ? join(",", var.spec.watch.namespaces) : null

      # Rendered as a STRING to match the chart's own declared type (its
      # default is the string "1"); the deployment template quotes the
      # value into the MAX_CONCURRENT_RECONCILES environment variable
      # either way, and the string keeps both engines byte-identical with
      # the chart's values file.
      maxConcurrentReconciles = try(var.spec.max_concurrent_reconciles, null) != null ? tostring(var.spec.max_concurrent_reconciles) : null

      # logStructured matches the chart default (false) — rendered only
      # when on. logLevel renders whenever the spec carries a value (the
      # chart default is "INFO"; re-stating it is harmless and keeps
      # rendering presence-driven, not value-driven).
      logStructured = try(var.spec.log.structured, false) ? true : null
      logLevel      = try(var.spec.log.level, "") != "" ? var.spec.log.level : null

      # Chart default false (telemetry on) — rendered only on explicit
      # opt-out.
      disableTelemetry = try(var.spec.disable_telemetry, false) ? true : null

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

      # Pull secrets are the raw Kubernetes object list ([{name: ...}]) —
      # the chart pipes imagePullSecrets straight into the pod spec with
      # toYaml.
      imagePullSecrets = length(try(var.spec.image_pull_secrets, [])) > 0 ? [
        for s in var.spec.image_pull_secrets : { name = s }
      ] : null
      image = length(local.image_values) > 0 ? local.image_values : null
    } : k => v if v != null
  }
}
