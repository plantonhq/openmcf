# KubernetesSolr Terraform module.
#
# Deploys one Apache Solr Operator-managed SolrCloud cluster:
#
#   1. the namespace (optional, create_namespace),
#   2. the solr.apache.org/v1beta1 SolrCloud CR itself.
#
# Everything else — the node StatefulSet, the provided ZooKeeper ensemble,
# services, PVCs, the basic-auth bootstrap Secret, Ingress exposure — is
# the operator's to create from the SolrCloud spec; the module renders the
# CR and exports the operator's deterministic names.
#
# The CR applies through kubectl_manifest (alekc/kubectl): unlike the
# hashicorp provider's kubernetes_manifest resource it needs no cluster
# connection at plan time — a SolrCloud can be PLANNED before the Solr
# operator's CRDs exist, which is what lets an infra chart deploy the
# operator and its clusters in one run (and lets offline plan proofs work).
#
# No wait_for block, deliberately: cluster readiness depends on the
# operator (image pulls, ZooKeeper quorum, node startup) that is not part
# of applying the resource — the same never-block-on-a-controller posture
# as the sibling operator-CR modules. Pulumi equivalent: the typed
# SolrCloud resource without await annotations.
#
# NAMING CONTRACT: every object the operator creates derives from
# metadata.name — StatefulSet `<name>-solrcloud`, common Service
# `<name>-solrcloud-common`, generated basic-auth Secret
# `<name>-solrcloud-basic-auth`, provided ZooKeeper client service
# `<name>-solrcloud-zookeeper-client:2181` — so the outputs derive them
# blind and both engines agree byte-for-byte.

# The optional namespace. Created before the cluster; deleted with the
# resource. Pre-existing-namespace deployments leave create_namespace
# false.
resource "kubernetes_namespace_v1" "namespace" {
  count = try(var.spec.create_namespace, false) ? 1 : 0

  metadata {
    name   = local.namespace
    labels = local.labels
  }
}

# ---- the SolrCloud ------------------------------------------------------------
resource "kubectl_manifest" "solr_cloud" {
  yaml_body = yamlencode(local.solrcloud_manifest)

  server_side_apply = true

  # BACKGROUND deletion, explicitly: the OPERATOR owns the SolrCloud's
  # cascade (its finalizer deletes the StatefulSet, services and the
  # provided ZooKeeper). Foreground propagation DEADLOCKS the teardown —
  # verified live: the foregroundDeletion finalizer waits for the child
  # ZookeeperCluster while the zookeeper-operator keeps reconciling it
  # back to life, and the fixture uninstall then hangs on the
  # zookeeper-operator's pre-delete hook. The provider would default to
  # Foreground whenever wait is set — never rely on that default here.
  # Pulumi twin: the pulumi.com/deletionPropagationPolicy annotation.
  delete_cascade = "Background"

  depends_on = [
    kubernetes_namespace_v1.namespace,
  ]
}
