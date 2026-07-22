---
title: "Presets"
description: "Ready-to-deploy configuration presets for KubernetesClusterSecretStore"
type: "preset-list"
componentSlug: "kubernetesclustersecretstore"
componentTitle: "KubernetesClusterSecretStore"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-aws-secrets-manager-irsa"
    rank: "01"
    title: "AWS Secrets Manager with IRSA (Keyless)"
    excerpt: "This preset creates a cluster-wide store connected to AWS Secrets Manager, authenticating through a referenced ServiceAccount whose IRSA (IAM Roles for Service Accounts) binding authorizes the reads...."
  - slug: "02-gcp-secret-manager-workload-identity"
    rank: "02"
    title: "GCP Secret Manager with Workload Identity (Keyless)"
    excerpt: "This preset creates a cluster-wide store connected to GCP Secret Manager, authenticating through a referenced ServiceAccount whose GKE Workload Identity binding authorizes the reads. No..."
  - slug: "03-azure-key-vault-workload-identity"
    rank: "03"
    title: "Azure Key Vault with Workload Identity (Keyless)"
    excerpt: "This preset creates a cluster-wide store connected to Azure Key Vault, authenticating through a referenced ServiceAccount whose AKS Workload Identity federation authorizes the reads. No client secret..."
  - slug: "04-vault-kubernetes-auth"
    rank: "04"
    title: "Vault KV with Kubernetes Auth (Keyless)"
    excerpt: "This preset creates a cluster-wide store connected to a HashiCorp Vault KV v2 engine, authenticating through Vault's Kubernetes auth method: the referenced ServiceAccount's token is exchanged for a..."
---

# KubernetesClusterSecretStore Presets

Ready-to-deploy configuration presets for KubernetesClusterSecretStore. Each preset is a complete manifest you can copy, customize, and deploy.
