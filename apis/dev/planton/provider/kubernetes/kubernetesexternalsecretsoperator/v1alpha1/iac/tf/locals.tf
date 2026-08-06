# Computed values for the KubernetesExternalSecretsOperator module. Every
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
  # Chart identity — must stay byte-identical with the Pulumi module's vars:
  # cross-engine chart-name drift deploys two different products from one
  # manifest.
  helm_chart_name = "external-secrets"
  helm_chart_repo = "https://charts.external-secrets.io"

  # Release name fixed to the chart name: one External Secrets Operator per
  # cluster is an upstream architectural constraint (cluster-scoped CRDs and
  # webhook configuration), so a manifest-derived name would only invite a
  # second broken install.
  release_name = local.helm_chart_name

  # Chart version resolved to the pinned default when unset, so both engines
  # install the same chart whether or not the platform's defaulting
  # middleware ran — mirror of the Pulumi module's DefaultChartVersion.
  chart_version = coalesce(var.spec.chart_version, "2.8.0")

  namespace = var.spec.namespace

  # Controller ServiceAccount name — the module pins it to the chart name
  # (serviceAccount.name) so cloud-side ambient-identity bindings have a
  # deterministic subject.
  controller_service_account = local.helm_chart_name

  # Resource-identity labels stamped on the namespace this module creates
  # (never injected into the chart's own resources — Helm owns those).
  labels = merge(
    {
      "planton.ai/resource"      = "true"
      "planton.ai/resource-name" = var.metadata.name
      "planton.ai/resource-kind" = "KubernetesExternalSecretsOperator"
    },
    var.metadata.id != null && var.metadata.id != "" ? { "planton.ai/resource-id" = var.metadata.id } : {},
    var.metadata.org != null && var.metadata.org != "" ? { "planton.ai/organization" = var.metadata.org } : {},
    var.metadata.env != null && var.metadata.env != "" ? { "planton.ai/environment" = var.metadata.env } : {}
  )

  # ---- CRD lifecycle -----------------------------------------------------
  # installCRDs matches the chart's own default (true). keep_on_uninstall
  # has NO chart knob — the chart templates its CRDs and Helm would DELETE
  # them on uninstall, cascading to every ESO object cluster-wide. Planton
  # default keeps them via the standard Helm resource-policy annotation,
  # which the chart forwards onto the CRDs (crds.annotations).
  crds_install = try(var.spec.crds.install, null) != null ? var.spec.crds.install : true
  crds_keep    = try(var.spec.crds.keep_on_uninstall, null) != null ? var.spec.crds.keep_on_uninstall : true

  # ---- workload identity annotations -----------------------------------
  # The chart creates the controller ServiceAccount; the identity annotation
  # rides serviceAccount.annotations. AKS additionally needs the
  # azure.workload.identity/use pod label.
  workload_identity_annotations = merge(
    try(var.spec.workload_identity.gke, null) != null ? {
      "iam.gke.io/gcp-service-account" = var.spec.workload_identity.gke.service_account_email
    } : {},
    try(var.spec.workload_identity.eks, null) != null ? {
      "eks.amazonaws.com/role-arn" = var.spec.workload_identity.eks.role_arn
    } : {},
    try(var.spec.workload_identity.aks, null) != null ? merge(
      {
        "azure.workload.identity/client-id" = var.spec.workload_identity.aks.client_id
      },
      try(var.spec.workload_identity.aks.tenant_id, null) != null ? {
        "azure.workload.identity/tenant-id" = var.spec.workload_identity.aks.tenant_id
      } : {}
    ) : {}
  )

  # ---- per-component resources (shared ContainerResources shape) --------
  controller_resources = try(var.spec.resources, null) == null ? null : {
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

  webhook_resources = try(var.spec.webhook.resources, null) == null ? null : {
    for k, v in {
      limits = try(var.spec.webhook.resources.limits, null) == null ? null : {
        for lk, lv in {
          cpu    = try(var.spec.webhook.resources.limits.cpu, "") != "" ? var.spec.webhook.resources.limits.cpu : null
          memory = try(var.spec.webhook.resources.limits.memory, "") != "" ? var.spec.webhook.resources.limits.memory : null
        } : lk => lv if lv != null
      }
      requests = try(var.spec.webhook.resources.requests, null) == null ? null : {
        for rk, rv in {
          cpu    = try(var.spec.webhook.resources.requests.cpu, "") != "" ? var.spec.webhook.resources.requests.cpu : null
          memory = try(var.spec.webhook.resources.requests.memory, "") != "" ? var.spec.webhook.resources.requests.memory : null
        } : rk => rv if rv != null
      }
    } : k => v if v != null && length(v) > 0
  }

  cert_controller_resources = try(var.spec.cert_controller.resources, null) == null ? null : {
    for k, v in {
      limits = try(var.spec.cert_controller.resources.limits, null) == null ? null : {
        for lk, lv in {
          cpu    = try(var.spec.cert_controller.resources.limits.cpu, "") != "" ? var.spec.cert_controller.resources.limits.cpu : null
          memory = try(var.spec.cert_controller.resources.limits.memory, "") != "" ? var.spec.cert_controller.resources.limits.memory : null
        } : lk => lv if lv != null
      }
      requests = try(var.spec.cert_controller.resources.requests, null) == null ? null : {
        for rk, rv in {
          cpu    = try(var.spec.cert_controller.resources.requests.cpu, "") != "" ? var.spec.cert_controller.resources.requests.cpu : null
          memory = try(var.spec.cert_controller.resources.requests.memory, "") != "" ? var.spec.cert_controller.resources.requests.memory : null
        } : rk => rv if rv != null
      }
    } : k => v if v != null && length(v) > 0
  }

  # ---- typed chart values (twin of the Pulumi module's buildHelmValues) --
  helm_values = {
    for k, v in {
      # installCRDs is always rendered (the chart's own default is true);
      # the keep annotation only makes sense when this release owns the
      # CRDs, so it rides along only when install && keep.
      installCRDs = local.crds_install
      crds = local.crds_install && local.crds_keep ? {
        annotations = { "helm.sh/resource-policy" = "keep" }
      } : null

      # ---- controller ---------------------------------------------------
      replicaCount    = try(var.spec.replicas, null)
      leaderElect     = var.spec.leader_elect ? true : null
      concurrent      = try(var.spec.concurrent, null)
      controllerClass = var.spec.controller_class != "" ? var.spec.controller_class : null
      scopedNamespace = var.spec.scoped_namespace != "" ? var.spec.scoped_namespace : null
      scopedRBAC      = var.spec.scoped_rbac ? true : null
      log             = try(var.spec.log_level, null) != null ? { level = var.spec.log_level } : null
      resources       = local.controller_resources

      # ---- workload identity --------------------------------------------
      # The chart creates the controller ServiceAccount; the module pins its
      # name (deterministic identity subject) and rides ambient-identity
      # annotations on it. Per-store identities (store auth blocks) need
      # nothing here.
      serviceAccount = {
        for sk, sv in {
          name        = local.controller_service_account
          annotations = length(local.workload_identity_annotations) > 0 ? local.workload_identity_annotations : null
        } : sk => sv if sv != null
      }
      # The azure-workload-identity webhook only injects the federated token
      # volume into pods carrying this label.
      podLabels = try(var.spec.workload_identity.aks, null) != null ? {
        "azure.workload.identity/use" = "true"
      } : null

      # ---- scheduling -----------------------------------------------------
      nodeSelector = length(var.spec.node_selector) > 0 ? var.spec.node_selector : null

      tolerations = length(var.spec.tolerations) > 0 ? [
        for t in var.spec.tolerations : {
          for tk, tv in {
            key               = t.key != "" ? t.key : null
            operator          = t.operator != "" ? t.operator : null
            value             = t.value != "" ? t.value : null
            effect            = t.effect != "" ? t.effect : null
            tolerationSeconds = try(t.toleration_seconds, null)
          } : tk => tv if tv != null
        }
      ] : null

      priorityClassName   = var.spec.priority_class_name != "" ? var.spec.priority_class_name : null
      podDisruptionBudget = var.spec.pod_disruption_budget ? { enabled = true, minAvailable = 1 } : null

      # ---- observability --------------------------------------------------
      serviceMonitor = try(var.spec.prometheus.service_monitor, false) ? {
        for sk, sv in {
          enabled          = true
          interval         = try(var.spec.prometheus.service_monitor_interval, null)
          additionalLabels = length(try(var.spec.prometheus.service_monitor_labels, {})) > 0 ? var.spec.prometheus.service_monitor_labels : null
        } : sk => sv if sv != null
      } : null

      # ---- webhook ---------------------------------------------------------
      # create is rendered ONLY when enabled is explicitly false — the chart
      # default (true) stays untouched otherwise, mirroring the Pulumi
      # module's `Enabled != nil && !GetEnabled()` guard.
      webhook = try(var.spec.webhook, null) == null ? null : {
        for wk, wv in {
          create       = try(var.spec.webhook.enabled, null) == false ? false : null
          replicaCount = try(var.spec.webhook.replicas, null)
          resources    = local.webhook_resources
        } : wk => wv if wv != null
      }

      # ---- cert-controller ---------------------------------------------------
      certController = try(var.spec.cert_controller, null) == null ? null : {
        for ck, cv in {
          create       = try(var.spec.cert_controller.enabled, null) == false ? false : null
          replicaCount = try(var.spec.cert_controller.replicas, null)
          resources    = local.cert_controller_resources
        } : ck => cv if cv != null
      }

      # ---- image ----------------------------------------------------------------
      image = var.spec.image_repository != "" ? { repository = var.spec.image_repository } : null
    } : k => v if v != null && v != {}
  }
}
