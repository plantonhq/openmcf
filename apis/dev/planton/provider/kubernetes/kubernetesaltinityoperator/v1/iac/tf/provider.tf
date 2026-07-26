# Provider requirements for the KubernetesAltinityOperator module.
#
# helm installs the operator release; kubernetes creates the optional
# installation namespace. No kubectl provider: the CRDs are chart-owned
# (crds/ directory + the crdHook upgrade job — see main.tf), so the
# module applies no manifests of its own.
#
# Providers are configured by the calling workspace/environment (the same
# kubeconfig environment contract).

terraform {
  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.35"
    }
    helm = {
      source  = "hashicorp/helm"
      version = "~> 3.0"
    }
  }
}

provider "kubernetes" {
}

provider "helm" {
}
