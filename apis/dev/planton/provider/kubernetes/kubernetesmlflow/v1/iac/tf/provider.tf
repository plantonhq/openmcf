# The hashicorp kubernetes provider renders the core objects (Deployment,
# Service, Secrets, PVCs, CronJob). kubectl (alekc/kubectl) applies the
# optional ServiceMonitor CR: unlike the hashicorp provider's
# kubernetes_manifest resource, kubectl_manifest needs no cluster
# connection at plan time, so the monitor can be planned before the
# Prometheus operator's CRDs exist (single-run infra charts, offline plan
# proofs).
terraform {
  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.35"
    }
    kubectl = {
      source  = "alekc/kubectl"
      version = ">= 2.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.6"
    }
  }
}

provider "kubernetes" {
}

provider "kubectl" {
}
