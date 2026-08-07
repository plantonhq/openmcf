# KubernetesGhaRunnerScaleSetController Terraform module.
#
# Installs the GitHub Actions Runner Scale Set controller from the
# official OCI chart as a real Helm release. The typed spec renders into
# chart values (locals.typed_helm_values); the helm_values escape hatch
# is passed as a SECOND values document, which the provider merges over
# the first with Helm -f semantics — the exact semantic twin of the
# Pulumi module's buildHelmValues + mergeMaps.
#
# OCI WIRING: the Terraform provider takes repository = the OCI registry
# path plus the bare chart name and joins them internally; the Pulumi
# twin passes the joined "oci://.../<chart>" string as the chart
# reference. Same chart bytes, different wiring — keep both sides of the
# split in lockstep.
#
# CRD LIFECYCLE: the chart installs the actions.github.com CRDs with the
# release and they are REMOVED with it — destroying the controller
# cascade-deletes every runner scale set on the cluster (the spec's CRD
# note carries the warning).

# The optional installation namespace. Created before the release; deleted
# with the resource.
resource "kubernetes_namespace_v1" "gha_runner_scale_set_controller" {
  count = var.spec.create_namespace ? 1 : 0

  metadata {
    name   = local.namespace
    labels = local.labels
  }
}

resource "helm_release" "gha_runner_scale_set_controller" {
  name       = local.release_name
  repository = local.helm_oci_repo
  chart      = local.helm_chart_name
  version    = local.chart_version
  namespace  = local.namespace

  # The module owns namespace creation (create_namespace flag).
  create_namespace = false

  # Wait for the controller to become Ready — a manager that never starts
  # should fail THIS apply, not the first runner scale set install.
  wait            = true
  atomic          = true
  cleanup_on_fail = true
  timeout         = 300

  # Documents merged in order by the provider (helm -f semantics): the
  # typed rendering first, the user's escape hatch second — and
  # fullnameOverride re-pinned LAST, the one deliberate exception to the
  # escape hatch's last-word contract (twin of the Pulumi module). The
  # ServiceAccount name output derives from the fullname; letting an
  # override move it would break the scale-set discovery handle.
  values = concat(
    [yamlencode(local.typed_helm_values)],
    try(var.spec.helm_values, "") != "" ? [var.spec.helm_values] : [],
    [yamlencode({ fullnameOverride = local.release_name })]
  )

  depends_on = [
    kubernetes_namespace_v1.gha_runner_scale_set_controller,
  ]
}
