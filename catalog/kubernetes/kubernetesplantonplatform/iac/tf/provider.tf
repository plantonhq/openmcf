terraform {
  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.35"
    }
    # alekc/kubectl for the PlantonPlatform CR: kubectl_manifest needs no
    # cluster connection at plan time — the CRD installed by the
    # prerequisite operator may not exist yet when a composed infra chart
    # plans this resource.
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
