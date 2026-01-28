# Cloudflare Load Balancer

A OpenMCF deployment component for deploying and managing Cloudflare Load Balancers - a Global Server Load Balancing (GSLB) solution that provides intelligent traffic steering, health monitoring, and automatic failover across geographically distributed origins.

## Overview

Cloudflare Load Balancer operates at the DNS level and is tightly integrated with Cloudflare's global network of 330+ data centers. Unlike static DNS round-robin, it continuously monitors origin health and dynamically routes traffic away from failures and toward healthy endpoints - providing high availability without managing hardware or installing software.

### Key Features

- **Global Traffic Management**: Distribute traffic across geographically dispersed data centers
- **Active-Passive Failover**: Automatic traffic rerouting when health checks fail (failover in seconds)
- **Geo-Steering**: Route users to the nearest origin for reduced latency
- **Health Monitoring**: Continuous health checks from multiple global locations
- **Session Affinity**: Sticky sessions for stateful applications
- **Traffic Steering Policies**: Support for failover, geo-routing, and weighted distribution
- **Proxied Mode**: Layer 7 routing with WAF, CDN caching, and advanced features
- **DNS-Only Mode**: Direct DNS responses for non-HTTP traffic

### Use Cases

Use Cloudflare Load Balancer when you need:

- **Multi-region high availability**: Failover between primary and backup data centers
- **Global performance optimization**: Route users to the closest origin server
- **Multi-cloud abstraction**: Balance traffic across AWS, GCP, Azure, and on-premises
- **Blue-green deployments**: Gradually shift traffic to new application versions
- **Disaster recovery**: Automatic failover to DR sites with configurable health checks

## Quick Start

### Basic Load Balancer with Two Origins

```yaml
apiVersion: cloudflare.openmcf.org/v1
kind: CloudflareLoadBalancer
metadata:
  name: api-lb
spec:
  hostname: api.example.com
  zone_id:
    value: "abc123def456"  # Your Cloudflare zone ID
  origins:
    - name: primary
      address: 203.0.113.10
      weight: 1
    - name: secondary  
      address: 198.51.100.20
      weight: 1
  proxied: true
  health_probe_path: /health
  session_affinity: SESSION_AFFINITY_COOKIE
  steering_policy: STEERING_OFF  # Active-passive failover
```

This configuration creates:
1. A health monitor that checks `/health` on each origin
2. A pool containing both origins
3. A load balancer at `api.example.com` with failover enabled
4. Sticky sessions via cookies

Traffic flows to the primary origin. If it fails health checks, traffic automatically moves to the secondary.

## Specification Fields

### Required Fields

- **`hostname`** (string): The DNS hostname for the load balancer (e.g., `api.example.com`)
  - Must be a valid DNS name within the specified zone
  
- **`zone_id`** (StringValueOrRef): Reference to the Cloudflare DNS zone
  - Can be a direct value or a reference to a CloudflareDnsZone resource
  - Example: `zone_id: { value: "abc123" }` or `zone_id: { ref: "my-zone" }`

- **`origins`** (list): List of origin servers (minimum 1 required)
  - **`name`** (string): Unique identifier for the origin
  - **`address`** (string): IP address or DNS hostname of the origin
  - **`weight`** (int32): Traffic weight for weighted distribution (default: 1)

### Optional Fields

- **`proxied`** (bool): Enable Cloudflare proxy (orange cloud) - Default: `true`
  - `true`: Traffic flows through Cloudflare's Layer 7 proxy (enables WAF, caching, session affinity)
  - `false`: DNS-only mode - Cloudflare returns origin IP directly

- **`health_probe_path`** (string): HTTP path for health checks - Default: `"/"`
  - Should point to an endpoint that returns HTTP 200 when healthy
  - Examples: `/health`, `/healthz`, `/api/status`

- **`session_affinity`** (enum): Session stickiness setting - Default: `SESSION_AFFINITY_NONE`
  - `SESSION_AFFINITY_NONE` (0): No sticky sessions
  - `SESSION_AFFINITY_COOKIE` (1): Cookie-based sticky sessions (requires proxied mode)

- **`steering_policy`** (enum): Traffic distribution policy - Default: `STEERING_OFF`
  - `STEERING_OFF` (0): Active-passive failover (uses priority order)
  - `STEERING_GEO` (1): Geographic routing (route to nearest origin)
  - `STEERING_RANDOM` (2): Weighted distribution (for A/B testing)

## Common Patterns

### Active-Passive Failover

Route all traffic to primary origin, fail over to secondary if primary is unhealthy:

```yaml
spec:
  hostname: app.example.com
  steering_policy: STEERING_OFF  # Priority-based failover
  origins:
    - name: primary-us-east
      address: 203.0.113.10
    - name: backup-us-west
      address: 198.51.100.20
```

### Geographic Routing

Route users to the nearest origin for optimal latency:

```yaml
spec:
  hostname: app.example.com
  steering_policy: STEERING_GEO
  origins:
    - name: us-origin
      address: 203.0.113.10
    - name: eu-origin
      address: 198.51.100.20
    - name: asia-origin
      address: 192.0.2.30
```

### Weighted A/B Testing

Send 80% of traffic to the stable version, 20% to the new version:

```yaml
spec:
  hostname: app.example.com
  steering_policy: STEERING_RANDOM
  origins:
    - name: stable-v1
      address: 203.0.113.10
      weight: 80
    - name: canary-v2
      address: 198.51.100.20
      weight: 20
```

