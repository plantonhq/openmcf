# Provider requirements for the KubernetesOtelOperator module.
#
# helm installs the operator release and renders the pinned chart to
# derive its CRDs; kubernetes creates the optional installation namespace
# and reads the CRDs already on the cluster for the never-downgrade check;
# kubectl (alekc/kubectl) applies the module-owned CRDs (server-side
# apply + apply_only — the keep-on-uninstall mechanism); http is the
# bundle-fetch branch of the shared helm_crds.tf, declared because init
# resolves every data source in that file even when this chart never
# uses it.
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
