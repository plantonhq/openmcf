---
title: "Presets"
description: "Ready-to-deploy configuration presets for Karpenter"
type: "preset-list"
componentSlug: "karpenter"
componentTitle: "Karpenter"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-eks-standard"
    rank: "01"
    title: "EKS Standard"
    excerpt: "This preset installs the Karpenter controller into `kube-system` on an EKS cluster with the two integrations every production installation carries: an IRSA role for keyless AWS API access and an SQS..."
  - slug: "02-eks-isolated-vpc"
    rank: "02"
    title: "EKS Isolated VPC"
    excerpt: "This preset installs Karpenter on an EKS cluster running in an isolated VPC — one without internet-reachable AWS endpoints beyond the provisioned VPC endpoints — with VPC CNI custom networking. It is..."
  - slug: "03-ha-tuned"
    rank: "03"
    title: "HA Tuned"
    excerpt: "This preset hardens the EKS-standard installation for clusters where Karpenter is load-bearing: explicit two-replica sizing with the documented large-cluster resource starting point, batching tuned..."
---

# Karpenter Presets

Ready-to-deploy configuration presets for Karpenter. Each preset is a complete manifest you can copy, customize, and deploy.
