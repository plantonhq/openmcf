# KubernetesSignoz Terraform module.
#
# Installs SigNoz from the official Helm chart as a real Helm release. The
# typed spec renders into chart values (locals.helm_values); exactly one
# database arm renders (bundled ClickHouse XOR external); the bundled
# ClickHouse admin password is module-GENERATED, delivered through
# set_sensitive (never a values document) and exported through the
# module-owned "<name>-clickhouse-auth" Secret; the helm_values escape
# hatch is passed as a SECOND values document, with fullnameOverride and
# the ClickHouse fullname re-pinned in a THIRD — the exact semantic twin
# of the Pulumi module's buildHelmValues + mergeMaps + re-pin.

# The optional installation namespace. Created before the release; deleted
# with the resource.
resource "kubernetes_namespace_v1" "signoz" {
  count = var.spec.create_namespace ? 1 : 0

  metadata {
    name   = local.namespace
    labels = local.labels
  }
}

# The bundled ClickHouse admin password, generated once and held in state
# as a sensitive value. Created only on the bundled arm — the external arm
# brings its own Secret. The chart's publicly-documented default password
# NEVER ships.
resource "random_password" "clickhouse" {
  count = local.is_external ? 0 : 1

  length  = 24
  special = false

  # The generation-shape arguments are ignored after creation so an
  # IMPORTED credential never silently regenerates: random_password's
  # import carries only the VALUE, the importer assumes the provider's
  # own generation defaults (special=true — verified live), and every
  # argument is ForceNew — without this, the first plan after an import
  # proposes replacing (rotating) the live credential. Rotation stays
  # an explicit verb (taint / destroy-recreate), never plan fallout.
  lifecycle {
    ignore_changes = [length, special, upper, lower, numeric, min_lower, min_numeric, min_special, min_upper, override_special]
  }
}

# The composition handle for the bundled credential (keys
# username/password, username "admin" — the chart's bundled-arm user).
# The chart itself offers no Secret for the bundled arm (it renders the
# password into the installation object and container env — the upstream
# grain); this module-owned Secret is what downstream kinds and operators
# reference. Pulumi twin: random.RandomPassword + core.Secret in
# clickhouse_auth_secret.go with the same name and keys.
resource "kubernetes_secret_v1" "clickhouse_auth" {
  count = local.is_external ? 0 : 1

  metadata {
    name      = local.clickhouse_auth_secret_name
    namespace = local.namespace
    labels    = local.labels
  }

  type = "Opaque"

  data = {
    username = "admin"
    password = sensitive(random_password.clickhouse[0].result)
  }

  depends_on = [kubernetes_namespace_v1.signoz]
}

resource "helm_release" "signoz" {
  name       = local.release_name
  repository = local.helm_chart_repo
  chart      = local.helm_chart_name
  version    = local.chart_version
  namespace  = local.namespace

  # The module owns namespace creation (create_namespace flag).
  create_namespace = false

  # Wait for the whole platform to become Ready — a SigNoz whose
  # ClickHouse never binds storage or whose schema migration fails should
  # fail THIS apply, not the first trace query. The budget covers the
  # ClickHouse operator reconcile + schema migration on a cold cluster.
  wait            = true
  atomic          = true
  cleanup_on_fail = true
  timeout         = 1200

  # Documents merged in order by the provider (helm -f semantics): the
  # typed rendering first, the user's escape hatch second, and the two
  # fullname pins re-asserted LAST — the one deliberate exception to the
  # escape hatch's last-word contract (twin of the Pulumi module). Every
  # child name — and the exported outputs built from them — derives from
  # these fullnames; letting an override move them would break every
  # output.
  values = concat(
    [yamlencode(local.helm_values)],
    try(var.spec.helm_values, "") != "" ? [var.spec.helm_values] : [],
    [yamlencode(merge(
      { fullnameOverride = local.release_name },
      local.is_external ? {} : { clickhouse = { fullnameOverride = local.clickhouse_fullname } },
    ))]
  )

  # The bundled admin password travels OUTSIDE the values documents so it
  # never appears in a plan diff. On the external arm nothing is set —
  # the chart reads the referenced Secret directly. set_sensitive is a
  # list-of-object ATTRIBUTE in helm provider v3 (the framework rewrite),
  # not a legacy block — verified against the provider source.
  set_sensitive = local.is_external ? [] : [
    { name = "clickhouse.password", value = random_password.clickhouse[0].result, type = "string" }
  ]

  # The chart wraps the ClickHouseInstallation's operator-generated
  # StatefulSet names in ~27 characters of scaffolding within Kubernetes'
  # 63-character cap — an over-long resource name corrupts the naming
  # contract the outputs promise. Fail THIS plan loudly instead (twin:
  # the Pulumi module's MaxNameLength guard).
  lifecycle {
    precondition {
      condition     = local.name_within_budget
      error_message = "metadata.name '${local.release_name}' is ${length(local.release_name)} characters — the signoz chart's child-name budget allows at most ${local.max_name_length} (the bundled ClickHouse composes names like chi-<name>-clickhouse-cluster-0-0 within Kubernetes' 63-character cap)."
    }
  }

  depends_on = [
    kubernetes_namespace_v1.signoz,
    kubernetes_secret_v1.clickhouse_auth,
  ]
}
