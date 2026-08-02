# KubernetesAirflow Terraform module.
#
# Installs Apache Airflow from the official Helm chart as a real Helm
# release. Every credential travels through module-owned Secrets composed
# BEFORE the release — connection URIs, security keys and the admin
# bootstrap password never appear in rendered values (only Secret NAMES
# do); the helm_values escape hatch is passed as a SECOND values
# document, which the provider merges over the first with Helm -f
# semantics — the exact semantic twin of the Pulumi module's
# buildHelmValues + mergeMaps.

# The optional installation namespace. Created before the release; deleted
# with the resource.
resource "kubernetes_namespace_v1" "airflow" {
  count = var.spec.create_namespace ? 1 : 0

  metadata {
    name   = local.namespace
    labels = local.labels
  }
}

# ---------------------------------------------------------------------------
# Module-generated security keys. Generation-shape arguments are ignored
# after creation so an IMPORTED credential never silently regenerates:
# rotation stays an explicit verb, never plan fallout (Pulumi twin:
# IgnoreChanges on the same arguments).
# ---------------------------------------------------------------------------

# Fernet requires EXACTLY 32 random bytes in URL-SAFE base64 (a
# password-charset string would not decode to 32 bytes); random_bytes +
# the +/→-_ substitution produces the exact shape.
resource "random_bytes" "fernet_key" {
  count  = local.fernet_key_secret_byo ? 0 : 1
  length = 32

  lifecycle {
    ignore_changes = [length]
  }
}

resource "random_password" "api_secret_key" {
  count   = local.api_secret_key_byo ? 0 : 1
  length  = 32
  special = false

  lifecycle {
    ignore_changes = [length, special, upper, lower, numeric, min_lower, min_numeric, min_special, min_upper, override_special]
  }
}

resource "random_password" "webserver_secret_key" {
  length  = 32
  special = false

  lifecycle {
    ignore_changes = [length, special, upper, lower, numeric, min_lower, min_numeric, min_special, min_upper, override_special]
  }
}

resource "random_password" "jwt_secret" {
  count   = local.jwt_secret_byo ? 0 : 1
  length  = 64
  special = false

  lifecycle {
    ignore_changes = [length, special, upper, lower, numeric, min_lower, min_numeric, min_special, min_upper, override_special]
  }
}

resource "random_password" "admin_password" {
  count       = local.admin_secret_module_owned ? 1 : 0
  length      = 24
  special     = false
  min_upper   = 2
  min_lower   = 2
  min_numeric = 2

  lifecycle {
    ignore_changes = [length, special, upper, lower, numeric, min_lower, min_numeric, min_special, min_upper, override_special]
  }
}

# Letters+digits only: the password embeds into the broker URL where
# reserved characters invite quoting bugs; the larger length compensates
# the smaller alphabet.
resource "random_password" "redis_password" {
  count   = local.bundled_redis_enabled ? 1 : 0
  length  = 40
  special = false

  lifecycle {
    ignore_changes = [length, special, upper, lower, numeric, min_lower, min_numeric, min_special, min_upper, override_special]
  }
}

# ---------------------------------------------------------------------------
# Module-owned key Secrets (the chart's *SecretName contracts). The
# chart's own path would draw NEW random values on every upgrade render
# (randAlphaNum in the templates), logging out every session — the
# module-generated values are stable.
# ---------------------------------------------------------------------------

resource "kubernetes_secret_v1" "fernet_key" {
  count = local.fernet_key_secret_byo ? 0 : 1

  metadata {
    name      = local.fernet_key_secret_name
    namespace = local.namespace
    labels    = local.labels
  }
  type = "Opaque"
  data = {
    # URL-safe base64 of 32 random bytes — the exact Fernet key shape.
    "fernet-key" = replace(replace(random_bytes.fernet_key[0].base64, "+", "-"), "/", "_")
    # The companion STANDARD-base64 form of the same bytes — the import
    # handle for the random_bytes resource, whose import format is
    # standard base64 while Fernet requires the URL-safe alphabet.
    # Never consumed by the chart.
    "fernet-key-std-b64" = random_bytes.fernet_key[0].base64
  }

  depends_on = [kubernetes_namespace_v1.airflow]
}

resource "kubernetes_secret_v1" "api_secret_key" {
  count = local.api_secret_key_byo ? 0 : 1

  metadata {
    name      = local.api_secret_key_name
    namespace = local.namespace
    labels    = local.labels
  }
  type = "Opaque"
  data = {
    "api-secret-key" = random_password.api_secret_key[0].result
  }

  depends_on = [kubernetes_namespace_v1.airflow]
}

# The FAB webserver session key has no BYO field: always module-owned.
resource "kubernetes_secret_v1" "webserver_secret_key" {
  metadata {
    name      = local.webserver_secret_key_name
    namespace = local.namespace
    labels    = local.labels
  }
  type = "Opaque"
  data = {
    "webserver-secret-key" = random_password.webserver_secret_key.result
  }

  depends_on = [kubernetes_namespace_v1.airflow]
}

