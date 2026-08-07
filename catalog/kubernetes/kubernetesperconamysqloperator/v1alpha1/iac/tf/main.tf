# KubernetesPerconaMysqlOperator Terraform module.
#
# Installs the Percona Operator for MySQL (based on Percona XtraDB
# Cluster) from the official pxc-operator Helm chart as a single Helm
# release named after metadata.name. The operator reconciles
# PerconaXtraDBCluster custom resources (declared through KubernetesMysql)
# into Galera clusters with automated failover, HAProxy/ProxySQL routing,
# and scheduled XtraBackup backups.
#
# CRD LIFECYCLE: the chart ships the PerconaXtraDBCluster CRDs in its
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
resource "kubernetes_namespace_v1" "percona_mysql_operator" {
  count = try(var.spec.create_namespace, false) ? 1 : 0

  metadata {
    name   = local.namespace
    labels = local.labels
  }
}

# The operator release.
resource "helm_release" "percona_mysql_operator" {
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
  # surface later as PerconaXtraDBCluster resources that mysteriously
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

  depends_on = [
    kubernetes_namespace_v1.percona_mysql_operator,
    kubernetes_validating_webhook_configuration_v1.pxc_validation,
  ]
}

# The operator's CR-validation webhook — MODULE-OWNED in the widened-watch
# arms, deliberately. Upstream behavior: an operator with cluster-scoped
# RBAC (which widened watch grants) registers this ONE fixed-name,
# cluster-scoped webhook at startup, pointing at a Service in its own
# namespace with failurePolicy Fail — and NOTHING ever removes it. The
# object cannot ride the operator Deployment's ownerReference (Kubernetes
# never garbage-collects a cluster-scoped dependent of a namespaced
# owner), and an operator that finds the object already present updates
# only the CA bundle, never the service pointer. Left to the operator,
# uninstalling it therefore strands a Fail-closed webhook whose service no
# longer exists — bricking every future PerconaXtraDBCluster admission in
# the cluster.
#
# The module renders the object FIRST (the release depends on it), so the
# operator's startup create hits AlreadyExists and merely refreshes the CA
# bundle into the module-declared shape — and destroy removes the webhook
# with the resource. The CA bundle and the operator-stamped ownerReference
# are the operator's to manage; ignore_changes keeps the module from
# fighting them.
#
# Own-namespace installs render nothing: the chart's namespaced Role
# carries no admissionregistration permissions, the operator's own
# registration attempt is Forbidden (which it treats as a soft skip), and
# the webhook simply does not exist — the upstream posture, preserved.
#
# One cluster carries at most ONE widened-watch operator: the webhook name
# is fixed upstream, so a second widened installation would contend for
# the same object (documented in the kind's README).
resource "kubernetes_validating_webhook_configuration_v1" "pxc_validation" {
  count = local.watch_widened ? 1 : 0

  metadata {
    name   = "percona-xtradbcluster-webhook"
    labels = local.labels
  }

  webhook {
    name                      = "validationwebhook.pxc.percona.com"
    admission_review_versions = ["v1"]

    client_config {
      service {
        name      = "percona-xtradb-cluster-operator"
        namespace = local.namespace
        path      = "/validate-percona-xtradbcluster"
        port      = 443
      }
    }

    side_effects   = "None"
    failure_policy = "Fail"

    rule {
      api_groups   = ["pxc.percona.com"]
      api_versions = ["*"]
      resources    = ["perconaxtradbclusters/*"]
      operations   = ["CREATE", "UPDATE"]
    }
  }

  lifecycle {
    # The operator patches the CA bundle in at startup (it issues the
    # serving certificate); the module never carries certificate material.
    ignore_changes = [webhook[0].client_config[0].ca_bundle]
  }
}
