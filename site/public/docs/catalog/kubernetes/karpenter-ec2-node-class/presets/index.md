---
title: "Presets"
description: "Ready-to-deploy configuration presets for Karpenter EC2 Node Class"
type: "preset-list"
componentSlug: "karpenter-ec2-node-class"
componentTitle: "Karpenter EC2 Node Class"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-al2023-standard"
    rank: "01"
    title: "AL2023 Standard"
    excerpt: "This preset declares the standard EKS machine template: EKS-optimized Amazon Linux 2023 AMIs, subnets and security groups selected by the `karpenter.sh/discovery` tag convention, and a node IAM role..."
  - slug: "02-bottlerocket-encrypted"
    rank: "02"
    title: "Bottlerocket Encrypted"
    excerpt: "This preset declares a security-hardened machine template: Bottlerocket (the container-optimized, minimal-surface OS) with both of its volumes on encrypted gp3 EBS under a customer-managed KMS key,..."
  - slug: "03-custom-ami-pipeline"
    rank: "03"
    title: "Custom AMI Pipeline"
    excerpt: "This preset declares a machine template for organizations that build their own node AMIs: the image is resolved from an SSM parameter your pipeline maintains, so AMI rollout is controlled by the..."
---

# Karpenter EC2 Node Class Presets

Ready-to-deploy configuration presets for Karpenter EC2 Node Class. Each preset is a complete manifest you can copy, customize, and deploy.
