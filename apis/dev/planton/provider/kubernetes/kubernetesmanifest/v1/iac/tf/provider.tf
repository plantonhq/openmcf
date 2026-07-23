# Provider requirements for the KubernetesManifest module.
#
# kubectl (alekc/kubectl) applies the raw documents: unlike the hashicorp
# kubernetes provider's kubernetes_manifest resource, kubectl_manifest needs
# no cluster connection at plan time and handles CRDs and their custom
# resources in the same apply. The hashicorp kubernetes provider is used only
# for the optional anchor-namespace creation.
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
