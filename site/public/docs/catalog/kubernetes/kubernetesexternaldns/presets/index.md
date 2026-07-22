---
title: "Presets"
description: "Ready-to-deploy configuration presets for KubernetesExternalDns"
type: "preset-list"
componentSlug: "kubernetesexternaldns"
componentTitle: "KubernetesExternalDns"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-aws-route53-eks-keyless"
    rank: "01"
    title: "AWS Route 53 on EKS (Keyless via IRSA)"
    excerpt: "This preset installs ExternalDNS on an EKS cluster publishing to a Route 53 hosted zone, authenticating keylessly through IRSA — no static AWS keys anywhere. It scopes the instance to one zone,..."
  - slug: "02-google-cloud-dns-gke"
    rank: "02"
    title: "Google Cloud DNS on GKE (Keyless via Workload Identity)"
    excerpt: "This preset installs ExternalDNS on a GKE cluster publishing to a Cloud DNS zone, authenticating keylessly through GKE Workload Identity — no service-account key anywhere. It scopes the instance to..."
  - slug: "03-azure-dns-aks"
    rank: "03"
    title: "Azure DNS on AKS (Keyless via Workload Identity)"
    excerpt: "This preset installs ExternalDNS on an AKS cluster publishing to an Azure DNS zone, authenticating keylessly through Azure AD Workload Identity — no service-principal secret anywhere. The module..."
  - slug: "04-cloudflare-any-cluster"
    rank: "04"
    title: "Cloudflare DNS from Any Cluster"
    excerpt: "This preset installs ExternalDNS publishing to Cloudflare — the canonical cross-cloud arm. Cloudflare has no workload-identity federation with Kubernetes clusters, so authentication is always an API..."
  - slug: "05-webhook-provider"
    rank: "05"
    title: "Webhook Provider (Out-of-Tree DNS)"
    excerpt: "This preset installs ExternalDNS with the webhook provider — upstream's extension architecture for every DNS provider that is not in-tree (Akamai, OVH, Scaleway, RFC2136, Hetzner, and many more). The..."
---

# KubernetesExternalDns Presets

Ready-to-deploy configuration presets for KubernetesExternalDns. Each preset is a complete manifest you can copy, customize, and deploy.
