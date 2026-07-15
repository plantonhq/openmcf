---
title: "EBS CSI Driver with Pod Identity"
description: "This preset installs the EBS CSI storage driver with its own IAM identity wired through EKS Pod Identity -- the modern, no-OIDC-provider way to give an add-on AWS permissions."
type: "preset"
rank: "02"
presetSlug: "02-ebs-csi-pod-identity"
componentSlug: "eks-add-on"
componentTitle: "EKS Add-on"
provider: "aws"
icon: "package"
order: 2
---

# EBS CSI Driver with Pod Identity

This preset installs the EBS CSI storage driver with its own IAM
identity wired through EKS Pod Identity -- the modern, no-OIDC-provider
way to give an add-on AWS permissions.

## When to Use

- Any cluster whose workloads use EBS-backed PersistentVolumes (the
  default StorageClass path on EKS)
- As the identity-wiring template for other AWS-calling add-ons
  (EFS CSI, CloudWatch observability, ...)

## Key Configuration Choices

- **`podIdentityAssociations` over `serviceAccountRoleArn`** -- Pod
  Identity needs no per-cluster OIDC provider and the same role works
  across clusters; prefer it on new clusters. IRSA
  (`serviceAccountRoleArn`) remains for clusters already invested in
  OIDC-provider wiring.
- **`serviceAccount: ebs-csi-controller-sa`** -- the driver's
  documented controller service account; each add-on documents its own.
- **The role carries its own policies** -- attach
  `service-role/AmazonEBSCSIDriverPolicy` on the referenced
  `AwsIamRole` (`managedPolicyArns`) and trust
  `pods.eks.amazonaws.com`; this add-on never modifies the role.

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<addon-resource-name>` | Name for this add-on resource (e.g. `platform-ebs-csi`) | Your naming convention |
| `<aws-region>` | AWS region code (e.g., `us-west-2`) | Your deployment region |
| `<cluster-resource-name>` | Name of the AwsEksCluster resource | Your cluster manifest's `metadata.name` |
| `<ebs-csi-role-resource-name>` | Name of the AwsIamRole with the EBS CSI policy | Your role manifest's `metadata.name` |

## Common Additions

- An `eks-pod-identity-agent` AwsEksAddon on the same cluster (the
  agent that makes Pod Identity work)
- `configurationValues` to tune the driver (e.g. enable volume
  snapshots alongside the snapshot-controller add-on)

## Related Presets

- **01-core-networking** -- the baseline trio every cluster runs
- **03-pinned-version** -- version-pinned add-on for controlled
  upgrades
