# Bottlerocket Encrypted

This preset declares a security-hardened machine template: Bottlerocket
(the container-optimized, minimal-surface OS) with both of its volumes on
encrypted gp3 EBS under a customer-managed KMS key, and the IMDS security
defaults stated explicitly. Discovery-tag selection and the
Karpenter-managed node role match the AL2023 standard preset.

> **NOTE on `bottlerocket@latest`:** the alias drifts nodes whenever a
> new Bottlerocket AMI ships. Pin a release version in production and
> roll AMI updates deliberately.

## When to Use

- Clusters with compliance requirements for encryption-at-rest under a
  customer-managed key
- Teams preferring a minimal, container-only host OS over a
  general-purpose Linux node image

## Key Configuration Choices

- **Bottlerocket via alias** — an alias term must be the only AMI
  selector term; the AMI family is inferred, so no `amiFamily` is set
- **Two block-device mappings** — Bottlerocket separates a small OS
  volume (`/dev/xvda`, 4Gi) from the data volume (`/dev/xvdb`, 100Gi)
  holding container images and pod storage; size the data volume for
  your image and ephemeral-storage footprint
- **Encrypted gp3 + customer-managed KMS** on both volumes, with
  `deleteOnTermination: true` so volumes never outlive their node
- **IMDS defaults made explicit** (`metadataOptions`) — IMDSv2 required,
  hop limit 1 so pods cannot reach the node's instance credentials, and
  the IPv6 metadata endpoint disabled; these are the CRD defaults and
  the EKS security best practice, stated here so a security review reads
  them off the manifest
- **`role` (not `instance_profile`)** — Karpenter manages the instance
  profile

## Placeholders to Replace

| Placeholder                       | Description                                | Where to Find                                 |
| --------------------------------- | ------------------------------------------ | --------------------------------------------- |
| `<eks-cluster-name>`              | Cluster name in the discovery tag value    | EKS console; tags on your VPC subnets and SGs |
| `<karpenter-node-role-name>`      | IAM role name nodes assume                 | IAM console                                   |
| `<ebs-kms-key-id-alias-or-arn>`   | Customer-managed KMS key (id/alias/ARN)    | KMS console or `AwsKmsKey` outputs            |

## Related Presets

- **01-al2023-standard** — the general-purpose Amazon Linux 2023 baseline
- **03-custom-ami-pipeline** — AMIs resolved from an SSM parameter
