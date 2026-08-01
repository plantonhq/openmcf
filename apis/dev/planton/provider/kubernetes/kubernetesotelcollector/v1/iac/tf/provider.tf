# Provider requirements for the KubernetesOtelCollector module.
#
# kubectl (alekc/kubectl) applies the OpenTelemetryCollector CR (no
# cluster connection needed at plan time — offline plan proofs and
# single-run infra charts both depend on that); kubernetes creates the
# optional deployment namespace.
#
# Providers are configured by the calling workspace/environment (the same
# kubeconfig environment contract).

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
  }
}

provider "kubernetes" {
}

provider "kubectl" {
}
