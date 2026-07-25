# Computed values for the KubernetesSignoz module. Every resolution here
# has an exact twin in the Pulumi module's locals.go / values.go — keep
# them in lockstep.
#
# SECRET DISCIPLINE (load-bearing): the bundled ClickHouse admin password
# is module-GENERATED (random_password in main.tf) and reaches the chart
# through set_sensitive — never through these values, so it never appears
# in a plan diff or a values document. KNOW THIS about the upstream grain:
# the chart renders that password into the ClickHouseInstallation object's
# user section and into literal container env — anyone with read access on
# the namespace's CHI/pods can see it. The module cannot change that; what
# it guarantees is that the chart's publicly-documented default password
# NEVER ships and the credential is unique per install. The external arm
# is fully secret-native (existingSecret → secretKeyRef). The SMTP
# password rides a valueFrom secretKeyRef env entry (the chart's flexible
# env structure) — a reference, never material.
#
# HCL DISCIPLINE (applies to every conditional object below): conditional
# entries are written as `key = cond ? value : null` inside ONE object
# literal, pruned with `{ for k, v in {...} : k => v if v != null }`.
# `cond ? {...} : {}` ternaries fail plan-time type unification when
# branches carry different attributes. Optional nested blocks are read
# with try() (HCL's && does NOT short-circuit).

