---
title: "ALB"
description: "ALB deployment documentation"
icon: "package"
order: 100
componentName: "awsalb"
---

# AWS ALB

Deploys an Application Load Balancer — the Layer-7 entry point that terminates HTTP/HTTPS and hands requests to the routing graph. The ALB itself carries no routing configuration by design: listeners (AwsLbListener) attach to it and own ports, TLS material, and default actions; listener rules (AwsLbListenerRule) own per-service routing; target groups (AwsLbTargetGroup) receive the traffic. This component owns what is truly load-balancer-wide — placement, security groups, and the HTTP behavior attributes.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Application Load Balancer** -- with its scheme, subnets, security groups, address family, and the configured HTTP behavior attributes (timeouts, HTTP/2, header handling, desync mitigation, WAF fail mode, zonal shift)
- **S3 log delivery** -- for whichever of the three streams (access, connection, health-check logs) has a bucket configured
- **Route53 alias A records** -- created only when DNS is enabled, one alias record per hostname pointing at the ALB's DNS name
- **AWS Tags** -- resource metadata tags (organization, environment, resource kind, resource ID) applied to the load balancer

Listeners, rules, and target groups are **not** created here — attach AwsLbListener resources through this ALB's `load_balancer_arn` output; TLS certificates live on the listeners.

## Before You Deploy

### Planton Setup

- **AWS Provider Connection** -- an active connection in the Connect module with credentials for the target AWS account. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Subnets** -- at least two AwsSubnet resources in different Availability Zones, referenced by their `subnet_id` outputs (public for internet-facing, private for internal).
- **Security Groups** -- AwsSecurityGroup resources opening exactly the listener ports; without them AWS attaches the VPC's default group (fine for a first boot, wrong for production).
- **S3 Bucket / Route53 Zone** -- optional AwsS3Bucket and AwsRoute53Zone resources for log delivery and alias DNS.

### AWS Account

- **ELB permissions** -- the credentials used by the Provider Connection must have `elasticloadbalancing:CreateLoadBalancer`, `DescribeLoadBalancers`, `ModifyLoadBalancerAttributes`, and `DeleteLoadBalancer`, plus `ec2:Describe*` on the subnets.
- **Two-AZ minimum** -- AWS requires ALB subnets in at least two Availability Zones.
- **Log-delivery bucket policy** -- when enabling any log stream, the S3 bucket must grant the regional ELB log-delivery service write access; delivery fails silently otherwise.

## Deploy

### Console

