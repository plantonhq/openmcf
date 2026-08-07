# Provider requirements for the KubernetesTekton module.
#
# kubectl (alekc/kubectl) renders the TektonConfig custom resource:
# unlike the hashicorp kubernetes provider's kubernetes_manifest
# resource, kubectl_manifest needs no cluster connection at plan time —
# the CRD (installed by the KubernetesTektonOperator prerequisite) may
# not exist yet when a composed infra chart plans this resource.
#
# Providers are configured by the calling workspace/environment (the same
# kubeconfig environment contract).

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
