# KubernetesValkey Terraform module.
#
# Installs Valkey from the official Helm chart as a real Helm release. The
# typed spec renders into chart values (locals.helm_values); declared ACL
# passwords materialize as the "<name>-auth" Kubernetes Secret the chart
# consumes via auth.usersExistingSecret; the helm_values escape hatch is
# passed as a SECOND values document, which the provider merges over the
# first with Helm -f semantics — the exact semantic twin of the Pulumi
# module's buildHelmValues + mergeMaps.
#
# The release is named after metadata.name (NOT a fixed chart name) and the
# chart's fullname is pinned to the same value: several Valkey instances
# coexist in one cluster, each rendering its own `<name>`,
# `<name>-headless`, and (replication) `<name>-read` Services.

# The optional installation namespace. Created before the release; deleted
# with the resource.
resource "kubernetes_namespace_v1" "valkey" {
  count = var.spec.create_namespace ? 1 : 0

  metadata {
    name   = local.namespace
    labels = local.labels
  }
}

# The declared ACL passwords, materialized as an Opaque Secret — ONE KEY
# PER USERNAME, each key's value that user's password. That layout is the
# chart's contract for auth.usersExistingSecret: its init script reads
# /valkey-users-secret/<passwordKey> where passwordKey defaults to the
# username (the module leaves passwordKey unset), and its metrics exporter
# reads the "default" key the same way. Because the rendered aclUsers carry
# no inline passwords, the chart renders no auth Secret of its own — this
# Secret is the only place the credentials land, and it never transits
# chart values.
resource "kubernetes_secret_v1" "auth" {
  count = local.auth_enabled ? 1 : 0

  metadata {
    name      = local.auth_secret_name
    namespace = local.namespace
    labels    = local.labels
  }

  type = "Opaque"

  data = {
    for user in var.spec.auth.users : user.name => user.password
  }

  depends_on = [kubernetes_namespace_v1.valkey]
}

resource "helm_release" "valkey" {
  name       = local.release_name
  repository = local.helm_chart_repo
  chart      = local.helm_chart_name
  version    = local.chart_version
  namespace  = local.namespace

  # The module owns namespace creation (create_namespace flag).
  create_namespace = false

  # Wait for the workload to become Ready — a store that never starts (bad
  # image, unschedulable pod, unbindable volume) should fail THIS apply,
  # not the first client connection. Replication starts pods one at a time
  # (OrderedReady) and each replica full-syncs before Ready, so the budget
  # is sized for a multi-pod StatefulSet, not a single Deployment.
  wait            = true
  atomic          = true
  cleanup_on_fail = true
  timeout         = 600

  # Two documents, merged in order by the provider (helm -f semantics):
  # the typed rendering first, the user's escape hatch last.
  values = [
    yamlencode(local.helm_values),
    try(var.spec.helm_values, ""),
  ]

  depends_on = [
    kubernetes_namespace_v1.valkey,
    kubernetes_secret_v1.auth,
  ]
}
