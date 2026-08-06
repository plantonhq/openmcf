---
title: "LB Target Group"
description: "LB Target Group deployment documentation"
icon: "package"
order: 100
componentName: "awslbtargetgroup"
---

# AWS LB Target Group

Deploys an ELBv2 target group: the routing destination that load balancer
listeners and listener rules forward traffic to. One kind serves both
Application and Network Load Balancers -- the protocol decides the family --
with health checks, stickiness, draining, and static target registration
managed together.

## What Gets Created

When you deploy an AwsLbTargetGroup resource, Planton provisions:

- **Target group** — an `aws_lb_target_group` / `lb.TargetGroup` with the name
  taken from `metadata.name` (truncated to AWS's 32-character limit when
  necessary), the specified target type, port/protocol, health check,
  stickiness, and connection-handling attributes
- **Target registrations** — one target-group attachment per entry in
  `targets`, for statically registered EC2 instances, IP addresses, a Lambda
  function, or an inner ALB

Dynamic registrars are **not** configured here — ECS services, auto-scaling
groups, and Kubernetes controllers register their own targets against this
group's `target_group_arn` output.

## Prerequisites

- **AWS credentials** configured via the Planton provider config (keyless SSO/OIDC).
- **A VPC** for `instance`, `ip`, and `alb` target types (not needed for `lambda`).
- **A load balancer** (`AwsAlb` or `AwsNlb`) whose listener will forward to this group — deployable before or after the group.

## Quick Start

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsLbTargetGroup
metadata:
  name: api
spec:
  region: us-west-2
  vpcId:
    valueFrom:
      kind: AwsVpc
      name: main-vpc
      fieldPath: status.outputs.vpc_id
  targetType: ip
  port: 8080
  protocol: HTTP
  healthCheck:
    path: /healthz
```

```shell
planton apply -f target-group.yaml
```

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
| --- | --- | --- | --- |
| `region` | `string` | AWS region for the target group. Must match the VPC and the forwarding load balancer. | Required; non-empty |
| `vpcId` | `string \| valueFrom` | VPC the targets live in. Defaults to referencing an `AwsVpc`'s `vpc_id` output. Immutable. | Required unless `targetType` is `lambda` |
| `port` | `int` | Port targets receive traffic on. Individual registrations can override it. Immutable. | Required unless `targetType` is `lambda`; 1–65535 |
| `protocol` | `string` | Protocol between the load balancer and targets; decides the family. Immutable. | Required unless `targetType` is `lambda`. One of: `HTTP`, `HTTPS` (ALB), `TCP`, `UDP`, `TCP_UDP`, `TLS` (NLB) |

### Optional Fields

| Field | Type | Default | Description |
| --- | --- | --- | --- |
| `targetType` | `string` | `instance` | How targets register: `instance`, `ip`, `lambda`, or `alb`. Immutable. |
| `protocolVersion` | `string` | `HTTP1` | ALB only: `HTTP1`, `HTTP2`, or `GRPC`. Immutable. |
| `ipAddressType` | `string` | `ipv4` | Address family of targets: `ipv4` or `ipv6` (requires a dualstack load balancer). Immutable. |
| `healthCheck` | `object` | protocol-appropriate AWS defaults | Probe protocol, port, path, thresholds, interval/timeout, and matcher. Health-check protocol is independent of the traffic protocol. |
| `stickiness` | `object` | disabled | `lb_cookie` / `app_cookie` (ALB) or `source_ip` (NLB) session affinity. |
| `deregistrationDelaySeconds` | `int` | `300` | Seconds a draining target keeps in-flight requests. Range: 0–3600. |
| `slowStartSeconds` | `int` | `0` (disabled) | ALB only: ramp-up window for newly registered targets, 30–900. Incompatible with stickiness and `least_outstanding_requests`. |
| `loadBalancingAlgorithmType` | `string` | `round_robin` | ALB only: `round_robin`, `least_outstanding_requests`, or `weighted_random`. |
| `loadBalancingAnomalyMitigation` | `string` | `off` | ALB only: automatic anomaly mitigation; requires `weighted_random`. |
| `loadBalancingCrossZoneEnabled` | `string` | `use_load_balancer_configuration` | Cross-zone override: `true`, `false`, or inherit from the load balancer. |
| `preserveClientIp` | `bool` | AWS default per target type | NLB only: targets see the original client IP. |
| `proxyProtocolV2` | `bool` | `false` | NLB only: send the Proxy Protocol v2 header. Targets must parse it — enabling against an unaware backend breaks connections. |
| `connectionTermination` | `bool` | `false` | NLB only: close established connections when the deregistration delay expires. |
| `lambdaMultiValueHeadersEnabled` | `bool` | `false` | Lambda targets only: deliver multi-value headers and query parameters as arrays. |
| `targetGroupHealth` | `object` | — | Group-level DNS failover and unhealthy-state-routing thresholds. |
| `targetHealthState` | `object` | — | NLB TCP/TLS only: established-connection handling for unhealthy targets. |
| `targets` | `object[]` | `[]` | Static registrations: `targetId` (value or `valueFrom`), optional per-target `port`, optional `availabilityZone: all` for out-of-VPC IPs. |

## Examples

### NLB TCP target group with HTTP health checks

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsLbTargetGroup
metadata:
  name: postgres-pool
spec:
  region: us-west-2
  vpcId:
    valueFrom:
      kind: AwsVpc
      name: main-vpc
      fieldPath: status.outputs.vpc_id
  targetType: ip
  port: 5432
  protocol: TCP
  preserveClientIp: true
  connectionTermination: true
  healthCheck:
    protocol: HTTP
    port: "8081"
    path: /readyz
```

### Lambda function target

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsLbTargetGroup
metadata:
  name: webhook-handler
spec:
  region: us-west-2
  targetType: lambda
  targets:
    - targetId:
        valueFrom:
          kind: AwsLambda
          name: webhook-handler
          fieldPath: status.outputs.function_arn
```

### Statically registered EC2 instances

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsLbTargetGroup
metadata:
  name: legacy-web
spec:
  region: us-west-2
  vpcId:
    valueFrom:
      kind: AwsVpc
      name: main-vpc
      fieldPath: status.outputs.vpc_id
  targetType: instance
  port: 80
  protocol: HTTP
  targets:
    - targetId:
        valueFrom:
          kind: AwsEc2Instance
          name: web-1
          fieldPath: status.outputs.instance_id
    - targetId:
        valueFrom:
          kind: AwsEc2Instance
          name: web-2
          fieldPath: status.outputs.instance_id
```

## Stack Outputs

| Output | Description |
| --- | --- |
| `target_group_arn` | ARN of the target group — what listener forward actions, ECS services, and ASG attachments reference |
| `target_group_name` | Friendly name of the group (`metadata.name`, truncated to the 32-character AWS limit when necessary) |
| `arn_suffix` | ARN suffix used as the `TargetGroup` dimension in CloudWatch metrics |

## Related Components

- [AwsLbListener](/docs/catalog/aws/lb-listener) — forwards listener traffic to this group via its default actions
- [AwsLbListenerRule](/docs/catalog/aws/lb-listener-rule) — routes matching ALB requests to this group
- [AwsAlb](/docs/catalog/aws/alb) — the Application Load Balancer in front of HTTP/HTTPS groups
- [AwsNlb](/docs/catalog/aws/nlb) — the Network Load Balancer in front of TCP/UDP/TLS groups
- [AwsEcsService](/docs/catalog/aws/ecs-service) — registers task IPs into this group as deployments roll
- [AwsVpc](/docs/catalog/aws/vpc) — provides the VPC targets live in