locals {
  # Chart identity — must stay byte-identical with the Pulumi module's
  # vars. The signoz chart's version tracks the SigNoz application
  # version in lockstep (chart 0.133.0 = app v0.133.0).
  helm_chart_name = "signoz"
  helm_chart_repo = "https://charts.signoz.io"

  # Release name — metadata.name, NOT a fixed chart name: several SigNoz
  # instances can coexist (per-team platforms, a staging instance).
  # fullnameOverride pins every chart child name to it.
  release_name = var.metadata.name

  chart_version = coalesce(var.spec.chart_version, "0.133.0")

  namespace = var.spec.namespace

  # Child names compose the ClickHouseInstallation's operator-generated
  # StatefulSets: `chi-<name>-clickhouse-cluster-0-0` is 27 characters of
  # scaffolding around the resource name inside Kubernetes' 63-character
  # cap. The helm_release precondition (main.tf) fails the plan loudly;
  # this flag is its condition (twin: the Pulumi module's MaxNameLength
  # guard).
  max_name_length    = 30
  name_within_budget = length(local.release_name) <= local.max_name_length

  labels = merge(
    {
      "planton.ai/resource"      = "true"
      "planton.ai/resource-name" = var.metadata.name
      "planton.ai/resource-kind" = "KubernetesSignoz"
    },
    var.metadata.id != null && var.metadata.id != "" ? { "planton.ai/resource-id" = var.metadata.id } : {},
    var.metadata.org != null && var.metadata.org != "" ? { "planton.ai/organization" = var.metadata.org } : {},
    var.metadata.env != null && var.metadata.env != "" ? { "planton.ai/environment" = var.metadata.env } : {}
  )

  # ---- database arm resolution -------------------------------------------
  # Empty oneof = the bundled arm with defaults (the appliance posture).
  is_external = try(var.spec.external_clickhouse, null) != null

  # The bundled installation's fullname — pins the CHI name AND the client
  # Service name the chart's own helpers derive (clickhouse.servicename
  # honors clickhouse.fullnameOverride).
  clickhouse_fullname = "${local.release_name}-clickhouse"

  # The module-owned Secret exporting the generated bundled credential
  # (keys username/password). The chart itself offers no Secret for the
  # bundled arm — this is the composition handle the outputs promise.
  clickhouse_auth_secret_name = "${local.release_name}-clickhouse-auth"

  # ---- bundled clickhouse rendering ----------------------------------------
  managed = try(var.spec.managed_clickhouse, null)

  managed_resources = local.is_external || try(local.managed.resources, null) == null ? null : { for rk, rv in {
    requests = try(local.managed.resources.requests, null) == null ? null : { for qk, qv in {
      cpu    = try(local.managed.resources.requests.cpu, "") != "" ? local.managed.resources.requests.cpu : null
      memory = try(local.managed.resources.requests.memory, "") != "" ? local.managed.resources.requests.memory : null
    } : qk => qv if qv != null }
    limits = try(local.managed.resources.limits, null) == null ? null : { for lk, lv in {
      cpu    = try(local.managed.resources.limits.cpu, "") != "" ? local.managed.resources.limits.cpu : null
      memory = try(local.managed.resources.limits.memory, "") != "" ? local.managed.resources.limits.memory : null
    } : lk => lv if lv != null }
  } : rk => rv if rv != null && rv != {} }

  zookeeper_resources = local.is_external || try(local.managed.zookeeper.resources, null) == null ? null : { for rk, rv in {
    requests = try(local.managed.zookeeper.resources.requests, null) == null ? null : { for qk, qv in {
      cpu    = try(local.managed.zookeeper.resources.requests.cpu, "") != "" ? local.managed.zookeeper.resources.requests.cpu : null
      memory = try(local.managed.zookeeper.resources.requests.memory, "") != "" ? local.managed.zookeeper.resources.requests.memory : null
    } : qk => qv if qv != null }
    limits = try(local.managed.zookeeper.resources.limits, null) == null ? null : { for lk, lv in {
      cpu    = try(local.managed.zookeeper.resources.limits.cpu, "") != "" ? local.managed.zookeeper.resources.limits.cpu : null
      memory = try(local.managed.zookeeper.resources.limits.memory, "") != "" ? local.managed.zookeeper.resources.limits.memory : null
    } : lk => lv if lv != null }
  } : rk => rv if rv != null && rv != {} }

  # Cold storage: keyless IRSA (role annotations on the ClickHouse service
  # account) XOR declared keys (rendered by the chart into ClickHouse's
  # storage configuration — the upstream grain, taught on the spec field).
  cold_s3  = local.is_external ? null : try(local.managed.cold_storage.s3, null)
  cold_gcs = local.is_external ? null : try(local.managed.cold_storage.gcs, null)

  cold_storage_block = local.cold_s3 == null && local.cold_gcs == null ? null : { for k, v in {
    enabled  = true
    type     = local.cold_s3 != null ? "s3" : "gcs"
    endpoint = local.cold_s3 != null ? local.cold_s3.endpoint : local.cold_gcs.endpoint
    accessKey = (
      local.cold_s3 != null && try(local.cold_s3.access_key, "") != "" ? local.cold_s3.access_key :
      local.cold_gcs != null ? local.cold_gcs.access_key : null
    )
    secretAccess = (
      local.cold_s3 != null && try(local.cold_s3.secret_key, "") != "" ? local.cold_s3.secret_key :
      local.cold_gcs != null ? local.cold_gcs.secret_key : null
    )
    role = local.cold_s3 != null && try(local.cold_s3.irsa_role_arn, "") != "" ? {
      enabled     = true
      annotations = { "eks.amazonaws.com/role-arn" = local.cold_s3.irsa_role_arn }
    } : null
  } : k => v if v != null }

  # ONE pruned object serves both arms: `enabled` is always present and
  # every bundled-arm key is null on the external arm — never a ternary
  # whose branches carry different attribute sets (the HCL
  # type-unification class this module's own first plan caught).
  clickhouse_block = { for k, v in {
    enabled          = !local.is_external
    fullnameOverride = local.is_external ? null : local.clickhouse_fullname
    layout = local.is_external ? null : {
      shardsCount   = coalesce(try(local.managed.shards, null), 1)
      replicasCount = coalesce(try(local.managed.replicas, null), 1)
    }
    persistence = local.is_external ? null : { for pk, pv in {
      enabled      = true
      size         = coalesce(try(local.managed.disk_size, null), "20Gi")
      storageClass = try(local.managed.storage_class, "") != "" ? local.managed.storage_class : null
    } : pk => pv if pv != null }
    resources         = local.managed_resources
    allowedNetworkIps = local.is_external ? null : (length(try(local.managed.allowed_network_ips, [])) > 0 ? local.managed.allowed_network_ips : null)
    zookeeper = local.is_external ? null : { for zk, zv in {
      replicaCount = coalesce(try(local.managed.zookeeper.replicas, null), 1)
      resources    = local.zookeeper_resources
    } : zk => zv if zv != null }
    coldStorage = local.cold_storage_block
  } : k => v if v != null }

  # ---- external clickhouse rendering ---------------------------------------
  external = try(var.spec.external_clickhouse, null)

  external_tcp_port  = local.is_external ? coalesce(try(local.external.tcp_port, null), 9000) : 9000
  external_http_port = local.is_external ? coalesce(try(local.external.http_port, null), 8123) : 8123

  external_clickhouse_block = !local.is_external ? null : { for k, v in {
    host                      = local.external.host
    cluster                   = try(local.external.cluster_name, "") != "" ? local.external.cluster_name : "cluster"
    tcpPort                   = local.external_tcp_port
    httpPort                  = local.external_http_port
    user                      = local.external.username
    existingSecret            = local.external.password_secret.secret_name
    existingSecretPasswordKey = local.external.password_secret.secret_key
    secure                    = try(local.external.secure, false) ? true : null
    verify                    = try(local.external.verify, false) ? true : null
  } : k => v if v != null }

  # ---- signoz server env ----------------------------------------------------
  # Typed entries WIN over the user's advanced env map (the spec's
  # documented contract). Env keys follow SigNoz's own derivation
  # (signoz_<section>_<key>, embedded underscores doubled). The SMTP
  # password rides a valueFrom secretKeyRef object — the chart's flexible
  # env structure renders it as a proper env-from-secret.
  smtp = try(var.spec.server.smtp, null)

  typed_env = merge(
    try(var.spec.server.external_url, "") != "" ? {
      signoz_alertmanager_signoz_external__url = var.spec.server.external_url
    } : {},
    local.smtp != null ? merge(
      {
        signoz_emailing_enabled      = "true"
        signoz_emailing_smtp_address = local.smtp.address
        signoz_emailing_smtp_from    = local.smtp.from
      },
      try(local.smtp.username, "") != "" ? { signoz_emailing_smtp_auth_username = local.smtp.username } : {},
      try(local.smtp.tls_enabled, false) ? { signoz_emailing_smtp_tls_enabled = "true" } : {},
    ) : {},
  )

  smtp_password_env = local.smtp != null && try(local.smtp.password_secret, null) != null ? {
    signoz_emailing_smtp_auth_password = {
      valueFrom = { secretKeyRef = {
        name = local.smtp.password_secret.name
        key  = local.smtp.password_secret.key
      } }
    }
  } : {}

  # merge() unifies the plain-string entries and the valueFrom object —
  # the chart's env helper accepts both shapes per key.
  signoz_env = merge(try(var.spec.server.env, {}), local.typed_env, local.smtp_password_env)

  # ---- scheduling (server + collector + migrator) ---------------------------
  node_selector = try(var.spec.scheduling.node_selector, {})
  tolerations = [
    for t in try(var.spec.scheduling.tolerations, []) : {
      for tk, tv in {
        key               = t.key != "" ? t.key : null
        operator          = t.operator != "" ? t.operator : null
        value             = t.value != "" ? t.value : null
        effect            = t.effect != "" ? t.effect : null
        tolerationSeconds = try(t.toleration_seconds, null)
      } : tk => tv if tv != null
    }
  ]
  priority_class = try(var.spec.scheduling.priority_class_name, "")

  server_resources = try(var.spec.server.resources, null) == null ? null : { for rk, rv in {
    requests = try(var.spec.server.resources.requests, null) == null ? null : { for qk, qv in {
      cpu    = try(var.spec.server.resources.requests.cpu, "") != "" ? var.spec.server.resources.requests.cpu : null
      memory = try(var.spec.server.resources.requests.memory, "") != "" ? var.spec.server.resources.requests.memory : null
    } : qk => qv if qv != null }
    limits = try(var.spec.server.resources.limits, null) == null ? null : { for lk, lv in {
      cpu    = try(var.spec.server.resources.limits.cpu, "") != "" ? var.spec.server.resources.limits.cpu : null
      memory = try(var.spec.server.resources.limits.memory, "") != "" ? var.spec.server.resources.limits.memory : null
    } : lk => lv if lv != null }
  } : rk => rv if rv != null && rv != {} }

  signoz_block = { for k, v in {
    persistence = { for pk, pv in {
      enabled      = true
      size         = coalesce(try(var.spec.server.disk_size, null), "1Gi")
      storageClass = try(var.spec.server.storage_class, "") != "" ? var.spec.server.storage_class : null
    } : pk => pv if pv != null }
    resources         = local.server_resources
    env               = length(local.signoz_env) > 0 ? local.signoz_env : null
    nodeSelector      = length(local.node_selector) > 0 ? local.node_selector : null
    tolerations       = length(local.tolerations) > 0 ? local.tolerations : null
    priorityClassName = local.priority_class != "" ? local.priority_class : null
  } : k => v if v != null }

  # ---- otel collector --------------------------------------------------------
  collector = try(var.spec.otel_collector, null)

  collector_resources = try(local.collector.resources, null) == null ? null : { for rk, rv in {
    requests = try(local.collector.resources.requests, null) == null ? null : { for qk, qv in {
      cpu    = try(local.collector.resources.requests.cpu, "") != "" ? local.collector.resources.requests.cpu : null
      memory = try(local.collector.resources.requests.memory, "") != "" ? local.collector.resources.requests.memory : null
    } : qk => qv if qv != null }
    limits = try(local.collector.resources.limits, null) == null ? null : { for lk, lv in {
      cpu    = try(local.collector.resources.limits.cpu, "") != "" ? local.collector.resources.limits.cpu : null
      memory = try(local.collector.resources.limits.memory, "") != "" ? local.collector.resources.limits.memory : null
    } : lk => lv if lv != null }
  } : rk => rv if rv != null && rv != {} }

  # Receiver toggles. jaeger + http-logs default ON (the chart's grain),
  # zipkin defaults OFF. The pipeline receiver LISTS are always rendered
  # from the toggles (lists replace under Helm merge — rendering them from
  # one derivation is what keeps the Service ports and the collector
  # pipelines in agreement by construction).
  jaeger_enabled    = try(local.collector.jaeger_receiver_enabled, null) != null ? local.collector.jaeger_receiver_enabled : true
  zipkin_enabled    = try(local.collector.zipkin_receiver_enabled, false)
  http_logs_enabled = try(local.collector.http_logs_receivers_enabled, null) != null ? local.collector.http_logs_receivers_enabled : true

  traces_receivers = concat(
    ["otlp"],
    local.jaeger_enabled ? ["jaeger"] : [],
    local.zipkin_enabled ? ["zipkin"] : [],
  )
  logs_receivers = concat(
    ["otlp"],
    local.http_logs_enabled ? ["httplogreceiver/heroku", "httplogreceiver/json"] : [],
  )

  collector_ports = { for k, v in {
    zipkin           = local.zipkin_enabled ? { enabled = true } : null
    "jaeger-thrift"  = local.jaeger_enabled ? null : { enabled = false }
    "jaeger-grpc"    = local.jaeger_enabled ? null : { enabled = false }
    logsheroku       = local.http_logs_enabled ? null : { enabled = false }
    logsjson         = local.http_logs_enabled ? null : { enabled = false }
  } : k => v if v != null }

  collector_config = { for k, v in {
    receivers = local.zipkin_enabled ? { zipkin = { endpoint = "0.0.0.0:9411" } } : null
    service = {
      pipelines = {
        traces = { receivers = local.traces_receivers }
        logs   = { receivers = local.logs_receivers }
      }
    }
  } : k => v if v != null }

  autoscaling = try(local.collector.autoscaling, null)
  autoscaling_block = try(local.autoscaling.enabled, false) ? { for k, v in {
    enabled                        = true
    minReplicas                    = coalesce(try(local.autoscaling.min_replicas, null), 1)
    maxReplicas                    = coalesce(try(local.autoscaling.max_replicas, null), 11)
    targetCPUUtilizationPercentage = try(local.autoscaling.target_cpu_utilization_percent, null)
    targetMemoryUtilizationPercentage = try(local.autoscaling.target_memory_utilization_percent, null)
  } : k => v if v != null } : null

  otel_collector_block = { for k, v in {
    replicaCount                     = coalesce(try(local.collector.replicas, null), 1)
    resources                        = local.collector_resources
    autoscaling                      = local.autoscaling_block
    lowCardinalityExceptionGrouping  = try(local.collector.low_cardinality_exception_grouping, false) ? true : null
    ports                            = length(local.collector_ports) > 0 ? local.collector_ports : null
    config                           = local.collector_config
    nodeSelector                     = length(local.node_selector) > 0 ? local.node_selector : null
    tolerations                      = length(local.tolerations) > 0 ? local.tolerations : null
    priorityClassName                = local.priority_class != "" ? local.priority_class : null
  } : k => v if v != null }

  migrator_block = { for k, v in {
    nodeSelector = length(local.node_selector) > 0 ? local.node_selector : null
    tolerations  = length(local.tolerations) > 0 ? local.tolerations : null
  } : k => v if v != null }

  # ---- globals ---------------------------------------------------------------
  # Every chart image here is the SPLIT registry+repository form and its
  # registry key defers to global.imageRegistry — one override reaches the
  # SigNoz server, the collector, ClickHouse, the bundled operator and its
  # metrics exporter, and ZooKeeper. The bundled arm's UDF init container
  # image (alpine) follows its own chart key — helm_values territory,
  # documented on the spec field.
  global_block = { for k, v in {
    imageRegistry    = var.spec.image_registry != "" ? var.spec.image_registry : null
    clusterName      = var.spec.cluster_name != "" ? var.spec.cluster_name : null
    imagePullSecrets = length(var.spec.image_pull_secrets) > 0 ? var.spec.image_pull_secrets : null
  } : k => v if v != null }

  # ---- typed chart values (twin of the Pulumi module's buildHelmValues) -----
  # The bundled admin password deliberately ABSENT: it travels via
  # set_sensitive (main.tf) so it never appears in a values document.
  helm_values = { for k, v in {
    fullnameOverride = local.release_name
    clusterName      = var.spec.cluster_name != "" ? var.spec.cluster_name : null
    global           = length(local.global_block) > 0 ? local.global_block : null

    clickhouse          = local.clickhouse_block
    externalClickhouse  = local.external_clickhouse_block

    signoz                 = local.signoz_block
    otelCollector          = local.otel_collector_block
    telemetryStoreMigrator = length(local.migrator_block) > 0 ? local.migrator_block : null
  } : k => v if v != null }
}
