# AwsRoute53HealthCheck

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `aws.planton.dev/v1`

AwsRoute53HealthCheckSpec defines an Amazon Route 53 health check — the
availability probe that DNS records reference (via health_check_id) to keep
unhealthy endpoints out of DNS answers. Health checks power failover
routing (primary/secondary), health-aware weighted routing, and multivalue
answers.

The check_type selects one of three fundamentally different monitoring
models, and the rest of the surface is gated by it:
- ENDPOINT checks (HTTP, HTTPS, HTTP_STR_MATCH, HTTPS_STR_MATCH, TCP):
  Route 53's global checker fleet probes an address you specify (fqdn
  and/or ip_address, port, path). The endpoint must be reachable from the
  public internet.
- CALCULATED: no probing of its own — aggregates the states of child
  health checks (healthy when at least child_health_threshold children are
  healthy). Composes per-endpoint checks into service-level health.
- CLOUDWATCH_METRIC: mirrors the state of a CloudWatch alarm — the way to
  health-check PRIVATE resources the checker fleet cannot reach, or to
  gate DNS on application-level metrics.
- RECOVERY_CONTROL: mirrors an Application Recovery Controller routing
  control — a manual/automated switch for disaster-recovery runbooks.

check_type, request_interval, measure_latency, and routing_control_arn are
create-time immutable (ForceNew).

## Example

```yaml
# AWS Route 53 health check — examples
#
# Usage:
#   planton apply -f manifest.yaml

apiVersion: aws.planton.dev/v1
kind: AwsRoute53HealthCheck
metadata:
  name: app-https-check
spec:
  region: us-west-2
  checkType: HTTPS
  fqdn: app.example.com
  resourcePath: /healthz
  requestInterval: 30
  failureThreshold: 3

---
# TCP check against a static address (a database listener, a bastion).

apiVersion: aws.planton.dev/v1
kind: AwsRoute53HealthCheck
metadata:
  name: db-tcp-check
spec:
  region: us-west-2
  checkType: TCP
  ipAddress: 192.0.2.10
  port: 5432
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.checkType` | `string` | yes |  |  |
| `spec.fqdn` | `string` |  |  |  |
| `spec.ipAddress` | `string` |  |  |  |
| `spec.port` | `int32` |  |  |  |
| `spec.resourcePath` | `string` |  |  |  |
| `spec.searchString` | `string` |  |  |  |
| `spec.requestInterval` | `int32` |  | `30` |  |
| `spec.failureThreshold` | `int32` |  | `3` |  |
| `spec.measureLatency` | `bool` |  |  |  |
| `spec.enableSni` | `bool` |  |  |  |
| `spec.regions` | `[]string` |  |  |  |
| `spec.invertHealthcheck` | `bool` |  |  |  |
| `spec.disabled` | `bool` |  |  |  |
| `spec.childHealthChecks` | `[]string \| valueFrom` |  |  | AwsRoute53HealthCheck (`status.outputs.health_check_id`) |
| `spec.childHealthThreshold` | `int32` |  |  |  |
| `spec.cloudwatchAlarmName` | `string` |  |  |  |
| `spec.cloudwatchAlarmRegion` | `string` |  |  |  |
| `spec.insufficientDataHealthStatus` | `string` |  |  |  |
| `spec.routingControlArn` | `string` |  |  |  |

## Field Details

### spec.region

`string` · required

The AWS region where the resource will be created.
Route 53 health checks are global objects; this selects the region used
for provider API calls.
Example: "us-west-2", "eu-west-1"

- rule: {"string":{"minLen":"1"}}

### spec.checkType

`string` · required

The monitoring model (create-time immutable, ForceNew):
- "HTTP" / "HTTPS": healthy when the endpoint answers with 2xx/3xx
  within the timeout.
- "HTTP_STR_MATCH" / "HTTPS_STR_MATCH": additionally requires
  search_string to appear in the first 5,120 bytes of the response body.
- "TCP": healthy when a TCP connection can be established.
- "CALCULATED": aggregates child_health_checks.
- "CLOUDWATCH_METRIC": mirrors a CloudWatch alarm's state.
- "RECOVERY_CONTROL": mirrors an Application Recovery Controller routing
  control state.

- rule: {"required":true,"string":{"in":["HTTP","HTTPS","HTTP_STR_MATCH","HTTPS_STR_MATCH","TCP","CALCULATED","CLOUDWATCH_METRIC","RECOVERY_CONTROL"]}}

### spec.fqdn

`string`

Domain name of the endpoint to probe. For HTTP(S) checks without
ip_address, Route 53 resolves this name and probes the result (also sent
as the Host header). When ip_address is also set, the probe goes to the
IP and this value is only the Host header. Max 255 characters.

- rule: {"string":{"maxLen":"255"}}

### spec.ipAddress

`string`

IPv4 or IPv6 address of the endpoint to probe. Use for endpoints whose
address is static; use fqdn alone when the address changes (e.g. behind
DNS-based scaling).

- rule: {"string":{"maxLen":"45"}}

### spec.port

`int32`

TCP port of the endpoint. Defaults: 80 for HTTP/HTTP_STR_MATCH, 443 for
HTTPS/HTTPS_STR_MATCH. Required for TCP checks (there is no default).

- rule: {"int32":{"lte":65535,"gte":0}}

### spec.resourcePath

`string`

Path to probe for HTTP(S) checks (e.g. "/healthz"). Defaults to "/".
Not applicable to TCP.

- rule: {"string":{"maxLen":"255"}}

### spec.searchString

`string`

