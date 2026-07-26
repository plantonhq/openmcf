# Computed values for the KubernetesGhaRunnerScaleSet module. Every
# resolution here has an exact twin in the Pulumi module — keep them in
# lockstep: same rendered chart values, same materialized Secret, same
# outputs.
#
# HCL DISCIPLINE: conditional keys are contributed via merge() of
# `cond ? { key = value } : {}` singleton maps — a ternary whose branches
# are differently-shaped objects fails plan-time type unification.
# Optional nested blocks are read with try(): HCL's && does NOT
# short-circuit.

locals {
  # Pinned OCI registry path and chart name; chart_version resolves to
  # the pinned default when unset. Keep it EQUAL to the controller
  # kind's — GitHub supports controller and scale set charts only at
  # matching versions.
  helm_oci_repo         = "oci://ghcr.io/actions/actions-runner-controller-charts"
  helm_chart_name       = "gha-runner-scale-set"
  default_chart_version = "0.14.2"
  chart_version         = try(var.spec.chart_version, "") != "" ? var.spec.chart_version : local.default_chart_version

  namespace    = var.spec.namespace
  release_name = var.metadata.name

  # The GitHub-visible fleet name — the exact `runs-on:` value.
  # spec.runner_scale_set_name, falling back to metadata.name. The chart
  # fails installs past 45 characters (main.tf precondition catches it
  # at plan).
  runner_scale_set_name = try(var.spec.runner_scale_set_name, "") != "" ? var.spec.runner_scale_set_name : var.metadata.name

  # The credential Secret the chart reads BY NAME: the user's own Secret
  # (existing-Secret arm) or the module-materialized
  # `<name>-github-auth` (declared arms) — credential material never
  # rides rendered chart values (Pulumi twin: githubAuthSecret).
  materialize_auth_secret = try(var.spec.auth.existing_secret_name, "") == ""
  github_auth_secret_name = local.materialize_auth_secret ? "${var.metadata.name}-github-auth" : var.spec.auth.existing_secret_name

  # The materialized Secret's data — the chart's pre-defined-secret key
  # contract (github_token for a PAT; github_app_* for an App). Exactly
  # one arm is populated (spec CEL).
  github_auth_secret_data = merge(
    try(var.spec.auth.pat, null) != null ? {
      github_token = var.spec.auth.pat.token
    } : {},
    try(var.spec.auth.github_app, null) != null ? {
      github_app_id              = var.spec.auth.github_app.app_id
      github_app_installation_id = var.spec.auth.github_app.installation_id
      github_app_private_key     = var.spec.auth.github_app.private_key
    } : {}
  )

  # Planton governance labels for the module-created satellites (the
  # namespace, the materialized Secret — never injected into the chart's
  # own resources; Helm owns those).
  labels = merge(
    {
      "planton.ai/resource"      = "true"
      "planton.ai/resource-name" = var.metadata.name
      "planton.ai/resource-kind" = "KubernetesGhaRunnerScaleSet"
    },
    try(var.metadata.id, "") != "" ? { "planton.ai/resource-id" = var.metadata.id } : {},
    try(var.metadata.org, "") != "" ? { "planton.ai/organization" = var.metadata.org } : {},
    try(var.metadata.env, "") != "" ? { "planton.ai/environment" = var.metadata.env } : {}
  )

  # ---- container mode ---------------------------------------------------------------
  container_mode_block = try(var.spec.container_mode, null) == null ? null : merge(
    { type = var.spec.container_mode.mode },
    try(var.spec.container_mode.kubernetes_work_volume, null) != null ? {
      kubernetesModeWorkVolumeClaim = {
        accessModes      = ["ReadWriteOnce"]
        storageClassName = var.spec.container_mode.kubernetes_work_volume.storage_class
        resources = {
          requests = {
            storage = var.spec.container_mode.kubernetes_work_volume.size
          }
        }
      }
    } : {}
  )

  # ---- the runner container -----------------------------------------------------------
  # Rendered ONLY when customized: Helm values LISTS replace (never
  # merge), so any override must re-state the chart's own container
  # contract — name `runner` (the chart applies its mode wiring to the
  # container with exactly that name) and the run.sh command.
  runner_customized = try(var.spec.runner, null) != null && (try(var.spec.runner.image, "") != "" || try(var.spec.runner.resources, null) != null)

  runner_resources = try(var.spec.runner.resources, null) == null ? null : {
    for k, v in {
      requests = try(var.spec.runner.resources.requests, null) == null ? null : {
        for rk, rv in {
          cpu    = var.spec.runner.resources.requests.cpu
          memory = var.spec.runner.resources.requests.memory
        } : rk => rv if rv != null && rv != ""
      }
      limits = try(var.spec.runner.resources.limits, null) == null ? null : {
        for rk, rv in {
          cpu    = var.spec.runner.resources.limits.cpu
          memory = var.spec.runner.resources.limits.memory
        } : rk => rv if rv != null && rv != ""
      }
    } : k => v if v != null
  }

  runner_template_block = !local.runner_customized ? null : {
    spec = {
      containers = [
        merge(
          {
            name    = "runner"
            image   = try(var.spec.runner.image, "") != "" ? var.spec.runner.image : "ghcr.io/actions/actions-runner:latest"
            command = ["/home/runner/run.sh"]
          },
          local.runner_resources != null ? { resources = local.runner_resources } : {}
        )
      ]
    }
  }

  # ---- proxy ------------------------------------------------------------------------------
  proxy_block = try(var.spec.proxy, null) == null ? null : merge(
    try(var.spec.proxy.http, null) != null ? {
      http = merge(
        { url = var.spec.proxy.http.url },
        try(var.spec.proxy.http.credential_secret_name, "") != "" ? { credentialSecretRef = var.spec.proxy.http.credential_secret_name } : {}
      )
    } : {},
    try(var.spec.proxy.https, null) != null ? {
      https = merge(
        { url = var.spec.proxy.https.url },
        try(var.spec.proxy.https.credential_secret_name, "") != "" ? { credentialSecretRef = var.spec.proxy.https.credential_secret_name } : {}
      )
    } : {},
    length(try(var.spec.proxy.no_proxy, [])) > 0 ? { noProxy = var.spec.proxy.no_proxy } : {}
  )

  # ---- GitHub Enterprise Server private CA -------------------------------------------------
  github_server_tls_block = try(var.spec.github_server_tls, null) == null ? null : merge(
    {
      certificateFrom = {
        configMapKeyRef = {
          name = var.spec.github_server_tls.config_map_name
          key  = try(var.spec.github_server_tls.key, "") != "" ? var.spec.github_server_tls.key : "ca.crt"
        }
      }
    },
    try(var.spec.github_server_tls.runner_mount_path, "") != "" ? { runnerMountPath = var.spec.github_server_tls.runner_mount_path } : {}
  )

  # ---- typed chart values (Pulumi twin: buildHelmValues) -------------------------------------
  # githubConfigSecret always renders as a Secret NAME (the chart's
  # pre-defined-secret form) — the secret discipline. The runner scale
  # set name is rendered explicitly (never left to the release-name
  # default) so the exported runs-on handle and the rendered chart agree
  # by construction.
  typed_helm_values = merge(
    {
      githubConfigUrl    = var.spec.github_config_url
      githubConfigSecret = local.github_auth_secret_name
      runnerScaleSetName = local.runner_scale_set_name
    },
    try(var.spec.runner_group, "") != "" ? { runnerGroup = var.spec.runner_group } : {},
    try(var.spec.min_runners, null) != null ? { minRunners = var.spec.min_runners } : {},
    try(var.spec.max_runners, null) != null ? { maxRunners = var.spec.max_runners } : {},
    local.container_mode_block != null ? { containerMode = local.container_mode_block } : {},
    local.runner_template_block != null ? { template = local.runner_template_block } : {},
    local.proxy_block != null ? { proxy = local.proxy_block } : {},
    local.github_server_tls_block != null ? { githubServerTLS = local.github_server_tls_block } : {},
    try(var.spec.controller_service_account, null) != null ? {
      controllerServiceAccount = {
        namespace = var.spec.controller_service_account.namespace
        name      = var.spec.controller_service_account.name
      }
    } : {}
  )
}
