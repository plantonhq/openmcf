# Provider requirements for the KubernetesSolrOperator module.
#
# kubectl (alekc/kubectl) applies the module-owned CRDs: it supports
# server-side apply (the SolrCloud CRD's schema exceeds the client-side
# last-applied-configuration annotation size limit) and apply_only (the
# keep-on-uninstall mechanism — Delete is a no-op in the provider
# source). The hashicorp kubernetes provider materializes the optional
# namespace; helm installs the operator release.
#
# All three providers are configured by the calling workspace/environment
# (the same kubeconfig environment contract).

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
