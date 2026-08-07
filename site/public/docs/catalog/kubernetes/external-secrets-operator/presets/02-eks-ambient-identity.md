---
title: "EKS Ambient Identity (IRSA)"
description: "This preset installs the External Secrets Operator on an EKS cluster with the controller ServiceAccount bound to an IAM role via IRSA. Every store that leaves its auth block empty authenticates..."
type: "preset"
rank: "02"
presetSlug: "02-eks-ambient-identity"
componentSlug: "external-secrets-operator"
componentTitle: "External Secrets Operator"
provider: "kubernetes"
icon: "package"
order: 2
---

# EKS Ambient Identity (IRSA)

This preset installs the External Secrets Operator on an EKS cluster with
the controller ServiceAccount bound to an IAM role via IRSA. Every store
that leaves its auth block empty authenticates through this ambient
identity — the simplest posture when one cloud identity may read every
secret the cluster syncs. No static AWS keys anywhere.

## When to Use

- EKS clusters syncing from AWS Secrets Manager or SSM Parameter Store
- Single-team clusters (or platform-owned secrets) where one IAM role
  reading everything is acceptable
- When you want stores to stay auth-free and inherit identity from the
  operator

For multi-team clusters where different stores must NOT share read access,
skip `workloadIdentity` here and give each store its own identity in its
auth block instead (see preset 01).

## Key Configuration Choices

- **Ambient identity** (`workloadIdentity.eks.roleArn`) — rendered as the
  `eks.amazonaws.com/role-arn` annotation on the controller ServiceAccount
  (`external-secrets` in the installation namespace); stores without their
  own auth fall back to it
- **Everything else at chart defaults** — CRDs installed and kept,
  single-replica components

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<external-secrets-irsa-role-arn>` | IAM role with read access to the synced secrets, trusting the cluster OIDC provider for `system:serviceaccount:external-secrets:external-secrets` | IAM console or `AwsIamRole` outputs |

## Related Presets

- **01-minimal** — no ambient identity; per-store auth only
- **03-tuned-multi-team** — sizing and concurrency for large clusters
