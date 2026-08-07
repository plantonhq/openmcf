# Provider requirements for the KubernetesCertificate module.
#
# kubectl (alekc/kubectl) applies the Certificate CR: unlike the hashicorp
# kubernetes provider's kubernetes_manifest resource, kubectl_manifest needs
# no cluster connection at plan time, so the certificate can be planned
# before cert-manager's CRDs exist (single-run infra charts, offline plan
# proofs). This module creates no other Kubernetes objects, so kubectl is
# its only provider.
#
# The provider is configured by the calling workspace/environment (the same
# kubeconfig environment contract).

terraform {
  required_version = ">= 1.0"

  required_providers {
    kubectl = {
      source  = "alekc/kubectl"
      version = ">= 2.0"
    }
  }
}

provider "kubectl" {
}
