# Provider requirements for the KubernetesSparkOperator module.
#
# helm installs the operator release; kubernetes creates the optional
# installation namespace. No kubectl provider: the chart's crds/-directory
# CRDs are Helm-installed (install-once, kept on uninstall) and the module
# owns no raw manifests.
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
