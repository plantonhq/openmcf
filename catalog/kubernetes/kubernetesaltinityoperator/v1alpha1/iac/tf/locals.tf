# Computed values for the KubernetesAltinityOperator module. Every
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
# short-circuit, so chained null checks still dereference the null.

locals {
  # Chart identity — must stay byte-identical with the Pulumi module's
  # vars: cross-engine chart drift deploys two different products from one
  # manifest.
  helm_chart_repo = "https://docs.altinity.com/clickhouse-operator/"
  helm_chart_name = "altinity-clickhouse-operator"

  # Release name = metadata.name.
  release_name = var.metadata.name

  # Chart version resolved to the pinned default when unset, so both
  # engines install the same chart whether or not the platform's
  # defaulting middleware ran — mirror of the Pulumi module's
  # DefaultChartVersion. Chart versions track operator releases
  # one-to-one (chart 0.27.2 = operator image 0.27.2).
  chart_version = coalesce(try(var.spec.chart_version, null), "0.27.2")

  namespace = var.spec.namespace

  # Resource-identity labels stamped on the namespace this module creates
  # (never injected into the chart's own resources — Helm owns those).
  labels = merge(
    {
      "planton.ai/resource"      = "true"
      "planton.ai/resource-name" = var.metadata.name
      "planton.ai/resource-kind" = "KubernetesAltinityOperator"
    },
    var.metadata.id != null && var.metadata.id != "" ? { "planton.ai/resource-id" = var.metadata.id } : {},
    var.metadata.org != null && var.metadata.org != "" ? { "planton.ai/organization" = var.metadata.org } : {},
    var.metadata.env != null && var.metadata.env != "" ? { "planton.ai/environment" = var.metadata.env } : {}
  )

  # ---- deployment name (twin of the Pulumi module's deploymentName) ---------
  # The module pins the chart's fullnameOverride to the resource name (see
  # typed_values), and the chart names its Deployment exactly the fullname
  # (templates/generated/Deployment-clickhouse-operator.yaml) — so the
  # Deployment IS the release name.
  deployment_name = local.release_name

  # ---- credentials Secret name (twin of the Pulumi module) ------------------
  # The chart names its credentials Secret exactly the fullname
  # (templates/generated/Secret-clickhouse-operator.yaml). The Secret is
  # created unconditionally (chart secret.create defaults true) — with the
  # publicly documented default credentials when operator_credentials is
  # unset, which is why the spec calls unset UNSAFE FOR PRODUCTION.
  credentials_secret_name = local.release_name

  # Metrics-exporter presence — the proto default (true) mirrors the chart
  # default, so an absent metrics block means the sidecar runs.
  metrics_enabled = coalesce(try(var.spec.metrics.enabled, null), true)

  # ---- metrics endpoint (twin of the Pulumi module's metricsEndpoint) -------
  # The chart's metrics Service is named "<fullname>-metrics" and carries
  # the 8888 ch-metrics port ONLY while metrics.enabled (verified in
  # templates/generated/Service-clickhouse-operator-metrics.yaml) — so the
  # endpoint is empty when the exporter is disabled.
  metrics_endpoint = local.metrics_enabled ? "http://${local.release_name}-metrics.${local.namespace}.svc.cluster.local:8888/metrics" : ""

  # ---- operator container resources (shared ContainerResources) -------------
  # Twin of the Pulumi module. The chart ships NO default requests/limits;
  # the key renders only when the spec sets values, so the rendered values
  # carry no empty resources section. Helm deep-merges per key, so a
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

  # ---- metrics-exporter sidecar resources (same shape and rules) ------------
  metrics_resources = try(var.spec.metrics.resources, null) == null ? null : {
    for k, v in {
      limits = try(var.spec.metrics.resources.limits, null) == null ? null : {
        for lk, lv in {
          cpu    = try(var.spec.metrics.resources.limits.cpu, "") != "" ? var.spec.metrics.resources.limits.cpu : null
          memory = try(var.spec.metrics.resources.limits.memory, "") != "" ? var.spec.metrics.resources.limits.memory : null
        } : lk => lv if lv != null
      }
      requests = try(var.spec.metrics.resources.requests, null) == null ? null : {
        for rk, rv in {
          cpu    = try(var.spec.metrics.resources.requests.cpu, "") != "" ? var.spec.metrics.resources.requests.cpu : null
          memory = try(var.spec.metrics.resources.requests.memory, "") != "" ? var.spec.metrics.resources.requests.memory : null
        } : rk => rv if rv != null
      }
    } : k => v if v != null && length(v) > 0
  }

  # ---- operator.image (rendered per half — deep-merges over the chart
  # defaults, leaving registry, pullPolicy and the appVersion-derived tag
  # intact) --------------------------------------------------------------------
  operator_image = {
    for k, v in {
      repository = try(var.spec.image.repo, "") != "" ? var.spec.image.repo : null
      tag        = try(var.spec.image.tag, "") != "" ? var.spec.image.tag : null
    } : k => v if v != null
  }

  # ---- crdHook.image (same per-half rendering; the chart default is
  # bitnami/kubectl:latest — see the spec's KNOW THIS note) --------------------
  crd_hook_image = {
    for k, v in {
      repository = try(var.spec.crd_hook.image.repo, "") != "" ? var.spec.crd_hook.image.repo : null
      tag        = try(var.spec.crd_hook.image.tag, "") != "" ? var.spec.crd_hook.image.tag : null
    } : k => v if v != null
  }

  # ---- the operator block (renders only when the spec sets something) -------
  operator_values = {
    for k, v in {
      image     = length(local.operator_image) > 0 ? local.operator_image : null
      resources = local.operator_resources != null && length(local.operator_resources) > 0 ? local.operator_resources : null
    } : k => v if v != null
  }

  # ---- the metrics block -----------------------------------------------------
  # enabled renders on presence — an explicit true re-states the chart
  # default harmlessly, an explicit false is the actual opt-out.
  metrics_values = {
    for k, v in {
      enabled   = try(var.spec.metrics.enabled, null)
      resources = local.metrics_resources != null && length(local.metrics_resources) > 0 ? local.metrics_resources : null
    } : k => v if v != null
  }

  # ---- the crdHook block (same presence rule for enabled) --------------------
  crd_hook_values = {
    for k, v in {
      enabled = try(var.spec.crd_hook.enabled, null)
      image   = length(local.crd_hook_image) > 0 ? local.crd_hook_image : null
    } : k => v if v != null
  }

  # ---- operator credentials (chart secret.username / secret.password) --------
  # Omitted ENTIRELY when the message is absent (the chart then provisions
  # its documented defaults); the proto's username default resolves here so
  # a password-only spec still names the standard operator user.
  credentials = try(var.spec.operator_credentials, null) == null ? null : {
    username = coalesce(try(var.spec.operator_credentials.username, null), "clickhouse_operator")
    password = var.spec.operator_credentials.password
  }

  # ---- typed chart values (twin of the Pulumi module's buildHelmValues) ------
  # Chart-default-matching values render only on divergence, so the
  # rendered values stay minimal on both engines.
  typed_values = {
    for k, v in {
      # fullnameOverride pins the chart's fullname to the resource name
      # (the catalog's Helm-kind identity convention). Load-bearing: the
      # Deployment / credentials Secret / metrics Service names — and the
      # exported outputs built from them — all derive from the fullname.
      # Twin of the Pulumi module; re-pinned after the helm_values merge
      # in main.tf.
      fullnameOverride = local.release_name

      # Entries are regexps; empty (the chart default) = the operator
      # watches only its own namespace.
      watchNamespaces = length(try(var.spec.watch_namespaces, [])) > 0 ? var.spec.watch_namespaces : null

      # Plain bool (no presence): false IS the chart default
      # (cluster-wide RBAC), so only true renders.
      rbac = try(var.spec.namespace_scoped_rbac, false) ? { namespaceScoped = true } : null

      secret = local.credentials

      metrics = length(local.metrics_values) > 0 ? local.metrics_values : null

      crdHook = length(local.crd_hook_values) > 0 ? local.crd_hook_values : null

      operator = length(local.operator_values) > 0 ? local.operator_values : null

      # Plain bool (no presence): false IS the chart default, so only
      # true renders. Requires the Prometheus Operator CRDs on the
      # cluster.
      serviceMonitor = try(var.spec.service_monitor_enabled, false) ? { enabled = true } : null

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

      # Pull secrets ride the chart's top-level imagePullSecrets list
      # (raw Kubernetes object list, piped into the operator pod spec).
      imagePullSecrets = length(try(var.spec.image_pull_secrets, [])) > 0 ? [for s in var.spec.image_pull_secrets : { name = s }] : null
    } : k => v if v != null
  }
}
