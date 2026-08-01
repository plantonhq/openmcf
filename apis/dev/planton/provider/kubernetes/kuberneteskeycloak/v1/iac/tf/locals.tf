# Computed values for the KubernetesKeycloak module. Every resolution
# here has an exact twin in the Pulumi module (locals.go /
# keycloak_cr.go) — keep them in lockstep: same rendered CR body, same
# outputs.
#
# HCL DISCIPLINE: conditional keys are contributed via merge() of
# `cond ? { key = value } : {}` singleton maps, or one object literal
# pruned with `{ for k, v in {...} : k => v if v != null }` — a ternary
# whose branches are differently-shaped objects fails plan-time type
# unification. Optional nested blocks are read with try(): HCL's && does
# NOT short-circuit. Optional scalars in string templates resolve with
# try(coalesce(...)).

locals {
  # metadata.name is the CR name — the operator's naming root: the
  # StatefulSet is the name itself, the main Service takes `-service`,
  # the headless JGroups discovery Service `-discovery`, and the
  # generated one-time bootstrap-admin Secret `-initial-admin`.
  keycloak_name = var.metadata.name
  namespace     = var.spec.namespace
  api_version   = "k8s.keycloak.org/v2beta1"

  # Planton governance labels on the module-created objects.
  labels = merge(
    {
      "planton.ai/resource"      = "true"
      "planton.ai/resource-name" = var.metadata.name
      "planton.ai/resource-kind" = "KubernetesKeycloak"
    },
    try(var.metadata.id, "") != "" ? { "planton.ai/resource-id" = var.metadata.id } : {},
    try(var.metadata.org, "") != "" ? { "planton.ai/organization" = var.metadata.org } : {},
    try(var.metadata.env, "") != "" ? { "planton.ai/environment" = var.metadata.env } : {}
  )

  # ---- resolved listener posture --------------------------------------------------
  # tls_secret_name arrives pre-resolved (literal or KubernetesCertificate
  # reference). Empty means the HTTPS listener is off and the server runs
  # plain HTTP — the spec's TLS-or-HTTP validation guarantees http_enabled
  # then (the server otherwise refuses to start; upstream surfaces that
  # only as a CrashLoopBackOff).
  tls_secret_name = try(var.spec.http.tls_secret_name, "") != "" ? var.spec.http.tls_secret_name : ""

  http_port       = try(coalesce(var.spec.http.http_port, 8080), 8080)
  https_port      = try(coalesce(var.spec.http.https_port, 8443), 8443)
  management_port = try(coalesce(var.spec.http_management_port, 9000), 9000)

  api_scheme = local.tls_secret_name != "" ? "https" : "http"
  api_port   = local.tls_secret_name != "" ? local.https_port : local.http_port

  # ---- spec.db ---------------------------------------------------------------------
  # The dev-file/dev-mem SANDBOX vendors run Keycloak's embedded H2 on
  # each pod's own ephemeral storage — no connection details apply, so
  # ONLY the vendor renders; everything else would be dead configuration
  # the CR should not carry.
  db_is_sandbox = contains(["dev-file", "dev-mem"], var.spec.db.vendor)

  db_username_secret = {
    for k, v in {
      name = try(var.spec.db.username_secret.name, null)
      key  = try(var.spec.db.username_secret.key, null)
    } : k => v if v != null
  }

  db_password_secret = {
    for k, v in {
      name = try(var.spec.db.password_secret.name, null)
      key  = try(var.spec.db.password_secret.key, null)
    } : k => v if v != null
  }

  # The CRD calls the JDBC URL override `url`; when set the server ignores
  # host/port/database.
  db_connection = {
    for k, v in {
      host        = !local.db_is_sandbox && try(var.spec.db.host, "") != "" ? var.spec.db.host : null
      port        = local.db_is_sandbox ? null : try(var.spec.db.port, null)
      database    = !local.db_is_sandbox && try(var.spec.db.database, "") != "" ? var.spec.db.database : null
      schema      = !local.db_is_sandbox && try(var.spec.db.schema, "") != "" ? var.spec.db.schema : null
      url         = !local.db_is_sandbox && try(var.spec.db.jdbc_url, "") != "" ? var.spec.db.jdbc_url : null
      poolMinSize = local.db_is_sandbox ? null : try(var.spec.db.pool_min_size, null)
      poolMaxSize = local.db_is_sandbox ? null : try(var.spec.db.pool_max_size, null)
    } : k => v if v != null
  }

  db_block = merge(
    { vendor = var.spec.db.vendor },
    local.db_connection,
    !local.db_is_sandbox && length(local.db_username_secret) > 0 ? { usernameSecret = local.db_username_secret } : {},
    !local.db_is_sandbox && length(local.db_password_secret) > 0 ? { passwordSecret = local.db_password_secret } : {}
  )

  # ---- spec.http / spec.hostname ---------------------------------------------------
  http_block = {
    for k, v in {
      tlsSecret   = local.tls_secret_name != "" ? local.tls_secret_name : null
      httpEnabled = try(var.spec.http.http_enabled, false) ? true : null
      httpPort    = try(var.spec.http.http_port, null)
      httpsPort   = try(var.spec.http.https_port, null)
    } : k => v if v != null
  }

  # strict is tri-state: unset leaves the server default (true); declared
  # renders either way — `strict: false` is the meaningful behind-a-proxy
  # posture.
  # backchannelDynamic and admin render verbatim; both are legal ONLY
  # with a FULL-URL hostname — server startup rules the spec's CEL
  # enforces (verified live: any other pairing crash-loops the pod
  # with no apply-time error from the operator).
  hostname_block = {
    for k, v in {
      hostname           = try(var.spec.hostname.hostname, "") != "" ? var.spec.hostname.hostname : null
      admin              = try(var.spec.hostname.admin, "") != "" ? var.spec.hostname.admin : null
      strict             = try(var.spec.hostname.strict, null)
      backchannelDynamic = try(var.spec.hostname.backchannel_dynamic, false) ? true : null
    } : k => v if v != null
  }

  # ---- spec.features / spec.cache / spec.truststores / additionalOptions ------------
  features_block = {
    for k, v in {
      enabled  = length(try(var.spec.features.enabled, [])) > 0 ? var.spec.features.enabled : null
      disabled = length(try(var.spec.features.disabled, [])) > 0 ? var.spec.features.disabled : null
    } : k => v if v != null
  }

  cache_config_map_file = {
    for k, v in {
      name = try(var.spec.cache_config.config_map_name, null)
      key  = try(var.spec.cache_config.key, null)
    } : k => v if v != null
  }

  # truststores is a CRD map keyed by an arbitrary handle — the module
  # keys each entry by its own Secret name.
  truststores = {
    for secret_name in try(var.spec.truststore_secret_names, []) :
    secret_name => { secret = { name = secret_name } }
  }

  # The spec validates value XOR secret — exactly one arm renders.
  additional_options = [
    for o in try(var.spec.additional_options, []) : merge(
      { name = o.name },
      try(o.value, "") != "" ? { value = o.value } : {},
      try(o.secret, null) != null ? { secret = { name = try(o.secret.name, ""), key = try(o.secret.key, "") } } : {}
    )
  ]

  # ---- spec.resources ----------------------------------------------------------------
  resources_requests = {
    for k, v in {
      cpu    = try(var.spec.resources.requests.cpu, "") != "" ? var.spec.resources.requests.cpu : null
      memory = try(var.spec.resources.requests.memory, "") != "" ? var.spec.resources.requests.memory : null
    } : k => v if v != null
  }

  resources_limits = {
    for k, v in {
      cpu    = try(var.spec.resources.limits.cpu, "") != "" ? var.spec.resources.limits.cpu : null
      memory = try(var.spec.resources.limits.memory, "") != "" ? var.spec.resources.limits.memory : null
    } : k => v if v != null
  }

  resources_block = merge(
    length(local.resources_requests) > 0 ? { requests = local.resources_requests } : {},
    length(local.resources_limits) > 0 ? { limits = local.resources_limits } : {}
  )

  # ---- spec.scheduling ----------------------------------------------------------------
  # The Keycloak CR models affinity rather than nodeSelector, so
  # node_selector translates to REQUIRED node affinity — one In-expression
  # per label, the catalog's established translation. Keys sort so both
  # engines render the same expression order.
  node_affinity_match_expressions = [
    for label_key in sort(keys(try(var.spec.scheduling.node_selector, {}))) : {
      key      = label_key
      operator = "In"
      values   = [var.spec.scheduling.node_selector[label_key]]
    }
  ]

  scheduling_tolerations = [
    for t in try(var.spec.scheduling.tolerations, []) : { for k, v in {
      key               = t.key != "" ? t.key : null
      operator          = t.operator != "" ? t.operator : null
      value             = t.value != "" ? t.value : null
      effect            = t.effect != "" ? t.effect : null
      tolerationSeconds = t.toleration_seconds
    } : k => v if v != null }
  ]

  scheduling_block = merge(
    length(local.node_affinity_match_expressions) > 0 ? {
      affinity = {
        nodeAffinity = {
          requiredDuringSchedulingIgnoredDuringExecution = {
            nodeSelectorTerms = [{ matchExpressions = local.node_affinity_match_expressions }]
          }
        }
      }
    } : {},
    length(local.scheduling_tolerations) > 0 ? { tolerations = local.scheduling_tolerations } : {},
    try(var.spec.scheduling.priority_class_name, "") != "" ? { priorityClassName = var.spec.scheduling.priority_class_name } : {}
  )

  # ---- spec probes ---------------------------------------------------------------------
  # Thresholds and periods only — the operator owns the probe endpoints.
  # The operator's startup default (600 × 1s) gives first boots a full
  # 10 minutes for schema migrations; the spec tightens it only when
  # declared.
  liveness_probe = {
    for k, v in {
      failureThreshold = try(var.spec.probes.liveness_failure_threshold, null)
      periodSeconds    = try(var.spec.probes.liveness_period_seconds, null)
    } : k => v if v != null
  }

  readiness_probe = {
    for k, v in {
      failureThreshold = try(var.spec.probes.readiness_failure_threshold, null)
      periodSeconds    = try(var.spec.probes.readiness_period_seconds, null)
    } : k => v if v != null
  }

  startup_probe = {
    for k, v in {
      failureThreshold = try(var.spec.probes.startup_failure_threshold, null)
      periodSeconds    = try(var.spec.probes.startup_period_seconds, null)
    } : k => v if v != null
  }

  # ---- networkPolicy / serviceMonitor -----------------------------------------------
  # Rendered ALWAYS, with the spec's effective value: the operator
  # defaults both to true, and an honest declaration states the value it
  # relies on instead of leaning on a default that can change under it.
  network_policy_enabled  = try(coalesce(var.spec.network_policy_enabled, true), true)
  service_monitor_enabled = try(coalesce(var.spec.service_monitor_enabled, true), true)

  # ---- spec.update / spec.tracing ----------------------------------------------------
  # KNOW THIS: with the default RecreateOnImageChange strategy, changing
  # the image takes a full scale-to-zero recreate — an outage window by
  # design (two Keycloak versions cannot share one cache cluster/schema).
  update_block = {
    for k, v in {
      strategy = try(var.spec.update.strategy, null)
      revision = try(var.spec.update.revision, "") != "" ? var.spec.update.revision : null
    } : k => v if v != null
  }

  # The CRD types samplerRatio as a NUMBER; the spec carries it as a
  # pattern-validated string — tonumber() renders what the schema wants.
  tracing_block = {
    for k, v in {
      enabled      = try(var.spec.tracing.enabled, false) ? true : null
      endpoint     = try(var.spec.tracing.endpoint, "") != "" ? var.spec.tracing.endpoint : null
      protocol     = try(var.spec.tracing.protocol, null)
      samplerRatio = try(tonumber(var.spec.tracing.sampler_ratio), null)
    } : k => v if v != null
  }

  # ---- the Keycloak CR spec body -----------------------------------------------------
  # Field names are the CRD's own JSON keys (verified against the pinned
  # v2beta1 schema at operator 26.7.0); values render ONLY when declared
  # so the operator's defaulting stays authoritative — except
  # networkPolicy/serviceMonitor (see above) and ingress (see below).
  # Pulumi twin: keycloakSpecBody.
  keycloak_spec = merge(
    {
      for k, v in {
        instances      = try(var.spec.instances, null)
        image          = try(var.spec.image, "") != "" ? var.spec.image : null
        startOptimized = try(var.spec.start_optimized, false) ? true : null
      } : k => v if v != null
    },
    { db = local.db_block },
    length(local.http_block) > 0 ? { http = local.http_block } : {},
    length(local.hostname_block) > 0 ? { hostname = local.hostname_block } : {},
    try(var.spec.proxy_headers, "") != "" ? { proxy = { headers = var.spec.proxy_headers } } : {},
    length(local.features_block) > 0 ? { features = local.features_block } : {},
    try(var.spec.transaction_xa_enabled, false) ? { transaction = { xaEnabled = true } } : {},
    length(local.cache_config_map_file) > 0 ? { cache = { configMapFile = local.cache_config_map_file } } : {},
    length(local.truststores) > 0 ? { truststores = local.truststores } : {},
    length(local.additional_options) > 0 ? { additionalOptions = local.additional_options } : {},
    try(var.spec.bootstrap_admin_secret_name, "") != "" ? {
      bootstrapAdmin = { user = { secret = var.spec.bootstrap_admin_secret_name } }
    } : {},
    length(local.resources_block) > 0 ? { resources = local.resources_block } : {},
    length(local.scheduling_block) > 0 ? { scheduling = local.scheduling_block } : {},
    length(local.liveness_probe) > 0 ? { livenessProbe = local.liveness_probe } : {},
    length(local.readiness_probe) > 0 ? { readinessProbe = local.readiness_probe } : {},
    length(local.startup_probe) > 0 ? { startupProbe = local.startup_probe } : {},
    try(var.spec.http_management_port, null) != null ? { httpManagement = { port = var.spec.http_management_port } } : {},
    { networkPolicy = { enabled = local.network_policy_enabled } },
    { serviceMonitor = { enabled = local.service_monitor_enabled } },
    length(local.update_block) > 0 ? { update = local.update_block } : {},
    length(local.tracing_block) > 0 ? { tracing = local.tracing_block } : {},
    # ALWAYS disable the operator's default Ingress — an ABSENT ingress
    # block means ENABLED (verified in operator source: absence defaults
    # the Ingress on), so the block must render explicitly. Exposure
    # composes from Gateway API kinds referencing the exported service
    # handles, never from this component.
    { ingress = { enabled = false } }
  )

  # ---- outputs -----------------------------------------------------------------------
  # All derived blind from the operator's naming contract.
  stateful_set      = local.keycloak_name
  service_name      = "${local.keycloak_name}-service"
  discovery_service = "${local.keycloak_name}-discovery"

  api_endpoint        = "${local.api_scheme}://${local.service_name}.${local.namespace}.svc.cluster.local:${local.api_port}"
  management_endpoint = "${local.api_scheme}://${local.service_name}.${local.namespace}.svc.cluster.local:${local.management_port}"

  # The bootstrap-admin credential handle: the user-provided Secret when
  # declared, else the operator-generated create-once `<name>-initial-admin`
  # (username "temp-admin") — seeded at FIRST start only, never rotated by
  # the operator; break-glass material.
  initial_admin_secret_name = try(var.spec.bootstrap_admin_secret_name, "") != "" ? var.spec.bootstrap_admin_secret_name : "${local.keycloak_name}-initial-admin"

  port_forward_command = "kubectl port-forward -n ${local.namespace} svc/${local.service_name} ${local.api_port}:${local.api_port}"
}
