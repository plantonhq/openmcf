# Computed values for the KubernetesGhaRunnerScaleSetController module.
# Every resolution here has an exact twin in the Pulumi module — keep
# them in lockstep: same rendered chart values, same outputs.
#
# HCL DISCIPLINE: conditional keys are contributed via merge() of
# `cond ? { key = value } : {}` singleton maps — a ternary whose branches
# are differently-shaped objects fails plan-time type unification.
# Optional nested blocks are read with try(): HCL's && does NOT
# short-circuit.

locals {
  # Pinned OCI registry path and chart name; chart_version resolves to
  # the pinned default when unset so both engines install the same chart
  # whether or not the platform's defaulting middleware ran. Chart and
  # controller image move in lockstep (0.14.2 = appVersion 0.14.2).
  helm_oci_repo         = "oci://ghcr.io/actions/actions-runner-controller-charts"
  helm_chart_name       = "gha-runner-scale-set-controller"
  default_chart_version = "0.14.2"
  chart_version         = try(var.spec.chart_version, "") != "" ? var.spec.chart_version : local.default_chart_version

  namespace    = var.spec.namespace
  release_name = var.metadata.name

  # The controller's ServiceAccount equals metadata.name:
  # fullnameOverride pins the chart fullname to the release name, and the
  # chart names the created ServiceAccount exactly the fullname.
  service_account_name = var.metadata.name

  # Planton governance labels for the module-created namespace (never
  # injected into the chart's own resources; Helm owns those).
  labels = merge(
    {
      "planton.ai/resource"      = "true"
      "planton.ai/resource-name" = var.metadata.name
      "planton.ai/resource-kind" = "KubernetesGhaRunnerScaleSetController"
    },
    try(var.metadata.id, "") != "" ? { "planton.ai/resource-id" = var.metadata.id } : {},
    try(var.metadata.org, "") != "" ? { "planton.ai/organization" = var.metadata.org } : {},
    try(var.metadata.env, "") != "" ? { "planton.ai/environment" = var.metadata.env } : {}
  )

  # ---- image override (air-gap) ----------------------------------------------------
  # The chart takes the image reference COMBINED (image.repository holds
  # the full mirror path; no separate registry value — verified in the
  # chart's values + deployment template at the pin). Every conditional
  # block is ONE object literal of try(..., null) values pruned with the
  # for-comprehension — NO outer ternary: `cond ? {} : {for…}` fails HCL
  # type unification (the empty object cannot gain the typed branch's
  # attributes; the class failed the KubernetesTekton plan gate).
  image_block = {
    for k, v in {
      repository = try(var.spec.image.repo, "") != "" ? var.spec.image.repo : null
      tag        = try(var.spec.image.tag, "") != "" ? var.spec.image.tag : null
    } : k => v if v != null
  }

  # ---- resources -------------------------------------------------------------------
  resources_block = {
    for k, v in {
      requests = try(var.spec.resources.requests, null) == null ? null : {
        for rk, rv in {
          cpu    = var.spec.resources.requests.cpu
          memory = var.spec.resources.requests.memory
        } : rk => rv if rv != null && rv != ""
      }
      limits = try(var.spec.resources.limits, null) == null ? null : {
        for rk, rv in {
          cpu    = var.spec.resources.limits.cpu
          memory = var.spec.resources.limits.memory
        } : rk => rv if rv != null && rv != ""
      }
    } : k => v if v != null
  }

  # ---- image pull secrets (joined + deduplicated, Pulumi twin: pullSecretNames) -----
  image_pull_secret_names = distinct(concat(
    try(var.spec.image_pull_secrets, []),
    try(var.spec.image.pull_secret_name, "") != "" ? [var.spec.image.pull_secret_name] : []
  ))

  # ---- flags ------------------------------------------------------------------------
  # rateLimiter is a STRUCTURED chart block ({name}) — the template reads
  # `.Values.flags.rateLimiter.name`; a bare string would break rendering.
  flags_block = {
    for k, v in {
      logLevel                        = try(var.spec.flags.log_level, "") != "" ? var.spec.flags.log_level : null
      logFormat                       = try(var.spec.flags.log_format, "") != "" ? var.spec.flags.log_format : null
      watchSingleNamespace            = try(var.spec.flags.watch_single_namespace, "") != "" ? var.spec.flags.watch_single_namespace : null
      runnerMaxConcurrentReconciles   = try(var.spec.flags.runner_max_concurrent_reconciles, null)
      updateStrategy                  = try(var.spec.flags.update_strategy, "") != "" ? var.spec.flags.update_strategy : null
      excludeLabelPropagationPrefixes = length(try(var.spec.flags.exclude_label_propagation_prefixes, [])) > 0 ? var.spec.flags.exclude_label_propagation_prefixes : null
      k8sClientRateLimiterQPS         = try(var.spec.flags.k8s_client_rate_limiter_qps, null)
      k8sClientRateLimiterBurst       = try(var.spec.flags.k8s_client_rate_limiter_burst, null)
      rateLimiter                     = try(var.spec.flags.rate_limiter, "") != "" ? { name = var.spec.flags.rate_limiter } : null
      healthProbeBindAddress          = try(var.spec.flags.health_probe_bind_address, "") != "" ? var.spec.flags.health_probe_bind_address : null
    } : k => v if v != null
  }

  # ---- scheduling --------------------------------------------------------------------
  scheduling_tolerations = [
    for t in try(var.spec.scheduling.tolerations, []) : { for k, v in {
      key               = t.key != "" ? t.key : null
      operator          = t.operator != "" ? t.operator : null
      value             = t.value != "" ? t.value : null
      effect            = t.effect != "" ? t.effect : null
      tolerationSeconds = t.toleration_seconds
    } : k => v if v != null }
  ]

  # ---- typed chart values (Pulumi twin: buildHelmValues) ------------------------------
  # fullnameOverride pins the chart fullname (and the ServiceAccount
  # name — the scale-set discovery handle) to the resource name.
  # Declaring metrics ENABLES metrics; absent = the chart passes empty
  # flags (disabled).
  typed_helm_values = merge(
    {
      fullnameOverride = local.release_name
      replicaCount     = try(var.spec.replicas, null) != null ? var.spec.replicas : 1
    },
    length(local.image_block) > 0 ? { image = local.image_block } : {},
    length(local.resources_block) > 0 ? { resources = local.resources_block } : {},
    length(local.image_pull_secret_names) > 0 ? {
      imagePullSecrets = [for n in local.image_pull_secret_names : { name = n }]
    } : {},
    length(local.flags_block) > 0 ? { flags = local.flags_block } : {},
    try(var.spec.flags.priority_class_name, "") != "" ? { priorityClassName = var.spec.flags.priority_class_name } : {},
    try(var.spec.metrics, null) != null ? {
      metrics = {
        controllerManagerAddr = var.spec.metrics.controller_manager_addr
        listenerAddr          = var.spec.metrics.listener_addr
        listenerEndpoint      = var.spec.metrics.listener_endpoint
      }
    } : {},
    length(try(var.spec.scheduling.node_selector, {})) > 0 ? { nodeSelector = var.spec.scheduling.node_selector } : {},
    length(local.scheduling_tolerations) > 0 ? { tolerations = local.scheduling_tolerations } : {}
  )
}
