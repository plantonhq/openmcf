---
title: "NLB"
description: "NLB deployment documentation"
icon: "package"
order: 100
componentName: "awsnlb"
---

# AWS NLB

Deploys an AWS Network Load Balancer: the Layer-4 entry point that owns node
placement with optional Elastic-IP static addresses, optional security
groups, traffic distribution behavior (cross-zone, DNS client routing, zonal
shift), TLS-listener access logs, and optional Route53 alias DNS. Routing
lives in separate components -- listeners and target groups attach to this
NLB by ARN.

## What Gets Created

When you deploy an AwsNlb resource, Planton provisions:

- **Network Load Balancer** — an `aws_lb` / `lb.LoadBalancer` of type
  `network`, with one node per subnet mapping (optionally pinned to an
  Elastic IP or a specific private IPv4 address), optional security groups,
  and the configured attributes (cross-zone load balancing, DNS client
  routing policy, zonal shift, PrivateLink security-group enforcement)
- **S3 access-log delivery** — enabled by the presence of the `accessLogs`
  block; captures TLS-listener traffic only (an AWS limitation)
- **Route53 A records** — created only when DNS is enabled, one alias record
  per hostname pointing to the NLB's DNS name

Listeners and target groups are **not** created here — attach
`AwsLbListener` resources to this NLB's `load_balancer_arn` output
(protocols `TCP`, `UDP`, `TCP_UDP`, `TLS`; forward-only actions) and
`AwsLbTargetGroup` resources as destinations. Listener rules do not apply to
NLBs — routing is purely by port and protocol.

## Prerequisites

- **AWS credentials** configured via the Planton provider config (keyless SSO/OIDC).
- **At least one subnet** (public for internet-facing, private for internal);
  two or more across Availability Zones for high availability.
- **Elastic IP allocations** (e.g. `AwsElasticIp` resources) if you need
  static public IPs — one per subnet mapping, internet-facing only.
- **An S3 bucket with the ELB log-delivery bucket policy** if enabling
  access logs — delivery fails silently otherwise.
- **A Route53 hosted zone** if enabling DNS management.

## Quick Start

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsNlb
metadata:
  name: edge-nlb
spec:
  region: us-west-2
  subnetMappings:
    - subnetId:
        value: subnet-0a1b2c3d4e5f00001
    - subnetId:
        value: subnet-0a1b2c3d4e5f00002
```

```shell
planton apply -f nlb.yaml
```

This creates an internet-facing NLB with one node in each subnet. Add an
`AwsLbListener` against its `load_balancer_arn` output to start accepting
traffic.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
| --- | --- | --- | --- |
| `region` | `string` | AWS region where the NLB is created (e.g., `us-west-2`). | Required; non-empty |
| `subnetMappings` | `object[]` | One entry per Availability Zone the NLB places a node in. | Required; minimum 1 item (two or more recommended for HA) |
| `subnetMappings[].subnetId` | `string \| valueFrom` | Subnet for this node — public for internet-facing, private for internal. References `AwsSubnet.subnet_id` by default. | Required |

### Optional Fields

| Field | Type | Default | Description |
| --- | --- | --- | --- |
| `subnetMappings[].allocationId` | `string \| valueFrom` | AWS-assigned | Elastic IP for a static public IP on this node. Internet-facing only; at most one per mapping. References `AwsElasticIp.allocation_id` by default. |
| `subnetMappings[].privateIpv4Address` | `string` | AWS-assigned | Specific private IPv4 address for this node, from the subnet's CIDR. Internal NLBs only. |
| `securityGroups` | `(string \| valueFrom)[]` | none | Optional for NLB (unlike ALB). One-way door: once attached, the last group can never be removed, only replaced. References `AwsSecurityGroup.security_group_id` by default. |
| `internal` | `bool` | `false` | When `true`, the NLB is reachable only inside the VPC. Immutable — changing it replaces the load balancer. |
| `deleteProtectionEnabled` | `bool` | `false` | Prevents deletion while enabled. Deleting an NLB orphans its listeners, and pinned Elastic IPs start billing as unattached. |
| `crossZoneLoadBalancingEnabled` | `bool` | `false` | Distribute traffic evenly across targets in all AZs. Off by AWS default for NLB (unlike ALB) because inter-AZ transfer is billed. |
| `ipAddressType` | `string` | `ipv4` | `ipv4` or `dualstack`. |
| `dnsRecordClientRoutingPolicy` | `string` | `any_availability_zone` | How DNS routes clients to NLB nodes: `any_availability_zone`, `availability_zone_affinity` (resolver's AZ — least cross-zone traffic), or `partial_availability_zone_affinity` (85% affinity, 15% spillover). |
| `zonalShiftEnabled` | `bool` | `false` | Allows Amazon Application Recovery Controller to shift traffic away from an impaired Availability Zone. |
| `enforceSecurityGroupInboundRulesOnPrivateLinkTraffic` | `string` | `on` (AWS, with SGs) | Whether inbound security-group rules apply to traffic arriving through PrivateLink VPC endpoints: `on` or `off`. Only meaningful when security groups are attached. |
| `accessLogs` | `object` | off | TLS-listener access logs to S3: `bucket` (references `AwsS3Bucket.bucket_id` by default) and optional `prefix`. Plain TCP/UDP flows are not logged — an AWS limitation. |
| `dns` | `object` | off | Route53 alias DNS: `enabled`, `route53ZoneId` (references `AwsRoute53Zone.zone_id` by default), `hostnames` (unique). |

## Examples

### Internal NLB with a pinned private address

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsNlb
metadata:
  name: internal-tcp-nlb
spec:
  region: us-west-2
  internal: true
  subnetMappings:
    - subnetId:
        value: subnet-private-az1
      privateIpv4Address: 10.0.1.50
    - subnetId:
        value: subnet-private-az2
```

