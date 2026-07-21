---
title: "Presets"
description: "Ready-to-deploy configuration presets for ServiceAccount"
type: "preset-list"
componentSlug: "serviceaccount"
componentTitle: "ServiceAccount"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-basic"
    rank: "01"
    title: "Basic ServiceAccount"
    excerpt: "This preset creates a plain ServiceAccount — a dedicated in-cluster identity for a workload, with no cloud federation and no pull secrets. Point workloads at it with `spec.serviceAccountName` and..."
  - slug: "02-workload-identity-gke"
    rank: "02"
    title: "GKE Workload Identity"
    excerpt: "This preset creates a ServiceAccount bound to a GCP service account via GKE Workload Identity. Pods running as this identity call Google Cloud APIs keylessly — the cluster's OIDC issuer vouches for..."
  - slug: "03-workload-identity-eks-irsa"
    rank: "03"
    title: "EKS IRSA (IAM Roles for Service Accounts)"
    excerpt: "This preset creates a ServiceAccount bound to an AWS IAM role via IRSA. Pods running as this identity call AWS APIs keylessly — the AWS SDK inside the pod exchanges the projected ServiceAccount token..."
  - slug: "04-image-pull-secrets"
    rank: "04"
    title: "Image Pull Secrets with Automount Hardening"
    excerpt: "This preset creates a ServiceAccount that carries private-registry pull credentials and disables the automatic API token mount. Every pod running as this identity pulls images with the attached..."
---

# ServiceAccount Presets

Ready-to-deploy configuration presets for ServiceAccount. Each preset is a complete manifest you can copy, customize, and deploy.
