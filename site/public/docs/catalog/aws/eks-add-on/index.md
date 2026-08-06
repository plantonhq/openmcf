---
title: "EKS Add-on"
description: "EKS Add-on deployment documentation"
icon: "package"
order: 100
componentName: "awseksaddon"
---

# AWS EKS Add-on

Installs an EKS managed add-on on an `AwsEksCluster`: cluster software
(vpc-cni, coredns, kube-proxy, CSI drivers, Pod Identity agent,
marketplace add-ons) whose install, upgrades, configuration, and
removal AWS manages through the EKS control plane -- with IAM identity
wired by reference through IRSA or EKS Pod Identity.

## What Gets Created

When you deploy an AwsEksAddon resource, Planton provisions:

- **Managed add-on** — an `aws_eks_addon` / `eks.Addon` keyed on the
  referenced cluster and the spec's `addonName`, at the AWS default
  version or your pin, with your configuration values and conflict
  handling

The IAM role an add-on assumes is never modified: attach the add-on's
required policies on the referenced `AwsIamRole` itself.

## Prerequisites

- **AWS credentials** configured via the Planton provider config (keyless SSO/OIDC).
- **An EKS cluster** (`AwsEksCluster`) for the add-on to install on.
- **An IAM role** (`AwsIamRole`) if the add-on needs its own identity — trusting the cluster's OIDC provider (IRSA) or `pods.eks.amazonaws.com` (Pod Identity).
- **The eks-pod-identity-agent add-on** on the cluster if you use `podIdentityAssociations`.

## Quick Start

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsEksAddon
metadata:
  name: platform-coredns
spec:
  region: us-west-2
  clusterName:
    valueFrom:
      kind: AwsEksCluster
      name: platform
      fieldPath: status.outputs.name
  addonName: coredns
  resolveConflictsOnCreate: OVERWRITE
```

```shell
planton apply -f addon.yaml
```

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
| --- | --- | --- | --- |
| `region` | `string` | AWS region; must match the cluster's. | Required; non-empty |
| `clusterName` | `string \| valueFrom` | The cluster the add-on installs on. Defaults to referencing an `AwsEksCluster` `name` output. | Required |
| `addonName` | `string` | The EKS catalog name (`vpc-cni`, `coredns`, `aws-ebs-csi-driver`, ...). One add-on name per cluster. Create-only. | Required; ≤100 chars |

### Optional Fields

| Field | Type | Default | Description |
| --- | --- | --- | --- |
| `addonVersion` | `string` | AWS default | Full version like `v1.18.1-eksbuild.3`. Empty follows the AWS default for the cluster's Kubernetes version. |
| `resolveConflictsOnCreate` | `string` | `NONE` | `NONE` or `OVERWRITE`. Use `OVERWRITE` to adopt bootstrap self-managed installs of vpc-cni/coredns/kube-proxy. |
| `resolveConflictsOnUpdate` | `string` | `NONE` | `NONE`, `OVERWRITE`, or `PRESERVE` — what happens to out-of-band edits on update. |
| `configurationValues` | `string` | `{}` | The add-on's own JSON configuration (`aws eks describe-addon-configuration` shows the schema). |
| `serviceAccountRoleArn` | `string \| valueFrom` | node role | IRSA role for the add-on's service account; requires the cluster's OIDC provider. |
| `podIdentityAssociations` | `object[]` | `[]` | Pod Identity bindings (`roleArn` + `serviceAccount`); the modern no-OIDC-provider path. |
| `preserve` | `bool` | `false` | Keep the add-on's Kubernetes resources running (self-managed) when this resource is deleted. |
| `namespaceConfig.namespace` | `string` | add-on default | Custom install namespace (RFC 1123 label). Create-only. |

## Examples

### EBS CSI driver with Pod Identity

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsEksAddon
metadata:
  name: platform-ebs-csi
spec:
  region: us-west-2
  clusterName:
    valueFrom: { kind: AwsEksCluster, name: platform, fieldPath: status.outputs.name }
  addonName: aws-ebs-csi-driver
  podIdentityAssociations:
    - roleArn:
        valueFrom: { kind: AwsIamRole, name: ebs-csi-driver, fieldPath: status.outputs.role_arn }
      serviceAccount: ebs-csi-controller-sa
```

### Pinned coredns with custom replica count

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsEksAddon
metadata:
  name: platform-coredns
spec:
  region: us-west-2
  clusterName:
    valueFrom: { kind: AwsEksCluster, name: platform, fieldPath: status.outputs.name }
  addonName: coredns
  addonVersion: v1.11.3-eksbuild.2
  configurationValues: '{"replicaCount":3}'
  resolveConflictsOnCreate: OVERWRITE
  resolveConflictsOnUpdate: OVERWRITE
```

## Stack Outputs

| Output | Description |
| --- | --- |
| `addon_arn` | The add-on's ARN |
| `addon_name` | The EKS catalog name it was installed under |
| `addon_version` | The version actually running — the resolved AWS default when the spec pinned nothing |

## Related Components

- [AwsEksCluster](/docs/catalog/aws/eks-cluster) — the control plane the add-on installs on
- [AwsIamRole](/docs/catalog/aws/iam-role) — the add-on's IAM identity (IRSA or Pod Identity)
- [AwsIamOidcProvider](/docs/catalog/aws/iam-oidc-provider) — required for IRSA wiring
- [AwsEksNodeGroup](/docs/catalog/aws/eks-node-group) — the compute the add-on's pods run on
