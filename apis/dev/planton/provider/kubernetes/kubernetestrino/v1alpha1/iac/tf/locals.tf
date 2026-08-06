# Every resolution here has an exact twin in the Pulumi module's
# locals.go / values.go — keep them in lockstep.
#
# HCL DISCIPLINE (the type-unification class): every conditional block
# below is a SINGLE-ATTRIBUTE ternary against {} or rides the null-prune
# merge idiom — a ternary whose true branch bundles attributes of more
# than one type cannot unify against {}. Mixed-shape LIST entries (the
# env list carries both value and valueFrom objects) ride the
# jsonencode/jsondecode seam, which erases the object-type distinction.

locals {
  # Chart identity — must stay byte-identical with the Pulumi module's
  # vars (cross-engine chart-name drift deploys two different products
  # from one manifest).
  helm_chart_name    = "trino"
  helm_chart_repo    = "https://trinodb.github.io/charts"
  helm_chart_version = "1.42.2"

  release_name = var.metadata.name
  namespace    = var.spec.namespace

  # Resource-identity labels stamped on the module-created satellites
  # (the namespace and the module-owned Secrets — never injected into
  # the chart's own resources; Helm owns those).
  labels = merge(
    {
      "planton.ai/resource"      = "true"
      "planton.ai/resource-name" = var.metadata.name
      "planton.ai/resource-kind" = "KubernetesTrino"
    },
    try(var.metadata.id, "") != "" ? { "planton.ai/resource-id" = var.metadata.id } : {},
    try(var.metadata.org, "") != "" ? { "planton.ai/organization" = var.metadata.org } : {},
    try(var.metadata.env, "") != "" ? { "planton.ai/environment" = var.metadata.env } : {},
  )

  # ----------------------------- name budget ------------------------------
  # Chart truth at the pin: `<fullname>-schemas-volume-coordinator`
  # (27-char suffix) renders unconditionally; the resource-groups
  # ConfigMap suffix (36 chars) renders only when resource groups are
  # declared. Both must fit the 63-char DNS bound — checked fail-loud
  # in main.tf preconditions (Pulumi twin: buildHelmValues errors).
  name_budget                 = 36
  name_budget_resource_groups = 27

  # --------------------------- authentication -----------------------------
  # The secured default: an ABSENT auth block means PASSWORD auth with
  # a module-generated admin — the chart's own default (no
  # authentication at all) never ships.
  auth_enabled = try(var.spec.auth.enabled, null) == null ? true : var.spec.auth.enabled

  # On the bring-your-own password-file arm no module-generated admin
  # exists — the exported credential handles stay EMPTY (honest).
  password_db_byo          = try(var.spec.auth.existing_password_db_secret, null) != null
  module_owned_password_db = local.auth_enabled && !local.password_db_byo
  admin_username           = local.module_owned_password_db ? coalesce(try(var.spec.auth.admin_username, null), "trino") : ""

  password_db_secret_name = local.auth_enabled ? (
    local.password_db_byo ? var.spec.auth.existing_password_db_secret.secret_name : "${local.release_name}-auth"
  ) : ""
  groups_secret_name   = try(var.spec.auth.groups_secret.secret_name, "")
  internal_secret_name = local.auth_enabled ? "${local.release_name}-internal" : ""

  shared_secret_env_var = "TRINO_INTERNAL_SHARED_SECRET"

  https_enabled      = try(var.spec.https.enabled, false)
  https_keystore_key = coalesce(try(var.spec.https.keystore_secret.secret_key, null), "keystore.jks")

  # --------------------- additional config properties ---------------------
  # Module-owned lines FIRST (the security spine), then the spec's
  # escape-hatch lines. Re-pinned after the helm_values merge —
  # additional_config_properties is the supported extension point,
  # never helm_values. The shared secret and any catalog password ride
  # `${ENV:VAR}` (Trino's own secrets substitution) — never literals in
  # this ConfigMap-rendered list.
  config_properties = concat(
    local.auth_enabled ? ["internal-communication.shared-secret=$${ENV:${local.shared_secret_env_var}}"] : [],
    # Password auth engages ONLY on secure requests. Verified in the
    # server's AuthenticationFilter at the pin: allow-insecure-over-http
    # does NOT run password auth over HTTP — it routes plain-HTTP
    # requests to the username-trust authenticator, so the password
    # file would guard nothing. process-forwarded is upstream's
    # TLS-terminating-proxy recipe: requests arriving with
    # X-Forwarded-Proto: https (what composed exposure kinds send) are
    # treated secure and the PASSWORD authenticator ENFORCES the file,
    # while plain HTTP data-plane requests fail CLOSED (403). Health
    # probes are unaffected (/v1/info and /v1/status are PUBLIC routes).
    local.auth_enabled && !local.https_enabled ? ["http-server.process-forwarded=true"] : [],
    try(var.spec.fault_tolerant_execution, null) != null ? ["retry-policy=${var.spec.fault_tolerant_execution.retry_policy}"] : [],
    try(var.spec.additional_config_properties, []),
  )

  # ------------------------------ catalogs --------------------------------
  # Each catalog renders as one properties document; credentials are
  # `${ENV:...}` references resolved from Secret-sourced env vars at
  # process start (env entries below). Catalog names are [a-z][a-z0-9_]*
  # by CEL, so upper(name) is env-var safe.
  postgres_catalogs = {
    for c in try(var.spec.catalogs.postgres, []) : c.name => join("\n", concat([
      "connector.name=postgresql",
      "connection-url=jdbc:postgresql://${c.host}:${coalesce(try(c.port, null), 5432)}/${c.database}",
      "connection-user=${coalesce(try(c.username, null), "app")}",
      "connection-password=$${ENV:TRINO_CATALOG_${upper(c.name)}_PASSWORD}",
    ], try(c.additional_properties, [])))
  }

  # MySQL exposes databases as Trino schemas — the JDBC URL carries no
  # database segment (connector truth at the pin).
  mysql_catalogs = {
    for c in try(var.spec.catalogs.mysql, []) : c.name => join("\n", concat([
      "connector.name=mysql",
      "connection-url=jdbc:mysql://${c.host}:${coalesce(try(c.port, null), 3306)}",
      "connection-user=${coalesce(try(c.username, null), "root")}",
      "connection-password=$${ENV:TRINO_CATALOG_${upper(c.name)}_PASSWORD}",
    ], try(c.additional_properties, [])))
  }

  sample_catalogs_enabled = try(var.spec.catalogs.sample_catalogs_enabled, null) == null ? true : var.spec.catalogs.sample_catalogs_enabled

  # Helm null-deletes the chart's default tpch/tpcds map entries when
  # the samples are disabled.
  catalogs_block = merge(
    local.postgres_catalogs,
    local.mysql_catalogs,
    try(var.spec.catalogs.custom, {}),
    local.sample_catalogs_enabled ? {} : { tpch = null, tpcds = null },
  )

  # ----------------- secret-sourced environment variables -----------------
  # The delivery vehicle for every ${ENV:...} reference rendered above:
  # the internal shared secret, per-catalog passwords, and the user's
  # extra entries — Secret NAMES only ever render here. The list mixes
  # valueFrom and value entry shapes → the jsonencode/jsondecode seam.
  env_entries = jsondecode(jsonencode(concat(
    local.auth_enabled ? [{
      name = local.shared_secret_env_var
      valueFrom = {
        secretKeyRef = {
          name = local.internal_secret_name
          key  = "shared-secret"
        }
      }
    }] : [],
    [for c in try(var.spec.catalogs.postgres, []) : {
      name = "TRINO_CATALOG_${upper(c.name)}_PASSWORD"
      valueFrom = {
        secretKeyRef = {
          name = c.password_secret.secret_name
          key  = coalesce(try(c.password_secret.secret_key, null), "password")
        }
      }
    }],
    [for c in try(var.spec.catalogs.mysql, []) : {
      name = "TRINO_CATALOG_${upper(c.name)}_PASSWORD"
      valueFrom = {
        secretKeyRef = {
          name = c.password_secret.secret_name
          key  = coalesce(try(c.password_secret.secret_key, null), "password")
        }
      }
    }],
    [for name, ref in try(var.spec.extra_env_from_secret, {}) : {
      name = name
      valueFrom = {
        secretKeyRef = {
          name = ref.secret_name
          key  = ref.secret_key
        }
      }
    }],
    [for name, value in try(var.spec.extra_env, {}) : {
      name  = name
      value = value
    }],
  )))

  # ------------------------------- server ---------------------------------
  worker_replicas = coalesce(try(var.spec.workers.replicas, null), 2)

  server_config_block = merge(
    {
      query = {
        maxMemory = coalesce(try(var.spec.max_query_memory, null), "4GB")
      }
    },
    local.auth_enabled ? { authenticationType = "PASSWORD" } : {},
    local.https_enabled ? {
      https = {
        enabled = true
        port    = coalesce(try(var.spec.https.port, null), 8443)
        keystore = {
          # The keystore Secret mounts through secretMounts; the chart
          # wires this path into config.properties.
          path = "/etc/trino/keystore/${local.https_keystore_key}"
        }
      }
    } : {},
  )

  # Worker autoscaling — HPA XOR KEDA (the spec's oneof). The chart
  # disables a metric when its target is an EMPTY STRING — 0 in the
  # spec maps to that disable contract; the mixed string/number targets
  # ride the jsonencode/jsondecode seam.
  hpa_cpu_target    = coalesce(try(var.spec.workers.hpa.target_cpu_utilization_percent, null), 50)
  hpa_memory_target = coalesce(try(var.spec.workers.hpa.target_memory_utilization_percent, null), 80)

  autoscaling_block = try(var.spec.workers.hpa, null) == null ? {} : {
    autoscaling = jsondecode(jsonencode({
      enabled                           = true
      maxReplicas                       = var.spec.workers.hpa.max_replicas
      targetCPUUtilizationPercentage    = local.hpa_cpu_target == 0 ? "" : local.hpa_cpu_target
      targetMemoryUtilizationPercentage = local.hpa_memory_target == 0 ? "" : local.hpa_memory_target
    }))
  }

  keda_block = try(var.spec.workers.keda, null) == null ? {} : {
    keda = merge(
      {
        enabled         = true
        maxReplicaCount = var.spec.workers.keda.max_replicas
        triggers        = yamldecode(var.spec.workers.keda.triggers)
      },
      try(var.spec.workers.keda.min_replicas, null) != null ? { minReplicaCount = var.spec.workers.keda.min_replicas } : {},
      try(var.spec.workers.keda.polling_interval_seconds, null) != null ? { pollingInterval = var.spec.workers.keda.polling_interval_seconds } : {},
      try(var.spec.workers.keda.cooldown_period_seconds, null) != null ? { cooldownPeriod = var.spec.workers.keda.cooldown_period_seconds } : {},
    )
  }

  # Fault-tolerant execution: the exchange manager spools exchange data
  # durably; the retry policy rides config_properties above.
  exchange_manager_block = try(var.spec.fault_tolerant_execution, null) == null ? {} : {
    exchangeManager = {
      name    = "filesystem"
      baseDir = var.spec.fault_tolerant_execution.exchange_manager.base_directories
    }
  }

  server_block = merge(
    {
      workers = local.worker_replicas
      node = {
        environment = coalesce(try(var.spec.node_environment, null), "production")
      }
      log = {
        trino = {
          level = coalesce(try(var.spec.log_level, null), "INFO")
        }
      }
      config = local.server_config_block
    },
    local.exchange_manager_block,
    local.autoscaling_block,
    local.keda_block,
  )

  # -------------------- coordinator / worker blocks -----------------------
  # JVM heap: percent-based sizing only works when the fixed -Xmx is
  # UNSET (chart truth). The chart's 8G default is disabled with an
  # EMPTY STRING, never null: the chart guards the -Xmx line with
  # `{{- if .Values.<node>.jvm.maxHeapSize }}` and "" is falsy there,
  # while a null's survival to Helm differs by engine (verified live:
  # the Pulumi seam dropped the null and the chart's -Xmx8G default
  # silently overrode the percent sizing inside the container limit).
  coordinator_jvm_block = try(var.spec.coordinator.jvm.max_heap_percent, null) != null ? {
    jvm = {
      maxHeapSize    = ""
      maxHeapPercent = var.spec.coordinator.jvm.max_heap_percent
    }
    } : (try(var.spec.coordinator.jvm.max_heap_size, "") != "" ? {
      jvm = {
        maxHeapSize = var.spec.coordinator.jvm.max_heap_size
      }
  } : {})

  worker_jvm_block = try(var.spec.workers.jvm.max_heap_percent, null) != null ? {
    jvm = {
      maxHeapSize    = ""
      maxHeapPercent = var.spec.workers.jvm.max_heap_percent
    }
    } : (try(var.spec.workers.jvm.max_heap_size, "") != "" ? {
      jvm = {
        maxHeapSize = var.spec.workers.jvm.max_heap_size
      }
  } : {})

  coordinator_resources = {
    requests = {
      cpu    = try(var.spec.coordinator.resources.requests.cpu, "")
      memory = try(var.spec.coordinator.resources.requests.memory, "")
    }
    limits = {
      cpu    = try(var.spec.coordinator.resources.limits.cpu, "")
      memory = try(var.spec.coordinator.resources.limits.memory, "")
    }
  }
  coordinator_resources_block = (
    local.coordinator_resources.requests.cpu != "" || local.coordinator_resources.requests.memory != "" ||
    local.coordinator_resources.limits.cpu != "" || local.coordinator_resources.limits.memory != ""
    ) ? {
    resources = merge(
      local.coordinator_resources.requests.cpu != "" || local.coordinator_resources.requests.memory != "" ? {
        requests = merge(
          local.coordinator_resources.requests.cpu != "" ? { cpu = local.coordinator_resources.requests.cpu } : {},
          local.coordinator_resources.requests.memory != "" ? { memory = local.coordinator_resources.requests.memory } : {},
        )
      } : {},
      local.coordinator_resources.limits.cpu != "" || local.coordinator_resources.limits.memory != "" ? {
        limits = merge(
          local.coordinator_resources.limits.cpu != "" ? { cpu = local.coordinator_resources.limits.cpu } : {},
          local.coordinator_resources.limits.memory != "" ? { memory = local.coordinator_resources.limits.memory } : {},
        )
      } : {},
    )
  } : {}

  worker_resources = {
    requests = {
      cpu    = try(var.spec.workers.resources.requests.cpu, "")
      memory = try(var.spec.workers.resources.requests.memory, "")
    }
    limits = {
      cpu    = try(var.spec.workers.resources.limits.cpu, "")
      memory = try(var.spec.workers.resources.limits.memory, "")
    }
  }
  worker_resources_block = (
    local.worker_resources.requests.cpu != "" || local.worker_resources.requests.memory != "" ||
    local.worker_resources.limits.cpu != "" || local.worker_resources.limits.memory != ""
    ) ? {
    resources = merge(
      local.worker_resources.requests.cpu != "" || local.worker_resources.requests.memory != "" ? {
        requests = merge(
          local.worker_resources.requests.cpu != "" ? { cpu = local.worker_resources.requests.cpu } : {},
          local.worker_resources.requests.memory != "" ? { memory = local.worker_resources.requests.memory } : {},
        )
      } : {},
      local.worker_resources.limits.cpu != "" || local.worker_resources.limits.memory != "" ? {
        limits = merge(
          local.worker_resources.limits.cpu != "" ? { cpu = local.worker_resources.limits.cpu } : {},
          local.worker_resources.limits.memory != "" ? { memory = local.worker_resources.limits.memory } : {},
        )
      } : {},
    )
  } : {}

  coordinator_tolerations = [for t in try(var.spec.coordinator.scheduling.tolerations, []) : jsondecode(jsonencode(merge(
    t.key != "" ? { key = t.key } : {},
    t.operator != "" ? { operator = t.operator } : {},
    t.value != "" ? { value = t.value } : {},
    t.effect != "" ? { effect = t.effect } : {},
    try(t.toleration_seconds, null) != null ? { tolerationSeconds = t.toleration_seconds } : {},
  )))]

  worker_tolerations = [for t in try(var.spec.workers.scheduling.tolerations, []) : jsondecode(jsonencode(merge(
    t.key != "" ? { key = t.key } : {},
    t.operator != "" ? { operator = t.operator } : {},
    t.value != "" ? { value = t.value } : {},
    t.effect != "" ? { effect = t.effect } : {},
    try(t.toleration_seconds, null) != null ? { tolerationSeconds = t.toleration_seconds } : {},
  )))]

  coordinator_block = merge(
    {
      config = merge(
        {
          query = {
            maxMemoryPerNode = coalesce(try(var.spec.coordinator.max_query_memory_per_node, null), "1GB")
          }
          nodeScheduler = {
            includeCoordinator = try(var.spec.coordinator.include_in_scheduling, false)
          }
        },
        try(var.spec.coordinator.heap_headroom_per_node, "") != "" ? {
          memory = { heapHeadroomPerNode = var.spec.coordinator.heap_headroom_per_node }
        } : {},
      )
    },
    local.coordinator_jvm_block,
    local.coordinator_resources_block,
    length(try(var.spec.coordinator.scheduling.node_selector, {})) > 0 ? { nodeSelector = var.spec.coordinator.scheduling.node_selector } : {},
    length(local.coordinator_tolerations) > 0 ? { tolerations = local.coordinator_tolerations } : {},
  )

  # Chart rule: the pod termination budget must be at least TWICE the
  # drain window — set exactly 2× so the drain always fits.
  worker_grace_period = coalesce(try(var.spec.workers.graceful_shutdown.grace_period_seconds, null), 120)

  worker_block = merge(
    {
      config = merge(
        {
          query = {
            maxMemoryPerNode = coalesce(try(var.spec.workers.max_query_memory_per_node, null), "1GB")
          }
        },
        try(var.spec.workers.heap_headroom_per_node, "") != "" ? {
          memory = { heapHeadroomPerNode = var.spec.workers.heap_headroom_per_node }
        } : {},
      )
    },
    local.worker_jvm_block,
    local.worker_resources_block,
    length(try(var.spec.workers.scheduling.node_selector, {})) > 0 ? { nodeSelector = var.spec.workers.scheduling.node_selector } : {},
    length(local.worker_tolerations) > 0 ? { tolerations = local.worker_tolerations } : {},
    try(var.spec.workers.graceful_shutdown.enabled, false) ? {
      gracefulShutdown = {
        enabled            = true
        gracePeriodSeconds = local.worker_grace_period
      }
    } : {},
    try(var.spec.workers.graceful_shutdown.enabled, false) ? {
      terminationGracePeriodSeconds = local.worker_grace_period * 2
    } : {},
  )

  # ------------------------------- image ----------------------------------
  # The chart's SPLIT image form: registry (empty = Docker Hub) +
  # repository + tag. The module always pins the tag explicitly so the
  # deployed engine version is declared, never inherited.
  image_block = merge(
    {
      repository = coalesce(try(var.spec.image.repository, null), "trinodb/trino")
      tag        = coalesce(try(var.spec.image.tag, null), "480")
    },
    try(var.spec.image.registry, "") != "" ? { registry = var.spec.image.registry } : {},
  )

  # -------------------------------- auth ----------------------------------
  auth_block = merge(
    { passwordAuthSecret = local.password_db_secret_name },
    local.groups_secret_name != "" ? { groupsAuthSecret = local.groups_secret_name } : {},
  )

  # ------------------------------- metrics --------------------------------
  metrics_enabled         = try(var.spec.metrics.enabled, false)
  service_monitor_enabled = try(var.spec.metrics.service_monitor_enabled, false)

  # The standalone JMX exporter FATALS without a hostPort/jmxUrl in its
  # config (verified live: "you must configure 'jmxUrl' or 'hostPort'"),
  # and the chart's default configProperties is EMPTY — enabling the
  # sidecar without composing the config ships a crash-loop. The module
  # renders the chart's own documented pairing; the `tpl` reference
  # keeps the port single-sourced from the chart's jmx.registryPort.
  jmx_exporter_config = <<-EOT
    hostPort: localhost:{{- .Values.jmx.registryPort }}
    startDelaySeconds: 0
    ssl: false
  EOT

  jmx_block = local.metrics_enabled ? {
    jmx = {
      enabled = true
      exporter = merge(
        {
          enabled          = true
          configProperties = local.jmx_exporter_config
        },
        try(var.spec.metrics.exporter_image, "") != "" ? { image = var.spec.metrics.exporter_image } : {},
      )
    }
  } : {}

  # ---------------------------- typed values ------------------------------
  helm_values = merge(
    {
      # Deterministic child names (`<name>-coordinator`, …) — the
      # release name never double-prefixes and the import map stays
      # exact (the fullname re-pin discipline).
      fullnameOverride = local.release_name
      image            = local.image_block
      server           = local.server_block
      coordinator      = local.coordinator_block
      worker           = local.worker_block
      service = merge(
        { type = coalesce(try(var.spec.service.type, null), "ClusterIP") },
        length(try(var.spec.service.annotations, {})) > 0 ? { annotations = var.spec.service.annotations } : {},
      )
    },
    length(local.config_properties) > 0 ? { additionalConfigProperties = local.config_properties } : {},
    length(try(var.spec.fault_tolerant_execution.exchange_manager.additional_properties, [])) > 0 ? {
      additionalExchangeManagerProperties = var.spec.fault_tolerant_execution.exchange_manager.additional_properties
    } : {},
    local.auth_enabled ? { auth = local.auth_block } : {},
    local.https_enabled ? {
      secretMounts = [{
        name       = "trino-keystore"
        secretName = var.spec.https.keystore_secret.secret_name
        path       = "/etc/trino/keystore"
      }]
    } : {},
    length(local.catalogs_block) > 0 ? { catalogs = local.catalogs_block } : {},
    length(local.env_entries) > 0 ? { env = local.env_entries } : {},
    try(var.spec.access_control_rules, "") != "" ? {
      accessControl = {
        type       = "configmap"
        configFile = "rules.json"
        rules = {
          "rules.json" = var.spec.access_control_rules
        }
      }
    } : {},
    try(var.spec.resource_groups_config, "") != "" ? {
      resourceGroups = {
        type                 = "configmap"
        resourceGroupsConfig = var.spec.resource_groups_config
      }
    } : {},
    try(var.spec.session_properties_config, "") != "" ? {
      sessionProperties = {
        type                    = "configmap"
        sessionPropertiesConfig = var.spec.session_properties_config
      }
    } : {},
    length(try(var.spec.event_listener_properties, [])) > 0 ? { eventListenerProperties = var.spec.event_listener_properties } : {},
    local.jmx_block,
    local.service_monitor_enabled ? { serviceMonitor = { enabled = true } } : {},
    try(var.spec.network_policy_enabled, false) ? { networkPolicy = { enabled = true } } : {},
    length(try(var.spec.image_pull_secrets, [])) > 0 ? {
      imagePullSecrets = [for s in var.spec.image_pull_secrets : { name = s }]
    } : {},
  )

  # ------------------------------- outputs --------------------------------
  coordinator_service  = "${local.release_name}-coordinator"
  worker_service       = "${local.release_name}-worker"
  coordinator_endpoint = "http://${local.coordinator_service}.${local.namespace}.svc.cluster.local:8080"
  port_forward_command = "kubectl port-forward svc/${local.coordinator_service} -n ${local.namespace} 8080:8080"

  admin_password_secret_output_name = local.module_owned_password_db ? local.password_db_secret_name : ""
  admin_password_secret_output_key  = local.module_owned_password_db ? "password" : ""
}