Open the deployment store, find **AWS ALB**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Internet Facing** preset in the [Presets](#presets) tab.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsAlb
metadata:
  name: api-gateway
  org: acme-corp
  env: prod
spec:
  region: us-east-1
  subnets:
    - valueFrom:
        kind: AwsSubnet
        name: public-subnet-a
        fieldPath: status.outputs.subnet_id
    - valueFrom:
        kind: AwsSubnet
        name: public-subnet-b
        fieldPath: status.outputs.subnet_id
  securityGroups:
    - valueFrom:
        kind: AwsSecurityGroup
        name: alb-ingress
        fieldPath: status.outputs.security_group_id
  deleteProtectionEnabled: true
```

```shell
planton apply -f alb.yaml
```

This creates an internet-facing ALB across two zones with explicit security groups and deletion protection. Attach AwsLbListener resources to route traffic. A Stack Job tracks the provisioning in real time.

## Key Configuration

These are the most important decisions when configuring an ALB. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Routing lives elsewhere** -- listeners own ports, TLS certificates, and default actions; rules own per-service routing; target groups receive traffic. The ALB spec is placement and behavior only.

**The scheme is a one-way door** -- `internal` decides internet-facing vs VPC-only and cannot change; AWS replaces the load balancer. Everything else — subnets, security groups, timeouts, header handling — updates in place.

**Security groups: explicit beats default** -- omitted groups fall back to the VPC's default security group. Production ALBs attach explicit groups opening exactly the listener ports.

**Idle timeout causes the classic 504** -- a response slower than the idle timeout (default 60s) times out at the ALB while the target still works. Raise it above the slowest legitimate response, and keep it below the targets' own keep-alive timeouts.

**Header handling matters behind proxies** -- an ALB behind CloudFront or a corporate gateway should usually preserve (not append to) the X-Forwarded-For chain; applications that route by hostname need the Host header preserved.

**Hardening is editable** -- desync mitigation ("defensive" by default; "monitor" to observe before enforcing "strictest"), invalid-header dropping, and the WAF fail-open/fail-closed call all update in place.

**WAF attaches here** -- `webAclArn` associates a REGIONAL-scope WAFv2 web ACL (same region, at most one per ALB) so every request is inspected before listener rules run; the fail-open toggle decides behavior while the WAF is unreachable.

**Deletion protection is recommended** -- deleting an ALB silently orphans every listener and rule attached to it. The wizard preselects protection ON.

**Capacity can be reserved for events** -- `minimumLoadBalancerCapacityUnits` pre-provisions LCUs for a dated traffic surge (the declarative replacement for pre-warming tickets). The reservation bills while set: size it from the event estimate, then remove the field to release it. `ipv4IpamPoolId` pairs with it for BYOIP addressing, keeping the ALB's public addresses inside ranges your clients have allowlisted.

## Outputs and Dependencies

### What This Component Consumes

| Field | References | Via |
|-------|-----------|-----|
| `subnets[]` | AwsSubnet | `status.outputs.subnet_id` |
| `securityGroups[]` | AwsSecurityGroup | `status.outputs.security_group_id` |
| `accessLogs.bucket` / `connectionLogs.bucket` / `healthCheckLogs.bucket` | AwsS3Bucket | `status.outputs.bucket_id` |
| `dns.route53ZoneId` | AwsRoute53Zone | `status.outputs.zone_id` |
| `webAclArn` | AwsWafWebAcl | `status.outputs.web_acl_arn` |

### What This Component Provides

After provisioning, `status.outputs` contains:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `load_balancer_arn` | ARN of the ALB | AwsLbListener resources attach through it |
| `load_balancer_name` | The ALB's name (truncated to AWS's 32-character limit) | Console URLs and CLI queries |
| `load_balancer_dns_name` | The AWS-assigned DNS name | CNAME targets and Route53 alias records |
| `load_balancer_hosted_zone_id` | Route53 hosted zone ID for the DNS name | Required when creating alias records manually |
| `arn_suffix` | The ARN's final segment | The CloudWatch LoadBalancer metric dimension — alarms, dashboards, request-count autoscaling |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Public HTTP entry point** -- an internet-facing ALB with explicit security groups; listeners add HTTPS with ACM certificates. Start from the **Internet Facing** preset.

**Internal service front** -- an internal ALB for service-to-service HTTP inside the VPC. Start from the **Internal Hardened** preset.

**Static IPs in front** -- put an AwsNlb with Elastic IPs in front and register this ALB in an `alb`-type AwsLbTargetGroup — static Layer-4 addresses, full Layer-7 routing.

## Works With

- **AwsLbListener** -- attaches to this ALB's `load_balancer_arn` output and owns ports, TLS certificates, and default actions.
- **AwsLbListenerRule** -- attaches to listeners for path/host/header routing.
- **AwsLbTargetGroup** -- receives the routed traffic; an `alb`-type group also lets an NLB front this ALB.
- **AwsSubnet / AwsSecurityGroup** -- placement and traffic filtering, referenced by `subnets` and `securityGroups`.
- **AwsS3Bucket** -- the log destinations, referenced per stream.
- **AwsRoute53Zone** -- the hosted zone for alias records, referenced by `dns.route53ZoneId`.
- **AwsWafWebAcl** -- a REGIONAL-scope web ACL inspecting requests in front of the listeners, referenced by `webAclArn`.
