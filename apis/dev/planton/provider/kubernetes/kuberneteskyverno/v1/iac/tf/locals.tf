# Computed values for the KubernetesKyverno module.
# Every resolution here has an exact twin in the Pulumi module — keep
# them in lockstep: same rendered chart values, same outputs.
#
# HCL DISCIPLINE: conditional keys are contributed via merge() of
# `cond ? { key = value } : {}` singleton maps — a ternary whose branches
# are differently-shaped objects fails plan-time type unification.
# Optional nested blocks are read with try(): HCL's && does NOT
# short-circuit.

locals {
  # Pinned chart identity; chart_version resolves to the pinned default
  # when unset so both engines install the same chart whether or not the
  # platform's defaulting middleware ran. The chart's appVersion pins
  # every controller image (chart 3.8.2 = Kyverno v1.18.2).
  helm_chart_name       = "kyverno"
  helm_chart_repo       = "https://kyverno.github.io/kyverno"
  default_chart_version = "3.8.2"
  chart_version         = try(var.spec.chart_version, "") != "" ? var.spec.chart_version : local.default_chart_version

  namespace    = var.spec.namespace
  release_name = var.metadata.name

  # Chart-derived name outputs (fullnameOverride pins the fullname to the
  # release name): the admission webhook Service is "<fullname>-svc" and
  # the runtime ConfigMap is the fullname itself — chart-truth from the
  # admission-controller serviceName and config configMapName helpers.
  admission_service_name = "${var.metadata.name}-svc"
  config_map_name        = var.metadata.name

  # Planton governance labels for the module-created namespace (never
  # injected into the chart's own resources; Helm owns those).
  labels = merge(
    {
      "planton.ai/resource"      = "true"
      "planton.ai/resource-name" = var.metadata.name
      "planton.ai/resource-kind" = "KubernetesKyverno"
    },
    try(var.metadata.id, "") != "" ? { "planton.ai/resource-id" = var.metadata.id } : {},
    try(var.metadata.org, "") != "" ? { "planton.ai/organization" = var.metadata.org } : {},
    try(var.metadata.env, "") != "" ? { "planton.ai/environment" = var.metadata.env } : {}
  )

  # ---- CRDs (the crds subchart) -------------------------------------------------
  # The CRDs are chart-TEMPLATED (no keep by default): they install and
  # DELETE with the release, cascade-deleting every policy on the
  # cluster. keep_on_uninstall injects the resource-policy annotation.
  # install and the migration hook are rendered EXPLICITLY whenever the
  # block is declared (the chart defaults both true; explicit rendering
  # keeps the manifest's intent visible in the diff).
  # Single null-pruned comprehension — never `cond ? {} : {for…}` or a
  # merge() of differently-shaped ternary branches (both fail HCL type
  # unification; this module's own plan gate caught them here).
  crds_block = {
    for k, v in {
      install     = try(var.spec.crds, null) == null ? null : try(var.spec.crds.install, true)
      migration   = try(var.spec.crds, null) == null ? null : { enabled = try(var.spec.crds.migration_enabled, true) }
      annotations = try(var.spec.crds.keep_on_uninstall, false) ? { "helm.sh/resource-policy" = "keep" } : null
    } : k => v if v != null
  }

  # ---- runtime config (resource filters, webhook selectors) -----------------------
  # Declaring webhook_exclude_namespaces REPLACES the chart's webhooks
  # value, so the chart's own kube-system exclusion is re-included here
  # by construction — dropping it would put the control plane on the
  # policy path.
  webhook_exclude_namespaces = distinct(concat(
    ["kube-system"],
    try(var.spec.config.webhook_exclude_namespaces, [])
  ))

  config_block = {
    for k, v in {
      webhooks = length(try(var.spec.config.webhook_exclude_namespaces, [])) > 0 ? {
        namespaceSelector = {
          matchExpressions = [{
            key      = "kubernetes.io/metadata.name"
            operator = "NotIn"
            values   = local.webhook_exclude_namespaces
          }]
        }
      } : null
      resourceFiltersInclude        = length(try(var.spec.config.resource_filters_include, [])) > 0 ? var.spec.config.resource_filters_include : null
      resourceFiltersExclude        = length(try(var.spec.config.resource_filters_exclude, [])) > 0 ? var.spec.config.resource_filters_exclude : null
      excludeGroups                 = length(try(var.spec.config.exclude_groups, [])) > 0 ? var.spec.config.exclude_groups : null
      excludeUsernames              = length(try(var.spec.config.exclude_usernames, [])) > 0 ? var.spec.config.exclude_usernames : null
      defaultRegistry               = try(var.spec.config.default_registry, "") != "" ? var.spec.config.default_registry : null
      enableDefaultRegistryMutation = try(var.spec.config.enable_default_registry_mutation, null)
    } : k => v if v != null
  }

  # ---- feature flags ----------------------------------------------------------------
  # Every chart feature is the nested {enabled: bool} shape (chart-truth
  # at the pin) — a bare bool would break template rendering.
  background_scan_block = try(var.spec.features.background_scan, null) == null ? null : {
    for k, v in {
      enabled                = try(var.spec.features.background_scan.enabled, true)
      backgroundScanWorkers  = try(var.spec.features.background_scan.workers, null)
      backgroundScanInterval = try(var.spec.features.background_scan.interval, "") != "" ? var.spec.features.background_scan.interval : null
    } : k => v if v != null
  }

  logging_block = try(var.spec.features, null) == null ? null : {
    for k, v in {
      format    = try(var.spec.features.logging_format, "") != "" ? var.spec.features.logging_format : null
      verbosity = try(var.spec.features.logging_verbosity, null)
    } : k => v if v != null
  }

  features_block = {
    for k, v in {
      forceFailurePolicyIgnore          = try(var.spec.features.force_failure_policy_ignore, false) ? { enabled = true } : null
      backgroundScan                    = local.background_scan_block
      generateValidatingAdmissionPolicy = try(var.spec.features.generate_validating_admission_policy, null) == null ? null : { enabled = var.spec.features.generate_validating_admission_policy }
      admissionReports                  = try(var.spec.features.admission_reports, null) == null ? null : { enabled = var.spec.features.admission_reports }
      aggregateReports                  = try(var.spec.features.aggregate_reports, null) == null ? null : { enabled = var.spec.features.aggregate_reports }
      policyReports                     = try(var.spec.features.policy_reports, null) == null ? null : { enabled = var.spec.features.policy_reports }
      logging                           = local.logging_block != null && length(coalesce(local.logging_block, {})) > 0 ? local.logging_block : null
      omitEvents                        = length(try(var.spec.features.omit_event_types, [])) > 0 ? { eventTypes = var.spec.features.omit_event_types } : null
    } : k => v if v != null
  }

  # ---- shared renderers ----------------------------------------------------------------
  # One resources renderer per controller (identical shape); tolerations
  # rendered with the null-prune idiom.
  admission_resources_block = {
    for k, v in {
      requests = try(var.spec.admission_controller.resources.requests, null) == null ? null : {
        for rk, rv in {
          cpu    = var.spec.admission_controller.resources.requests.cpu
          memory = var.spec.admission_controller.resources.requests.memory
        } : rk => rv if rv != null && rv != ""
      }
      limits = try(var.spec.admission_controller.resources.limits, null) == null ? null : {
        for rk, rv in {
          cpu    = var.spec.admission_controller.resources.limits.cpu
          memory = var.spec.admission_controller.resources.limits.memory
        } : rk => rv if rv != null && rv != ""
      }
    } : k => v if v != null
  }

  # ---- admission controller ---------------------------------------------------------------
  # Container resources sit under `container.resources` (the admission
  # controller runs an init container + the main container); the other
  # three controllers take `resources` directly — chart-truth from the
  # values shapes at the pin.
  admission_autoscaling = try(var.spec.admission_controller.autoscaling, null) == null ? {} : {
    autoscaling = merge(
      {
        enabled     = true
        maxReplicas = var.spec.admission_controller.autoscaling.max_replicas
      },
      try(var.spec.admission_controller.autoscaling.min_replicas, null) != null ? { minReplicas = var.spec.admission_controller.autoscaling.min_replicas } : {},
      try(var.spec.admission_controller.autoscaling.target_cpu_utilization_percentage, null) != null ? { targetCPUUtilizationPercentage = var.spec.admission_controller.autoscaling.target_cpu_utilization_percentage } : {}
    )
  }

  admission_controller_block = merge(
    try(var.spec.admission_controller.replicas, null) != null ? { replicas = var.spec.admission_controller.replicas } : {},
    length(local.admission_resources_block) > 0 ? { container = { resources = local.admission_resources_block } } : {},
    local.scheduling_entries["admission_controller"],
    local.admission_autoscaling,
    local.cert_manager_block != null ? { certManager = local.cert_manager_block } : {},
    local.service_monitor_enabled ? { serviceMonitor = { enabled = true } } : {}
  )

  # ---- optional controllers (background / cleanup / reports) -------------------------------
  controller_specs = {
    background_controller = try(var.spec.background_controller, null)
    cleanup_controller    = try(var.spec.cleanup_controller, null)
    reports_controller    = try(var.spec.reports_controller, null)
  }

  controller_resources = {
    for name, controller in local.controller_specs :
    name => {
      for k, v in {
        requests = try(controller.resources.requests, null) == null ? null : {
          for rk, rv in {
            cpu    = controller.resources.requests.cpu
            memory = controller.resources.requests.memory
          } : rk => rv if rv != null && rv != ""
        }
        limits = try(controller.resources.limits, null) == null ? null : {
          for rk, rv in {
            cpu    = controller.resources.limits.cpu
            memory = controller.resources.limits.memory
          } : rk => rv if rv != null && rv != ""
        }
      } : k => v if v != null
    }
  }

  scheduling_entries = {
    for name, block in {
      admission_controller  = try(var.spec.admission_controller.scheduling, null)
      background_controller = try(var.spec.background_controller.scheduling, null)
      cleanup_controller    = try(var.spec.cleanup_controller.scheduling, null)
      reports_controller    = try(var.spec.reports_controller.scheduling, null)
    } :
    name => merge(
      block == null ? {} : (
        length(try(block.node_selector, {})) > 0 ? { nodeSelector = block.node_selector } : {}
      ),
      block == null ? {} : (
        length(try(block.tolerations, [])) > 0 ? {
          tolerations = [
            for t in block.tolerations : { for k, v in {
              key               = t.key != "" ? t.key : null
              operator          = t.operator != "" ? t.operator : null
              value             = t.value != "" ? t.value : null
              effect            = t.effect != "" ? t.effect : null
              tolerationSeconds = t.toleration_seconds
            } : k => v if v != null }
          ]
        } : {}
      )
    )
  }

  # Null-prune idiom per entry (a `controller == null ? {} : merge(...)`
  # ternary fails HCL type unification — this module's own plan gate
  # caught that shape here).
  optional_controller_blocks = {
    for name, controller in local.controller_specs :
    name => merge(
      {
        for k, v in {
          enabled   = controller == null ? null : try(controller.enabled, null)
          replicas  = controller == null ? null : try(controller.replicas, null)
          resources = length(local.controller_resources[name]) > 0 ? local.controller_resources[name] : null
        } : k => v if v != null
      },
      local.scheduling_entries[name],
      local.service_monitor_enabled ? { serviceMonitor = { enabled = true } } : {}
    )
  }

  # ---- webhook certificates (the cert-manager arm) -------------------------------------------
  # Omitted = Kyverno-managed runtime certificates with rotation (the
  # chart default; nothing to render). The arm applies to BOTH webhook
  # servers — admission AND cleanup controllers (each exposes its own
  # certManager block; missing one would leave that webhook on runtime
  # certs, a split-brain trust posture — the Pulumi twin fans out
  # identically).
  # Null-prune idiom (NOT a `cond ? {a,b} : {}` ternary — differently
  # shaped branches fail HCL type unification; this module's own plan
  # gate caught exactly that shape here).
  cert_manager_block = try(var.spec.certificates.cert_manager, null) == null ? null : {
    for k, v in {
      enabled                = true
      createSelfSignedIssuer = try(var.spec.certificates.cert_manager.issuer_name, "") != "" ? false : null
      issuerRef = try(var.spec.certificates.cert_manager.issuer_name, "") == "" ? null : {
        name  = var.spec.certificates.cert_manager.issuer_name
        kind  = try(var.spec.certificates.cert_manager.issuer_kind, "") != "" ? var.spec.certificates.cert_manager.issuer_kind : "ClusterIssuer"
        group = "cert-manager.io"
      }
    } : k => v if v != null
  }

  # ---- metrics: ServiceMonitor fan-out ---------------------------------------------------------
  # All four controllers expose the toggle — enabling only some would
  # silently blind the others (the fan-out completeness lesson).
  service_monitor_enabled = try(var.spec.metrics.service_monitor, false)

  cleanup_controller_block = merge(
    local.optional_controller_blocks["cleanup_controller"],
    local.cert_manager_block != null ? { certManager = local.cert_manager_block } : {}
  )

  # ---- typed chart values (Pulumi twin: buildHelmValues) ----------------------------------------
  # fullnameOverride pins the chart fullname to the resource name; the
  # webhooksCleanup pre-delete hook is rendered EXPLICITLY in both states
  # (it is the designed uninstall path for the runtime-registered webhook
  # configurations). existingImagePullSecrets takes NAMES of pre-existing
  # secrets; the chart's imagePullSecrets map (which CREATES secrets from
  # credentials) is deliberately not modeled.
  typed_helm_values = merge(
    {
      fullnameOverride = local.release_name
      webhooksCleanup = {
        enabled = try(var.spec.webhooks_cleanup_enabled, null) == null ? true : var.spec.webhooks_cleanup_enabled
      }
    },
    length(local.crds_block) > 0 ? { crds = local.crds_block } : {},
    length(local.config_block) > 0 ? { config = local.config_block } : {},
    length(local.features_block) > 0 ? { features = local.features_block } : {},
    length(local.admission_controller_block) > 0 ? { admissionController = local.admission_controller_block } : {},
    length(local.optional_controller_blocks["background_controller"]) > 0 ? { backgroundController = local.optional_controller_blocks["background_controller"] } : {},
    length(local.cleanup_controller_block) > 0 ? { cleanupController = local.cleanup_controller_block } : {},
    length(local.optional_controller_blocks["reports_controller"]) > 0 ? { reportsController = local.optional_controller_blocks["reports_controller"] } : {},
    try(var.spec.image_registry, "") != "" ? { global = { image = { registry = var.spec.image_registry } } } : {},
    length(try(var.spec.image_pull_secrets, [])) > 0 ? { existingImagePullSecrets = var.spec.image_pull_secrets } : {}
  )
}
