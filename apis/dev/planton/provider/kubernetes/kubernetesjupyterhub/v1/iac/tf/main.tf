# KubernetesJupyterHub Terraform module.
#
# Installs JupyterHub from the official Zero to JupyterHub Helm chart as a
# real Helm release. Credentials travel through module-owned Secrets
# composed BEFORE the release — the external database password and the
# sign-in secrets never appear in rendered values (which this chart embeds
# READABLE inside its own hub Secret); the helm_values escape hatch is
# passed as a SECOND values document, which the provider merges over the
# first with Helm -f semantics — the exact semantic twin of the Pulumi
# module's buildHelmValues + mergeMaps.

# The optional installation namespace. Created before the release; deleted
# with the resource. KNOW THIS: deleting the namespace also deletes every
# user's home PVC living in it — back user homes up before tearing an
# instance down.
resource "kubernetes_namespace_v1" "jupyterhub" {
  count = var.spec.create_namespace ? 1 : 0

  metadata {
    name   = local.namespace
    labels = local.labels
  }
}

# ---------------------------------------------------------------------------
# The module-generated shared sign-in password (the secured default: the
# chart's own default authenticator accepts ANY username with NO password
# — that never ships). Generation-shape arguments are ignored after
# creation so an IMPORTED credential never silently regenerates (Pulumi
# twin: IgnoreChanges on the same arguments).
# ---------------------------------------------------------------------------

# Letters+digits only: the password reaches JupyterHub through an env var
# and users type it at the login form — symbol-free avoids both quoting
# bugs and login-typing friction; the larger length compensates the
# smaller alphabet.
resource "random_password" "shared_password" {
  count       = local.shared_password_module_owned ? 1 : 0
  length      = 24
  special     = false
  min_upper   = 2
  min_lower   = 2
  min_numeric = 2

  lifecycle {
    ignore_changes = [length, special, upper, lower, numeric, min_lower, min_numeric, min_special, min_upper, override_special]
  }
}

resource "kubernetes_secret_v1" "shared_password" {
  count = local.shared_password_module_owned ? 1 : 0

  metadata {
    name      = local.shared_password_secret_name
    namespace = local.namespace
    labels    = local.labels
  }
  type = "Opaque"
  data = {
    "password" = random_password.shared_password[0].result
  }

  depends_on = [kubernetes_namespace_v1.jupyterhub]
}

# ---------------------------------------------------------------------------
# The referenced database-credential read. The referenced Secret is
# created by ANOTHER component (e.g. the KubernetesPostgres app Secret),
# so it exists before this module applies. depends_on a module-created
# resource DEFERS the read to apply time on a fresh plan — an offline
# plan (no cluster) stays green with "(known after apply)" while a real
# apply reads the live value (Pulumi twin: the DryRun-gated GetSecret).
# The dependency anchor is a marker ConfigMap because the module creates
# no unconditional Secret of its own.
# ---------------------------------------------------------------------------

resource "kubernetes_config_map_v1" "read_anchor" {
  count = local.hub_secret_enabled ? 1 : 0

  metadata {
    name      = "${local.release_name}-read-anchor"
    namespace = local.namespace
    labels    = local.labels
  }
  data = {
    # Marker only — this ConfigMap exists to defer the credential
    # data-source read below to apply time; it carries no configuration.
    purpose = "defers-referenced-secret-reads-to-apply"
  }

  depends_on = [kubernetes_namespace_v1.jupyterhub]
}

data "kubernetes_secret_v1" "db_password" {
  count = local.hub_secret_enabled ? 1 : 0

  metadata {
    name      = local.db_password_secret
    namespace = local.namespace
  }

  depends_on = [kubernetes_config_map_v1.read_anchor]
}

# The hub's existing-secret (`<name>-hub-secret`): the chart mounts it at
# /usr/local/etc/jupyterhub/existing-secret/ and the hub exports
# PGPASSWORD/MYSQL_PWD from the `hub.db.password` key at startup — the
# password never rides hub.db.url or any rendered value.
resource "kubernetes_secret_v1" "hub_secret" {
  count = local.hub_secret_enabled ? 1 : 0

  metadata {
    name      = local.hub_secret_name
    namespace = local.namespace
    labels    = local.labels
  }
  type = "Opaque"
  data = {
    "hub.db.password" = data.kubernetes_secret_v1.db_password[0].data[local.db_password_secret_key]
  }

  depends_on = [kubernetes_namespace_v1.jupyterhub]
}

# ---------------------------------------------------------------------------
# The Helm release.
# ---------------------------------------------------------------------------

resource "helm_release" "jupyterhub" {
  name       = local.release_name
  repository = local.helm_chart_repo
  chart      = local.helm_chart_name
  version    = local.chart_version
  namespace  = local.namespace

  # The module owns namespace creation (create_namespace flag).
  create_namespace = false

  # Wait for the hub, proxy and scheduling machinery to become Ready —
  # an install whose hub cannot reach its database or whose image-pull
  # hook cannot finish should fail THIS apply, not the first user's
  # login. With the pre-puller hook on (the chart default) the
  # notebook-image pull to every node runs inside this budget.
  wait            = true
  atomic          = true
  cleanup_on_fail = true
  timeout         = 1200

  # Documents merged in order by the provider (helm -f semantics): the
  # typed rendering first, the user's escape hatch second — and the
  # deliberate re-pin LAST, the exception to the escape hatch's
  # last-word contract (twin of the Pulumi module): fullnameOverride
  # stays "" — the chart-fixed bare names (hub, proxy-public…) ARE the
  # exported outputs and the per-namespace singleton contract.
  values = concat(
    [yamlencode(local.helm_values)],
    try(var.spec.helm_values, "") != "" ? [var.spec.helm_values] : [],
    [yamlencode({
      fullnameOverride = ""
    })]
  )

  depends_on = [
    kubernetes_namespace_v1.jupyterhub,
    kubernetes_secret_v1.shared_password,
    kubernetes_secret_v1.hub_secret,
  ]
}
