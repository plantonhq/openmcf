# KubernetesOpenSearchOperator Terraform module.
#
# Installs the OpenSearch Kubernetes Operator from the official
# opensearch-operator Helm chart as a single Helm release named after
# metadata.name. The operator reconciles OpenSearchCluster custom
# resources (declared through KubernetesOpenSearch) into running search
# clusters with managed TLS, security bootstrap, rolling upgrades and
# Dashboards.
#
# CRD LIFECYCLE: the chart templates its CRDs as release-owned resources
# with NO keep-on-uninstall knob upstream, so a Helm-owned install would
# cascade-delete every OpenSearchCluster (and its data) on uninstall. The
# module therefore OWNS the CRDs through the catalog's derive-branch
# primitive (the generated helm_crds.tf): the pinned chart is rendered at
# plan time with the release's own values plus installCRDs turned on, each
# CustomResourceDefinition is applied keyed by its own name as a kept
# resource (retained on destroy unless crds.keep_on_uninstall is false;
# re-adopted on reinstall; refused when the manifest lowers chart_version
# below what the cluster carries), and the release below installs with
# skip_crds = true and installCRDs pinned false so Helm never touches them.
# The release depends on the CRDs so the operator never starts against an
# unregistered API group.
#
# The typed spec renders into chart values (locals.typed_values); the
# helm_values escape hatch is the SECOND values document and the
# load-bearing re-pin the THIRD (locals.helm_release_values) — the exact
# semantic twin of the Pulumi module's buildHelmValues + mergeMaps.

# The optional installation namespace. Created before the release; deleted
# with the resource (pre-existing-namespace installs leave create_namespace
# false).
resource "kubernetes_namespace_v1" "opensearch_operator" {
  count = try(var.spec.create_namespace, false) ? 1 : 0

  metadata {
    name   = local.namespace
    labels = local.labels
  }
}

# The operator release.
resource "helm_release" "opensearch_operator" {
  name       = local.release_name
  repository = local.helm_chart_repo
  chart      = local.helm_chart_name
  version    = local.chart_version
  namespace  = local.namespace

  # The module owns namespace creation (create_namespace flag).
  create_namespace = false

  # Wait for the operator to become Available — an operator that never
  # becomes ready (an unpullable image from a private mirror is the
  # classic case) should fail THIS apply with a readiness timeout, not
  # surface later as OpenSearch clusters that mysteriously never
  # reconcile.
  wait            = true
  atomic          = true
  cleanup_on_fail = true
  timeout         = 600

  # The same values list the CRD render consumed (see locals.tf for the
  # merge order and the re-pin).
  values = local.helm_release_values

  # The module owns the CRDs (helm_crds.tf); Helm must never install its
  # own copy of them, whichever way installCRDs is set.
  skip_crds = true

  depends_on = [
    kubernetes_namespace_v1.opensearch_operator,
    kubectl_manifest.helm_crds,
  ]
}
