# KubernetesSuperset Terraform module.
#
# Installs Apache Superset from the official ASF Helm chart as a real Helm
# release. The chart consumes ALL runtime credentials through ONE
# environment Secret — the chart's own copy is turned OFF
# (secretEnv.create=false) and this module composes `<name>-env` itself:
# non-secret connection facts and module-GENERATED material render into
# it; REFERENCED material (the database/cache passwords, bring-your-own
# keys) arrives in the pods as extraEnvRaw secretKeyRef entries — never
# copied, never read at apply time. The helm_values escape hatch is
# passed as a SECOND values document (helm -f semantics) and the security
# spine is re-pinned in a THIRD — the exact semantic twin of the Pulumi
# module's buildHelmValues + mergeMaps + re-pins.

# The optional installation namespace. Created before the release; deleted
# with the resource.
resource "kubernetes_namespace_v1" "superset" {
  count = var.spec.create_namespace ? 1 : 0

  metadata {
    name   = local.namespace
    labels = local.labels
  }
}

# ---------------------------------------------------------------------------
# The session-signing SECRET_KEY (module-generated arm). Superset REFUSES
# to start on its insecure default; the key also encrypts datasource
# credentials stored in the metadata database — the random is generated
# ONCE and shape-ignored (rotating it without Superset's own re-encrypt
# procedure orphans every stored connection). Generation-shape arguments
# are ignored after creation so an IMPORTED credential never silently
# regenerates (Pulumi twin: IgnoreChanges on the same arguments).
# ---------------------------------------------------------------------------

resource "random_password" "secret_key" {
  count   = local.secret_key_module_owned ? 1 : 0
  length  = 42
  special = false

  lifecycle {
    ignore_changes = [length, special, upper, lower, numeric, min_lower, min_numeric, min_special, min_upper, override_special]
  }
}

# The dedicated handle Secret (`<name>-secret-key`) — the exported
# reference operators rotate against (`superset re-encrypt-secrets`).
resource "kubernetes_secret_v1" "secret_key" {
  count = local.secret_key_module_owned ? 1 : 0

  metadata {
    name      = local.secret_key_secret_name
    namespace = local.namespace
    labels    = local.labels
  }
  type = "Opaque"
  data = {
    "secret_key" = random_password.secret_key[0].result
  }

  depends_on = [kubernetes_namespace_v1.superset]
}

# ---------------------------------------------------------------------------
# The bootstrap admin password (module-generated arm) — the chart's
# documented admin/admin default never ships. Letters+digits only:
# operators type this at the login form.
# ---------------------------------------------------------------------------

resource "random_password" "admin" {
  count       = local.admin_module_owned ? 1 : 0
  length      = 24
  special     = false
  min_upper   = 2
  min_lower   = 2
  min_numeric = 2

  lifecycle {
    ignore_changes = [length, special, upper, lower, numeric, min_lower, min_numeric, min_special, min_upper, override_special]
  }
}

# The exported credential handle (`<name>-admin-auth`).
resource "kubernetes_secret_v1" "admin_auth" {
  count = local.admin_module_owned ? 1 : 0

  metadata {
    name      = local.admin_password_secret_name
    namespace = local.namespace
    labels    = local.labels
  }
  type = "Opaque"
  data = {
    "password" = random_password.admin[0].result
  }

  depends_on = [kubernetes_namespace_v1.superset]
}

# ---------------------------------------------------------------------------
# The websocket JWT (websockets arm) — one key serves both sides: the ws
# server reads JWT_SECRET from its environment natively, and the module's
# configOverrides snippet points Superset's async-queries JWT at the same
# variable.
# ---------------------------------------------------------------------------

resource "random_password" "ws_jwt" {
  count   = local.websockets_enabled ? 1 : 0
  length  = 48
  special = false

  lifecycle {
    ignore_changes = [length, special, upper, lower, numeric, min_lower, min_numeric, min_special, min_upper, override_special]
  }
}

# ---------------------------------------------------------------------------
# The module-owned environment Secret (`<name>-env`) — the chart's
# runtime-credential contract: every component (web, worker, beat,
# flower, websocket, MCP, the init Job and its wait-for initContainers)
# envFroms it. Plain connection facts + module-GENERATED material only;
# referenced credentials arrive via extraEnvRaw and are never copied
# here.
# ---------------------------------------------------------------------------

resource "kubernetes_secret_v1" "env" {
  metadata {
    name      = local.env_secret_name
    namespace = local.namespace
    labels    = local.labels
  }
  type = "Opaque"
  data = merge(
    local.env_plain,
    local.secret_key_module_owned ? { SUPERSET_SECRET_KEY = random_password.secret_key[0].result } : {},
    local.admin_module_owned ? { ADMIN_PASSWORD = random_password.admin[0].result } : {},
    local.websockets_enabled ? { JWT_SECRET = random_password.ws_jwt[0].result } : {},
  )

  depends_on = [kubernetes_namespace_v1.superset]
}

# ---------------------------------------------------------------------------
# The Helm release.
# ---------------------------------------------------------------------------

resource "helm_release" "superset" {
  name       = local.release_name
  repository = local.helm_chart_repo
  chart      = local.helm_chart_name
  version    = local.helm_chart_version
  namespace  = local.namespace

  # The module owns namespace creation (create_namespace flag).
  create_namespace = false

  # Wait for the web/worker rollouts and the post-install init Job
  # (schema migration + admin bootstrap run inside this budget against
  # the composed database) — an install whose migration cannot reach
  # its database should fail THIS apply, not the first login.
  wait            = true
  atomic          = true
  cleanup_on_fail = true
  timeout         = 900

  # Documents merged in order by the provider (helm -f semantics): the
  # typed rendering first, the user's escape hatch second — and the
  # deliberate re-pins LAST, the exception to the escape hatch's
  # last-word contract (twin of the Pulumi module): the deterministic
  # names, the module-owned env-Secret contract, the dead bundled
  # subcharts and the env-driven admin bootstrap cannot be silently
  # re-enabled or redirected.
  values = concat(
    [yamlencode(local.helm_values)],
    try(var.spec.helm_values, "") != "" ? [var.spec.helm_values] : [],
    [yamlencode({
      fullnameOverride = local.release_name
      secretEnv        = { create = false }
      envFromSecret    = local.env_secret_name
      postgresql       = { enabled = false }
      redis            = { enabled = false }
      init = {
        createAdmin = false
        command     = local.init_command
      }
      configOverrides = local.config_overrides
    })]
  )

  # Fail loud on the chart's name-derivation budget before anything
  # renders (twin of the Pulumi module's buildHelmValues error).
  lifecycle {
    precondition {
      condition     = length(local.release_name) <= local.name_budget
      error_message = "metadata.name '${local.release_name}' is ${length(local.release_name)} characters — the chart derives '<name>-celerybeat' (11-char suffix), so the name must be at most ${local.name_budget} characters."
    }
  }

  depends_on = [
    kubernetes_namespace_v1.superset,
    kubernetes_secret_v1.env,
    kubernetes_secret_v1.secret_key,
    kubernetes_secret_v1.admin_auth,
  ]
}
