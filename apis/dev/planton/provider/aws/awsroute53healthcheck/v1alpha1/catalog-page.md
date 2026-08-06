# AWS Route53 Health Check

Deploys an Amazon Route 53 health check — the availability probe DNS records reference to keep unhealthy endpoints out of DNS answers. Supports all eight check types: internet-facing endpoint probes (HTTP, HTTPS, string-match variants, TCP), calculated aggregation over child checks, CloudWatch alarm mirroring for private resources, and Application Recovery Controller routing-control mirroring.

## What Gets Created

When you deploy an AwsRoute53HealthCheck resource, Planton provisions:

- **Route 53 Health Check** — the health check with the chosen monitoring model, probe tuning (interval, threshold, checker regions, latency measurement, SNI), state shaping (inversion, administrative disable), and Planton resource tags (the `Name` tag is the console display name)

## Prerequisites

- **AWS credentials** configured via environment variables or Planton provider config
- **A publicly reachable endpoint** for the probing check types (the checker fleet probes from the internet)
- **A CloudWatch alarm** for `CLOUDWATCH_METRIC` checks (the private-resource pattern)

## Quick Start

Create a file `health-check.yaml`:

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsRoute53HealthCheck
metadata:
  name: app-https-check
spec:
  region: us-west-2
  checkType: HTTPS
  fqdn: app.example.com
  resourcePath: /healthz
```

Deploy:

```shell
planton apply -f health-check.yaml
```

## Configuration Reference

### Required Fields

| Field | Type | Description |
|-------|------|-------------|
| `region` | `string` | The AWS region used for provider API calls (health checks are global objects). |
| `checkType` | `string` | The monitoring model: `HTTP`, `HTTPS`, `HTTP_STR_MATCH`, `HTTPS_STR_MATCH`, `TCP`, `CALCULATED`, `CLOUDWATCH_METRIC`, or `RECOVERY_CONTROL`. Create-time immutable. |

### Endpoint Check Fields (HTTP / HTTPS / *_STR_MATCH / TCP)

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `fqdn` | `string` | — | Domain name to probe (also the Host header when `ipAddress` is set). One of `fqdn`/`ipAddress` required. |
| `ipAddress` | `string` | — | IPv4/IPv6 address to probe. |
| `port` | `int32` | 80 / 443 | Endpoint port. Required for `TCP` (no default). |
| `resourcePath` | `string` | `/` | Path to probe (HTTP-shaped types only). |
| `searchString` | `string` | — | Body substring that must appear. Required for, and only valid with, the `*_STR_MATCH` types. |
| `requestInterval` | `int32` | `30` | Seconds between probes: `10` or `30`. Create-time immutable. |
| `failureThreshold` | `int32` | `3` | Consecutive results (1–10) required to flip the health state. |
| `measureLatency` | `bool` | `false` | Latency graphs in the console. Create-time immutable. |
| `enableSni` | `bool` | AWS default | Send SNI in the TLS handshake (HTTPS types). |
| `regions` | `list(string)` | all | Checker regions to probe from (minimum 3 when set). |

### Aggregation, Mirroring, and State Shaping

| Field | Type | Description |
|-------|------|-------------|
| `childHealthChecks` | `list(string \| ref)` | Child check IDs for `CALCULATED` (max 256; references other AwsRoute53HealthCheck resources). |
| `childHealthThreshold` | `int32` | Minimum healthy children (defaults to all when omitted). |
| `cloudwatchAlarmName` / `cloudwatchAlarmRegion` | `string` | The alarm a `CLOUDWATCH_METRIC` check mirrors. |
| `insufficientDataHealthStatus` | `string` | `Healthy`, `Unhealthy`, or `LastKnownStatus` while the alarm lacks data. |
| `routingControlArn` | `string` | The ARC routing control a `RECOVERY_CONTROL` check mirrors. Create-time immutable. |
| `invertHealthcheck` | `bool` | Report the opposite of the underlying result. |
| `disabled` | `bool` | Stop probing; the check reports healthy (the maintenance-window switch). |

## Examples

### Failover pair's primary probe

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsRoute53HealthCheck
metadata:
  name: primary-region-check
spec:
  region: us-east-1
  checkType: HTTPS
  fqdn: api.example.com
  resourcePath: /healthz
  requestInterval: 10
  failureThreshold: 2
```

### String-match content check

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsRoute53HealthCheck
metadata:
  name: status-content-check
spec:
  region: us-east-1
  checkType: HTTPS_STR_MATCH
  fqdn: status.example.com
  resourcePath: /status
  searchString: '"healthy":true'
```

### Private resource via CloudWatch alarm

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsRoute53HealthCheck
metadata:
  name: internal-api-check
spec:
  region: us-west-2
  checkType: CLOUDWATCH_METRIC
  cloudwatchAlarmName: internal-api-5xx-rate
  cloudwatchAlarmRegion: us-west-2
  insufficientDataHealthStatus: LastKnownStatus
```

### Service-level calculated check

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsRoute53HealthCheck
metadata:
  name: service-health
spec:
  region: us-west-2
  checkType: CALCULATED
  childHealthChecks:
    - valueFrom:
        kind: AwsRoute53HealthCheck
        name: web-check
        fieldPath: status.outputs.health_check_id
    - valueFrom:
        kind: AwsRoute53HealthCheck
        name: api-check
        fieldPath: status.outputs.health_check_id
  childHealthThreshold: 1
```

## Stack Outputs

| Output | Type | Description |
|--------|------|-------------|
| `health_check_id` | `string` | The health check ID — referenced by DNS records (`healthCheckId`) and calculated parents. |
| `health_check_arn` | `string` | The health check ARN — for IAM policies and the CloudWatch `HealthCheckStatus` metric dimension. |

## Related Components

- [AwsRoute53DnsRecord](/docs/catalog/aws/awsroute53dnsrecord) — references the check to gate failover, weighted, and multivalue answers
- [AwsRoute53Zone](/docs/catalog/aws/awsroute53zone) — the hosted zone the gated records live in
- [AwsCloudwatchAlarm](/docs/catalog/aws/awscloudwatchalarm) — the alarm a CLOUDWATCH_METRIC check mirrors
