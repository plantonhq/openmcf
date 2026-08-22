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
  }

  required_version = ">= 1.0"
}

# Both providers are configured by the calling workspace/environment (the
# same kubeconfig environment contract every kubernetes-provider module
# rides). Keep these blocks empty -- do not wire cluster coordinates here.
provider "kubernetes" {
}

provider "helm" {
}
