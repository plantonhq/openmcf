# KubernetesHarbor Terraform module.
#
# Installs Harbor from the official chart as a real Helm release. The
# typed spec renders into chart values (locals.typed_helm_values); every
# module-generated credential is materialized into module-owned Secrets
# BEFORE the release and referenced by NAME through the chart's
# existingSecret sites — the chart's publicly documented default
# credentials (Harbor12345 / changeit / not-a-secure-key /
# harbor_registry_password) never ship. Exact twin of the Pulumi module
# (values.go / auth_secrets.go).
#
# EXPOSURE COMPOSES: the module always renders one of the chart's
# in-cluster Service exposure types (ClusterIP/NodePort/LoadBalancer)
# with the chart's nginx front door terminating client traffic; the
# chart's ingress and Gateway API route types are never rendered —
# north-south exposure is a separate composed resource pointed at the
# exported front-door Service.
#
# DESTROY TRUTH (chart-verified at 1.19.1): with the default
# keep_volumes_on_uninstall the registry and jobservice PVCs carry
# `helm.sh/resource-policy: keep` and survive uninstall for a reinstall
# to adopt; the INTERNAL database and Redis volumes are
# StatefulSet-template PVCs Helm never deletes regardless. Retiring an
# install for good means sweeping those PVCs explicitly.

# The optional installation namespace. Created before the release;
# deleted with the resource.
resource "kubernetes_namespace_v1" "harbor" {
  count = var.spec.create_namespace ? 1 : 0

  metadata {
    name   = local.namespace
    labels = local.labels
  }
}

# ---------------------------------------------------------------------------
# Generated credentials. Letters+digits only (several consumers embed
# these into URLs, htpasswd lines, and env values where shell/URL
# structural characters invite quoting bugs); complexity minimums keep
# every generated credential valid against Harbor's password policy.
# The generation-shape arguments are ignored after creation so an
# IMPORTED credential never silently regenerates: rotation stays an
# explicit verb, never plan fallout. Twin: the Pulumi module's
# RandomPassword resources with the same shapes.
# ---------------------------------------------------------------------------

resource "random_password" "admin" {
  count = local.admin_generated ? 1 : 0

  length      = 32
  special     = false
  min_upper   = 2
  min_lower   = 2
  min_numeric = 2

  lifecycle {
    ignore_changes = [
      length, special, upper, lower, numeric,
      min_lower, min_numeric, min_special, min_upper, override_special,
    ]
  }
}

# The 16-char encryption key (chart contract: exactly 16 characters).
resource "random_password" "encryption_key" {
  length      = 16
  special     = false
  min_upper   = 2
  min_lower   = 2
  min_numeric = 2

  lifecycle {
    ignore_changes = [
      length, special, upper, lower, numeric,
      min_lower, min_numeric, min_special, min_upper, override_special,
    ]
  }
}

resource "random_password" "core_secret" {
  length      = 16
  special     = false
  min_upper   = 2
  min_lower   = 2
  min_numeric = 2

  lifecycle {
    ignore_changes = [
      length, special, upper, lower, numeric,
      min_lower, min_numeric, min_special, min_upper, override_special,
    ]
  }
}

# The CSRF key (chart contract: exactly 32 characters).
resource "random_password" "csrf_key" {
  length      = 32
  special     = false
  min_upper   = 2
  min_lower   = 2
  min_numeric = 2

  lifecycle {
    ignore_changes = [
      length, special, upper, lower, numeric,
      min_lower, min_numeric, min_special, min_upper, override_special,
    ]
  }
}

resource "random_password" "jobservice_secret" {
  length      = 16
  special     = false
  min_upper   = 2
  min_lower   = 2
  min_numeric = 2

  lifecycle {
    ignore_changes = [
      length, special, upper, lower, numeric,
      min_lower, min_numeric, min_special, min_upper, override_special,
    ]
  }
}

resource "random_password" "registry_http_secret" {
  length      = 16
  special     = false
  min_upper   = 2
  min_lower   = 2
  min_numeric = 2

  lifecycle {
    ignore_changes = [
      length, special, upper, lower, numeric,
      min_lower, min_numeric, min_special, min_upper, override_special,
    ]
  }
}

resource "random_password" "registry_credential" {
  length      = 32
  special     = false
  min_upper   = 2
  min_lower   = 2
  min_numeric = 2

  lifecycle {
    ignore_changes = [
      length, special, upper, lower, numeric,
      min_lower, min_numeric, min_special, min_upper, override_special,
    ]
  }
}

resource "random_password" "internal_database" {
  count = local.internal_database ? 1 : 0

  length      = 24
  special     = false
  min_upper   = 2
  min_lower   = 2
  min_numeric = 2

  lifecycle {
    ignore_changes = [
      length, special, upper, lower, numeric,
      min_lower, min_numeric, min_special, min_upper, override_special,
    ]
  }
}

# ---------------------------------------------------------------------------
# Module-owned credential Secrets — created BEFORE the release (the
# chart reads several at install time via template-time lookups and
# wires the rest as secretKeyRef env).
# ---------------------------------------------------------------------------

# The exported credential handle (generated arm only).
resource "kubernetes_secret_v1" "admin_auth" {
  count = local.admin_generated ? 1 : 0

  metadata {
    name      = local.admin_secret_name
    namespace = local.namespace
    labels    = local.labels
  }

  data = {
    HARBOR_ADMIN_PASSWORD = random_password.admin[0].result
  }

  depends_on = [
    kubernetes_namespace_v1.harbor,
  ]
}