### Static IPs via Elastic IP references

Internet-facing NLB whose public IPs never change — the addresses partners
allowlist:

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsNlb
metadata:
  name: partner-ingress
spec:
  region: us-west-2
  subnetMappings:
    - subnetId:
        valueFrom:
          kind: AwsSubnet
          name: public-az1
          fieldPath: status.outputs.subnet_id
      allocationId:
        valueFrom:
          kind: AwsElasticIp
          name: ingress-eip-az1
          fieldPath: status.outputs.allocation_id
    - subnetId:
        valueFrom:
          kind: AwsSubnet
          name: public-az2
          fieldPath: status.outputs.subnet_id
      allocationId:
        valueFrom:
          kind: AwsElasticIp
          name: ingress-eip-az2
          fieldPath: status.outputs.allocation_id
  deleteProtectionEnabled: true
```

### Hardened NLB with security groups and access logs

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsNlb
metadata:
  name: prod-nlb
spec:
  region: us-west-2
  subnetMappings:
    - subnetId:
        value: subnet-private-az1
    - subnetId:
        value: subnet-private-az2
  internal: true
  securityGroups:
    - valueFrom:
        kind: AwsSecurityGroup
        name: nlb-sg
        fieldPath: status.outputs.security_group_id
  enforceSecurityGroupInboundRulesOnPrivateLinkTraffic: "on"
  deleteProtectionEnabled: true
  accessLogs:
    bucket:
      value: prod-nlb-logs-bucket
    prefix: nlb/prod
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
| --- | --- | --- |
| `load_balancer_arn` | `string` | ARN of the NLB — what `AwsLbListener` resources attach through, and what Global Accelerator endpoints reference |
| `load_balancer_name` | `string` | Name assigned to the NLB (`metadata.name`, truncated to AWS's 32-character limit when necessary) |
| `load_balancer_dns_name` | `string` | DNS name automatically assigned by AWS (e.g., `edge-nlb-abc123.elb.us-west-2.amazonaws.com`) |
| `load_balancer_hosted_zone_id` | `string` | Route53 hosted zone ID for the NLB's DNS name, used for alias records |

## Related Components

- [AwsLbListener](/docs/catalog/aws/lb-listener) — attaches to this NLB by ARN; owns port, protocol, and TLS material (forward-only actions)
- [AwsLbTargetGroup](/docs/catalog/aws/lb-target-group) — the destination of forward actions; `targetType: alb` expresses the NLB-in-front-of-ALB pattern
- [AwsElasticIp](/docs/catalog/aws/elastic-ip) — provides `allocationId` values for static public IPs in subnet mappings
- [AwsSubnet](/docs/catalog/aws/subnet) — provides the subnets for node placement
- [AwsSecurityGroup](/docs/catalog/aws/security-group) — optional inbound traffic filtering for the NLB
- [AwsS3Bucket](/docs/catalog/aws/s3-bucket) — receives TLS-listener access logs
- [AwsRoute53Zone](/docs/catalog/aws/route53-zone) — hosts the DNS zone for alias records
- [AwsAlb](/docs/catalog/aws/alb) — the Layer-7 alternative; combine both by registering an ALB as an NLB target
