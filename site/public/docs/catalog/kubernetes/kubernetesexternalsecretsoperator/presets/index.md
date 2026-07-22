---
title: "Presets"
description: "Ready-to-deploy configuration presets for KubernetesExternalSecretsOperator"
type: "preset-list"
componentSlug: "kubernetesexternalsecretsoperator"
componentTitle: "KubernetesExternalSecretsOperator"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-minimal"
    rank: "01"
    title: "Minimal Installation"
    excerpt: "This preset installs the External Secrets Operator with chart defaults: one controller replica, CRDs installed with the release and kept on uninstall, webhook and cert-controller enabled. No ambient..."
  - slug: "02-eks-ambient-identity"
    rank: "02"
    title: "EKS Ambient Identity (IRSA)"
    excerpt: "This preset installs the External Secrets Operator on an EKS cluster with the controller ServiceAccount bound to an IAM role via IRSA. Every store that leaves its auth block empty authenticates..."
  - slug: "03-tuned-multi-team"
    rank: "03"
    title: "Tuned Multi-Team Installation"
    excerpt: "This preset sizes the External Secrets Operator for clusters where many teams sync many secrets: reconcile concurrency raised to 5, explicit controller resources, two webhook replicas (every ESO..."
---

# KubernetesExternalSecretsOperator Presets

Ready-to-deploy configuration presets for KubernetesExternalSecretsOperator. Each preset is a complete manifest you can copy, customize, and deploy.
