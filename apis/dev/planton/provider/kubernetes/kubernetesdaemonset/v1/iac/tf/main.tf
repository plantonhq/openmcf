# KubernetesDaemonSet Terraform module.
#
# Deploys a node agent: optional namespace (this file), env and image-pull
# satellite Secrets (secret.tf), and the apps/v1 DaemonSet (daemonset.tf).
# There is no replica count — node membership IS the replica count — and no
# Service or ingress: clients that must reach an agent do so on its node via
# per-container host_port or pod.host_network.
#
# Identity and permissions are composed, not created: pods run as the
# ServiceAccount referenced in spec.pod.service_account, and API permissions
# come from KubernetesRbac grants targeting that identity. This module never
# creates ServiceAccounts, RBAC objects, certificates, gateways, or routes.

resource "kubernetes_namespace" "this" {
  count = var.spec.create_namespace ? 1 : 0

  metadata {
    name   = local.namespace
    labels = local.final_labels
  }
}
