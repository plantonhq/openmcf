# KubernetesSolrOperator Terraform module.
#
# Installs the Apache Solr Operator from the official solr-operator Helm
# chart as a single Helm release named after metadata.name. The operator
# reconciles SolrCloud custom resources (declared through KubernetesSolr)
# into running Solr clusters, plus SolrBackup and SolrPrometheusExporter
# resources.
#
# CRD LIFECYCLE: unlike most operator charts, the solr-operator chart
# ships NO CRDs — they are separate release artifacts. The module OWNS
# them: the four staged files at ../crds (the three solr.apache.org CRDs
# plus the ZookeeperCluster CRD of the bundled zookeeper-operator
# dependency) are applied before the release with apply_only, so a
# destroy removes the operator but NEVER the CRDs and SolrCloud resources
# are never cascade-deleted. The bundled subchart's own CRD switch is
# pinned off (zookeeper-operator.crd.create = false) so the
# ZookeeperCluster CRD never falls under Helm's delete-on-uninstall
# lifecycle.
#
# The typed spec renders into chart values (locals.typed_values); the
# helm_values escape hatch is passed as a SECOND values document, which
# the provider merges over the first with Helm -f semantics — the exact
# semantic twin of the Pulumi module's buildHelmValues + mergeMaps.

# The optional installation namespace. Created before the release; deleted
# with the resource (pre-existing-namespace installs leave create_namespace
# false).
resource "kubernetes_namespace_v1" "solr_operator" {
  count = try(var.spec.create_namespace, false) ? 1 : 0

  metadata {
    name   = local.namespace
    labels = local.labels
  }
}

# The module-owned CRDs, one resource per CRD keyed by the CRD's OWN
# metadata.name (locals.crd_documents does the per-file document
# splitting and comment-header filtering).
#
# KEEP-ON-UNINSTALL: apply_only = true makes the provider's Delete a
# NO-OP (verified in the alekc/kubectl provider source) — `terraform
# destroy` removes the operator release but leaves these CRDs (and
# therefore every SolrCloud/SolrBackup/ZookeeperCluster resource
# cluster-wide) untouched. This is the exact twin of the Pulumi module's
# retainOnDelete on each CRD.
#
# server_side_apply matters twice over: the SolrCloud CRD's schema blows
# past the client-side last-applied-configuration annotation size limit,
# and SSA field ownership lets a restaged (upgraded) CRD file apply as an
# in-place update.
resource "kubectl_manifest" "solr_operator_crds" {
  for_each = local.crd_documents

  yaml_body = each.value

  server_side_apply = true
  apply_only        = true
}

# The operator release.
resource "helm_release" "solr_operator" {
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
  # surface later as SolrCloud resources that mysteriously never
  # reconcile.
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

  # The CRDs must exist before the operator starts (its controllers
  # watch these types immediately, and the bundled zookeeper-operator
  # refuses to start without the ZookeeperCluster CRD present).
  depends_on = [
    kubernetes_namespace_v1.solr_operator,
    kubectl_manifest.solr_operator_crds,
  ]
}
