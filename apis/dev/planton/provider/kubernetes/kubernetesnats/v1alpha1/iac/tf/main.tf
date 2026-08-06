# KubernetesNats Terraform module.
#
# Installs NATS from the official Helm chart as a real Helm release. The
# typed spec renders into chart values (locals.helm_values); declared
# users' passwords are generated below into the `<name>-auth` Secret
# BEFORE the release and reach the server as secretKeyRef env vars the
# rendered config references (`$NATS_PW_<i>`) — nothing credential-bearing
# transits values; the helm_values escape hatch is passed as a SECOND
# values document, which the provider merges over the first with Helm -f
# semantics — the exact semantic twin of the Pulumi module's
# buildHelmValues + mergeMaps.

# The optional installation namespace. Created before the release; deleted
# with the resource.
resource "kubernetes_namespace_v1" "nats" {
  count = var.spec.create_namespace ? 1 : 0

  metadata {
    name   = local.namespace
    labels = local.labels
  }
}

# One generated password per declared user (flat and account users alike),
# keyed by USERNAME — the credential's identity. Usernames are unique
# across the whole auth block (CEL-enforced), the password follows its
# user through spec reorders (an index/env keying would silently SWAP
# passwords between users when the list order changes), and a renamed
# user is honestly a NEW credential. Twin: the Pulumi module's
# per-username RandomPassword resources.
resource "random_password" "user" {
  for_each = { for u in local.all_users_meta : u.username => u }

  # LETTERS ONLY, and longer to compensate the smaller alphabet
  # (40 letters ≈ 228 bits — far past any practical bar). The server
  # RESOLVES the $NATS_PW_<i> env reference and RE-PARSES the resolved
  # value through its own config parser (verified live: a generated
  # password with a digit/symbol prefix crash-loops every server with
  # "variable reference for 'NATS_PW_0' ... could not be parsed").
  # Digits can lex as numbers, and '-' '#' '$' '{' quotes are all
  # structural tokens — a pure-letter password is the only shape the
  # parser can never misread. Twin: the Pulumi module's RandomPassword
  # with the same alphabet.
  length  = 40
  special = false
  numeric = false

  # The generation-shape arguments are ignored after creation so an
  # IMPORTED credential never silently regenerates: rotation stays an
  # explicit verb, never plan fallout. Twin: the Pulumi module's
  # IgnoreChanges on the same argument set.
  lifecycle {
    ignore_changes = [
      length, special, upper, lower, numeric,
      min_lower, min_numeric, min_special, min_upper, override_special,
    ]
  }
}

# The credential Secret — ONE KEY PER USERNAME, each key's value that
# user's password. The ONLY place the credentials land: the chart values
# carry secretKeyRef env wiring plus `$NATS_PW_<i>` config references the
# server resolves from environment at load. Clients read the same Secret
# (the exported auth_secret_name handle).
resource "kubernetes_secret_v1" "auth" {
  count = local.auth_enabled ? 1 : 0

  metadata {
    name      = local.auth_secret_name
    namespace = local.namespace
    labels    = local.labels
  }

  type = "Opaque"

  data = {
    for u in local.all_users_meta : u.username => random_password.user[u.username].result
  }

  depends_on = [
    kubernetes_namespace_v1.nats,
  ]
}

resource "helm_release" "nats" {
  name       = local.release_name
  repository = local.helm_chart_repo
  chart      = local.helm_chart_name
  version    = local.chart_version
  namespace  = local.namespace

  # The module owns namespace creation (create_namespace flag).
  create_namespace = false

  # Wait for the StatefulSet to become Ready — a server that never
  # starts (bad TLS Secret name, malformed config merge) should fail
  # THIS apply, not the first client connection.
  wait            = true
  atomic          = true
  cleanup_on_fail = true
  timeout         = 300

  # Documents merged in order by the provider (helm -f semantics): the
  # typed rendering first, the user's escape hatch second — and
  # fullnameOverride re-pinned LAST, the one deliberate exception to the
  # escape hatch's last-word contract (twin of the Pulumi module). Every
  # child name (`<name>`, `<name>-headless`, `<name>-box`, ...) — and
  # the exported outputs built from them — derive from the fullname;
  # letting an override move it would break every output.
  values = concat(
    [yamlencode(local.helm_values)],
    try(var.spec.helm_values, "") != "" ? [var.spec.helm_values] : [],
    [yamlencode({ fullnameOverride = local.release_name })]
  )

  lifecycle {
    # FAIL LOUDLY on names past the chart's fullname budget: derived
    # child names (`<fullname>-box-contents` is the longest, +13 chars)
    # truncate SILENTLY at 63 characters past a 50-character fullname,
    # breaking the naming contract the exported outputs are built on.
    # Twin: the Pulumi module's Resources() guard.
    precondition {
      condition     = length(var.metadata.name) <= 50
      error_message = "The nats chart derives child names from the resource name and silently truncates past 50 characters (63 minus its longest derived suffix), which would break the naming contract — use a name of at most 50 characters."
    }
  }

  depends_on = [
    kubernetes_namespace_v1.nats,
    kubernetes_secret_v1.auth,
  ]
}
