# Provider requirements for the KubernetesHelmRelease module.
#
# helm >= 3.1: floor set by take_ownership (added in provider 3.1.0); helm
# also renders the pinned chart to derive the CRDs in its crds/ directory.
# kubernetes creates the optional namespace and reads the CRDs already on
# the cluster for the never-downgrade check; kubectl (alekc/kubectl) applies
# the module-owned CRDs (server-side apply + apply_only — the
# keep-on-uninstall mechanism); http is the index read and the bundle-fetch
# branch of the shared helm_crds.tf, declared because init resolves every
# data source in that file.
#
# Providers are configured by the calling workspace/environment (the same
# kubeconfig environment contract).

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
    kubectl = {
      source  = "alekc/kubectl"
      version = ">= 2.0"
    }
    http = {
      source  = "hashicorp/http"
      version = "~> 3.4"
    }
  }
}

provider "kubernetes" {
}

provider "helm" {
}

provider "kubectl" {
}
