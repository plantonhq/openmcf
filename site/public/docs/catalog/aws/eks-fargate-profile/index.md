---
title: "EKS Fargate Profile"
description: "EKS Fargate Profile deployment documentation"
icon: "package"
order: 100
componentName: "awseksfargateprofile"
---

# AWS EKS Fargate Profile

Declares which Kubernetes pods of an `AwsEksCluster` run on AWS
Fargate -- serverless, per-pod compute with no EC2 nodes to size,
patch, or scale -- selected by namespace (with wildcards) and optional
pod labels.

## What Gets Created

When you deploy an AwsEksFargateProfile resource, Planton provisions:

- **Fargate profile** — an `aws_eks_fargate_profile` /
  `eks.FargateProfile` named from `metadata.name`, attached to the
  referenced cluster, running matched pods as the referenced pod
  execution role inside the referenced private subnets

The pod execution role is never modified: attach
`AmazonEKSFargatePodExecutionRolePolicy` on the referenced `AwsIamRole`
itself, with a trust policy for `eks-fargate-pods.amazonaws.com`.

## Prerequisites

- **AWS credentials** configured via the Planton provider config (keyless SSO/OIDC).
- **An EKS cluster** (`AwsEksCluster`) for the profile to attach to.
- **A pod execution role** (`AwsIamRole`) trusting `eks-fargate-pods.amazonaws.com` with the Fargate pod execution policy.
- **Private subnets** (`AwsSubnet`) — AWS rejects subnets with an internet-gateway route; use a NAT gateway for egress.

## Quick Start

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsEksFargateProfile
metadata:
  name: serverless
spec:
  region: us-west-2
  clusterName:
    valueFrom:
      kind: AwsEksCluster
      name: platform
      fieldPath: status.outputs.name
  podExecutionRoleArn:
    valueFrom:
      kind: AwsIamRole
      name: fargate-pod-execution
      fieldPath: status.outputs.role_arn
  subnetIds:
    - valueFrom:
        kind: AwsSubnet
        name: private-a
        fieldPath: status.outputs.subnet_id
    - valueFrom:
        kind: AwsSubnet
        name: private-b
        fieldPath: status.outputs.subnet_id
  selectors:
    - namespace: serverless
```

```shell
planton apply -f fargate-profile.yaml
```

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
| --- | --- | --- | --- |
| `region` | `string` | AWS region; must match the cluster's. | Required; non-empty |
| `clusterName` | `string \| valueFrom` | The cluster the profile attaches to. Defaults to referencing an `AwsEksCluster` `name` output. | Required |
| `podExecutionRoleArn` | `string \| valueFrom` | The role Fargate runs matched pods as. Defaults to referencing an `AwsIamRole` `role_arn` output. | Required |
| `subnetIds` | `string[] \| valueFrom` | Private subnets pods launch into. | Required; ≥1 entry |
| `selectors` | `object[]` | Which pods run on Fargate: `namespace` (wildcards `*`/`?` allowed) + optional `labels` (AND semantics, ≤5 pairs). A pod matches when ANY selector matches. | Required; 1–5 entries |

Every field is create-time immutable in AWS — changing anything
replaces the profile.

## Examples

### Namespace-per-team wildcard

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsEksFargateProfile
metadata:
  name: team-namespaces
spec:
  region: us-west-2
  clusterName:
    valueFrom: { kind: AwsEksCluster, name: platform, fieldPath: status.outputs.name }
  podExecutionRoleArn:
    valueFrom: { kind: AwsIamRole, name: fargate-pod-execution, fieldPath: status.outputs.role_arn }
  subnetIds:
    - valueFrom: { kind: AwsSubnet, name: private-a, fieldPath: status.outputs.subnet_id }
    - valueFrom: { kind: AwsSubnet, name: private-b, fieldPath: status.outputs.subnet_id }
  selectors:
    - namespace: team-*
```

### Label-scoped batch workloads

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsEksFargateProfile
metadata:
  name: batch-fargate
spec:
  region: us-west-2
  clusterName:
    valueFrom: { kind: AwsEksCluster, name: platform, fieldPath: status.outputs.name }
  podExecutionRoleArn:
    valueFrom: { kind: AwsIamRole, name: fargate-pod-execution, fieldPath: status.outputs.role_arn }
  subnetIds:
    - valueFrom: { kind: AwsSubnet, name: private-a, fieldPath: status.outputs.subnet_id }
    - valueFrom: { kind: AwsSubnet, name: private-b, fieldPath: status.outputs.subnet_id }
  selectors:
    - namespace: batch
      labels:
        compute: fargate
```

## Stack Outputs

| Output | Description |
| --- | --- |
| `fargate_profile_arn` | The profile's ARN |
| `fargate_profile_name` | The profile's name |
| `status` | The profile's state after provisioning — `ACTIVE` on success |

## Related Components

- [AwsEksCluster](/docs/catalog/aws/eks-cluster) — the control plane the profile attaches to
- [AwsIamRole](/docs/catalog/aws/iam-role) — the pod execution role (carries its own policy)
- [AwsSubnet](/docs/catalog/aws/subnet) — the private subnets pods launch into
- [AwsNatGateway](/docs/catalog/aws/nat-gateway) — outbound internet for Fargate pods
- [AwsEksNodeGroup](/docs/catalog/aws/eks-node-group) — EC2-backed compute for everything else
