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
    # alekc/kubectl for the module-owned CRD: kubectl_manifest needs no
    # cluster connection at plan time (the CRD may not exist when a
    # composed infra chart plans this resource) and its server-side apply
    # natively ADOPTS a CRD retained by a previous install's destroy.
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
