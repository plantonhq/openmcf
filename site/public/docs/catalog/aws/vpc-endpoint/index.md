---
title: "VPC Endpoint"
description: "VPC Endpoint deployment documentation"
icon: "package"
order: 100
componentName: "awsvpcendpoint"
---

# AWS VPC Endpoint

Creates a private connection from a VPC to an AWS service, a PrivateLink
service, or a VPC Lattice resource -- S3/DynamoDB gateway endpoints on
referenced route tables, ENI-based interface endpoints with private DNS
on referenced subnets and security groups -- keeping traffic on the AWS
network instead of the internet.

## What Gets Created

When you deploy an AwsVpcEndpoint resource, Planton provisions:

- **VPC endpoint** — an `aws_vpc_endpoint` / `ec2.VpcEndpoint` of the
  chosen type in the referenced VPC: a Gateway endpoint injecting the
  service's prefix-list route into your route tables, or an
  Interface/GatewayLoadBalancer/Lattice endpoint placing one ENI per
  referenced subnet, with your policy, DNS options, and IP address type

The endpoint never modifies a resource it merely references: route
tables keep their own routes (AWS manages the endpoint's prefix-list
route as part of the endpoint), and security groups carry their own
rules.

## Prerequisites

- **AWS credentials** configured via the Planton provider config (keyless SSO/OIDC).
- **A VPC** (`AwsVpc`) for the endpoint to live in.
- **Route tables** for gateway endpoints — an `AwsSubnet`'s `route_table_id` output (when the subnet owns its table) or the VPC's `main_route_table_id` / `default_route_table_id` outputs.
- **Subnets** (`AwsSubnet`) for interface endpoints — one ENI per subnet, spread across AZs.
- **DNS support and DNS hostnames enabled on the VPC** if you turn on `privateDnsEnabled`.

## Quick Start

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsVpcEndpoint
metadata:
  name: platform-s3-gateway
spec:
  region: us-west-2
  vpcId:
    valueFrom:
      kind: AwsVpc
      name: platform-vpc
      fieldPath: status.outputs.vpc_id
  serviceName: com.amazonaws.us-west-2.s3
  routeTableIds:
    - valueFrom:
        kind: AwsSubnet
        name: platform-private-1
        fieldPath: status.outputs.route_table_id
```

```shell
planton apply -f endpoint.yaml
```

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
| --- | --- | --- | --- |
| `region` | `string` | AWS region; must match the VPC's. | Required; non-empty |
| `vpcId` | `string \| valueFrom` | The VPC the endpoint lives in. Defaults to referencing an `AwsVpc` `vpc_id` output. Create-only. | Required |
| service target | `string` | Exactly one of `serviceName`, `resourceConfigurationArn`, or `serviceNetworkArn`. Create-only. | Exactly one |

### Optional Fields

| Field | Type | Default | Description |
| --- | --- | --- | --- |
| `endpointType` | `string` | `Gateway` | `Gateway`, `Interface`, `GatewayLoadBalancer`, `Resource`, or `ServiceNetwork`. Create-only. |
| `routeTableIds` | `(string \| valueFrom)[]` | `[]` | Route tables a Gateway endpoint injects its prefix-list route into. Gateway only. |
| `subnetIds` | `(string \| valueFrom)[]` | `[]` | Subnets an ENI-based endpoint places interfaces in (one per subnet). Not for Gateway. |
| `securityGroupIds` | `(string \| valueFrom)[]` | VPC default SG | Security groups on an Interface endpoint's ENIs — must allow the service port (443) inbound. |
| `privateDnsEnabled` | `bool` | `false` | Resolve the service's public DNS name to the endpoint inside the VPC. Interface only; needs VPC DNS support + hostnames. |
| `dnsOptions.dnsRecordIpType` | `string` | AWS picks | `ipv4`, `ipv6`, `dualstack`, or `service-defined`. |
| `dnsOptions.privateDnsOnlyForInboundResolverEndpoint` | `bool` | `false` | The S3 dual-stack pattern: in-VPC traffic rides the gateway endpoint; on-premises resolver traffic reaches this interface endpoint. |
| `dnsOptions.privateDnsPreference` | `string` | AWS default | Lattice types: which private domains get hosted zones. Create-only. |
| `dnsOptions.privateDnsSpecifiedDomains` | `string[]` | `[]` | Required exactly when the preference includes specified domains. 1–10 domains. Create-only. |
| `ipAddressType` | `string` | AWS picks | `ipv4`, `dualstack`, or `ipv6`; the service must support it. |
| `policy` | `string` | full access | IAM policy document scoping who may use the endpoint for what. |
| `subnetConfigurations` | `object[]` | `[]` | Pin static IPv4/IPv6 ENI addresses per subnet (each subnet must also be in `subnetIds`). |
| `serviceRegion` | `string` | endpoint's region | Cross-region target (Interface only). Create-only. |
| `autoAccept` | `bool` | `false` | Auto-accept the connection for a same-account PrivateLink service. |

## Examples

### Interface endpoint for ECR with private DNS

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsVpcEndpoint
metadata:
  name: platform-ecr-dkr
spec:
  region: us-west-2
  vpcId:
    valueFrom: { kind: AwsVpc, name: platform-vpc, fieldPath: status.outputs.vpc_id }
  endpointType: Interface
  serviceName: com.amazonaws.us-west-2.ecr.dkr
  subnetIds:
    - valueFrom: { kind: AwsSubnet, name: platform-private-1, fieldPath: status.outputs.subnet_id }
    - valueFrom: { kind: AwsSubnet, name: platform-private-2, fieldPath: status.outputs.subnet_id }
  securityGroupIds:
    - valueFrom: { kind: AwsSecurityGroup, name: platform-endpoints-sg, fieldPath: status.outputs.security_group_id }
  privateDnsEnabled: true
```

