---
title: "EKS Cluster"
description: "EKS Cluster deployment documentation"
icon: "package"
order: 100
componentName: "awsekscluster"
---

# AWS EKS Cluster

Deploys an EKS cluster control plane: the managed Kubernetes API server
and etcd, with networking exposure, access-entry authentication, secrets
encryption, control-plane logging, upgrade policy, and optional EKS Auto
Mode -- the foundation node that `AwsEksNodeGroup` compute and
IRSA identity compose onto.

## What Gets Created

When you deploy an AwsEksCluster resource, Planton provisions:

- **EKS cluster** — an `aws_eks_cluster` / `eks.Cluster` named from
  `metadata.name`, attached to the referenced subnets and cluster role,
  with the endpoint exposure, access config, encryption, logging,
  upgrade, and zonal-shift posture you declare
- **Auto Mode capabilities** (when `autoMode.enabled`) — AWS-managed
  compute, block storage, and load balancing, expanded from the single
  toggle into the three settings the EKS API requires to move together

The cluster role is never modified: attach `AmazonEKSClusterPolicy` on
the referenced `AwsIamRole` itself (`managedPolicyArns`).

## Prerequisites

- **AWS credentials** configured via the Planton provider config (keyless SSO/OIDC).
- **Subnets** (`AwsSubnet`) in at least two availability zones.
- **A cluster IAM role** (`AwsIamRole`) trusting `eks.amazonaws.com` with `AmazonEKSClusterPolicy` attached.
- **A KMS key** (`AwsKmsKey`) if you want customer-managed secrets encryption.

## Quick Start

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsEksCluster
metadata:
  name: platform
spec:
  region: us-west-2
  subnetIds:
    - valueFrom:
        kind: AwsSubnet
        name: private-a
        fieldPath: status.outputs.subnet_id
    - valueFrom:
        kind: AwsSubnet
        name: private-b
        fieldPath: status.outputs.subnet_id
  clusterRoleArn:
    valueFrom:
      kind: AwsIamRole
      name: eks-cluster-role
      fieldPath: status.outputs.role_arn
  accessConfig:
    authenticationMode: API
```

```shell
planton apply -f eks-cluster.yaml
```

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
| --- | --- | --- | --- |
| `region` | `string` | AWS region for the control plane. | Required; non-empty |
| `subnetIds` | `string[] \| valueFrom` | Subnets for the control plane's network interfaces, in ≥2 AZs. Defaults to referencing `AwsSubnet` `subnet_id` outputs. | Required; ≥2 entries |
| `clusterRoleArn` | `string \| valueFrom` | The IAM role EKS assumes (must trust `eks.amazonaws.com` and carry `AmazonEKSClusterPolicy`). Defaults to referencing an `AwsIamRole` `role_arn` output. | Required |

### Optional Fields

| Field | Type | Default | Description |
| --- | --- | --- | --- |
| `version` | `string` | AWS default | Kubernetes minor (`"1.31"`); ≥1.24; EKS never downgrades. |
| `securityGroupIds` | `string[] \| valueFrom` | `[]` | Extra control-plane security groups on top of the EKS-managed one. |
| `endpointPublicAccess` | `bool` | `true` | Internet reachability of the API server. Explicit `false` = private-only (pair with private access). |
| `endpointPrivateAccess` | `bool` | `false` | In-VPC API endpoint. Almost always wanted in production. |
| `publicAccessCidrs` | `string[]` | `0.0.0.0/0` | CIDRs allowed at the public endpoint — the cheapest hardening step. |
| `controlPlaneEgressMode` | `string` | `AWS_MANAGED` | `CUSTOMER_ROUTED` (your route tables) or `CUSTOMER_ISOLATED`; reverting CUSTOMER_ROUTED replaces the cluster. |
| `ipFamily` | `string` | `ipv4` | `ipv4` or `ipv6` pod/service networking. Create-only. |
| `serviceIpv4Cidr` | `string` | AWS default | Service CIDR (/12–/24, private ranges, non-overlapping). Create-only. |
| `enabledClusterLogTypes` | `string[]` | `[]` | Any of `api`, `audit`, `authenticator`, `controllerManager`, `scheduler` to CloudWatch. |
| `kmsKeyArn` | `string \| valueFrom` | AWS-owned keys | Customer-managed envelope encryption of secrets. One-way: cannot be disabled later. |
| `accessConfig.authenticationMode` | `string` | `API_AND_CONFIG_MAP` | `API` (access entries), `API_AND_CONFIG_MAP`, or `CONFIG_MAP`. Migration toward `API` is one-way. |
| `accessConfig.bootstrapClusterCreatorAdminPermissions` | `bool` | `true` | Whether the creating identity gets cluster-admin. Create-only. |
| `autoMode.enabled` | `bool` | `false` | EKS Auto Mode: AWS manages compute, block storage, and load balancing (the API's all-or-nothing trio, one toggle). |
| `autoMode.nodePools` | `string[]` | `[]` | Built-in pools: `general-purpose`, `system`. Empty = custom in-cluster NodePools only. |
| `autoMode.nodeRoleArn` | `string \| valueFrom` | — | Node identity for Auto Mode capacity. Required when `nodePools` is set. |
| `upgradeSupportType` | `string` | `EXTENDED` | `STANDARD` (upgrade on schedule, no surcharge) or `EXTENDED`. |
| `zonalShiftEnabled` | `bool` | `false` | Amazon Application Recovery Controller zonal shift. |
| `deletionProtection` | `bool` | `false` | EKS refuses deletion until disabled. |
| `bootstrapSelfManagedAddons` | `bool` | `true` | Install default vpc-cni/kube-proxy/CoreDNS at creation. `false` = bring-your-own add-ons. Create-only. |
| `forceUpdateVersion` | `bool` | `false` | Force version updates past unsatisfiable pod disruption budgets. |

## Examples

### Private, encrypted, deletion-protected control plane

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsEksCluster
metadata:
  name: prod
spec:
  region: us-west-2
  subnetIds:
    - valueFrom: { kind: AwsSubnet, name: private-a, fieldPath: status.outputs.subnet_id }
    - valueFrom: { kind: AwsSubnet, name: private-b, fieldPath: status.outputs.subnet_id }
  clusterRoleArn:
    valueFrom: { kind: AwsIamRole, name: eks-cluster-role, fieldPath: status.outputs.role_arn }
  endpointPublicAccess: false
  endpointPrivateAccess: true
  kmsKeyArn:
    valueFrom: { kind: AwsKmsKey, name: eks-secrets, fieldPath: status.outputs.key_arn }
  accessConfig:
    authenticationMode: API
  enabledClusterLogTypes: [audit, authenticator]
  deletionProtection: true
```

