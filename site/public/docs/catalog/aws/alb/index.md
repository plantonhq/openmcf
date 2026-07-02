---
title: "ALB"
description: "ALB deployment documentation"
icon: "package"
order: 100
componentName: "awsalb"
---

# AWS ALB

Deploys an AWS Application Load Balancer: the Layer-7 entry point that owns
placement (two or more subnets across Availability Zones), security groups,
HTTP behavior attributes, S3 log delivery, and optional Route53 alias DNS.
Routing lives in separate components -- listeners, rules, and target groups
attach to this ALB by ARN.

## What Gets Created

When you deploy an AwsAlb resource, Planton provisions:

- **Application Load Balancer** — an `aws_lb` / `lb.LoadBalancer` of type
  `application`, placed in the specified subnets with the attached security
  groups and the configured load balancer attributes (timeouts, HTTP/2,
  desync mitigation, XFF handling, WAF fail-open, zonal shift)
- **S3 log delivery** — access, connection, and health-check log streams,
  each enabled by the presence of its block in the spec
- **Route53 A records** — created only when DNS is enabled, one alias record
  per hostname pointing to the ALB's DNS name

Listeners, listener rules, and target groups are **not** created here —
attach `AwsLbListener` resources to this ALB's `load_balancer_arn` output,
`AwsLbListenerRule` resources to those listeners, and `AwsLbTargetGroup`
resources as destinations.

## Prerequisites

- **AWS credentials** configured via the Planton provider config (keyless SSO/OIDC).
- **At least two subnets** in different Availability Zones (public subnets
  for internet-facing, private for internal).
- **A security group** allowing inbound traffic on the ports your listeners
  will use.
- **An S3 bucket with the ELB log-delivery bucket policy** if enabling any
  log stream — delivery fails silently otherwise.
- **A Route53 hosted zone** if enabling DNS management.

## Quick Start

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsAlb
metadata:
  name: main-alb
spec:
  region: us-west-2
  subnets:
    - value: subnet-0a1b2c3d4e5f00001
    - value: subnet-0a1b2c3d4e5f00002
  securityGroups:
    - value: sg-0a1b2c3d4e5f00001
```

```shell
planton apply -f alb.yaml
```

This creates an internet-facing ALB across two subnets. Add an
`AwsLbListener` against its `load_balancer_arn` output to start accepting
traffic.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
| --- | --- | --- | --- |
| `region` | `string` | AWS region where the ALB is created (e.g., `us-west-2`). | Required; non-empty |
| `subnets` | `(string \| valueFrom)[]` | Subnets for the ALB's nodes — public for internet-facing, private for internal. References `AwsSubnet.subnet_id` by default. | Required; minimum 2 items, in different Availability Zones |

### Optional Fields

| Field | Type | Default | Description |
| --- | --- | --- | --- |
| `securityGroups` | `(string \| valueFrom)[]` | VPC default SG | Security groups controlling traffic to and from the ALB. The VPC default is fine for a first boot, wrong for production. References `AwsSecurityGroup.security_group_id` by default. |
| `internal` | `bool` | `false` | When `true`, the ALB is reachable only inside the VPC. Immutable — changing it replaces the load balancer. |
| `ipAddressType` | `string` | `ipv4` | `ipv4`, `dualstack`, or `dualstack-without-public-ipv4` (public IPv6 with private IPv4 — avoids public-IPv4 charges). |
| `deleteProtectionEnabled` | `bool` | `false` | Prevents deletion while enabled. Recommended for production: deleting an ALB silently orphans every listener and rule attached to it. |
| `idleTimeoutSeconds` | `int` | `60` (AWS) | Seconds an idle connection stays open, 1–4000. Raise above the application's slowest response to avoid 504s. |
| `clientKeepAliveSeconds` | `int` | `3600` (AWS) | Seconds an HTTP client connection may stay alive across requests, 60–604800. |
| `http2Enabled` | `bool` | `true` (AWS) | Whether HTTP/2 is offered to clients. Optional so that an explicit `false` is distinguishable from "keep the AWS default". |
| `wafFailOpenEnabled` | `bool` | `false` | When an attached WAF is unreachable: `true` passes requests through, `false` rejects them. An availability-versus-security call. |
| `zonalShiftEnabled` | `bool` | `false` | Allows Amazon Application Recovery Controller to shift traffic away from an impaired Availability Zone. |
| `dropInvalidHeaderFields` | `bool` | `false` | Drop request headers with invalid names instead of forwarding them — hardens against header smuggling. |
| `preserveHostHeader` | `bool` | `false` | Forward the client's original Host header unchanged instead of rewriting it to the target address. |
| `xffClientPortEnabled` | `bool` | `false` | Append the client's source port to `X-Forwarded-For`. |
| `xffHeaderProcessingMode` | `string` | `append` | `append`, `preserve`, or `remove` — `preserve`/`remove` matter when the ALB sits behind another proxy whose XFF chain must win. |
| `desyncMitigationMode` | `string` | `defensive` | Protection against HTTP desync attacks: `monitor` (classify only), `defensive`, or `strictest` (block everything not RFC 7230 compliant). |
| `tlsVersionAndCipherSuiteHeadersEnabled` | `bool` | `false` | Inject `x-amzn-tls-version` / `x-amzn-tls-cipher-suite` headers toward targets. |
| `accessLogs` | `object` | off | Per-request logs to S3: `bucket` (references `AwsS3Bucket.bucket_id` by default) and optional `prefix`. |
| `connectionLogs` | `object` | off | Per-connection logs to S3 (TLS handshake details) — where negotiation failures that never become requests show up. Same shape as `accessLogs`. |
| `healthCheckLogs` | `object` | off | Per-probe health-check logs to S3, for debugging flapping targets without packet captures. Same shape as `accessLogs`. |
| `dns` | `object` | off | Route53 alias DNS: `enabled`, `route53ZoneId` (references `AwsRoute53Zone.zone_id` by default), `hostnames` (unique). |

## Examples

### Internal ALB

An ALB reachable only inside the VPC, for internal microservice routing:

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsAlb
metadata:
  name: internal-api-alb
spec:
  region: us-west-2
  subnets:
    - value: subnet-private-az1
    - value: subnet-private-az2
  securityGroups:
    - value: sg-internal-alb
  internal: true
```

