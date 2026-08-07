# KubernetesOpenSearch Terraform module.
#
# Deploys one operator-managed OpenSearch cluster:
#
#   1. the namespace (optional, create_namespace),
#   2. the opensearch.opster.io/v1 OpenSearchCluster CR itself — the ONLY
#      custom resource: node StatefulSets, Services, TLS Secrets, the
#      security-plugin admin bootstrap and the optional Dashboards
#      deployment are all operator-created from it. No ingress resources
#      — exposure composes from first-class kinds referencing the
#      exported handles.
#
# The CR applies through kubectl_manifest (alekc/kubectl): unlike the
# hashicorp provider's kubernetes_manifest resource it needs no cluster
# connection at plan time — a cluster can be PLANNED before the OpenSearch
# operator's CRDs exist, which is what lets an infra chart deploy the
# operator and its clusters in one run (and lets offline plan proofs work).
#
# No wait_for block, deliberately: cluster readiness depends on the
# operator (image pulls, quorum bootstrap, security initialization) that is
# not part of applying the resource — the same never-block-on-a-controller
# posture as the sibling operator-CR modules. Pulumi equivalent:
# the typed OpenSearchCluster resource without await annotations.
#
# NAMING CONTRACT: every object the operator creates derives from
# metadata.name — StatefulSets `<name>-<pool>`, the main Service `<name>`
# (this module pins general.serviceName to it), the discovery Service
# `<name>-discovery`, the admin credentials Secret `<name>-admin-password`,
# the Dashboards deployment/Service `<name>-dashboards`.

# The optional namespace. Created before the cluster (the CR is
# namespaced); deleted with the resource. Pre-existing-namespace
# deployments leave create_namespace false.
resource "kubernetes_namespace_v1" "namespace" {
  count = try(var.spec.create_namespace, false) ? 1 : 0

  metadata {
    name   = local.namespace
    labels = local.labels
  }
}

# ---- the OpenSearchCluster ------------------------------------------------
resource "kubectl_manifest" "opensearch_cluster" {
  yaml_body = yamlencode(local.cluster_manifest)

  server_side_apply = true

  depends_on = [
    kubernetes_namespace_v1.namespace,
  ]
}