### EKS Auto Mode (no node groups to operate)

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsEksCluster
metadata:
  name: hands-off
spec:
  region: us-west-2
  subnetIds:
    - valueFrom: { kind: AwsSubnet, name: private-a, fieldPath: status.outputs.subnet_id }
    - valueFrom: { kind: AwsSubnet, name: private-b, fieldPath: status.outputs.subnet_id }
  clusterRoleArn:
    valueFrom: { kind: AwsIamRole, name: eks-cluster-role, fieldPath: status.outputs.role_arn }
  autoMode:
    enabled: true
    nodePools: [general-purpose, system]
    nodeRoleArn:
      valueFrom: { kind: AwsIamRole, name: eks-auto-node-role, fieldPath: status.outputs.role_arn }
  bootstrapSelfManagedAddons: false
```

### IRSA in one reference

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsIamOidcProvider
metadata:
  name: platform-irsa
spec:
  region: us-west-2
  url:
    valueFrom:
      kind: AwsEksCluster
      name: platform
      fieldPath: status.outputs.oidc_issuer_url
  clientIdList:
    - sts.amazonaws.com
```

## Stack Outputs

| Output | Description |
| --- | --- |
| `endpoint` | Kubernetes API server URL |
| `cluster_ca_certificate` | Base64 cluster CA certificate, for kubeconfigs |
| `cluster_security_group_id` | The EKS-managed control-plane security group (node groups use it for control-plane communication) |
| `oidc_issuer_url` | The OIDC issuer — point an `AwsIamOidcProvider` here to enable IRSA |
| `cluster_arn` | Cluster ARN, for IAM policies and access entries |
| `name` | Cluster name — what `AwsEksNodeGroup.clusterName` references |
| `platform_version` | EKS platform revision of the control plane (e.g. `eks.12`) |

## Related Components

- [AwsEksNodeGroup](/docs/catalog/aws/eks-node-group) — managed worker fleets registered to this cluster
- [AwsIamRole](/docs/catalog/aws/iam-role) — the cluster role (and Auto Mode node role) this cluster assumes
- [AwsIamOidcProvider](/docs/catalog/aws/iam-oidc-provider) — IRSA trust anchor fed by `oidc_issuer_url`
- [AwsSubnet](/docs/catalog/aws/subnet) — where the control plane attaches its network interfaces
- [AwsKmsKey](/docs/catalog/aws/kms-key) — customer-managed secrets encryption
- [AwsSecurityGroup](/docs/catalog/aws/security-group) — additional control-plane security groups