### Production ALB with hardening and access logs

Deletion protection, strict desync mitigation, invalid-header dropping, and
per-request logging:

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsAlb
metadata:
  name: prod-alb
spec:
  region: us-west-2
  subnets:
    - value: subnet-public-az1
    - value: subnet-public-az2
  securityGroups:
    - value: sg-prod-alb
  deleteProtectionEnabled: true
  idleTimeoutSeconds: 120
  desyncMitigationMode: strictest
  dropInvalidHeaderFields: true
  accessLogs:
    bucket:
      value: prod-alb-logs-bucket
    prefix: alb/prod
```

### ALB with Route53 alias DNS and references

Everything by reference — subnets, security group, and hosted zone resolve
from other components' outputs at deploy time:

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsAlb
metadata:
  name: web-alb
spec:
  region: us-west-2
  subnets:
    - valueFrom:
        kind: AwsSubnet
        name: public-az1
        fieldPath: status.outputs.subnet_id
    - valueFrom:
        kind: AwsSubnet
        name: public-az2
        fieldPath: status.outputs.subnet_id
  securityGroups:
    - valueFrom:
        kind: AwsSecurityGroup
        name: alb-sg
        fieldPath: status.outputs.security_group_id
  dns:
    enabled: true
    route53ZoneId:
      valueFrom:
        kind: AwsRoute53Zone
        name: example-com
        fieldPath: status.outputs.zone_id
    hostnames:
      - app.example.com
      - api.example.com
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
| --- | --- | --- |
| `load_balancer_arn` | `string` | ARN of the ALB — what `AwsLbListener` resources attach through |
| `load_balancer_name` | `string` | Name assigned to the ALB (`metadata.name`, truncated to AWS's 32-character limit when necessary) |
| `load_balancer_dns_name` | `string` | DNS name automatically assigned by AWS (e.g., `main-alb-123456.us-west-2.elb.amazonaws.com`) |
| `load_balancer_hosted_zone_id` | `string` | Route53 hosted zone ID for the ALB's DNS entry, used for alias records |

## Related Components

- [AwsLbListener](/docs/catalog/aws/lb-listener) — attaches to this ALB by ARN; owns ports, TLS material, and default actions
- [AwsLbListenerRule](/docs/catalog/aws/lb-listener-rule) — per-service routing attached to a listener
- [AwsLbTargetGroup](/docs/catalog/aws/lb-target-group) — the destination of forward actions
- [AwsSubnet](/docs/catalog/aws/subnet) — provides the subnets for ALB placement
- [AwsSecurityGroup](/docs/catalog/aws/security-group) — controls inbound and outbound traffic to the ALB
- [AwsS3Bucket](/docs/catalog/aws/s3-bucket) — receives access, connection, and health-check logs
- [AwsRoute53Zone](/docs/catalog/aws/route53-zone) — hosts the DNS zone for alias records
