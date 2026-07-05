---
title: "MSK Serverless Cluster"
description: "MSK Serverless Cluster deployment documentation"
icon: "package"
order: 100
componentName: "awsmskserverlesscluster"
---

# AWS MSK Serverless Cluster

Deploys an Amazon MSK Serverless cluster — fully managed Apache Kafka with automatic capacity scaling and pay-per-throughput billing. There are no brokers, instance types, storage volumes, or Kafka version to configure: the declaration is where the cluster lives (subnets and security groups), and clients authenticate with AWS IAM (SASL/IAM) on port 9098.

## What Gets Created

When you deploy an AwsMskServerlessCluster resource, Planton provisions:

- **MSK Serverless Cluster** — an `aws_msk_serverless_cluster` with network interfaces in the referenced subnets, the referenced security groups attached, and SASL/IAM client authentication enabled (the only scheme serverless MSK supports)

The resource is effectively immutable: everything except tags is create-time, so changing networking replaces the cluster.

## Prerequisites

- **AWS credentials** configured via environment variables or Planton provider config
- **Private subnets** in 2+ Availability Zones (can be managed by AwsSubnet resources)
- **A security group** whose ingress opens the SASL/IAM port (9098) to your clients — optional; AWS attaches the VPC default group when omitted

## Quick Start

Create a file `msk-serverless.yaml`:

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsMskServerlessCluster
metadata:
  name: my-kafka
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AwsMskServerlessCluster.my-kafka
spec:
  region: us-west-2
  subnetIds:
    - value: subnet-0a1b2c3d4e5f00001
    - value: subnet-0a1b2c3d4e5f00002
  securityGroupIds:
    - value: sg-0a1b2c3d4e5f00001
```

Deploy:

```shell
planton apply -f msk-serverless.yaml
```

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `region` | `string` | AWS region where the cluster is created. | Required; non-empty |
| `subnetIds` | `StringValueOrRef[]` | VPC subnets for the cluster network interfaces. Changing subnets forces replacement. | Minimum 1 item |
| `subnetIds[].value` | `string` | Direct subnet ID value | — |
| `subnetIds[].valueFrom` | `object` | Foreign key reference to an AwsSubnet resource | Default field: `status.outputs.subnet_id` |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `securityGroupIds` | `StringValueOrRef[]` | VPC default group | Security groups attached to the cluster network interfaces (max 5). The port-9098 ingress rule lives on the referenced AwsSecurityGroup nodes. Changing forces replacement. |

SASL/IAM authentication is not a field: AWS requires it and offers no alternative, so both IaC engines enable it unconditionally.

## Examples

### Composed with References

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsMskServerlessCluster
metadata:
  name: events-kafka
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AwsMskServerlessCluster.events-kafka
spec:
  region: us-west-2
  subnetIds:
    - valueFrom:
        kind: AwsSubnet
        name: platform-private-a
        fieldPath: status.outputs.subnet_id
    - valueFrom:
        kind: AwsSubnet
        name: platform-private-b
        fieldPath: status.outputs.subnet_id
  securityGroupIds:
    - valueFrom:
        kind: AwsSecurityGroup
        name: kafka-broker-sg
        fieldPath: status.outputs.security_group_id
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `cluster_arn` | `string` | ARN of the cluster (also its resource identifier) — used in IAM policies and Lambda event source mappings |
| `cluster_name` | `string` | Human-readable cluster name |
| `cluster_uuid` | `string` | Unique identifier extracted from the cluster ARN |
| `bootstrap_brokers_sasl_iam` | `string` | Comma-separated SASL/IAM broker endpoint list (port 9098) — the only connection string serverless MSK exposes |

## Related Components

- [AwsMskCluster](/docs/catalog/aws/msk-cluster) — the provisioned sibling: choose it for SCRAM/mTLS auth, public access, PrivateLink, tiered storage, or sustained high throughput
- [AwsSubnet](/docs/catalog/aws/subnet) — provides the subnets where the cluster places network interfaces
- [AwsSecurityGroup](/docs/catalog/aws/security-group) — controls which clients can reach the SASL/IAM listener (port 9098)
- [AwsIamRole](/docs/catalog/aws/iam-role) — client workloads need `kafka-cluster:*` IAM permissions scoped to the exported `cluster_arn`