### Multi-Cloud Load Balancing

Balance traffic across AWS, GCP, and on-premises:

```yaml
spec:
  hostname: api.example.com
  origins:
    - name: aws-us-east-1
      address: 203.0.113.10
    - name: gcp-us-central1
      address: 198.51.100.20
    - name: on-prem-datacenter
      address: 192.0.2.30
```

## Health Monitoring

Cloudflare probes origins from three separate data centers per region. An origin is marked healthy only if the majority of probes succeed. Health check configuration:

- **Type**: HTTPS (or HTTP for non-TLS origins)
- **Path**: Configured via `health_probe_path` (default: `/`)
- **Expected Response**: HTTP 2xx status codes
- **Interval**: 60 seconds (Pro plan minimum)
- **Timeout**: 5 seconds
- **Retries**: 2 attempts before marking unhealthy

**Important**: Ensure origin firewalls allow Cloudflare health check IP ranges.

## Session Affinity

When `session_affinity: SESSION_AFFINITY_COOKIE` is set:

1. Cloudflare generates a `__cflb` cookie on the first request
2. Subsequent requests from the same client route to the same origin
3. Ensures stateful applications (e.g., shopping carts) work correctly

**Note**: Session affinity requires `proxied: true` (orange cloud mode).

## Failover Behavior

With `steering_policy: STEERING_OFF` (active-passive failover):

1. All traffic goes to the first origin in the list
2. If that origin fails health checks, traffic moves to the next origin
3. When the primary recovers, traffic automatically "fails back"
4. Failover time depends on health check interval (60s for Pro plan)

## Cost Considerations

Cloudflare Load Balancing is a paid add-on service:

- **Base Fee**: $5/month (includes 2 origins)
- **Additional Origins**: $5/month per origin
- **Geo-Routing**: $10/month (flat fee)
- **DNS Queries**: First 500,000/month free, then $0.50 per 500K

**Health check interval** is gated by plan tier:
- Pro Plan ($20/mo): 60-second minimum (60s RTO)
- Business Plan ($200/mo): 15-second minimum (15s RTO)
- Enterprise Plan: 10-second minimum (10s RTO)

## Integration with Other Components

### With CloudflareDnsZone

Reference a CloudflareDnsZone resource for the zone_id:

```yaml
---
apiVersion: cloudflare.openmcf.org/v1
kind: CloudflareDnsZone
metadata:
  name: example-zone
spec:
  domain_name: example.com
---
apiVersion: cloudflare.openmcf.org/v1
kind: CloudflareLoadBalancer
metadata:
  name: api-lb
spec:
  hostname: api.example.com
  zone_id:
    ref: example-zone  # References the CloudflareDnsZone above
  origins:
    - name: origin-1
      address: 203.0.113.10
```

### With Cloudflare Tunnel

Use Cloudflare Tunnel to load balance private origins without public IPs:

```yaml
spec:
  hostname: internal-api.example.com
  origins:
    - name: private-origin
      address: <tunnel-uuid>.cfargotunnel.com  # Cloudflare Tunnel hostname
```

## Best Practices

1. **Always define at least 2 origins** for true high availability
2. **Use health check paths that accurately reflect service health** (not just `/` which might always return 200)
3. **Separate environments by account** (dev and prod in different Cloudflare accounts)
4. **Enable session affinity** for stateful applications
5. **Use proxied mode** (default) to benefit from WAF and CDN features
6. **Configure fallback pool** (automatically set to the last origin in OpenMCF)
7. **Monitor health check results** in Cloudflare dashboard

## Anti-Patterns to Avoid

❌ **Single origin**: Provides no redundancy  
❌ **Missing health checks**: Can route to failed origins  
❌ **Firewall blocking health checks**: Causes false negatives  
❌ **Creating load balancer before pool is healthy**: Can cause immediate outage  
❌ **Sharing pools across dev/prod**: Environment isolation violation  

## Deployment

OpenMCF handles the complexity of creating and linking:
1. A Cloudflare Monitor (health check) - account-level resource
2. A Cloudflare Pool (origin grouping) - account-level resource  
3. A Cloudflare Load Balancer (hostname mapping) - zone-level resource

You define a simple, denormalized spec - OpenMCF manages the resource dependency graph.

## Examples

See [examples.md](./examples.md) for complete, working examples including:
- Basic active-passive failover
- Geographic routing
- Weighted A/B testing
- Multi-cloud load balancing
- Session affinity configuration

## Documentation

- **[Research Documentation](./docs/README.md)**: Deep dive into deployment methods, architectural decisions, and 80/20 scoping
- **[Pulumi Module](./iac/pulumi/README.md)**: Pulumi deployment guide
- **[Terraform Module](./iac/tf/README.md)**: Terraform deployment guide

## Support

- Cloudflare Load Balancer: [Official Docs](https://developers.cloudflare.com/load-balancing/)
- Health Checks: [Cloudflare Health Check Documentation](https://developers.cloudflare.com/load-balancing/monitors/)
- Traffic Steering: [Cloudflare Steering Policies](https://developers.cloudflare.com/load-balancing/understand-basics/traffic-steering/)

## What's Next?

- Explore [examples.md](./examples.md) for complete usage patterns
- Read [docs/README.md](./docs/README.md) for architectural deep dive
- Deploy using [iac/pulumi/](./iac/pulumi/) or [iac/terraform/](./iac/tf/)

