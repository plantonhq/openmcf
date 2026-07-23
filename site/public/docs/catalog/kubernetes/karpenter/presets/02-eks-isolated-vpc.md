---
title: "EKS Isolated VPC"
description: "This preset installs Karpenter on an EKS cluster running in an isolated VPC — one without internet-reachable AWS endpoints beyond the provisioned VPC endpoints — with VPC CNI custom networking. It is..."
type: "preset"
rank: "02"
presetSlug: "02-eks-isolated-vpc"
componentSlug: "karpenter"
componentTitle: "Karpenter"
provider: "kubernetes"
icon: "package"
order: 2
---

# EKS Isolated VPC

This preset installs Karpenter on an EKS cluster running in an isolated
VPC — one without internet-reachable AWS endpoints beyond the provisioned
VPC endpoints — with VPC CNI custom networking. It is the EKS-standard
posture plus the two flags those environments require: `isolatedVpc` and
`reservedEnis`.

## When to Use

- EKS clusters in private/air-gapped VPCs where AWS services are only
  reachable through VPC endpoints
- Clusters using VPC CNI custom networking, where pod ENIs live in
  separate subnets and one ENI per node belongs to the CNI

## Key Configuration Choices

- **`aws.isolatedVpc: true`** — Karpenter avoids AWS services without a
  VPC endpoint; note the pricing API has none, so price-aware
  consolidation falls back to static pricing data
- **`aws.reservedEnis: 1`** — one ENI per node is reserved outside
  Karpenter's max-pods and kube-reserved math for the CNI (chart default
  is 0; set it to match your custom-networking layout)
- **`cluster.eksControlPlane: true`** — endpoint/CA discovery via
  DescribeCluster still works through the EKS VPC endpoint
- **IRSA + interruption queue** — same keyless credential and
  drain-ahead-of-interruption posture as **01-eks-standard**
- **CRDs installed and kept on uninstall** (spec defaults, made explicit)

## Placeholders to Replace

| Placeholder                                           | Description                                 | Where to Find                             |
| ----------------------------------------------------- | ------------------------------------------- | ----------------------------------------- |
| `<eks-cluster-name>`                                  | EKS cluster name                            | EKS console or `AwsEksCluster` outputs    |
| `arn:aws:iam::123456789012:role/karpenter-controller` | IRSA role ARN — replace account id and name | IAM console (role per upstream IAM guide) |
| `karpenter-interruptions`                             | SQS interruption queue name                 | SQS console                               |

## Related Presets

- **01-eks-standard** — the same posture without the isolated-VPC and
  custom-networking flags
- **03-ha-tuned** — HA sizing, batching tuning, and telemetry
