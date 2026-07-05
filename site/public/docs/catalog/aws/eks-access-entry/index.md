---
title: "EKS Access Entry"
description: "EKS Access Entry deployment documentation"
icon: "package"
order: 100
componentName: "awseksaccessentry"
---

# AWS EKS Access Entry

Grants one IAM principal access to an `AwsEksCluster`'s Kubernetes API
through EKS access entries -- the modern replacement for the aws-auth
ConfigMap -- with authorization from AWS-managed access policies
(cluster- or namespace-scoped) and/or your own RBAC group mappings.

## What Gets Created

When you deploy an AwsEksAccessEntry resource, Planton provisions:

- **Access entry** — an `aws_eks_access_entry` / `eks.AccessEntry`
  keyed on the referenced cluster and principal, with your group
  mappings and username
- **Policy associations** — one `aws_eks_access_policy_association` /
  `eks.AccessPolicyAssociation` per `policyAssociations` entry, keyed
  by policy name so each one adds, re-scopes, or removes independently

## Prerequisites

- **AWS credentials** configured via the Planton provider config (keyless SSO/OIDC).
- **An EKS cluster** (`AwsEksCluster`) with API authentication enabled (`accessConfig.authenticationMode: API` or `API_AND_CONFIG_MAP`).
- **An IAM principal** (`AwsIamRole`, or a literal user/role ARN) to grant access to.

## Quick Start

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsEksAccessEntry
metadata:
  name: platform-viewers
spec:
  region: us-west-2
  clusterName:
    valueFrom:
      kind: AwsEksCluster
      name: platform
      fieldPath: status.outputs.name
  principalArn:
    valueFrom:
      kind: AwsIamRole
      name: team-viewer
      fieldPath: status.outputs.role_arn
  policyAssociations:
    - policyArn: arn:aws:eks::aws:cluster-access-policy/AmazonEKSViewPolicy
      accessScope:
        type: cluster
```

```shell
planton apply -f access-entry.yaml
```

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
| --- | --- | --- | --- |
| `region` | `string` | AWS region; must match the cluster's. | Required; non-empty |
| `clusterName` | `string \| valueFrom` | The cluster access is granted on. Defaults to referencing an `AwsEksCluster` `name` output. | Required |
| `principalArn` | `string \| valueFrom` | The IAM principal (role or user). Defaults to referencing an `AwsIamRole` `role_arn` output. One entry per principal per cluster. Create-only. | Required |

### Optional Fields

| Field | Type | Default | Description |
| --- | --- | --- | --- |
| `type` | `string` | `STANDARD` | `STANDARD` for humans/workloads; `EC2`, `EC2_LINUX`, `EC2_WINDOWS`, `FARGATE_LINUX`, `HYBRID_LINUX` for self-managed/hybrid node registration (forbid the fields below). Create-only. |
| `kubernetesGroups` | `string[]` | `[]` | Group names your own RBAC bindings reference; `system:` prefix rejected. Updates in place. |
| `userName` | `string` | AWS default | The Kubernetes username in audit logs; empty preserves AWS's session-templated default for roles. Updates in place. |
| `policyAssociations` | `object[]` | `[]` | AWS-managed access policies: `policyArn` + `accessScope` (`type: cluster` or `type: namespace` + `namespaces`). Each association updates independently. |

## Examples

### Namespace-scoped team admin

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsEksAccessEntry
metadata:
  name: team-a-admin
spec:
  region: us-west-2
  clusterName:
    valueFrom: { kind: AwsEksCluster, name: platform, fieldPath: status.outputs.name }
  principalArn:
    valueFrom: { kind: AwsIamRole, name: team-a, fieldPath: status.outputs.role_arn }
  policyAssociations:
    - policyArn: arn:aws:eks::aws:cluster-access-policy/AmazonEKSAdminPolicy
      accessScope:
        type: namespace
        namespaces: [team-a]
```

### RBAC group mapping (bring your own bindings)

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsEksAccessEntry
metadata:
  name: platform-operators
spec:
  region: us-west-2
  clusterName:
    valueFrom: { kind: AwsEksCluster, name: platform, fieldPath: status.outputs.name }
  principalArn:
    valueFrom: { kind: AwsIamRole, name: platform-operator, fieldPath: status.outputs.role_arn }
  kubernetesGroups:
    - platform-operators
  userName: "ops:{{SessionName}}"
```

## Stack Outputs

| Output | Description |
| --- | --- |
| `access_entry_arn` | The entry's ARN |
| `principal_arn` | The IAM principal granted access, as resolved at provisioning time |

## Related Components

- [AwsEksCluster](/docs/catalog/aws/eks-cluster) — the cluster access is granted on (needs API authentication mode)
- [AwsIamRole](/docs/catalog/aws/iam-role) — the principal being granted access
- [AwsEksNodeGroup](/docs/catalog/aws/eks-node-group) — managed nodes get their entries auto-created by EKS
