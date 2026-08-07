---
title: "Presets"
description: "Ready-to-deploy configuration presets for EKS Addon"
type: "preset-list"
componentSlug: "eks-addon"
componentTitle: "EKS Addon"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-core-networking"
    rank: "01"
    title: "Core Networking Add-on"
    excerpt: "This preset adopts one of the cluster's core networking add-ons (vpc-cni, coredns, kube-proxy) into AWS-managed lifecycle -- upgrades and configuration flow through the EKS control plane instead of..."
  - slug: "02-ebs-csi-pod-identity"
    rank: "02"
    title: "EBS CSI Driver with Pod Identity"
    excerpt: "This preset installs the EBS CSI storage driver with its own IAM identity wired through EKS Pod Identity -- the modern, no-OIDC-provider way to give an add-on AWS permissions."
  - slug: "03-pinned-version"
    rank: "03"
    title: "Pinned Version with Configuration"
    excerpt: "This preset pins an add-on to an exact version and carries custom configuration -- the controlled-upgrade pattern for fleets that must run byte-identical clusters."
---

# EKS Addon Presets

Ready-to-deploy configuration presets for EKS Addon. Each preset is a complete manifest you can copy, customize, and deploy.
