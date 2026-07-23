# KubernetesJob Terraform module.
#
# Runs work to completion: optional namespace (this file), env and image-pull
# satellite Secrets (secret.tf), and the batch/v1 Job itself (job.tf).
#
# Identity is composed, not created: pods run as the ServiceAccount referenced
# in spec.pod.service_account. This module never creates ServiceAccounts, RBAC
# objects, certificates, gateways, or routes — and Jobs front no Service, so
# nothing here creates exposure either.

resource "kubernetes_namespace" "this" {
  count = var.spec.create_namespace ? 1 : 0

  metadata {
    name   = local.namespace
    labels = local.final_labels
  }
}
