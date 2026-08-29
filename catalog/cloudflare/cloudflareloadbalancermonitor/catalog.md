# Load Balancer Monitor on Cloudflare

Deploys a reusable Cloudflare Load Balancing health monitor. A monitor probes the origins inside a load balancer pool and decides whether each origin (and the pool) is healthy. Monitors are account-scoped and reusable -- many pools can reference the same monitor -- and integrate with Planton's Provider Connections for Cloudflare credential management.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Load Balancer Monitor** -- a health check of the configured protocol (HTTP, HTTPS, TCP, UDP/ICMP, ICMP ping, or SMTP) with the chosen probe configuration, timing, and health-flip thresholds
- **Cloudflare Labels** -- resource metadata applied for organization and environment tracking

## Before You Deploy

### Planton Setup

- **Cloudflare Provider Connection** -- an active connection in the Connect module with a Cloudflare API token that has Load Balancing edit access. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline API token authentication.

### Cloudflare Account

- **Load Balancing add-on** -- Cloudflare Load Balancing is a paid feature and must be enabled on your account before deploying.
- **Reachable origins** -- the origins this monitor will eventually check (via a pool) must be reachable from Cloudflare's probing regions over the chosen protocol.

## Deploy

### Console

Open the deployment store, find **Load Balancer Monitor on Cloudflare**, and click **Deploy**. The creation wizard walks you through account and protocol selection, the protocol-specific probe configuration, and the timing thresholds. Start from the **HTTPS web** preset in the [Presets](#presets) tab for an application-layer check.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: cloudflare.planton.dev/v1
kind: CloudflareLoadBalancerMonitor
metadata:
  name: web-https
  org: acme-corp
  env: prod
spec:
  accountId: "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4"
  type: https
  path: /healthz
  expectedCodes: "2xx"
  method: GET
  interval: 60
  timeout: 5
  retries: 2
  headers:
    - name: Host
      values:
        - app.example.com
```

```shell
planton apply -f cloudflare-load-balancer-monitor.yaml
```

This creates an HTTPS monitor that probes `/healthz` on each origin and marks it healthy on a `2xx` response. A Stack Job tracks the provisioning in real time.

## Key Configuration

These are the most important decisions when configuring a monitor. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Protocol (`type`)** -- Determines what a check measures. `http`/`https` validate the application response (path, status codes, body, headers) and catch "up but broken" origins. `tcp`, `udp_icmp`, and `smtp` confirm a port accepts connections. `icmp_ping` checks raw reachability only. Prefer an application-layer check whenever clients use HTTP(S).

**Expected codes** -- For HTTP(S), `expectedCodes` is the response code or range (e.g. `200`, `2xx`, `200-299`) that marks an origin healthy. The endpoint must signal health through its HTTP status, not just the body. Cloudflare requires `expectedCodes` or `expectedBody` on every http/https monitor -- creation fails with a 400 (code 1002) when both are empty.

**Port** -- Required for `tcp`, `udp_icmp`, and `smtp` monitors; optional for HTTP(S) where it defaults to 80/443. Set it to a non-standard port only when origins listen elsewhere.

**Timing and thresholds** -- `interval`, `timeout`, `retries`, `consecutiveUp`, and `consecutiveDown` trade failover speed for stability. Leaving any at 0 uses Cloudflare's defaults (60s / 5s / 2). Effective detection time is roughly interval × consecutive-down.

## Outputs and Dependencies

### What This Component Consumes

This component has no foreign-key dependencies -- a monitor is defined entirely by its own probe configuration.

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `monitor_id` | The Cloudflare-assigned identifier of the monitor | Referenced by a CloudflareLoadBalancerPool's `monitor` field |
| `monitor_type` | The health-check protocol (echoed for convenience) | Verification, dashboards |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**HTTPS web health check** -- An HTTPS monitor probing `/healthz` for a `2xx` response, with a `Host` header for virtual-hosted origins. Use for web/API pools where application-level health matters. Start from the **HTTPS web** preset.

**TCP port check** -- A TCP monitor that confirms a port accepts connections, for non-HTTP origins such as databases or message brokers. Use when only connection-level health is needed. Start from the **TCP port** preset.

## Works With

- [**Load Balancer Pool on Cloudflare**](/cloud-catalog/cloudflare-load-balancer-pool) -- references this monitor (via `monitor`) to health-check its origins
