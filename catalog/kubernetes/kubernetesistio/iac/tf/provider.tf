# Provider requirements for the KubernetesIstio module.
#
# helm installs the four control-plane releases (base, istiod, and in ambient
# mode cni + ztunnel); kubernetes creates the optional anchor namespace;
# kubectl (alekc/kubectl) applies the module-owned CRD bundle (server-side
# apply — co-ownable, unlike Helm-owned CRDs); http fetches that bundle at
# the pinned version.
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
    http = {
      source  = "hashicorp/http"
      version = ">= 3.0"
    }
  }
}

provider "kubernetes" {
}

provider "helm" {
}

provider "kubectl" {
}
