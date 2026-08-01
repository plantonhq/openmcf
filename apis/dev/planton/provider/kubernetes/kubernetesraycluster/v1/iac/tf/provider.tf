# Provider requirements for the KubernetesRayCluster module.
#
# kubectl (alekc/kubectl) applies the RayCluster custom resource: unlike
# the hashicorp kubernetes provider's kubernetes_manifest resource,
# kubectl_manifest needs no cluster connection at plan time — the CRD
# (installed by the KubernetesKubeRayOperator prerequisite) may not
# exist yet when a composed infra chart plans this resource. The
# hashicorp kubernetes provider materializes the optional namespace.
#
# No helm provider: this module installs no chart — it renders one
# custom resource the already-installed operator reconciles.
#
# Both providers are configured by the calling workspace/environment
# (the same kubeconfig environment contract).

terraform {
  required_version = ">= 1.0"

  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.35"
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