resource "kubernetes_secret_v1" "jwt_secret" {
  count = local.jwt_secret_byo ? 0 : 1

  metadata {
    name      = local.jwt_secret_name
    namespace = local.namespace
    labels    = local.labels
  }
  type = "Opaque"
  data = {
    "jwt-secret" = random_password.jwt_secret[0].result
  }

  depends_on = [kubernetes_namespace_v1.airflow]
}

resource "kubernetes_secret_v1" "admin_auth" {
  count = local.admin_secret_module_owned ? 1 : 0

  metadata {
    name      = local.admin_secret_name
    namespace = local.namespace
    labels    = local.labels
  }
  type = "Opaque"
  data = {
    "password" = random_password.admin_password[0].result
  }

  depends_on = [kubernetes_namespace_v1.airflow]
}

resource "kubernetes_secret_v1" "redis_password" {
  count = local.bundled_redis_enabled ? 1 : 0

  metadata {
    name      = local.redis_password_secret_name
    namespace = local.namespace
    labels    = local.labels
  }
  type = "Opaque"
  data = {
    "password" = random_password.redis_password[0].result
  }

  depends_on = [kubernetes_namespace_v1.airflow]
}

# ---------------------------------------------------------------------------
# The referenced credential reads. The referenced Secrets are created by
# OTHER components (e.g. the KubernetesPostgres app Secret), so they
# exist before this module applies. depends_on a module-created resource
# DEFERS the read to apply time on a fresh plan — an offline plan (no
# cluster) stays green with "(known after apply)" while a real apply
# reads the live value (Pulumi twin: the DryRun-gated GetSecret).
# ---------------------------------------------------------------------------

data "kubernetes_secret_v1" "db_password" {
  metadata {
    name      = local.db_password_secret
    namespace = local.namespace
  }

  depends_on = [kubernetes_secret_v1.webserver_secret_key]
}

data "kubernetes_secret_v1" "broker_password" {
  count = local.valkey_broker != null && local.broker_password_secret != "" ? 1 : 0

  metadata {
    name      = local.broker_password_secret
    namespace = local.namespace
  }

  depends_on = [kubernetes_secret_v1.webserver_secret_key]
}

data "kubernetes_secret_v1" "log_backend_password" {
  count = local.log_backend != "" && local.log_backend_user != "" ? 1 : 0

  metadata {
    name      = local.log_backend_password_secret
    namespace = local.namespace
  }

  depends_on = [kubernetes_secret_v1.webserver_secret_key]
}

# ---------------------------------------------------------------------------
# Composed connection Secrets (the chart's `connection`-key contracts).
# ---------------------------------------------------------------------------

resource "kubernetes_secret_v1" "metadata_conn" {
  metadata {
    name      = local.metadata_conn_secret_name
    namespace = local.namespace
    labels    = local.labels
  }
  type = "Opaque"
  # `kedaConnection` rides along wherever a KEDA trigger could need the
  # direct-database form (see the locals comment) — the chart's worker
  # and triggerer autoscalers read it from THIS Secret by name.
  data = merge(
    { "connection" = local.metadata_conn_uri },
    local.keda_conn_needed ? { "kedaConnection" = local.keda_conn_uri } : {}
  )

  depends_on = [kubernetes_namespace_v1.airflow]
}

resource "kubernetes_secret_v1" "result_backend_conn" {
  count = local.celery_enabled ? 1 : 0

  metadata {
    name      = local.result_backend_conn_secret_name
    namespace = local.namespace
    labels    = local.labels
  }
  type = "Opaque"
  data = {
    "connection" = local.result_backend_uri
  }

  depends_on = [kubernetes_namespace_v1.airflow]
}

resource "kubernetes_secret_v1" "broker_url" {
  count = local.broker_url_secret_module_owned ? 1 : 0

  metadata {
    name      = local.broker_url_secret_name
    namespace = local.namespace
    labels    = local.labels
  }
  type = "Opaque"
  data = {
    # redis://[user]:<password>@<host>:<port>/<db> — the chart's own URL
    # shape for the broker; userinfo urlencoded (chart parity).
    "connection" = local.bundled_redis_enabled ? (
      "redis://:${random_password.redis_password[0].result}@${local.broker_host}:${local.broker_port}/${local.broker_db}"
      ) : (
      local.broker_password_secret != "" ? (
        "redis://${urlencode(local.broker_user)}:${urlencode(data.kubernetes_secret_v1.broker_password[0].data[local.broker_password_secret_key])}@${local.broker_host}:${local.broker_port}/${local.broker_db}"
        ) : (
        "redis://${urlencode(local.broker_user)}:@${local.broker_host}:${local.broker_port}/${local.broker_db}"
      )
    )
  }

  depends_on = [kubernetes_namespace_v1.airflow]
}

