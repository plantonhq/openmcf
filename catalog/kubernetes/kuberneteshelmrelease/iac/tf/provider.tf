# Provider requirements for the KubernetesHelmRelease module.
#
# helm >= 3.1: floor set by take_ownership (added in provider 3.1.0).
# Both providers are configured by the calling workspace/environment (the
# same kubeconfig environment contract).

terraform {
  required_version = ">= 1.0"

  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = ">= 2.20"
    }
    helm = {
      source  = "hashicorp/helm"
      version = ">= 3.1"
    }
  }
}

provider "kubernetes" {
}

provider "helm" {
}
