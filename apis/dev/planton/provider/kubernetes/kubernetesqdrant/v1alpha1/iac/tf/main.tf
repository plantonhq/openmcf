# KubernetesQdrant Terraform module.
#
# Installs Qdrant from the official Helm chart as a real Helm release. The
# typed spec renders into chart values (locals.helm_values); API keys stay
# chart-owned (generated once via the chart's lookup, or read from an
# existing Secret AT TEMPLATE TIME — the referenced Secret must exist before
# the apply); the helm_values escape hatch is passed as a SECOND values
# document, which the provider merges over the first with Helm -f semantics
# — the exact semantic twin of the Pulumi module's buildHelmValues +
# mergeMaps.

# The optional installation namespace. Created before the release; deleted
# with the resource.
resource "kubernetes_namespace_v1" "qdrant" {
  count = var.spec.create_namespace ? 1 : 0

  metadata {
    name   = local.namespace
    labels = local.labels
  }
}

resource "helm_release" "qdrant" {
  name       = local.release_name
  repository = local.helm_chart_repo
  chart      = local.helm_chart_name
  version    = local.chart_version
  namespace  = local.namespace

  # The module owns namespace creation (create_namespace flag).
  create_namespace = false

  # Wait for the cluster to become Ready — a database that never starts
  # (bad image, unschedulable pod, unbindable volume) should fail THIS
  # apply, not the first client connection. Multi-node consensus bootstrap
  # is quick; storage recovery on big volumes is what the generous budget
  # covers.
  wait            = true
  atomic          = true
  cleanup_on_fail = true
  timeout         = 600

  # Documents merged in order by the provider (helm -f semantics): the
  # typed rendering first, the user's escape hatch second — and
  # fullnameOverride re-pinned LAST, the one deliberate exception to the
  # escape hatch's last-word contract (twin of the Pulumi module). The
  # Service / `-headless` / `-apikey` Secret names — and the exported
  # outputs built from them — all derive from the fullname; letting an
  # override move it would break every output.
  values = concat(
    [yamlencode(local.helm_values)],
    try(var.spec.helm_values, "") != "" ? [var.spec.helm_values] : [],
    [yamlencode({ fullnameOverride = local.release_name })]
  )

  depends_on = [
    kubernetes_namespace_v1.qdrant,
  ]
}
