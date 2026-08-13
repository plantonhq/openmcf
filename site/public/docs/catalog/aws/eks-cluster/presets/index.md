---
title: "Presets"
description: "Ready-to-deploy configuration presets for EKS Cluster"
type: "preset-list"
componentSlug: "eks-cluster"
componentTitle: "EKS Cluster"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-standard"
    rank: "01"
    title: "Standard EKS Cluster"
    excerpt: "This preset creates a two-AZ EKS control plane with the modern access-entries authentication model, audit-grade control-plane logging, and standard upgrade support. The API endpoint stays publicly..."
  - slug: "02-private-endpoint"
    rank: "02"
    title: "Private-Endpoint EKS Cluster"
    excerpt: "This preset creates a hardened, fully private EKS control plane: the API server is reachable only from inside the VPC, Kubernetes secrets are envelope-encrypted with a customer-managed KMS key, and..."
  - slug: "03-auto-mode"
    rank: "03"
    title: "EKS Auto Mode Cluster"
    excerpt: "This preset creates a hands-off Kubernetes platform: EKS Auto Mode provisions and scales EC2 capacity, provisions EBS volumes, and manages load balancers for the cluster's workloads. There are no..."
  - slug: "04-hybrid-nodes"
    rank: "04"
    title: "EKS Hybrid Nodes Cluster"
    excerpt: "This preset creates a control plane that on-premises or edge machines join as workers over your VPN or Direct Connect: one Kubernetes API, cloud and on-prem capacity underneath it. AWS bills hybrid..."
---

# EKS Cluster Presets

Ready-to-deploy configuration presets for EKS Cluster. Each preset is a complete manifest you can copy, customize, and deploy.
