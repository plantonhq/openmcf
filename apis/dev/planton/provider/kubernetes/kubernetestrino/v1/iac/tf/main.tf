# KubernetesTrino Terraform module.
#
# Installs Trino from the official trinodb Helm chart as a real Helm
# release. Every properties surface in this chart renders into ConfigMaps,
# so NOTHING credential-bearing rides values: catalog passwords and the
# internal shared secret are `${ENV:VAR}` references (Trino's own secrets
# substitution) resolved from Secret-sourced environment variables; the
# helm_values escape hatch is passed as a SECOND values document, which
# the provider merges over the first with Helm -f semantics, and the
# security spine is re-pinned in a THIRD document — the exact semantic
# twin of the Pulumi module's buildHelmValues + mergeMaps + re-pins.

# The optional installation namespace. Created before the release; deleted
# with the resource.
resource "kubernetes_namespace_v1" "trino" {
  count = var.spec.create_namespace ? 1 : 0

  metadata {
    name   = local.namespace
    labels = local.labels
  }
}

# ---------------------------------------------------------------------------
# The module-generated admin credential (the secured default: the chart
# ships NO authentication — that never ships from this kind). Key
# `password` is the plaintext clients use; key `password.db` is the
# htpasswd-format bcrypt file the chart mounts through
# auth.passwordAuthSecret. Both keys derive from ONE random — a verified
# pairing by construction. Generation-shape arguments are ignored after
# creation so an IMPORTED credential never silently regenerates (Pulumi
# twin: IgnoreChanges on the same arguments). KNOW THIS (import
# semantics): random_password.bcrypt_hash re-salts on import — the import
# map carries an import_normalized tolerance for the `password.db` key
# (the Harbor htpasswd class).
# ---------------------------------------------------------------------------

# Letters+digits only: users type this at BI-tool connection forms and
# the trino CLI — symbol-free avoids quoting bugs; the larger length
# compensates the smaller alphabet.
resource "random_password" "admin" {
  count       = local.module_owned_password_db ? 1 : 0
  length      = 24
  special     = false
  min_upper   = 2
  min_lower   = 2
  min_numeric = 2

  lifecycle {
    ignore_changes = [length, special, upper, lower, numeric, min_lower, min_numeric, min_special, min_upper, override_special]
  }
}

resource "kubernetes_secret_v1" "auth" {
  count = local.module_owned_password_db ? 1 : 0

  metadata {
    name      = local.password_db_secret_name
    namespace = local.namespace
    labels    = local.labels
  }
  type = "Opaque"
  data = {
    "password"    = random_password.admin[0].result
    "password.db" = "${local.admin_username}:${random_password.admin[0].bcrypt_hash}\n"
  }

  depends_on = [kubernetes_namespace_v1.trino]
}

# ---------------------------------------------------------------------------
# The internal-communication shared secret. Trino REQUIRES it once
# authentication is on (internal-communication.md at the pin); it reaches
# config.properties as `${ENV:TRINO_INTERNAL_SHARED_SECRET}` — never
# rendered. Length 64, letters+digits: the value transits an env var into
# Trino's own config parser — symbol-free sidesteps quoting ambiguity
# (the NATS password-parser lesson, applied preemptively).
# ---------------------------------------------------------------------------

resource "random_password" "internal_shared_secret" {
  count   = local.auth_enabled ? 1 : 0
  length  = 64
  special = false

  lifecycle {
    ignore_changes = [length, special, upper, lower, numeric, min_lower, min_numeric, min_special, min_upper, override_special]
  }
}

resource "kubernetes_secret_v1" "internal" {
  count = local.auth_enabled ? 1 : 0

  metadata {
    name      = local.internal_secret_name
    namespace = local.namespace
    labels    = local.labels
  }
  type = "Opaque"
  data = {
    "shared-secret" = random_password.internal_shared_secret[0].result
  }

  depends_on = [kubernetes_namespace_v1.trino]
}

# ---------------------------------------------------------------------------
# The Helm release.
# ---------------------------------------------------------------------------

resource "helm_release" "trino" {
  name       = local.release_name
  repository = local.helm_chart_repo
  chart      = local.helm_chart_name
  version    = local.helm_chart_version
  namespace  = local.namespace

  # The module owns namespace creation (create_namespace flag).
  create_namespace = false

  # Wait for the coordinator and worker rollouts — an install whose
  # coordinator cannot start (bad catalog properties, a missing
  # referenced Secret) should fail THIS apply, not the first query.
  wait            = true
  atomic          = true
  cleanup_on_fail = true
  timeout         = 900

  # Documents merged in order by the provider (helm -f semantics): the
  # typed rendering first, the user's escape hatch second — and the
  # deliberate re-pins LAST, the exception to the escape hatch's
  # last-word contract (twin of the Pulumi module): the deterministic
  # names, the PASSWORD authentication wiring and the module-owned
  # config-properties list (which carries the shared secret and the
  # insecure-over-http pairing) cannot be silently disabled.
  # Single-attribute ternaries only (the HCL type-unification class —
  # a composite true branch cannot unify against {}).
  values = concat(
    [yamlencode(local.helm_values)],
    try(var.spec.helm_values, "") != "" ? [var.spec.helm_values] : [],
    [yamlencode(merge(
      { fullnameOverride = local.release_name },
      local.auth_enabled ? { server = { config = { authenticationType = "PASSWORD" } } } : {},
      local.auth_enabled ? { auth = local.auth_block } : {},
      length(local.config_properties) > 0 ? { additionalConfigProperties = local.config_properties } : {},
    ))]
  )

  # Fail loud on the chart's name-derivation budgets before anything
  # renders (twin of the Pulumi module's buildHelmValues errors).
  lifecycle {
    precondition {
      condition     = length(local.release_name) <= local.name_budget
      error_message = "metadata.name '${local.release_name}' is ${length(local.release_name)} characters — the chart derives '<name>-schemas-volume-coordinator' (27-char suffix), so the name must be at most ${local.name_budget} characters."
    }
    precondition {
      condition     = try(var.spec.resource_groups_config, "") == "" || length(local.release_name) <= local.name_budget_resource_groups
      error_message = "metadata.name '${local.release_name}' is ${length(local.release_name)} characters — with resource_groups_config set the chart derives '<name>-resource-groups-volume-coordinator' (36-char suffix), so the name must be at most ${local.name_budget_resource_groups} characters."
    }
  }

  depends_on = [
    kubernetes_namespace_v1.trino,
    kubernetes_secret_v1.auth,
    kubernetes_secret_v1.internal,
  ]
}
