# KubernetesLocust Terraform module.
#
# Installs Locust from the deliveryhero Helm chart (the OFFICIAL
# locustio/locust image) as a real Helm release. The typed spec renders
# into chart values (locals.helm_values); the test scripts and the
# module-owned login backend travel through ConfigMaps composed BEFORE
# the release; the login credential lives in a module-owned Secret
# projected as files — nothing credential-bearing appears in rendered
# values or process arguments; the helm_values escape hatch is passed as
# a SECOND values document, which the provider merges over the first
# with Helm -f semantics, and the security spine is re-pinned in a THIRD
# document — the exact semantic twin of the Pulumi module's
# buildHelmValues + mergeMaps + re-pins.
#
# OCI WIRING: the chart's serving home is the OCI registry (the classic
# index stalls at 0.31.6) — the Terraform provider takes repository =
# the OCI registry path plus the bare chart name and joins them
# internally; the Pulumi twin passes the joined "oci://.../locust"
# string as the chart reference. Same chart bytes, different wiring.

# The optional installation namespace. Created before the release;
# deleted with the resource.
resource "kubernetes_namespace_v1" "locust" {
  count = var.spec.create_namespace ? 1 : 0

  metadata {
    name   = local.namespace
    labels = local.labels
  }
}

# ---------------------------------------------------------------------------
# The module-owned script ConfigMaps (inline arm): the user's locustfile
# as `main.py` and the supporting modules. The chart mounts them at the
# locustfile path (+`/lib`); content changes roll the pods through the
# module's checksum annotation (the chart's own checksums cover only
# chart-rendered ConfigMaps). The existing-ConfigMaps arm creates
# nothing here — the chart values reference the user's own.
# ---------------------------------------------------------------------------

resource "kubernetes_config_map_v1" "locustfile" {
  count = local.module_owned_scripts ? 1 : 0

  metadata {
    name      = local.locustfile_config_map
    namespace = local.namespace
    labels    = local.labels
  }
  data = {
    (local.locustfile_name) = local.inline_scripts.locustfile_content
  }

  depends_on = [kubernetes_namespace_v1.locust]
}

resource "kubernetes_config_map_v1" "lib" {
  count = local.module_owned_scripts && length(local.lib_files) > 0 ? 1 : 0

  metadata {
    name      = local.lib_config_map
    namespace = local.namespace
    labels    = local.labels
  }
  data = local.lib_files

  depends_on = [kubernetes_namespace_v1.locust]
}

# ---------------------------------------------------------------------------
# The web-UI login (the secured default: the chart ships the UI OPEN —
# that never ships from this kind). The login backend is module-owned
# CODE the master loads alongside the locustfile (Locust's documented
# extension seam; see locals.web_auth_backend_py), and the credential
# lives in the `<name>-auth` Secret projected as FILES — never
# environment, rendered values or process arguments.
# ---------------------------------------------------------------------------

resource "kubernetes_config_map_v1" "web_auth_code" {
  count = local.web_login_enabled ? 1 : 0

  metadata {
    name      = local.web_auth_code_name
    namespace = local.namespace
    labels    = local.labels
  }
  data = {
    "planton_auth.py" = local.web_auth_backend_py
  }

  depends_on = [kubernetes_namespace_v1.locust]
}

# Letters+digits only: operators type this at the login form —
# symbol-free avoids quoting bugs; the larger length compensates the
# smaller alphabet. Generation-shape arguments are ignored after
# creation so an IMPORTED credential never silently regenerates (Pulumi
# twin: IgnoreChanges on the same arguments).
resource "random_password" "web_ui" {
  count       = local.web_login_enabled ? 1 : 0
  length      = 24
  special     = false
  min_upper   = 2
  min_lower   = 2
  min_numeric = 2

  lifecycle {
    ignore_changes = [length, special, upper, lower, numeric, min_lower, min_numeric, min_special, min_upper, override_special]
  }
}

# The Flask session-signing key the login backend loads — generated once
# and stable, so sessions survive pod restarts (a per-start random would
# log every user out on every roll).
resource "random_password" "flask_secret_key" {
  count   = local.web_login_enabled ? 1 : 0
  length  = 64
  special = false

  lifecycle {
    ignore_changes = [length, special, upper, lower, numeric, min_lower, min_numeric, min_special, min_upper, override_special]
  }
}

resource "kubernetes_secret_v1" "auth" {
  count = local.web_login_enabled ? 1 : 0

  metadata {
    name      = local.auth_secret_name
    namespace = local.namespace
    labels    = local.labels
  }
  type = "Opaque"
  data = {
    "username"         = local.web_username
    "password"         = random_password.web_ui[0].result
    "flask-secret-key" = random_password.flask_secret_key[0].result
  }

  depends_on = [kubernetes_namespace_v1.locust]
}

# ---------------------------------------------------------------------------
# The Helm release.
# ---------------------------------------------------------------------------

resource "helm_release" "locust" {
  name       = local.release_name
  repository = local.helm_oci_repo
  chart      = local.helm_chart_name
  version    = local.helm_chart_version
  namespace  = local.namespace

  # The module owns namespace creation (create_namespace flag).
  create_namespace = false

  # Wait for the master and worker rollouts — an install whose pods
  # cannot start (a broken locustfile import, a missing referenced
  # Secret, a failing pip install) should fail THIS apply, not the
  # first test run.
  wait            = true
  atomic          = true
  cleanup_on_fail = true
  timeout         = 600

  # Documents merged in order by the provider (helm -f semantics): the
  # typed rendering first, the user's escape hatch second — and the
  # deliberate re-pins LAST, the exception to the escape hatch's
  # last-word contract (twin of the Pulumi module): the deterministic
  # names, the script wiring, the login wiring, and the NULLED
  # environment_secret (the chart's values-rendered-Secret path can
  # never engage).
  values = concat(
    [yamlencode(local.helm_values)],
    try(var.spec.helm_values, "") != "" ? [var.spec.helm_values] : [],
    [yamlencode(local.helm_values_repins)],
  )

  # Fail loud before anything renders (twin of the Pulumi module's
  # buildHelmValues errors).
  lifecycle {
    precondition {
      condition     = length(local.release_name) <= local.name_budget
      error_message = "metadata.name '${local.release_name}' is ${length(local.release_name)} characters — the module derives '<name>-locustfile' (11-char suffix), so the name must be at most ${local.name_budget} characters."
    }
    precondition {
      condition     = !local.web_login_enabled || local.image_tag_login_capable
      error_message = "image.tag '${local.image_tag}' cannot prove Locust >= 2.21.0 — below that the chart renders the login credential as a literal pod argument, which this module refuses; use a numeric tag at or above 2.21.0, or disable web_ui_auth."
    }
  }

  depends_on = [
    kubernetes_namespace_v1.locust,
    kubernetes_config_map_v1.locustfile,
    kubernetes_config_map_v1.lib,
    kubernetes_config_map_v1.web_auth_code,
    kubernetes_secret_v1.auth,
  ]
}
