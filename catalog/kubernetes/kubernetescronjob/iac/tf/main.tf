# KubernetesCronJob Terraform module.
#
# Runs work on a recurring schedule: optional namespace (this file), env and
# image-pull satellite Secrets (secret.tf), and the batch/v1 CronJob itself
# (cron_job.tf) — scheduling controls at the top level, the Job stamped out at
# each run from spec.job_template.
#
# Identity is composed, not created: pods run as the ServiceAccount referenced
# in spec.job_template.pod.service_account. This module never creates
# ServiceAccounts, RBAC objects, certificates, gateways, or routes — and
# CronJobs front no Service, so nothing here creates exposure either.

resource "kubernetes_namespace" "this" {
  count = var.spec.create_namespace ? 1 : 0

  metadata {
    name   = local.namespace
    labels = local.final_labels
  }
}