resource "kubernetes_secret_v1" "log_read_conn" {
  count = local.log_backend != "" ? 1 : 0

  metadata {
    name      = local.log_read_conn_secret_name
    namespace = local.namespace
    labels    = local.labels
  }
  type = "Opaque"
  data = {
    "connection" = local.log_backend_user != "" ? (
      "${local.log_backend_scheme}://${urlencode(local.log_backend_user)}:${urlencode(data.kubernetes_secret_v1.log_backend_password[0].data[local.log_backend_password_secret_key])}@${local.log_backend_config.host}:${local.log_backend_port}"
      ) : (
      "${local.log_backend_scheme}://${local.log_backend_config.host}:${local.log_backend_port}"
    )
  }

  depends_on = [kubernetes_namespace_v1.airflow]
}

resource "kubernetes_secret_v1" "pgbouncer_config" {
  count = local.pgbouncer_enabled ? 1 : 0

  metadata {
    name      = local.pgbouncer_config_secret
    namespace = local.namespace
    labels    = local.labels
  }
  type = "Opaque"
  data = {
    # Byte-faithful to the chart's own pgbouncer_config helper: the
    # chart's rendering path would embed the database password in Helm
    # values, so the module composes both files instead.
    "pgbouncer.ini" = local.pgbouncer_ini
    "users.txt"     = local.pgbouncer_users_txt
  }

  depends_on = [kubernetes_namespace_v1.airflow]
}

resource "kubernetes_secret_v1" "pgbouncer_stats" {
  count = local.pgbouncer_enabled ? 1 : 0

  metadata {
    name      = local.pgbouncer_stats_secret
    namespace = local.namespace
    labels    = local.labels
  }
  type = "Opaque"
  data = {
    # The metrics-exporter sidecar's DSN (chart contract: key
    # `connection`, pgbouncer admin database over localhost). The
    # chart's own stats Secret would render split-values defaults the
    # auth_file never carries — see the pgbouncer_stats_uri comment.
    "connection" = local.pgbouncer_stats_uri
  }

  depends_on = [kubernetes_namespace_v1.airflow]
}

# ---------------------------------------------------------------------------
# The Helm release.
# ---------------------------------------------------------------------------

resource "helm_release" "airflow" {
  name       = local.release_name
  repository = local.helm_chart_repo
  chart      = local.helm_chart_name
  version    = local.chart_version
  namespace  = local.namespace

  # The module owns namespace creation (create_namespace flag).
  create_namespace = false

  # Wait for the components to become Ready — an install whose
  # migration Job cannot reach the database, whose credential Secrets
  # are misnamed, or whose scheduler crash-loops should fail THIS
  # apply, not the first pipeline run. The post-install migration +
  # create-user hook Jobs run inside this budget too.
  wait            = true
  atomic          = true
  cleanup_on_fail = true
  timeout         = 900

  # Documents merged in order by the provider (helm -f semantics): the
  # typed rendering first, the user's escape hatch second — and the
  # deliberate re-pins LAST, the two exceptions to the escape hatch's
  # last-word contract (twin of the Pulumi module): useStandardNaming
  # stays false (every child name and every exported output derives
  # from the release name) and postgresql.enabled stays false (the
  # bundled database is non-production by upstream's own definition
  # with a frozen image line).
  values = concat(
    [yamlencode(local.helm_values_with_log_backend)],
    try(var.spec.helm_values, "") != "" ? [var.spec.helm_values] : [],
    [yamlencode({
      useStandardNaming = false
      postgresql        = { enabled = false }
    })]
  )

  lifecycle {
    # FAIL LOUDLY on names past the chart's fullname budget: at the
    # default naming scheme the fullname IS the release name and child
    # names append fixed suffixes — the longest
    # ("-run-airflow-migrations", 23 chars) pushes names past the
    # Kubernetes 63-character limit when the resource name exceeds 40,
    # failing the deploy midway with API rejections. Twin: the Pulumi
    # module's Resources() guard.
    precondition {
      condition     = length(var.metadata.name) <= 40
      error_message = "The airflow chart appends child-name suffixes up to 23 characters (\"-run-airflow-migrations\") and Kubernetes caps names at 63, so the deploy would fail midway — use a name of at most 40 characters."
    }

    # The Celery pairing is CEL-enforced at the API; this precondition
    # is the engine-level backstop for manifests that bypassed it.
    precondition {
      condition     = !local.celery_enabled || local.broker_url_secret_name != ""
      error_message = "A Celery-family executor needs a message broker — declare spec.broker (bundled Redis, a composed KubernetesValkey, or an existing broker-URL Secret)."
    }
  }

  depends_on = [
    kubernetes_namespace_v1.airflow,
    kubernetes_secret_v1.fernet_key,
    kubernetes_secret_v1.api_secret_key,
    kubernetes_secret_v1.webserver_secret_key,
    kubernetes_secret_v1.jwt_secret,
    kubernetes_secret_v1.admin_auth,
    kubernetes_secret_v1.redis_password,
    kubernetes_secret_v1.metadata_conn,
    kubernetes_secret_v1.result_backend_conn,
    kubernetes_secret_v1.broker_url,
    kubernetes_secret_v1.log_read_conn,
    kubernetes_secret_v1.pgbouncer_config,
  ]
}
