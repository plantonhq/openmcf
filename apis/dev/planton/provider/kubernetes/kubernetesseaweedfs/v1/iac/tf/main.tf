# KubernetesSeaweedFs Terraform module.
#
# Installs SeaweedFS from the official Helm chart as a real Helm release.
# The typed spec renders into chart values (locals.helm_values); S3
# credentials stay chart-owned (`<name>-s3-secret`, generated once and kept
# on uninstall) or come from a referenced existing config Secret; the
# admin-console credentials materialize as the "<name>-admin-auth" Secret;
# the helm_values escape hatch is passed as a SECOND values document, which
# the provider merges over the first with Helm -f semantics — the exact
# semantic twin of the Pulumi module's buildHelmValues + mergeMaps.

# The optional installation namespace. Created before the release; deleted
# with the resource.
resource "kubernetes_namespace_v1" "seaweedfs" {
  count = var.spec.create_namespace ? 1 : 0

  metadata {
    name   = local.namespace
    labels = local.labels
  }
}

# The admin-console password, generated once and held in state as a
# sensitive value. Created only when the console is enabled without an
# existing Secret.
resource "random_password" "admin" {
  count = local.create_admin_secret ? 1 : 0

  length  = 24
  special = false
}

# The console credentials Secret (keys user/password, user "admin"). The
# chart consumes it via admin.secret.existingSecret + userKey/pwKey — its
# OWN adminPassword value path is never used, so the credential never
# transits chart values. Pulumi twin: random.RandomPassword +
# core.Secret in admin_secret.go with the same name and keys.
resource "kubernetes_secret_v1" "admin_auth" {
  count = local.create_admin_secret ? 1 : 0

  metadata {
    name      = local.admin_auth_secret_name
    namespace = local.namespace
    labels    = local.labels
  }

  type = "Opaque"

  data = {
    user     = "admin"
    password = sensitive(random_password.admin[0].result)
  }

  depends_on = [kubernetes_namespace_v1.seaweedfs]
}

resource "helm_release" "seaweedfs" {
  name       = local.release_name
  repository = local.helm_chart_repo
  chart      = local.helm_chart_name
  version    = local.chart_version
  namespace  = local.namespace

  # The module owns namespace creation (create_namespace flag).
  create_namespace = false

  # Wait for every tier to become Ready — a store that never starts (bad
  # image, unschedulable pod, unbindable PVC) should fail THIS apply, not
  # the first S3 request. The budget covers a three-tier cold start plus
  # the bucket-creation hook.
  wait            = true
  atomic          = true
  cleanup_on_fail = true
  timeout         = 600

  # Documents merged in order by the provider (helm -f semantics): the
  # typed rendering first, the user's escape hatch second — and
  # fullnameOverride re-pinned LAST, the one deliberate exception to the
  # escape hatch's last-word contract (twin of the Pulumi module). Every
  # componentName child (`-master`, `-filer`, `-s3`, `-admin`, `-volume`)
  # and the chart-generated `-s3-secret` — and the exported outputs built
  # from them — all derive from the fullname; letting an override move it
  # would break every output.
  values = concat(
    [yamlencode(local.helm_values)],
    try(var.spec.helm_values, "") != "" ? [var.spec.helm_values] : [],
    [yamlencode({ fullnameOverride = local.release_name })]
  )

  depends_on = [
    kubernetes_namespace_v1.seaweedfs,
    kubernetes_secret_v1.admin_auth,
  ]
}
