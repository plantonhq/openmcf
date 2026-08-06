# Provider requirements for the KubernetesKafkaTopic module.
#
# kubectl (alekc/kubectl) applies the KafkaTopic custom resource: unlike
# the hashicorp provider's kubernetes_manifest resource it needs no
# cluster connection at plan time, so the topic can be planned before the
# Strimzi operator's CRDs exist (single-run infra charts, offline plan
# proofs). This module creates no core resources — the namespace belongs
# to the KubernetesKafka resource's lifecycle.

terraform {
  required_version = ">= 1.0"

  required_providers {
    kubectl = {
      source  = "alekc/kubectl"
      version = ">= 2.0"
    }
  }
}

provider "kubectl" {
}
