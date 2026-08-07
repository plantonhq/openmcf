# KubernetesStatefulSet Terraform module.
#
# Deploys a stateful workload: optional namespace (this file), env and
# image-pull satellite Secrets (secret.tf), the headless governing Service
# (service.tf), the apps/v1 StatefulSet with its volume claim templates
# (statefulset.tf), and optional PDB (pdb.tf).
#
# Identity and exposure are composed, not created: pods run as the
# ServiceAccount referenced in spec.pod.service_account, and external exposure
# attaches through first-class ingress kinds referencing this workload's
# exported Service handle. This module never creates ServiceAccounts, RBAC
# objects, certificates, gateways, or routes. There is deliberately no HPA:
# stateful members join and leave through application-aware procedures.

resource "kubernetes_namespace" "this" {
  count = var.spec.create_namespace ? 1 : 0

  metadata {
    name   = local.namespace
    labels = local.final_labels
  }
}
