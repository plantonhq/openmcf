# KubernetesPerconaMongoOperator Terraform module.
#
# Installs the Percona Operator for MongoDB from the official
# psmdb-operator Helm chart as a single Helm release named after
# metadata.name. The operator reconciles PerconaServerMongoDB custom
# resources (declared through KubernetesMongodb) into replica sets and
# sharded clusters with automated failover and Percona Backup for MongoDB.
#
# CRD LIFECYCLE: the chart ships the PerconaServerMongoDB CRDs in its
# Helm-native crds/ directory — installed on first install, never upgraded
# or deleted by Helm. Uninstalling the release therefore NEVER
# cascade-deletes the database clusters (the upstream safety posture).
#
# The typed spec renders into chart values (locals.typed_values); the
# helm_values escape hatch is passed as a SECOND values document, which
# the provider merges over the first with Helm -f semantics — the exact
# semantic twin of the Pulumi module's buildHelmValues + mergeMaps.

# The optional installation namespace. Created before the release; deleted
# with the resource (pre-existing-namespace installs leave create_namespace
# false).
resource "kubernetes_namespace_v1" "percona_mongo_operator" {
  count = try(var.spec.create_namespace, false) ? 1 : 0

  metadata {
    name   = local.namespace
    labels = local.labels
  }
}

# The operator release.
resource "helm_release" "percona_mongo_operator" {
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
  # surface later as PerconaServerMongoDB resources that mysteriously
  # never reconcile.
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

  depends_on = [kubernetes_namespace_v1.percona_mongo_operator]
}
