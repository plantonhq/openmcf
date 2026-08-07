# OpenStackLoadBalancerMonitor

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `openstack.planton.dev/v1alpha1`

OpenStackLoadBalancerMonitorSpec defines the configuration for an Octavia health monitor
in OpenStack. A health monitor periodically checks the health of pool members and
removes unhealthy members from the pool's rotation until they recover.

Terraform resource: openstack_lb_monitor_v2
Pulumi resource:    loadbalancer.Monitor

Validations:
- url_path, http_method, and expected_codes are only valid for HTTP or HTTPS monitors.
- max_retries must be between 1 and 10.
- max_retries_down must be between 1 and 10 (when set).

Note: The Terraform provider does NOT support tags on health monitors.

## Example

```yaml
apiVersion: openstack.planton.dev/v1alpha1
kind: OpenStackLoadBalancerMonitor
metadata:
  name: test-monitor
spec:
  pool_id:
    value: "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
  type: "HTTP"
  delay: 5
  timeout: 10
  max_retries: 3
  url_path: "/healthz"
  http_method: "GET"
  expected_codes: "200"
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.poolId` | `string \| valueFrom` | yes |  | OpenStackLoadBalancerPool (`status.outputs.pool_id`) |
| `spec.type` | `string` | yes |  |  |
| `spec.delay` | `int32` | yes |  |  |
| `spec.timeout` | `int32` | yes |  |  |
| `spec.maxRetries` | `int32` | yes |  |  |
| `spec.maxRetriesDown` | `int32` |  |  |  |
| `spec.urlPath` | `string` |  |  |  |
| `spec.httpMethod` | `string` |  |  |  |
| `spec.expectedCodes` | `string` |  |  |  |
| `spec.adminStateUp` | `bool` |  | `true` |  |
| `spec.region` | `string` |  |  |  |

## Field Details

### spec.poolId

`string | valueFrom` · required

(Required) The pool to monitor.
ForceNew: changing this requires recreating the monitor.

FK: OpenStackLoadBalancerPool.status.outputs.pool_id

- references: OpenStackLoadBalancerPool (`status.outputs.pool_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: OpenStackLoadBalancerPool, name: <that resource's name>, fieldPath: status.outputs.pool_id}} -- a bare string does not parse

### spec.type

`string` · required

(Required) The type of health check to perform.
ForceNew: changing this requires recreating the monitor.

- HTTP: Send HTTP requests and check response codes
- HTTPS: Send HTTPS requests and check response codes
- PING: ICMP ping the member
- TCP: Attempt a TCP connection
- TLS-HELLO: Perform a TLS handshake
- UDP-CONNECT: Send a UDP datagram and check for ICMP errors

- rule: {"required":true,"string":{"in":["HTTP","HTTPS","PING","TCP","TLS-HELLO","UDP-CONNECT"]}}

### spec.delay

`int32` · required

(Required) The interval in seconds between health checks.
A check is sent every `delay` seconds to each member.

- rule: {"required":true}

### spec.timeout

`int32` · required

(Required) The maximum time in seconds to wait for a health check response.
If a member does not respond within this time, the check is considered failed.

- rule: {"required":true}

### spec.maxRetries

`int32` · required

(Required) The number of consecutive successful health checks required before
a member is considered healthy (brought back into rotation).
Must be between 1 and 10.

- rule: max_retries must be between 1 and 10
- rule: {"required":true}

### spec.maxRetriesDown

`int32` · optional (explicit presence)

(Optional) The number of consecutive failed health checks required before
a member is considered unhealthy (removed from rotation).
Must be between 1 and 10. Default: Octavia uses the same value as max_retries.

- rule: max_retries_down must be between 1 and 10

### spec.urlPath

`string`

(Optional) The URL path to request for HTTP/HTTPS health checks.
Typically "/" or "/healthz". Only applicable when type is HTTP or HTTPS.

### spec.httpMethod

`string`

(Optional) The HTTP method to use for HTTP/HTTPS health checks.
Only applicable when type is HTTP or HTTPS.

- rule: {"string":{"in":["","GET","HEAD","POST","PUT","DELETE","PATCH","OPTIONS","CONNECT","TRACE"]}}

### spec.expectedCodes

`string`

(Optional) Expected HTTP response codes for a healthy member.
Supports single codes ("200"), ranges ("200-299"), and comma-separated lists ("200,202").
Only applicable when type is HTTP or HTTPS.

### spec.adminStateUp

`bool` · optional (explicit presence)

(Optional) Administrative state of the monitor.
When false, the monitor stops checking members. Default: true.

- default: `true`

### spec.region

`string`

(Optional) Override the region from the provider configuration.

## Validation Rules

- `http_fields_require_http_type`: url_path, http_method, and expected_codes are only valid for HTTP or HTTPS monitor types

## Outputs

Reference an output from another manifest as `valueFrom: {kind: OpenStackLoadBalancerMonitor, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.monitor_id` | `string` | monitor_id is the unique identifier (UUID) of the health monitor. |
| `status.outputs.name` | `string` | name is the name of the monitor (derived from metadata.name). |
| `status.outputs.type` | `string` | type is the health check type (HTTP, HTTPS, PING, TCP, TLS-HELLO, UDP-CONNECT). |
| `status.outputs.pool_id` | `string` | pool_id is the ID of the monitored pool. |
| `status.outputs.region` | `string` | region is the OpenStack region where the monitor was created. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.poolId` | OpenStackLoadBalancerPool | `status.outputs.pool_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
