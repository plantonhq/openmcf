# Every resolution here has an exact twin in the Pulumi module's
# locals.go / values.go — keep them in lockstep.
#
# HCL DISCIPLINE (the type-unification class): every conditional block
# below is a SINGLE-ATTRIBUTE ternary against {} or rides the null-prune
# merge idiom; the env-Secret data map merges conditionally-present
# string-valued maps only.

locals {
  # Chart identity — must stay byte-identical with the Pulumi module's
  # vars (cross-engine chart-name drift deploys two different products
  # from one manifest).
  helm_chart_name    = "superset"
  helm_chart_repo    = "https://apache.github.io/superset"
  helm_chart_version = "0.22.4"

  release_name = var.metadata.name
  namespace    = var.spec.namespace

  labels = merge(
    {
      "planton.ai/resource"      = "true"
      "planton.ai/resource-name" = var.metadata.name
      "planton.ai/resource-kind" = "KubernetesSuperset"
    },
    try(var.metadata.id, "") != "" ? { "planton.ai/resource-id" = var.metadata.id } : {},
    try(var.metadata.org, "") != "" ? { "planton.ai/organization" = var.metadata.org } : {},
    try(var.metadata.env, "") != "" ? { "planton.ai/environment" = var.metadata.env } : {},
  )

  # ----------------------------- name budget ------------------------------
  # Chart truth at the pin: the longest derived child is
  # `<name>-celerybeat` (11-char suffix) — checked fail-loud in
  # main.tf preconditions (Pulumi twin: buildHelmValues errors).
  name_budget = 52

  # The module-owned environment Secret — the chart's runtime-credential
  # contract (every component envFroms it; the chart's own copy is OFF).
  env_secret_name = "${local.release_name}-env"

  config_mount_path = "/app/pythonpath"

  # --------------------------- metadata database --------------------------
  db_host = var.spec.metadata_database.host
  db_port = coalesce(try(var.spec.metadata_database.port, null), 5432)
  db_user = coalesce(try(var.spec.metadata_database.username, null), "superset")
  db_name = coalesce(try(var.spec.metadata_database.database_name, null), "superset")

  db_password_secret = var.spec.metadata_database.password_secret.secret_name
  db_password_key    = coalesce(try(var.spec.metadata_database.password_secret.secret_key, null), "password")

  db_ssl_enabled = try(var.spec.metadata_database.ssl.enabled, false)
  db_ssl_mode    = coalesce(try(var.spec.metadata_database.ssl.mode, null), "require")

  # -------------------------------- cache ---------------------------------
  cache_enabled = try(var.spec.cache, null) != null
  cache_host    = try(var.spec.cache.host, "")
  cache_port    = local.cache_enabled ? coalesce(try(var.spec.cache.port, null), 6379) : 6379
  cache_user    = try(var.spec.cache.username, "")

  cache_password_secret = try(var.spec.cache.password_secret.secret_name, "")
  cache_password_key    = coalesce(try(var.spec.cache.password_secret.secret_key, null), "password")

  # ------------------------------ SECRET_KEY -------------------------------
  secret_key_byo          = try(var.spec.secret_key_secret, null) != null
  secret_key_module_owned = !local.secret_key_byo
  secret_key_secret_name  = local.secret_key_byo ? var.spec.secret_key_secret.secret_name : "${local.release_name}-secret-key"

  # ------------------------------ admin user -------------------------------
  admin_username             = coalesce(try(var.spec.init.admin.username, null), "admin")
  admin_email                = coalesce(try(var.spec.init.admin.email, null), "admin@superset.local")
  admin_byo                  = try(var.spec.init.admin.password_secret, null) != null
  admin_module_owned         = !local.admin_byo
  admin_password_secret_name = local.admin_byo ? var.spec.init.admin.password_secret.secret_name : "${local.release_name}-admin-auth"
  admin_password_secret_key  = local.admin_byo ? var.spec.init.admin.password_secret.secret_key : "password"

  # -------------------------- component toggles ---------------------------
  worker_enabled     = local.cache_enabled && (try(var.spec.worker.enabled, null) == null ? true : var.spec.worker.enabled)
  beat_enabled       = try(var.spec.beat.enabled, false)
  flower_enabled     = try(var.spec.flower.enabled, false)
  websockets_enabled = try(var.spec.websockets.enabled, false)
  mcp_enabled        = try(var.spec.mcp.enabled, false)

  # -------------------- environment-Secret plain keys ---------------------
  # NON-SECRET connection facts + the admin identity. The GENERATED
  # material (SECRET_KEY, ADMIN_PASSWORD, the ws JWT) joins in main.tf
  # where the randoms live; REFERENCED material rides extraEnvRaw.
  env_plain = merge(
    {
      DB_HOST        = local.db_host
      DB_PORT        = tostring(local.db_port)
      DB_USER        = local.db_user
      DB_NAME        = local.db_name
      ADMIN_USERNAME = local.admin_username
      ADMIN_EMAIL    = local.admin_email
    },
    local.db_ssl_enabled ? { DB_SSL_MODE = local.db_ssl_mode } : {},
    local.cache_enabled ? {
      REDIS_HOST = local.cache_host
      REDIS_PORT = tostring(local.cache_port)
    } : {},
    local.cache_enabled && local.cache_user != "" ? { REDIS_USER = local.cache_user } : {},
  )

  # ---------------------- extraEnvRaw (references) -------------------------
  # Explicit env entries override the envFrom Secret — the chart's own
  # bring-your-own-credential mechanism. Secret NAMES only ever render.
  extra_env_raw = concat(
    [{
      name = "DB_PASS"
      valueFrom = {
        secretKeyRef = {
          name = local.db_password_secret
          key  = local.db_password_key
        }
      }
    }],
    local.cache_password_secret != "" ? [{
      name = "REDIS_PASSWORD"
      valueFrom = {
        secretKeyRef = {
          name = local.cache_password_secret
          key  = local.cache_password_key
        }
      }
    }] : [],
    local.secret_key_byo ? [{
      name = "SUPERSET_SECRET_KEY"
      valueFrom = {
        secretKeyRef = {
          name = var.spec.secret_key_secret.secret_name
          key  = var.spec.secret_key_secret.secret_key
        }
      }
    }] : [],
    local.admin_byo ? [{
      name = "ADMIN_PASSWORD"
      valueFrom = {
        secretKeyRef = {
          name = var.spec.init.admin.password_secret.secret_name
          key  = var.spec.init.admin.password_secret.secret_key
        }
      }
    }] : [],
    [for name, ref in try(var.spec.extra_env_from_secret, {}) : {
      name = name
      valueFrom = {
        secretKeyRef = {
          name = ref.secret_name
          key  = ref.secret_key
        }
      }
    }],
  )

  # ------------------------- init command override -------------------------
  # The chart's own rendered init script (schema migration + role init —
  # createAdmin=false keeps it admin-free) followed by an idempotent
  # create-admin step reading the admin identity FROM ENVIRONMENT — the
  # chart's literal-password rendering path is never used. Keep
  # byte-identical with the Pulumi module's initCommand.
  init_command = ["/bin/sh", "-c", join("; ", [
    ". ${local.config_mount_path}/superset_bootstrap.sh",
    ". ${local.config_mount_path}/superset_init.sh",
    "if superset fab list-users 2>/dev/null | grep -qF \"username:$${ADMIN_USERNAME}\"; then echo \"Admin user already exists, skipping.\"; else superset fab create-admin --username \"$ADMIN_USERNAME\" --firstname Superset --lastname Admin --email \"$ADMIN_EMAIL\" --password \"$ADMIN_PASSWORD\"; fi",
  ])]

  # ------------------- module-owned configOverrides ------------------------
  # Env-indirection replacements for the chart's password-from-values
  # config blocks. Keep byte-identical with the Pulumi module's
  # moduleConfigOverrides.
  module_config_overrides = merge(
    local.cache_enabled && local.cache_password_secret != "" ? {
      planton_redis_auth = join("\n", [
        "# Authed cache: the chart-rendered RESULTS_BACKEND and async-queries",
        "# backends carry no password (it never rides values) — redefine them",
        "# reading the environment the pods already carry.",
        "from flask_caching.backends.rediscache import RedisCache as _PlantonRedisCache",
        "RESULTS_BACKEND = _PlantonRedisCache(",
        "    host=env('REDIS_HOST'),",
        "    port=int(env('REDIS_PORT', '6379')),",
        "    password=env('REDIS_PASSWORD') or None,",
        "    key_prefix='superset_results',",
        ")",
        "GLOBAL_ASYNC_QUERIES_CACHE_BACKEND = {",
        "    'CACHE_TYPE': 'RedisCache',",
        "    'CACHE_REDIS_HOST': env('REDIS_HOST'),",
        "    'CACHE_REDIS_PORT': int(env('REDIS_PORT', '6379')),",
        "    'CACHE_REDIS_PASSWORD': env('REDIS_PASSWORD', ''),",
        "    'CACHE_REDIS_DB': int(env('REDIS_DB', '1')),",
        "    'CACHE_KEY_PREFIX': 'qc-',",
        "    'CACHE_DEFAULT_TIMEOUT': 86400,",
        "}",
        "GLOBAL_ASYNC_QUERIES_RESULTS_BACKEND = {",
        "    'backend': 'redis',",
        "    'host': env('REDIS_HOST'),",
        "    'port': int(env('REDIS_PORT', '6379')),",
        "    'password': env('REDIS_PASSWORD', ''),",
        "    'db': int(env('REDIS_DB', '1')),",
        "    'prefix': 'qc-',",
        "}",
      ])
    } : {},
    local.websockets_enabled ? {
      planton_ws_jwt = join("\n", [
        "# The async-queries JWT shared with the websocket server — both sides",
        "# read the same module-generated environment variable.",
        "GLOBAL_ASYNC_QUERIES_JWT_SECRET = env('JWT_SECRET')",
      ])
    } : {},
  )

  config_overrides = merge(local.module_config_overrides, try(var.spec.config_overrides, {}))

  # ----------------------------- components -------------------------------
  # Single-attribute ternaries only (the HCL type-unification class): a
  # branch-shaped conditional (autoscaling XOR replicas) splits into two
  # complementary ternaries against {}.
  web_block = merge(
    try(var.spec.web.hpa, null) != null ? {
      autoscaling = merge(
        {
          enabled     = true
          maxReplicas = var.spec.web.hpa.max_replicas
        },
        try(var.spec.web.hpa.min_replicas, null) != null ? { minReplicas = var.spec.web.hpa.min_replicas } : {},
        try(var.spec.web.hpa.target_cpu_utilization_percent, null) != null ? { targetCPUUtilizationPercentage = var.spec.web.hpa.target_cpu_utilization_percent } : {},
      )
    } : {},
    try(var.spec.web.hpa, null) == null ? {
      replicas = {
        enabled      = true
        replicaCount = coalesce(try(var.spec.web.replicas, null), 1)
      }
    } : {},
    local.web_resources_block,
  )

  web_resources = {
    requests_cpu    = try(var.spec.web.resources.requests.cpu, "")
    requests_memory = try(var.spec.web.resources.requests.memory, "")
    limits_cpu      = try(var.spec.web.resources.limits.cpu, "")
    limits_memory   = try(var.spec.web.resources.limits.memory, "")
  }
  web_resources_block = (
    local.web_resources.requests_cpu != "" || local.web_resources.requests_memory != "" ||
    local.web_resources.limits_cpu != "" || local.web_resources.limits_memory != ""
    ) ? {
    resources = merge(
      local.web_resources.requests_cpu != "" || local.web_resources.requests_memory != "" ? {
        requests = merge(
          local.web_resources.requests_cpu != "" ? { cpu = local.web_resources.requests_cpu } : {},
          local.web_resources.requests_memory != "" ? { memory = local.web_resources.requests_memory } : {},
        )
      } : {},
      local.web_resources.limits_cpu != "" || local.web_resources.limits_memory != "" ? {
        limits = merge(
          local.web_resources.limits_cpu != "" ? { cpu = local.web_resources.limits_cpu } : {},
          local.web_resources.limits_memory != "" ? { memory = local.web_resources.limits_memory } : {},
        )
      } : {},
    )
  } : {}

  worker_resources = {
    requests_cpu    = try(var.spec.worker.resources.requests.cpu, "")
    requests_memory = try(var.spec.worker.resources.requests.memory, "")
    limits_cpu      = try(var.spec.worker.resources.limits.cpu, "")
    limits_memory   = try(var.spec.worker.resources.limits.memory, "")
  }
  worker_resources_block = (
    local.worker_resources.requests_cpu != "" || local.worker_resources.requests_memory != "" ||
    local.worker_resources.limits_cpu != "" || local.worker_resources.limits_memory != ""
    ) ? {
    resources = merge(
      local.worker_resources.requests_cpu != "" || local.worker_resources.requests_memory != "" ? {
        requests = merge(
          local.worker_resources.requests_cpu != "" ? { cpu = local.worker_resources.requests_cpu } : {},
          local.worker_resources.requests_memory != "" ? { memory = local.worker_resources.requests_memory } : {},
        )
      } : {},
      local.worker_resources.limits_cpu != "" || local.worker_resources.limits_memory != "" ? {
        limits = merge(
          local.worker_resources.limits_cpu != "" ? { cpu = local.worker_resources.limits_cpu } : {},
          local.worker_resources.limits_memory != "" ? { memory = local.worker_resources.limits_memory } : {},
        )
      } : {},
    )
  } : {}

  # No cache (or explicitly disabled): the worker Deployment never
  # renders — it would crash-loop without a broker. Complementary
  # single-attribute ternaries: exactly one `replicas`/`autoscaling`
  # arm is non-empty.
  worker_block = merge(
    local.worker_enabled && try(var.spec.worker.hpa, null) != null ? {
      autoscaling = merge(
        {
          enabled     = true
          maxReplicas = var.spec.worker.hpa.max_replicas
        },
        try(var.spec.worker.hpa.min_replicas, null) != null ? { minReplicas = var.spec.worker.hpa.min_replicas } : {},
        try(var.spec.worker.hpa.target_cpu_utilization_percent, null) != null ? { targetCPUUtilizationPercentage = var.spec.worker.hpa.target_cpu_utilization_percent } : {},
      )
    } : {},
    local.worker_enabled && try(var.spec.worker.hpa, null) == null ? {
      replicas = {
        enabled      = true
        replicaCount = coalesce(try(var.spec.worker.replicas, null), 1)
      }
    } : {},
    local.worker_enabled ? local.worker_resources_block : {},
    !local.worker_enabled ? { replicas = { enabled = false } } : {},
  )

  # The two arms carry different attribute sets — merge of
  # single-attribute ternaries (the jsonencode seam does NOT rescue a
  # ternary: plan-time constant folding re-derives concrete object
  # types through jsondecode).
  websockets_block = merge(
    { enabled = local.websockets_enabled },
    local.websockets_enabled ? {
      image = {
        repository = coalesce(try(var.spec.websockets.image.repository, null), "oneacrefund/superset-websocket")
        tag        = coalesce(try(var.spec.websockets.image.tag, null), "latest")
        pullPolicy = "IfNotPresent"
      }
    } : {},
    local.websockets_enabled ? { replicaCount = coalesce(try(var.spec.websockets.replicas, null), 1) } : {},
  )

  tolerations = [for t in try(var.spec.scheduling.tolerations, []) : jsondecode(jsonencode(merge(
    t.key != "" ? { key = t.key } : {},
    t.operator != "" ? { operator = t.operator } : {},
    t.value != "" ? { value = t.value } : {},
    t.effect != "" ? { effect = t.effect } : {},
    try(t.toleration_seconds, null) != null ? { tolerationSeconds = t.toleration_seconds } : {},
  )))]

  # ---------------------------- typed values -------------------------------
  helm_values = merge(
    {
      # Deterministic child names (`<name>`, `<name>-worker`, the
      # `<name>-env`/`<name>-config` Secrets) — the release name never
      # double-prefixes and the import map stays exact.
      fullnameOverride = local.release_name
      # The module composes the environment Secret — the chart's copy
      # (which renders credentials from values) never ships.
      secretEnv     = { create = false }
      envFromSecret = local.env_secret_name
      # The bundled subcharts ride frozen legacy image lines and never
      # ship from this kind — the metadata database and the cache are
      # ALWAYS external (composition-first).
      postgresql = { enabled = false }
      redis      = { enabled = false }

      image = {
        repository = coalesce(try(var.spec.image.repository, null), "apachesuperset.docker.scarf.sh/apache/superset")
        tag        = coalesce(try(var.spec.image.tag, null), "6.1.0")
      }

      # NON-SECRET connection facts only — the password never rides
      # values (DB_PASS arrives via extraEnvRaw); the ssl block feeds
      # the rendered config's sslmode parameters.
      database = merge(
        {
          host = local.db_host
          port = tostring(local.db_port)
          user = local.db_user
          name = local.db_name
        },
        local.db_ssl_enabled ? {
          ssl = {
            enabled = true
            mode    = local.db_ssl_mode
          }
        } : {},
      )

      # Merge of single-attribute ternaries — the enabled/disabled arms
      # carry different attribute sets (the type-unification class).
      cache = merge(
        { enabled = local.cache_enabled },
        local.cache_enabled ? { host = local.cache_host } : {},
        local.cache_enabled ? { port = tostring(local.cache_port) } : {},
        local.cache_enabled && local.cache_user != "" ? { user = local.cache_user } : {},
        local.cache_enabled && try(var.spec.cache.cache_db, null) != null ? { cacheDb = var.spec.cache.cache_db } : {},
        local.cache_enabled && try(var.spec.cache.celery_db, null) != null ? { celeryDb = var.spec.cache.celery_db } : {},
      )

      supersetNode         = local.web_block
      supersetWorker       = local.worker_block
      supersetCeleryBeat   = { enabled = local.beat_enabled }
      supersetCeleryFlower = { enabled = local.flower_enabled }
      supersetWebsockets   = local.websockets_block
      supersetMcp          = { enabled = local.mcp_enabled }

      # createAdmin stays FALSE so the chart's literal-password
      # rendering path (and its config-template fail gate) never
      # engage; the module's command override appends the
      # create-admin-from-env step.
      init = {
        createAdmin  = false
        loadExamples = try(var.spec.init.load_examples, false)
        command      = local.init_command
      }

      extraEnvRaw = jsondecode(jsonencode(local.extra_env_raw))

      service = merge(
        {
          type = coalesce(try(var.spec.service.type, null), "ClusterIP")
          port = 8088
        },
        length(try(var.spec.service.annotations, {})) > 0 ? { annotations = var.spec.service.annotations } : {},
      )
    },
    length(local.config_overrides) > 0 ? { configOverrides = local.config_overrides } : {},
    try(var.spec.bootstrap_script, "") != "" ? { bootstrapScript = var.spec.bootstrap_script } : {},
    length(try(var.spec.extra_env, {})) > 0 ? { extraEnv = var.spec.extra_env } : {},
    length(try(var.spec.feature_flags, {})) > 0 ? { featureFlags = var.spec.feature_flags } : {},
    length(try(var.spec.scheduling.node_selector, {})) > 0 ? { nodeSelector = var.spec.scheduling.node_selector } : {},
    length(local.tolerations) > 0 ? { tolerations = local.tolerations } : {},
    length(try(var.spec.image_pull_secrets, [])) > 0 ? {
      imagePullSecrets = [for s in var.spec.image_pull_secrets : { name = s }]
    } : {},
  )

  # ------------------------------- outputs --------------------------------
  service              = local.release_name
  endpoint             = "http://${local.service}.${local.namespace}.svc.cluster.local:8088"
  port_forward_command = "kubectl port-forward svc/${local.service} -n ${local.namespace} 8088:8088"
}