# Every generated inter-component credential, one Secret with each
# chart site's CONTRACT KEY. The htpasswd line uses random_password's
# STABLE bcrypt hash — the chart's own `htpasswd` template function
# re-salts on every render, which would rotate the credential on every
# apply (the chart's values comment itself recommends a pre-computed
# line for CD tools).
resource "kubernetes_secret_v1" "internal_auth" {
  metadata {
    name      = local.internal_auth_secret_name
    namespace = local.namespace
    labels    = local.labels
  }

  data = {
    secretKey            = random_password.encryption_key.result
    secret               = random_password.core_secret.result
    CSRF_KEY             = random_password.csrf_key.result
    JOBSERVICE_SECRET    = random_password.jobservice_secret.result
    REGISTRY_HTTP_SECRET = random_password.registry_http_secret.result
    REGISTRY_PASSWD      = random_password.registry_credential.result
    REGISTRY_HTPASSWD    = "${local.registry_credential_username}:${random_password.registry_credential.bcrypt_hash}"
  }

  depends_on = [
    kubernetes_namespace_v1.harbor,
  ]
}

# Declared external-redis password → the chart's contract keys.
resource "kubernetes_secret_v1" "redis_auth" {
  count = local.redis_auth_secret_name != "" ? 1 : 0

  metadata {
    name      = local.redis_auth_secret_name
    namespace = local.namespace
    labels    = local.labels
  }

  data = merge(
    { REDIS_PASSWORD = var.spec.cache.external.password },
    try(coalesce(var.spec.cache.external.username), "") != "" ? { REDIS_USERNAME = var.spec.cache.external.username } : {},
  )

  depends_on = [
    kubernetes_namespace_v1.harbor,
  ]
}

# Declared s3/gcs/azure storage credentials → each driver's contract
# keys.
resource "kubernetes_secret_v1" "storage_auth" {
  count = local.storage_auth_secret_name != "" ? 1 : 0

  metadata {
    name      = local.storage_auth_secret_name
    namespace = local.namespace
    labels    = local.labels
  }

  data = merge(
    local.storage_s3 != null && try(coalesce(local.storage_s3.credentials.access_key), "") != "" ? {
      REGISTRY_STORAGE_S3_ACCESSKEY = local.storage_s3.credentials.access_key
      REGISTRY_STORAGE_S3_SECRETKEY = local.storage_s3.credentials.secret_key
    } : {},
    local.storage_gcs != null && try(coalesce(local.storage_gcs.key_data), "") != "" ? {
      GCS_KEY_DATA = local.storage_gcs.key_data
    } : {},
    local.storage_azure != null && try(coalesce(local.storage_azure.account_key), "") != "" ? {
      AZURE_STORAGE_ACCESS_KEY = local.storage_azure.account_key
    } : {},
  )

  depends_on = [
    kubernetes_namespace_v1.harbor,
  ]
}

resource "helm_release" "harbor" {
  name       = local.release_name
  repository = local.helm_chart_repo
  chart      = local.helm_chart_name
  version    = local.chart_version
  namespace  = local.namespace

  # The module owns namespace creation (create_namespace flag).
  create_namespace = false

  # Wait for the rollout: every Harbor component self-readies (core's
  # startup probe budgets 60 minutes for first-boot schema migrations —
  # the 900s timeout covers the normal case; genuinely slow clusters
  # surface honestly as a timeout). The Pulumi twin waits identically.
  # wait_for_jobs stays false in BOTH engines: the chart ships no hook
  # jobs at this pin (enableMigrateHelmHook defaults off and is not
  # modeled), so waiting on jobs is meaningless.
  wait            = true
  wait_for_jobs   = false
  atomic          = false
  cleanup_on_fail = false
  timeout         = 900

  # Documents merged in order by the provider (helm -f semantics): the
  # typed rendering first, the air-gap image mirror, the user's escape
  # hatch — and fullnameOverride plus the front-door Service name
  # re-pinned LAST, the one deliberate exception to the escape hatch's
  # last-word contract (twin of the Pulumi module). Every exported name
  # output derives from the fullname; letting an override move it would
  # break the composition handles.
  values = concat(
    [yamlencode(local.typed_helm_values)],
    length(local.image_mirror_values) > 0 ? [yamlencode(local.image_mirror_values)] : [],
    try(var.spec.helm_values, "") != "" ? [var.spec.helm_values] : [],
    [yamlencode({
      fullnameOverride = local.release_name
      expose           = { (local.expose_type) = { name = local.release_name } }
    })]
  )

  lifecycle {
    # NAME BUDGET (chart truth at 1.19.1): the chart truncates its
    # fullname at 63 and then APPENDS component suffixes — the longest,
    # `-jobservice-internal-tls` (24 chars), renders whenever
    # internalTLS runs in auto mode. The Pulumi twin enforces the same
    # budget.
    precondition {
      condition     = length(var.metadata.name) <= local.max_name_length
      error_message = "metadata.name exceeds the Harbor name budget: the chart derives object names by suffixing up to 24 characters onto it and Kubernetes caps names at 63 — use at most 39 characters."
    }
  }

  depends_on = [
    kubernetes_namespace_v1.harbor,
    kubernetes_secret_v1.admin_auth,
    kubernetes_secret_v1.internal_auth,
    kubernetes_secret_v1.redis_auth,
    kubernetes_secret_v1.storage_auth,
  ]
}
