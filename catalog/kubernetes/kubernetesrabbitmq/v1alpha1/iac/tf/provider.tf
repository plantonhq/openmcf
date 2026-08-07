# Provider requirements for the KubernetesRabbitMq module.
#
# kubectl (alekc/kubectl) applies the RabbitmqCluster custom resource:
# unlike the hashicorp kubernetes provider's kubernetes_manifest resource,
# kubectl_manifest needs no cluster connection at plan time, so the
# cluster can be planned before the RabbitMQ Cluster Operator's CRD exists
# (single-run infra charts, offline plan proofs). The hashicorp kubernetes
# provider materializes the optional namespace.
#
# Both providers are configured by the calling workspace/environment (the
# same kubeconfig environment contract).

terraform {
  required_version = ">= 1.0"

  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = ">= 2.20"
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
