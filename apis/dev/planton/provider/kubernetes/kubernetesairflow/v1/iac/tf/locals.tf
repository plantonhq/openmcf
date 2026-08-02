# Every resolution here has an exact twin in the Pulumi module's
# locals.go / values.go — keep them in lockstep. Conditional blocks use
# the null-prune idiom ({for k,v in {...} : k => v if v != null}) —
# NEVER a ternary between two differently-shaped objects (the HCL
# type-unification class).

locals {
  # Chart identity — byte-identical with the Pulumi module's vars.
  helm_chart_name = "airflow"
  helm_chart_repo = "https://airflow.apache.org"

  release_name = var.metadata.name
  namespace    = var.spec.namespace

  # Resource-identity labels for the module-created satellites (the
  # namespace and module-owned Secrets — never injected into the
  # chart's own resources; Helm owns those).
  labels = merge(
    {
      "planton.ai/resource"      = "true"
      "planton.ai/resource-name" = var.metadata.name
      "planton.ai/resource-kind" = "KubernetesAirflow"
    },
    var.metadata.id != "" ? { "planton.ai/resource-id" = var.metadata.id } : {},
    var.metadata.org != "" ? { "planton.ai/organization" = var.metadata.org } : {},
    var.metadata.env != "" ? { "planton.ai/environment" = var.metadata.env } : {}
  )

  chart_version   = try(coalesce(var.spec.chart_version), null) != null ? var.spec.chart_version : "1.22.0"
  airflow_version = try(coalesce(var.spec.airflow_version), null) != null ? var.spec.airflow_version : "3.2.2"

  # This kind's default is the Kubernetes-native executor (the chart's
  # own default is CeleryExecutor); the Celery test mirrors the chart's
  # check-values substring pairing.
  executor       = try(coalesce(var.spec.executor), null) != null ? var.spec.executor : "KubernetesExecutor"
  celery_enabled = strcontains(local.executor, "CeleryExecutor") || strcontains(local.executor, "CeleryKubernetesExecutor")

  # ---- database resolution --------------------------------------------------
  db_is_postgres = try(var.spec.database.postgres, null) != null
  db_protocol    = local.db_is_postgres ? "postgresql" : "mysql"
  db_host        = local.db_is_postgres ? var.spec.database.postgres.host : var.spec.database.mysql.host
  db_port = local.db_is_postgres ? (
    try(coalesce(var.spec.database.postgres.port), null) != null ? var.spec.database.postgres.port : 5432
    ) : (
    try(coalesce(var.spec.database.mysql.port), null) != null ? var.spec.database.mysql.port : 3306
  )
  db_name = local.db_is_postgres ? (
    try(coalesce(var.spec.database.postgres.database_name), null) != null ? var.spec.database.postgres.database_name : "airflow"
    ) : (
    try(coalesce(var.spec.database.mysql.database_name), null) != null ? var.spec.database.mysql.database_name : "airflow"
  )
  db_user = local.db_is_postgres ? (
    try(coalesce(var.spec.database.postgres.username), null) != null ? var.spec.database.postgres.username : "airflow"
    ) : (
    try(coalesce(var.spec.database.mysql.username), null) != null ? var.spec.database.mysql.username : "airflow"
  )
  db_ssl_mode = local.db_is_postgres ? (
    try(coalesce(var.spec.database.postgres.ssl_mode), null) != null ? var.spec.database.postgres.ssl_mode : "disable"
  ) : ""
  db_password_secret = local.db_is_postgres ? var.spec.database.postgres.password_secret.secret_name : var.spec.database.mysql.password_secret.secret_name
  db_password_secret_key = local.db_is_postgres ? (
    try(coalesce(var.spec.database.postgres.password_secret.secret_key), null) != null ? var.spec.database.postgres.password_secret.secret_key : "password"
    ) : (
    try(coalesce(var.spec.database.mysql.password_secret.secret_key), null) != null ? var.spec.database.mysql.password_secret.secret_key : "password"
  )

  # ---- pgbouncer routing ------------------------------------------------------
  # Mirrors the chart's own connection rewriting EXACTLY: with the
  # pooler on, Airflow dials `<name>-pgbouncer:6543` and the DATABASE
  # segment becomes the pgbouncer.ini alias, which the ini maps back to
  # the real host/database.
  pgbouncer_enabled           = try(var.spec.pgbouncer.enabled, false)
  pgbouncer_port              = 6543
  effective_db_host           = local.pgbouncer_enabled ? "${local.release_name}-pgbouncer" : local.db_host
  effective_db_port           = local.pgbouncer_enabled ? local.pgbouncer_port : local.db_port
  effective_metadata_db       = local.pgbouncer_enabled ? "${local.release_name}-metadata" : local.db_name
  effective_result_backend_db = local.pgbouncer_enabled ? "${local.release_name}-result-backend" : local.db_name
  pgbouncer_config_secret     = "${local.release_name}-pgbouncer-config"

  # ---- broker resolution ---------------------------------------------------------
  bundled_redis_enabled = local.celery_enabled && try(var.spec.broker.bundled_redis, null) != null
  valkey_broker         = local.celery_enabled ? try(var.spec.broker.valkey, null) : null
  broker_host           = local.bundled_redis_enabled ? "${local.release_name}-redis" : try(local.valkey_broker.host, "")
  broker_port = local.bundled_redis_enabled ? 6379 : (
    try(coalesce(local.valkey_broker.port), null) != null ? local.valkey_broker.port : 6379
  )
  broker_user = try(local.valkey_broker.username, "")
  broker_db = local.bundled_redis_enabled ? 0 : (
    try(coalesce(local.valkey_broker.database_number), null) != null ? local.valkey_broker.database_number : 0
  )
  broker_password_secret         = try(local.valkey_broker.password_secret.secret_name, "")
  broker_password_secret_key     = try(coalesce(try(local.valkey_broker.password_secret.secret_key, null)), null) != null ? local.valkey_broker.password_secret.secret_key : "password"
  broker_url_secret_module_owned = local.bundled_redis_enabled || local.valkey_broker != null
  broker_url_secret_name = local.celery_enabled ? (
    local.broker_url_secret_module_owned ? "${local.release_name}-broker-url" : var.spec.broker.existing_broker_url_secret.secret_name
  ) : ""
  redis_password_secret_name = "${local.release_name}-redis-password"

  # ---- log read path ---------------------------------------------------------------
  log_backend_config              = try(var.spec.logging.elasticsearch, null) != null ? var.spec.logging.elasticsearch : try(var.spec.logging.opensearch, null)
  log_backend                     = try(var.spec.logging.elasticsearch, null) != null ? "elasticsearch" : (try(var.spec.logging.opensearch, null) != null ? "opensearch" : "")
  log_backend_scheme              = try(coalesce(try(local.log_backend_config.scheme, null)), null) != null ? local.log_backend_config.scheme : "http"
  log_backend_port                = try(coalesce(try(local.log_backend_config.port, null)), null) != null ? local.log_backend_config.port : 9200
  log_backend_user                = try(local.log_backend_config.username, "")
  log_backend_password_secret     = try(local.log_backend_config.password_secret.secret_name, "")
  log_backend_password_secret_key = try(coalesce(try(local.log_backend_config.password_secret.secret_key, null)), null) != null ? local.log_backend_config.password_secret.secret_key : "password"
  log_read_conn_secret_name       = "${local.release_name}-log-read-conn"

  # ---- module-owned secret names (BYO respected) --------------------------------------
  metadata_conn_secret_name       = "${local.release_name}-metadata-conn"
  result_backend_conn_secret_name = "${local.release_name}-result-backend-conn"

  fernet_key_secret_byo     = try(coalesce(try(var.spec.security.fernet_key_secret_name, null)), null) != null && try(var.spec.security.fernet_key_secret_name, "") != ""
  fernet_key_secret_name    = local.fernet_key_secret_byo ? var.spec.security.fernet_key_secret_name : "${local.release_name}-fernet-key"
  api_secret_key_byo        = try(var.spec.security.api_secret_key_secret_name, "") != ""
  api_secret_key_name       = local.api_secret_key_byo ? var.spec.security.api_secret_key_secret_name : "${local.release_name}-api-secret-key"
  webserver_secret_key_name = "${local.release_name}-webserver-secret-key"
  jwt_secret_byo            = try(var.spec.security.jwt_secret_name, "") != ""
  jwt_secret_name           = local.jwt_secret_byo ? var.spec.security.jwt_secret_name : "${local.release_name}-jwt-secret"

  # ---- admin bootstrap user -----------------------------------------------------------
  admin_create      = try(coalesce(try(var.spec.admin_user.create, null)), null) != null ? var.spec.admin_user.create : true
  admin_username    = try(coalesce(try(var.spec.admin_user.username, null)), null) != null ? var.spec.admin_user.username : "admin"
  admin_email       = try(coalesce(try(var.spec.admin_user.email, null)), null) != null ? var.spec.admin_user.email : "admin@example.com"
  admin_secret_byo  = try(var.spec.admin_user.password_secret, null) != null
  admin_secret_name = local.admin_secret_byo ? var.spec.admin_user.password_secret.secret_name : "${local.release_name}-admin-auth"
  admin_secret_key = local.admin_secret_byo ? (
    try(coalesce(var.spec.admin_user.password_secret.secret_key), null) != null ? var.spec.admin_user.password_secret.secret_key : "password"
  ) : "password"
  admin_secret_module_owned = local.admin_create && !local.admin_secret_byo
  # The exported credential handle is honest: with the bootstrap user
  # disabled no admin credential exists, so the handle exports EMPTY
  # rather than a name that points at nothing (Pulumi twin exports the
  # same empties).
  admin_password_secret_output_name = local.admin_create ? local.admin_secret_name : ""
  admin_password_secret_output_key  = local.admin_create ? local.admin_secret_key : ""

  # ---- composed connection URIs (sensitive: password from the data read) ---------------
  # userinfo segments are urlencoded — the chart's own urlquery
  # treatment; sslmode rides the query only on postgresql (chart
  # parity).
  db_password       = data.kubernetes_secret_v1.db_password.data[local.db_password_secret_key]
  db_query          = local.db_protocol == "postgresql" ? "?sslmode=${local.db_ssl_mode}" : ""
  metadata_conn_uri = "${local.db_protocol}://${urlencode(local.db_user)}:${urlencode(local.db_password)}@${local.effective_db_host}:${local.effective_db_port}/${local.effective_metadata_db}${local.db_query}"
  # Celery's result backend needs SQLAlchemy's `db+` scheme prefix (the
  # chart's own result-backend template renders exactly this).
  result_backend_uri = "db+${local.db_protocol}://${urlencode(local.db_user)}:${urlencode(local.db_password)}@${local.effective_db_host}:${local.effective_db_port}/${local.effective_result_backend_db}${local.db_query}"

  # The chart's KEDA autoscalers read env KEDA_DB_CONN from THIS
  # Secret's `kedaConnection` key whenever the trigger cannot ride the
  # normal `connection` URI: always on mysql (KEDA's mysql scaler wants
  # the Go-DSN `user:pass@tcp(host:port)/db` form, never a URI) and on
  # the postgres pgbouncer-BYPASS posture (the scaler dials the real
  # database while Airflow rides the pooler). The chart gates the key on
  # KEDA being enabled; the module renders it whenever it COULD be
  # needed — the extra key is inert otherwise (nothing else reads this
  # module-owned Secret, and it carries the same credential material as
  # `connection`), and escape-hatch KEDA configs (triggerer.keda,
  # workers.keda.usePgbouncer=false) then work without a spec change.
  keda_conn_needed = !local.db_is_postgres || local.pgbouncer_enabled
  keda_conn_uri = local.db_is_postgres ? (
    # The DIRECT postgres URI (real host and database — bypassing the
    # pooler), mirroring the chart's own kedaConnection rendering.
    "postgresql://${urlencode(local.db_user)}:${urlencode(local.db_password)}@${local.db_host}:${local.db_port}/${local.db_name}${local.db_query}"
    ) : (
    # KEDA's mysql scaler DSN form (the chart urlJoins then trims the
    # scheme) — userinfo urlencoded, chart parity.
    "${urlencode(local.db_user)}:${urlencode(local.db_password)}@tcp(${local.db_host}:${local.db_port})/${local.db_name}"
  )

  # ---- pgbouncer config (module-composed; byte-faithful to the chart's helper) --------
  pgbouncer_metadata_pool_size       = try(coalesce(try(var.spec.pgbouncer.metadata_pool_size, null)), null) != null ? var.spec.pgbouncer.metadata_pool_size : 10
  pgbouncer_result_backend_pool_size = try(coalesce(try(var.spec.pgbouncer.result_backend_pool_size, null)), null) != null ? var.spec.pgbouncer.result_backend_pool_size : 5
  pgbouncer_max_client_conn          = try(coalesce(try(var.spec.pgbouncer.max_client_connections, null)), null) != null ? var.spec.pgbouncer.max_client_connections : 100

  pgbouncer_ini = <<-EOT
    [databases]
    ${local.release_name}-metadata = host=${local.db_host} dbname=${local.db_name} port=${local.db_port} pool_size=${local.pgbouncer_metadata_pool_size}
    ${local.release_name}-result-backend = host=${local.db_host} dbname=${local.db_name} port=${local.db_port} pool_size=${local.pgbouncer_result_backend_pool_size}

    [pgbouncer]
    pool_mode = transaction
    listen_port = ${local.pgbouncer_port}
    listen_addr = *
    auth_type = scram-sha-256
    auth_file = /etc/pgbouncer/users.txt
    stats_users = ${local.db_user}
    ignore_startup_parameters = extra_float_digits
    max_client_conn = ${local.pgbouncer_max_client_conn}
    verbose = 0
    log_disconnections = 0
    log_connections = 0

    server_tls_sslmode = prefer
    server_tls_ciphers = normal
  EOT

  # Two lines mirror the chart's own users.txt helper (the metadata and
  # result-backend users) — identical here because one declared user
  # serves both connections. %q escapes embedded quotes/backslashes the
  # same way the chart's `quote` (and the Pulumi twin's %q) does.
  pgbouncer_users_txt = format("%q %q\n%q %q\n",
    local.db_user, local.db_password,
    local.db_user, local.db_password,
  )

  # The pgbouncer metrics-exporter sidecar (unconditional whenever the
  # pooler deploys at this pin) reads its DSN from a stats Secret. The
  # chart's OWN stats Secret composes that DSN from split values
  # (data.metadataConnection.user/pass — defaults postgres/postgres),
  # which the secret-native design never sets: the exporter then dials
  # a user pgbouncer's auth_file does not carry and crash-loops
  # (verified live: "no such user: postgres" in the pooler log). The
  # module composes the stats DSN itself — the metadata user IS the ini
  # stats_users grant — through the chart's documented
  # statsSecretName seam (DSN shape from the chart's own Secret
  # template: 127.0.0.1:<port>/pgbouncer, sslmode disable default).
  pgbouncer_stats_secret = "${local.release_name}-pgbouncer-stats"
  pgbouncer_stats_uri    = "postgresql://${urlencode(local.db_user)}:${urlencode(local.db_password)}@127.0.0.1:${local.pgbouncer_port}/pgbouncer?sslmode=disable"

  # ---- resources helper renderings ----------------------------------------------------
  component_resources = {
    api_server    = try(var.spec.components.api_server.resources, null)
    scheduler     = try(var.spec.components.scheduler.resources, null)
    dag_processor = try(var.spec.components.dag_processor.resources, null)
    triggerer     = try(var.spec.components.triggerer.resources, null)
    workers       = try(var.spec.components.workers.resources, null)
    pgbouncer     = try(var.spec.pgbouncer.resources, null)
    bundled_redis = try(var.spec.broker.bundled_redis.resources, null)
    git_sync      = try(var.spec.dags.git_sync.resources, null)
  }

  rendered_resources = {
    for name, r in local.component_resources : name => r == null ? null : {
      for rk, rv in {
        requests = try(r.requests, null) == null ? null : {
          for qk, qv in {
            cpu    = r.requests.cpu != "" ? r.requests.cpu : null
            memory = r.requests.memory != "" ? r.requests.memory : null
          } : qk => qv if qv != null
        }
        limits = try(r.limits, null) == null ? null : {
          for lk, lv in {
            cpu    = r.limits.cpu != "" ? r.limits.cpu : null
            memory = r.limits.memory != "" ? r.limits.memory : null
          } : lk => lv if lv != null
        }
      } : rk => rv if rv != null && rv != {}
    }
  }

  # ---- component blocks ---------------------------------------------------------------
  api_server_block = {
    for k, v in {
      replicas  = try(coalesce(try(var.spec.components.api_server.replicas, null)), null)
      resources = local.rendered_resources.api_server
    } : k => v if v != null
  }

  scheduler_block = {
    for k, v in {
      replicas  = try(coalesce(try(var.spec.components.scheduler.replicas, null)), null)
      resources = local.rendered_resources.scheduler
    } : k => v if v != null
  }

  dag_processor_block = {
    for k, v in {
      replicas  = try(coalesce(try(var.spec.components.dag_processor.replicas, null)), null)
      resources = local.rendered_resources.dag_processor
    } : k => v if v != null
  }

  triggerer_persistence = {
    for k, v in {
      size = try(coalesce(try(var.spec.components.triggerer.persistence_size, null)), null)
    } : k => v if v != null
  }
  triggerer_block = {
    for k, v in {
      enabled     = try(coalesce(try(var.spec.components.triggerer.enabled, null)), null)
      replicas    = try(coalesce(try(var.spec.components.triggerer.replicas, null)), null)
      persistence = length(local.triggerer_persistence) > 0 ? local.triggerer_persistence : null
      resources   = local.rendered_resources.triggerer
    } : k => v if v != null
  }

  workers_keda = try(var.spec.components.workers.keda, null)
  workers_keda_block = try(local.workers_keda.enabled, false) ? {
    for k, v in {
      enabled         = true
      minReplicaCount = try(coalesce(try(local.workers_keda.min_replicas, null)), null)
      maxReplicaCount = try(coalesce(try(local.workers_keda.max_replicas, null)), null)
      pollingInterval = try(coalesce(try(local.workers_keda.polling_interval_seconds, null)), null)
      cooldownPeriod  = try(coalesce(try(local.workers_keda.cooldown_period_seconds, null)), null)
    } : k => v if v != null
  } : null

  workers_persistence = {
    for k, v in {
      enabled = try(coalesce(try(var.spec.components.workers.persistence_enabled, null)), null)
      size    = try(coalesce(try(var.spec.components.workers.persistence_size, null)), null)
    } : k => v if v != null
  }
  workers_block = {
    for k, v in {
      replicas    = try(coalesce(try(var.spec.components.workers.replicas, null)), null)
      resources   = local.rendered_resources.workers
      persistence = length(local.workers_persistence) > 0 ? local.workers_persistence : null
      keda        = local.workers_keda_block
    } : k => v if v != null
  }

  # ---- DAG delivery ----------------------------------------------------------------------
  git_sync = try(var.spec.dags.git_sync, null)
  git_sync_block = local.git_sync == null ? null : {
    for k, v in {
      enabled = true
      repo    = local.git_sync.repo
      # The chart renders BOTH env generations UNCONDITIONALLY —
      # GITSYNC_REF from `ref` (v4) and GIT_SYNC_BRANCH from `branch`
      # (legacy) — and git-sync v4 translates the deprecated --branch
      # OVER --ref, so a ref-only rendering silently syncs the chart's
      # default branch (verified live: the sidecar fetched v2-2-stable
      # while ref carried the declared value). Both keys always render
      # the spec value — INCLUDING the empty string, which neutralizes
      # the chart's v2-2-stable defaults so the spec's Empty = HEAD
      # promise holds (git-sync treats empty ref/branch as HEAD).
      ref               = local.git_sync.ref
      branch            = local.git_sync.ref
      subPath           = local.git_sync.sub_path
      period            = try(coalesce(local.git_sync.period_seconds), null) != null ? "${local.git_sync.period_seconds}s" : null
      depth             = try(coalesce(local.git_sync.depth), null)
      credentialsSecret = local.git_sync.credentials_secret != "" ? local.git_sync.credentials_secret : null
      sshKeySecret      = local.git_sync.ssh_key_secret != "" ? local.git_sync.ssh_key_secret : null
      knownHosts        = local.git_sync.known_hosts != "" ? local.git_sync.known_hosts : null
      resources         = local.rendered_resources.git_sync
    } : k => v if v != null
  }

  dags_persistence = try(var.spec.dags.persistence, null)
  dags_persistence_block = local.dags_persistence == null ? null : {
    for k, v in {
      enabled          = true
      size             = try(coalesce(local.dags_persistence.size), null)
      storageClassName = local.dags_persistence.storage_class != "" ? local.dags_persistence.storage_class : null
      existingClaim    = local.dags_persistence.existing_claim != "" ? local.dags_persistence.existing_claim : null
    } : k => v if v != null
  }

  # One null-pruned map — never `{gitSync=…}` vs `{persistence=…}` on
  # ternary branches (the type-unification class). At most one arm is
  # non-null (the spec's oneof).
  dags_block_entries = {
    for k, v in {
      gitSync     = local.git_sync_block
      persistence = local.dags_persistence_block
    } : k => v if v != null
  }
  dags_block = length(local.dags_block_entries) > 0 ? local.dags_block_entries : null

  # ---- logs persistence ---------------------------------------------------------------------
  logs_persistence = try(var.spec.logging.persistence, null)
  logs_block = try(local.logs_persistence.enabled, false) ? {
    persistence = {
      for k, v in {
        enabled          = true
        size             = try(coalesce(local.logs_persistence.size), null)
        storageClassName = local.logs_persistence.storage_class != "" ? local.logs_persistence.storage_class : null
      } : k => v if v != null
    }
  } : null

  # ---- broker (bundled redis) ------------------------------------------------------------------
  bundled_redis = try(var.spec.broker.bundled_redis, null)
  redis_persistence = local.bundled_redis == null ? {} : {
    for k, v in {
      size             = try(coalesce(local.bundled_redis.persistence_size), null)
      storageClassName = try(local.bundled_redis.storage_class, "") != "" ? local.bundled_redis.storage_class : null
    } : k => v if v != null
  }
  # Never `cond ? {…} : {enabled=false}` — the type-unification class;
  # one null-pruned map serves both shapes.
  redis_block = {
    for k, v in {
      enabled = local.bundled_redis_enabled
      # The module owns the password Secret — the chart's own path
      # draws a NEW random password on every render behind a
      # pre-install hook (upstream's admitted hack).
      passwordSecretName = local.bundled_redis_enabled ? local.redis_password_secret_name : null
      persistence        = local.bundled_redis_enabled && length(local.redis_persistence) > 0 ? local.redis_persistence : null
      resources          = local.bundled_redis_enabled ? local.rendered_resources.bundled_redis : null
    } : k => v if v != null
  }

  # ---- pgbouncer block ----------------------------------------------------------------------------
  pgbouncer_block = local.pgbouncer_enabled ? {
    for k, v in {
      enabled = true
      # The module composes pgbouncer.ini + users.txt — the chart's
      # own rendering path embeds the database password in Helm
      # values and is never used.
      configSecretName      = local.pgbouncer_config_secret
      metadataPoolSize      = try(coalesce(try(var.spec.pgbouncer.metadata_pool_size, null)), null)
      resultBackendPoolSize = try(coalesce(try(var.spec.pgbouncer.result_backend_pool_size, null)), null)
      maxClientConn         = try(coalesce(try(var.spec.pgbouncer.max_client_connections, null)), null)
      resources             = local.rendered_resources.pgbouncer
      # The module-composed stats DSN for the metrics-exporter sidecar
      # (see the pgbouncer_stats_uri comment) — with statsSecretName
      # set, the chart skips creating its split-values stats Secret.
      metricsExporterSidecar = { statsSecretName = local.pgbouncer_stats_secret }
    } : k => v if v != null
  } : null

  # ---- admin bootstrap Job -----------------------------------------------------------------------------
  # The chart's default renders the admin password as a LITERAL POD
  # ARGUMENT — this block replaces the password argument with a
  # job-scoped env var read from the admin Secret. One null-pruned map
  # serves the enabled and disabled shapes (the type-unification class).
  # useHelmHooks false for the same reason as migrateDatabaseJob: a
  # post-install hook only fires after the release wait, which never
  # completes on a fresh database (the migration-hook deadlock class).
  create_user_job_block = {
    for k, v in {
      useHelmHooks = false
      enabled      = local.admin_create
      env = local.admin_create ? [{
        name = "ADMIN_PASSWORD"
        valueFrom = {
          secretKeyRef = {
            name = local.admin_secret_name
            key  = local.admin_secret_key
          }
        }
      }] : null
      args = local.admin_create ? [
        "bash",
        "-c",
        "exec \\\nairflow users create \"$@\"",
        "--",
        "-r", "Admin",
        "-u", local.admin_username,
        "-e", local.admin_email,
        "-f", "admin",
        "-l", "user",
        "-p", "$(ADMIN_PASSWORD)",
      ] : null
    } : k => v if v != null
  }

  # ---- scheduling -------------------------------------------------------------------------------------------
  scheduling_block = {
    for k, v in {
      nodeSelector = length(try(var.spec.scheduling.node_selector, {})) > 0 ? var.spec.scheduling.node_selector : null
      tolerations = length(try(var.spec.scheduling.tolerations, [])) > 0 ? [
        for t in var.spec.scheduling.tolerations : {
          for tk, tv in {
            key               = t.key != "" ? t.key : null
            operator          = t.operator != "" ? t.operator : null
            value             = t.value != "" ? t.value : null
            effect            = t.effect != "" ? t.effect : null
            tolerationSeconds = try(t.toleration_seconds, null)
          } : tk => tv if tv != null
        }
      ] : null
    } : k => v if v != null
  }

  # ---- images -------------------------------------------------------------------------------------------------
  images = try(var.spec.images, null)
  airflow_image_block = local.images == null ? null : {
    for k, v in {
      repository = local.images.airflow_repository != "" ? local.images.airflow_repository : null
      tag        = local.images.airflow_tag != "" ? local.images.airflow_tag : null
      digest     = local.images.airflow_digest != "" ? local.images.airflow_digest : null
    } : k => v if v != null
  }
  images_block = local.images == null ? null : {
    for k, v in {
      airflow   = length(coalesce(local.airflow_image_block, {})) > 0 ? local.airflow_image_block : null
      statsd    = local.images.statsd_repository != "" ? { repository = local.images.statsd_repository } : null
      redis     = local.images.redis_repository != "" ? { repository = local.images.redis_repository } : null
      pgbouncer = local.images.pgbouncer_repository != "" ? { repository = local.images.pgbouncer_repository } : null
      gitSync   = local.images.git_sync_repository != "" ? { repository = local.images.git_sync_repository } : null
    } : k => v if v != null
  }

  statsd_enabled = try(coalesce(try(var.spec.statsd_enabled, null)), null) != null ? var.spec.statsd_enabled : true

  # ---- data block (connection secret names only — never values) --------------------------------------------------
  data_block = merge(
    { metadataSecretName = local.metadata_conn_secret_name },
    local.celery_enabled ? {
      # The result backend must carry the `db+` scheme, so it never
      # falls back to the metadata Secret; the broker URL Secret is
      # module-composed (bundled/valkey arms) or user-provided.
      resultBackendSecretName = local.result_backend_conn_secret_name
      brokerUrlSecretName     = local.broker_url_secret_name
    } : {}
  )

  # ---- typed chart values (twin of the Pulumi module's buildHelmValues) ----
  helm_values = {
    for k, v in {
      # The bundled PostgreSQL subchart NEVER ships (upstream marks it
      # non-production; frozen bitnamilegacy image). Re-pinned after
      # the escape-hatch merge in main.tf.
      postgresql = { enabled = false }

      executor = local.executor
      # airflowVersion gates version-specific rendering;
      # defaultAirflowTag keeps the image tag in lockstep (the chart's
      # airflow_image helper falls back to defaultAirflowTag, NOT to
      # airflowVersion).
      airflowVersion    = local.airflow_version
      defaultAirflowTag = try(local.images.airflow_tag, "") != "" ? local.images.airflow_tag : local.airflow_version

      defaultAirflowRepository = try(local.images.airflow_repository, "") != "" ? local.images.airflow_repository : null

      data = local.data_block

      fernetKeySecretName          = local.fernet_key_secret_name
      apiSecretKeySecretName       = local.api_secret_key_name
      webserverSecretKeySecretName = local.webserver_secret_key_name
      jwtSecretName                = local.jwt_secret_name

      redis = local.redis_block

      apiServer    = length(local.api_server_block) > 0 ? local.api_server_block : null
      scheduler    = length(local.scheduler_block) > 0 ? local.scheduler_block : null
      dagProcessor = length(local.dag_processor_block) > 0 ? local.dag_processor_block : null
      triggerer    = length(local.triggerer_block) > 0 ? local.triggerer_block : null
      workers      = length(local.workers_block) > 0 ? local.workers_block : null

      dags = local.dags_block
      logs = local.logs_block

      pgbouncer = local.pgbouncer_block

      statsd = { enabled = local.statsd_enabled }

      # The env-var form, NOT config.core: the official image BAKES
      # AIRFLOW__CORE__LOAD_EXAMPLES=False as a container env, and
      # Airflow's precedence puts env above airflow.cfg — a cfg-only
      # True is silently defeated (verified live: examples never
      # parsed; the chart's own docs prescribe the env route).
      env = var.spec.load_examples ? [
        { name = "AIRFLOW__CORE__LOAD_EXAMPLES", value = "True" }
      ] : null

      createUserJob = local.create_user_job_block

      # The chart's migration Job defaults to a post-install Helm HOOK,
      # and Helm runs post-install hooks only AFTER the release wait
      # completes — while every component's wait-for-airflow-migrations
      # init container blocks on the migrations that hook would apply.
      # Under any wait-style install (both engines wait) that is a
      # deadlock by construction: the pods never turn Ready, the wait
      # expires, the hook never fires (verified live: no Job existed
      # while every init container crash-looped on "unapplied
      # migrations"). Hook-less mode makes the Job an ordinary release
      # resource applied WITH the install; the chart's own
      # ttlSecondsAfterFinished: 300 default self-deletes the finished
      # Job, so day-2 applies recreate it cleanly.
      migrateDatabaseJob = { useHelmHooks = false }

      nodeSelector = try(local.scheduling_block.nodeSelector, null)
      tolerations  = try(local.scheduling_block.tolerations, null)

      images = local.images_block != null && length(coalesce(local.images_block, {})) > 0 ? local.images_block : null
    } : k => v if v != null
  }

  # Log read path folded separately — the chart key is the backend name
  # itself (elasticsearch:/opensearch:), so it cannot be a static key in
  # the map above. Unconditional merge with an empty map (never a
  # ternary between the merged and unmerged objects — the
  # type-unification class).
  helm_values_with_log_backend = merge(
    local.helm_values,
    {
      for k, v in {
        elasticsearch = local.log_backend == "elasticsearch" ? { enabled = true, secretName = local.log_read_conn_secret_name } : null
        opensearch    = local.log_backend == "opensearch" ? { enabled = true, secretName = local.log_read_conn_secret_name } : null
      } : k => v if v != null
    }
  )

  # ---- outputs -------------------------------------------------------------------------------------------------------
  api_server_service_name = "${local.release_name}-api-server"
  api_server_endpoint     = "http://${local.api_server_service_name}.${local.namespace}.svc.cluster.local:8080"
  port_forward_command    = "kubectl port-forward svc/${local.api_server_service_name} -n ${local.namespace} 8080:8080"
}
