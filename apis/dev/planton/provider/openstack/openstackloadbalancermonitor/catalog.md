# OpenStack Load Balancer Monitor

Deploys an Octavia health monitor on OpenStack that periodically checks the health of pool members and removes unhealthy members from rotation until they recover. The monitor supports HTTP, HTTPS, TCP, PING, TLS-HELLO, and UDP-CONNECT check types with configurable intervals, timeouts, and retry thresholds. ValueFromRef wiring connects the monitor to an OpenStackLoadBalancerPool Cloud Resource in InfraChart deployments.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Octavia Health Monitor** -- a monitor resource attached to the specified pool, periodically probing each member according to the configured check type, interval, and retry thresholds
- **HTTP Check Configuration** -- created only when type is `HTTP` or `HTTPS`; configures the URL path, HTTP method, and expected response codes for application-level health validation

## Before You Deploy

### OpenStack Account

- **Pool** -- an Octavia pool must exist to attach the monitor to. Provide the pool ID directly or reference an OpenStackLoadBalancerPool Cloud Resource via ValueFromRef.
- **Health endpoint** (for HTTP/HTTPS monitors) -- backend members should expose a health check endpoint (e.g., `/healthz`) that returns an expected HTTP status code when healthy.

## Deploy

### Console

Open the deployment store, find **OpenStack Load Balancer Monitor**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **HTTP Health Check Monitor** preset in the [Presets](#presets) tab to pre-populate a working configuration.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: openstack.planton.dev/v1
kind: OpenStackLoadBalancerMonitor
metadata:
  name: http-health
  org: acme-corp
  env: prod
spec:
  poolId:
    value: "<pool-id>"
  type: HTTP
  delay: 10
  timeout: 5
  maxRetries: 3
  urlPath: /healthz
  httpMethod: GET
  expectedCodes: "200"
```

```shell
planton apply -f monitor.yaml
```

This creates an HTTP health monitor that sends a GET request to `/healthz` every 10 seconds, expects a 200 response within 5 seconds, and requires 3 consecutive successes to mark a member healthy.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the monitor to a pool deployed in the same InfraPipeline:

```yaml
spec:
  poolId:
    valueFrom:
      kind: OpenStackLoadBalancerPool
      name: web-pool
      fieldPath: status.outputs.pool_id
```

The InfraPipeline resolves the dependency graph, deploys the pool first, then provisions the monitor with the resolved pool ID.

## Key Configuration

These are the most important decisions when configuring a health monitor. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Check type** -- Choose `HTTP` or `HTTPS` for application-level health validation with URL path and response code checking. Choose `TCP` for port-reachability checks on non-HTTP services. Choose `PING` for ICMP-based liveness checks. Changing the type requires recreating the monitor.

**Timing parameters** -- `delay` controls the interval between checks, `timeout` sets the maximum wait for a response. Set `timeout` lower than `delay` to avoid overlapping checks. Typical production values: 10-second delay, 5-second timeout.

**Retry thresholds** -- `maxRetries` controls how many consecutive successes are needed to mark a member healthy. `maxRetriesDown` controls failures needed to mark unhealthy (defaults to `maxRetries` when unset). Lower values detect failures faster but increase false positives.

**HTTP configuration** -- For HTTP/HTTPS monitors, set `urlPath` to your application's health endpoint, `httpMethod` to the request method (typically GET), and `expectedCodes` to the acceptable response codes (e.g., `"200"` or `"200-299"`).

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **OpenStackLoadBalancerPool** | `poolId` | `status.outputs.pool_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `monitor_id` | UUID of the health monitor | Resource identification, API operations |
| `name` | Name of the monitor | Monitoring labels |
| `type` | Health check type | Diagnostics, configuration auditing |
| `pool_id` | ID of the monitored pool | Cross-referencing pool relationships |
| `region` | OpenStack region where the monitor was created | Region-aware downstream resources |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**HTTP Health Check** -- Sends GET requests to `/healthz` every 10 seconds with a 5-second timeout. Validates application health by checking HTTP response codes. Best for web services and REST APIs. Start from the **HTTP Health Check Monitor** preset.

**TCP Health Check** -- Attempts a TCP connection every 10 seconds with a 5-second timeout. Validates port reachability without application-level protocol checks. Best for databases, message queues, and non-HTTP services. Start from the **TCP Health Check Monitor** preset.

## Works With

- [**OpenStack Load Balancer Pool**](/cloud-catalog/openstack-load-balancer-pool) -- provides the pool ID that this monitor checks members for