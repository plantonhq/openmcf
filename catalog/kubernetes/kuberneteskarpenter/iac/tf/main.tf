# KubernetesKarpenter Terraform module.
#
# Installs Karpenter from the official OCI-served Helm charts as TWO real
# Helm releases in the same namespace: the companion karpenter-crd chart
# (when crds.install — upstream's supported mechanism for keeping CRDs
# upgradable; Helm installs the copies bundled inside the main chart once
# and NEVER upgrades them) and the karpenter controller chart, which
# depends on it and skips its bundled CRDs unconditionally.
#
# Both release names are FIXED ("karpenter-crd" / "karpenter"): Karpenter
# owns the cluster-wide karpenter.sh label domain, its CRDs, and node
# lifecycle — one installation per cluster is an upstream constraint.
#
# The typed spec renders into controller-chart values
# (locals.typed_values); the helm_values escape hatch is passed as a SECOND
# values document, which the provider merges over the first with Helm -f
# semantics — the exact semantic twin of the Pulumi module's
# buildHelmValues + mergeMaps.

# The optional installation namespace. Created before the releases; deleted
# with the resource (pre-existing-namespace installs — always the case for
# kube-system, upstream's recommended home — leave create_namespace false).
resource "kubernetes_namespace_v1" "karpenter" {
  count = try(var.spec.create_namespace, false) ? 1 : 0

  metadata {
    name   = local.namespace
    labels = local.labels
  }
}

# The CRD release. Its whole values surface is one knob: the keep
# annotation stamped onto every CRD it templates (see locals.crd_values) —
# without it a plain uninstall cascade-deletes every NodePool/EC2NodeClass/
# NodeClaim in the cluster.
#
# OCI ENGINE ASYMMETRY: this provider takes repository =
# "oci://public.ecr.aws/karpenter" plus the bare chart name and joins them
# internally; Pulumi's helm.v3.Release needs the JOINED
# "oci://public.ecr.aws/karpenter/karpenter-crd" chart reference with no
# repository opts. Same chart bytes, different wiring — keep both sides of
# the split in lockstep.
resource "helm_release" "crds" {
  count = local.crds_install ? 1 : 0

  name       = local.crd_release_name
  repository = local.helm_oci_repo
  chart      = local.crd_chart_name
  version    = local.chart_version
  namespace  = local.namespace

  # The module owns namespace creation (create_namespace flag).
  create_namespace = false

  # Same atomic/wait posture as the controller release: a CRD chart that
  # fails to apply should fail THIS apply cleanly, never leave half the
  # CRDs behind for the controller release to trip over.
  wait            = true
  atomic          = true
  cleanup_on_fail = true
  timeout         = 600

  values = length(local.crd_values) > 0 ? [yamlencode(local.crd_values)] : []

  depends_on = [kubernetes_namespace_v1.karpenter]
}

# The controller release.
resource "helm_release" "controller" {
  name       = local.release_name
  repository = local.helm_oci_repo # see the OCI ENGINE ASYMMETRY note above
  chart      = local.helm_chart_name
  version    = local.chart_version
  namespace  = local.namespace

  # The module owns namespace creation (create_namespace flag).
  create_namespace = false

  # The CRD release owns the CRDs (upstream's upgradable path); skipping
  # the controller chart's bundled copies UNCONDITIONALLY keeps this
  # release's shape deterministic whether or not crds.install is on.
  # Identical SkipCrds on the Pulumi side.
  skip_crds = true

  # Wait for the controller to become Available — a Karpenter that never
  # becomes ready (a ServiceMonitor rendered without the Prometheus
  # operator CRDs, a bad IRSA trust policy) should fail THIS apply with a
  # readiness timeout, not surface later as pods that stay Pending forever
  # because no nodes ever appear.
  wait            = true
  atomic          = true
  cleanup_on_fail = true
  timeout         = 600

  # Two documents, merged in order by the provider (helm -f semantics):
  # the typed rendering first, the user's escape hatch last.
  values = concat(
    [yamlencode(local.typed_values)],
    try(var.spec.helm_values, "") != "" ? [var.spec.helm_values] : []
  )

  # The controller's pods reconcile NodePools/NodeClaims from the moment
  # they start — the CRDs must exist first.
  depends_on = [kubernetes_namespace_v1.karpenter, helm_release.crds]
}
