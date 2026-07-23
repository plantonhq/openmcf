---
title: "AL2023 Standard"
description: "This preset declares the standard EKS machine template: EKS-optimized Amazon Linux 2023 AMIs, subnets and security groups selected by the `karpenter.sh/discovery` tag convention, and a node IAM role..."
type: "preset"
rank: "01"
presetSlug: "01-al2023-standard"
componentSlug: "karpenter-ec2-node-class"
componentTitle: "Karpenter EC2 Node Class"
provider: "kubernetes"
icon: "package"
order: 1
---

# AL2023 Standard

This preset declares the standard EKS machine template: EKS-optimized
Amazon Linux 2023 AMIs, subnets and security groups selected by the
`karpenter.sh/discovery` tag convention, and a node IAM role whose
instance profile Karpenter manages. One NodeClass is typically shared by
several NodePools — the pools differ in constraints and taints; this
class is the common "how a node is built".

> **NOTE on `al2023@latest`:** the alias resolves to the newest
> EKS-optimized AL2023 release, which DRIFTS nodes whenever a new AMI
> ships. Pin a release version (e.g. `al2023@v20240807`) in production
> and roll AMI updates deliberately.

## When to Use

- The first EC2NodeClass of a new Karpenter installation on EKS
- Clusters whose subnets and security groups already carry the
  `karpenter.sh/discovery: <cluster>` tag (the standard EKS setup)

## Key Configuration Choices

- **AMI by alias** (`al2023@latest`) — the simplest selector arm; an
  alias term must be the only AMI selector term, and the AMI family is
  inferred from it (no `amiFamily` needed)
- **Discovery-tag selection** for subnets and security groups — new
  subnets/SGs enroll by tagging alone; terms are ORed, and Karpenter
  spreads nodes across the selected subnets' zones
- **`role` (not `instance_profile`)** — Karpenter creates and manages
  the instance profile; the controller role needs `iam:PassRole`
- **Everything else at CRD defaults** — including the IMDS security
  defaults (IMDSv2 required, hop limit 1) and the AMI family's default
  disk layout; see **02-bottlerocket-encrypted** for explicit
  block-device and IMDS blocks

## Placeholders to Replace

| Placeholder                 | Description                              | Where to Find                                 |
| --------------------------- | ---------------------------------------- | --------------------------------------------- |
| `<eks-cluster-name>`        | Cluster name in the discovery tag value  | EKS console; tags on your VPC subnets and SGs |
| `<karpenter-node-role-name>`| IAM role name nodes assume (name, not ARN) | IAM console (node role per upstream guidance) |

## Related Presets

- **02-bottlerocket-encrypted** — container-optimized OS with encrypted
  gp3 volumes and explicit IMDS posture
- **03-custom-ami-pipeline** — AMIs resolved from an SSM parameter
  maintained by your image pipeline