### S3 gateway endpoint with a scoping policy

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsVpcEndpoint
metadata:
  name: platform-s3-gateway
spec:
  region: us-west-2
  vpcId:
    valueFrom: { kind: AwsVpc, name: platform-vpc, fieldPath: status.outputs.vpc_id }
  serviceName: com.amazonaws.us-west-2.s3
  routeTableIds:
    - valueFrom: { kind: AwsSubnet, name: platform-private-1, fieldPath: status.outputs.route_table_id }
    - valueFrom: { kind: AwsSubnet, name: platform-private-2, fieldPath: status.outputs.route_table_id }
  policy: |
    {
      "Version": "2012-10-17",
      "Statement": [{
        "Effect": "Allow",
        "Principal": "*",
        "Action": ["s3:GetObject", "s3:PutObject", "s3:ListBucket"],
        "Resource": ["arn:aws:s3:::platform-*", "arn:aws:s3:::platform-*/*"]
      }]
    }
```

### Third-party PrivateLink service

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsVpcEndpoint
metadata:
  name: vendor-privatelink
spec:
  region: us-west-2
  vpcId:
    valueFrom: { kind: AwsVpc, name: platform-vpc, fieldPath: status.outputs.vpc_id }
  endpointType: Interface
  serviceName: com.amazonaws.vpce.us-west-2.vpce-svc-0a1b2c3d4e5f67890
  subnetIds:
    - valueFrom: { kind: AwsSubnet, name: platform-private-1, fieldPath: status.outputs.subnet_id }
```

## Stack Outputs

| Output | Description |
| --- | --- |
| `vpc_endpoint_id` | The endpoint's id (vpce-...) |
| `arn` | The endpoint's ARN |
| `state` | Lifecycle state — `available`, or `pendingAcceptance` when the service requires manual acceptance |
| `prefix_list_id` | The service's prefix list (gateway endpoints only) — for security-group and route rules scoped to the service |
| `dns_name` | The primary private DNS name (interface endpoints only) — the Route53 alias target |
| `hosted_zone_id` | The Route53 zone of `dns_name`, needed alongside it for alias records |
| `network_interface_ids` | The endpoint's ENIs, one per attached subnet |

## Related Components

- [AwsVpc](/docs/catalog/aws/vpc) — the VPC the endpoint lives in (and the main/default route-table outputs gateway endpoints attach to)
- [AwsSubnet](/docs/catalog/aws/subnet) — subnets for interface ENIs; subnet-owned route tables for gateway endpoints
- [AwsSecurityGroup](/docs/catalog/aws/security-group) — controls client access to an interface endpoint's ENIs
- [AwsRoute53Zone](/docs/catalog/aws/route53-zone) — private zones that alias onto the endpoint's DNS name
