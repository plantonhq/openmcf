---
title: "HTTP API VPC Link"
description: "HTTP API VPC Link deployment documentation"
icon: "package"
order: 100
componentName: "awshttpapivpclink"
---

# AWS HTTP API VPC Link

Deploys an API Gateway v2 VPC link -- the managed network attachment that lets HTTP APIs proxy traffic to private backends (internal ALBs, NLBs, or Cloud Map services) inside a VPC without exposing them to the internet.

## What Gets Created

When you deploy an AwsHttpApiVpcLink resource, Planton provisions:

- **VPC link** -- a set of AWS-managed elastic network interfaces in the subnets you choose, tagged with Planton identity tags

HTTP API private integrations then reference the link by ID (`connectionId` with `connectionType: VPC_LINK`).

## Prerequisites

- **AWS credentials** configured via environment variables or Planton provider config
- **Subnets** in the target VPC (two or more availability zones recommended -- the link can only reach targets in AZs it has an ENI in)
- (Optional) **A security group** allowing egress to the target ALB/NLB listener ports

## Quick Start

Create a file `vpc-link.yaml`:

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsHttpApiVpcLink
metadata:
  name: private-services-link
spec:
  region: us-east-1
  subnetIds:
    - value: subnet-0abc123def456
    - value: subnet-0def456abc789
  securityGroupIds:
    - value: sg-0abc123def456
```

Deploy:

```shell
planton apply -f vpc-link.yaml
```

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `region` | `string` | The AWS region where the VPC link will be created | Non-empty |
| `subnetIds` | `StringValueOrRef[]` | Subnets for the link's network interfaces. Immutable after creation (changing the set replaces the link). Can reference `AwsSubnet` via `valueFrom`. | Minimum 1 item |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `securityGroupIds` | `StringValueOrRef[]` | `[]` | Security groups on the link's ENIs. Immutable after creation. When omitted, AWS applies no filtering on the link side -- reachability is governed solely by the target's security groups. Can reference `AwsSecurityGroup` via `valueFrom`. |

## Examples

### Composed from Planton-Managed Networking

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsHttpApiVpcLink
metadata:
  name: internal-alb-link
spec:
  region: us-east-1
  subnetIds:
    - valueFrom:
        kind: AwsSubnet
        name: private-subnet-a
        fieldPath: status.outputs.subnet_id
    - valueFrom:
        kind: AwsSubnet
        name: private-subnet-b
        fieldPath: status.outputs.subnet_id
  securityGroupIds:
    - valueFrom:
        kind: AwsSecurityGroup
        name: vpc-link-egress
        fieldPath: status.outputs.security_group_id
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `vpc_link_id` | `string` | The VPC link ID -- what private integrations set as their `connectionId` |
| `vpc_link_arn` | `string` | ARN of the VPC link, for IAM policies and tag-based governance |

## Related Components

- [AwsHttpApiGateway](/docs/catalog/aws/http-api-gateway) — HTTP APIs whose private integrations route through this link
- [AwsSubnet](/docs/catalog/aws/subnet) — Subnets hosting the link's network interfaces
- [AwsSecurityGroup](/docs/catalog/aws/security-group) — Security groups governing the link's reach
- [AwsAlb](/docs/catalog/aws/alb) — Internal ALBs targeted by private integrations
- [AwsNlb](/docs/catalog/aws/nlb) — Internal NLBs targeted by private integrations