String that must appear in the first 5,120 bytes of the response body.
Required for (and only valid with) HTTP_STR_MATCH / HTTPS_STR_MATCH.

- rule: {"string":{"maxLen":"255"}}

### spec.requestInterval

`int32` · optional (explicit presence)

Seconds between probes from each checker: 10 or 30 (default 30).
Create-time immutable (ForceNew). Fast (10s) checks cost more but detect
failures ~3x sooner.

- default: `30`
- rule: {"int32":{"in":[10,30]}}

### spec.failureThreshold

`int32` · optional (explicit presence)

Consecutive probe results required to flip the health state (1–10,
default 3). Lower reacts faster; higher rides out blips.

- default: `3`
- rule: {"int32":{"lte":10,"gte":1}}

### spec.measureLatency

`bool`

Measure and graph endpoint latency in the Route 53 console.
Create-time immutable (ForceNew); small extra cost.

### spec.enableSni

`bool` · optional (explicit presence)

Send SNI (the fqdn value) in the TLS handshake for HTTPS checks —
required by most name-based virtual hosting endpoints. AWS defaults this
to true for HTTPS checks when fqdn is set.

### spec.regions

`[]string`

Subset of Route 53 checker regions to probe from (minimum 3 when set).
Valid values: us-east-1, us-west-1, us-west-2, eu-west-1, ap-southeast-1,
ap-southeast-2, ap-northeast-1, sa-east-1. Default: all checker regions.
Note: once set, AWS ignores removing the list (it keeps the last value).

- rule: {"repeated":{"items":{"string":{"in":["us-east-1","us-west-1","us-west-2","eu-west-1","ap-southeast-1","ap-southeast-2","ap-northeast-1","sa-east-1"]}}}}

### spec.invertHealthcheck

`bool`

Invert the result: report unhealthy when the underlying check is healthy
and vice versa. Occasionally useful for "route AWAY while X is up"
arrangements.

### spec.disabled

`bool`

Administratively disable probing. Route 53 then treats the check as
always healthy (unless inverted) — the maintenance-window switch that
stops failover from firing while you work on the endpoint.

### spec.childHealthChecks

`[]string | valueFrom`

The child health checks this calculated check aggregates (max 256).
Can reference other AwsRoute53HealthCheck resources.

- references: AwsRoute53HealthCheck (`status.outputs.health_check_id`)
- rule: {"repeated":{"maxItems":"256"}}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsRoute53HealthCheck, name: <that resource's name>, fieldPath: status.outputs.health_check_id}} -- a bare string does not parse

### spec.childHealthThreshold

`int32`

Minimum number of healthy children for this check to report healthy
(1–256). Defaults to the number of children when omitted (all must be
healthy).

- rule: {"int32":{"lte":256,"gte":0}}

### spec.cloudwatchAlarmName

`string`

Name of the CloudWatch alarm whose state this check mirrors.

### spec.cloudwatchAlarmRegion

`string`

Region the CloudWatch alarm lives in (alarms are regional even though
the health check is global). Example: "us-west-2".

### spec.insufficientDataHealthStatus

`string`

What to report while the alarm is in INSUFFICIENT_DATA state:
"Healthy", "Unhealthy", or "LastKnownStatus".

- rule: {"string":{"in":["","Healthy","Unhealthy","LastKnownStatus"]}}

### spec.routingControlArn

`string`

ARN of the Application Recovery Controller routing control whose state
this check mirrors. Create-time immutable (ForceNew).

## Validation Rules

- `endpoint_checks_require_target`: fqdn or ip_address is required for HTTP, HTTPS, HTTP_STR_MATCH, HTTPS_STR_MATCH, and TCP checks
- `target_requires_endpoint_check`: fqdn, ip_address, port, and resource_path only apply to endpoint checks (HTTP, HTTPS, HTTP_STR_MATCH, HTTPS_STR_MATCH, TCP)
- `tcp_requires_port`: port is required for TCP checks (HTTP defaults to 80 and HTTPS to 443, but TCP has no default)
- `resource_path_http_only`: resource_path only applies to HTTP, HTTPS, HTTP_STR_MATCH, and HTTPS_STR_MATCH checks
- `search_string_for_str_match`: search_string is required for HTTP_STR_MATCH/HTTPS_STR_MATCH checks and not valid for any other type
- `calculated_requires_children`: child_health_checks (at least one) is required for CALCULATED checks and not valid for any other type
- `cloudwatch_requires_alarm`: cloudwatch_alarm_name and cloudwatch_alarm_region are required for CLOUDWATCH_METRIC checks and not valid for any other type
- `recovery_control_requires_arn`: routing_control_arn is required for RECOVERY_CONTROL checks and not valid for any other type
- `probe_tuning_endpoint_only`: regions, measure_latency, and enable_sni only apply to endpoint checks (HTTP, HTTPS, HTTP_STR_MATCH, HTTPS_STR_MATCH, TCP)
- `regions_min_three`: regions must list at least 3 checker regions when set (AWS minimum for reliable quorum)

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AwsRoute53HealthCheck, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.healthCheckId` | `string` | The health check ID (a UUID, e.g. "abcdef11-2222-3333-4444-555555fedcba"). The identifier DNS records and calculated parent checks reference. |
| `status.outputs.healthCheckArn` | `string` | The Amazon Resource Name of the health check (arn:aws:route53:::healthcheck/<id>). Used in IAM policies and as the dimension for the CloudWatch HealthCheckStatus metric. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.childHealthChecks` | AwsRoute53HealthCheck | `status.outputs.health_check_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
