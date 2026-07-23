# KubernetesMysql Terraform module.
#
# Deploys one Percona-operator-managed MySQL (XtraDB Cluster):
#
#   1. the namespace (optional, create_namespace),
#   2. declared-credential Secrets (user passwords, backup-storage keys)
#      — secrets always travel via secret references, never inline in a
#      custom resource,
#   3. the pxc.percona.com/v1 PerconaXtraDBCluster CR itself.
#
# The CR applies through kubectl_manifest (alekc/kubectl): unlike the
# hashicorp provider's kubernetes_manifest resource it needs no cluster
# connection at plan time — the cluster can be PLANNED before the Percona
# operator's CRDs exist, which is what lets an infra chart deploy the
# operator and its databases in one run (and lets offline plan proofs
# work).
#
# No wait_for block, deliberately: cluster readiness depends on the
# operator (image pulls, Galera bootstrap, proxy rollout) that is not
# part of applying the resource — the same never-block-on-a-controller
# posture as the KubernetesPostgres exemplar. Pulumi equivalent:
# CustomResource without await annotations.

resource "kubernetes_namespace_v1" "namespace" {
  count = try(var.spec.create_namespace, false) ? 1 : 0

  metadata {
    name   = local.namespace
    labels = local.labels
  }
}

resource "kubernetes_secret_v1" "user_password_secret" {
  for_each = local.user_password_secrets

  metadata {
    name      = each.key
    namespace = local.namespace
    labels    = local.labels
  }

  data = {
    password = each.value
  }

  # The operator CO-OWNS this Secret's annotations: after applying a
  # password it stamps a `percona.com/<cluster>-<user>-hash` marker (its
  # rotation detector) onto the object. The module owns the data; the
  # annotations are the operator's — never fight them.
  lifecycle {
    ignore_changes = [metadata[0].annotations]
  }

  depends_on = [kubernetes_namespace_v1.namespace]
}

resource "kubernetes_secret_v1" "backup_credentials_secret" {
  for_each = local.backup_credential_secrets

  metadata {
    name      = each.key
    namespace = local.namespace
    labels    = local.labels
  }

  data = each.value

  depends_on = [kubernetes_namespace_v1.namespace]
}

resource "kubectl_manifest" "mysql" {
  yaml_body = yamlencode(local.mysql_manifest)

  server_side_apply = true

  depends_on = [
    kubernetes_namespace_v1.namespace,
    kubernetes_secret_v1.user_password_secret,
    kubernetes_secret_v1.backup_credentials_secret,
  ]
}
