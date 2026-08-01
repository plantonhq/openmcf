# Provider requirements for the KubernetesOtelOperator module.
#
# helm installs the operator release; kubernetes creates the optional
# installation namespace; kubectl (alekc/kubectl) applies the module-owned
# CRDs (server-side apply + apply_only — the keep-on-uninstall mechanism,
# see main.tf).
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
    kubectl = {
      source  = "alekc/kubectl"
      version = ">= 2.0"
    }
  }
}

provider "kubernetes" {
}

provider "helm" {
}

provider "kubectl" {
}
