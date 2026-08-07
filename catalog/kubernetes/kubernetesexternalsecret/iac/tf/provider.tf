# Provider requirements for the KubernetesExternalSecret module.
#
# kubectl (alekc/kubectl) applies the ExternalSecret CR: unlike the
# hashicorp kubernetes provider's kubernetes_manifest resource,
# kubectl_manifest needs no cluster connection at plan time, so the secret
# can be planned before the External Secrets Operator's CRDs exist
# (single-run infra charts, offline plan proofs). The hashicorp kubernetes
# provider is carried for the family's shared provider contract (the store
# twins materialize credential Secrets with it; this module creates none).
#
# Both providers are configured by the calling workspace/environment (the
# same kubeconfig environment contract).

terraform {
  required_version = ">= 1.0"

  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = ">= 2.20"
    }
    kubectl = {
      source  = "alekc/kubectl"
      version = ">= 2.0"
    }
  }
}

provider "kubernetes" {
}

provider "kubectl" {
}
