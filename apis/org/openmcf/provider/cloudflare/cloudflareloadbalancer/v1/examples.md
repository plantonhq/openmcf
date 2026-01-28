# Cloudflare Load Balancer Examples

Complete, working examples for common Cloudflare Load Balancer patterns using OpenMCF.

## Table of Contents

- [Example 1: Basic Active-Passive Failover](#example-1-basic-active-passive-failover)
- [Example 2: Geographic Routing (Multi-Region)](#example-2-geographic-routing-multi-region)
- [Example 3: Weighted A/B Testing](#example-3-weighted-ab-testing)
- [Example 4: Multi-Cloud Load Balancing](#example-4-multi-cloud-load-balancing)
- [Example 5: Session Affinity for Stateful Apps](#example-5-session-affinity-for-stateful-apps)
- [Example 6: DNS-Only Mode (Non-Proxied)](#example-6-dns-only-mode-non-proxied)
- [Example 7: Custom Health Check Path](#example-7-custom-health-check-path)
- [Example 8: Three-Tier Failover](#example-8-three-tier-failover)

---

## Example 1: Basic Active-Passive Failover

**Use Case**: Route all traffic to a primary origin. If it fails, automatically fail over to a backup origin.

**Features**:
- Single pool with two origins
- Health checks every 60 seconds
- Automatic failover and failback
- Session affinity disabled

```yaml
apiVersion: cloudflare.openmcf.org/v1
kind: CloudflareLoadBalancer
metadata:
  name: api-failover-lb
spec:
  hostname: api.example.com
  zoneId:
    value: "abc123def456"  # Replace with your Cloudflare zone ID
  
  # Active-passive failover using priority order
  steeringPolicy: STEERING_OFF
  
  # Origins are tried in order - first healthy origin receives traffic
  origins:
    - name: primary-us-east-1
      address: 203.0.113.10
      weight: 1
    
    - name: backup-us-west-2
      address: 198.51.100.20
      weight: 1
  
  # Orange cloud mode - traffic flows through Cloudflare proxy
  proxied: true
  
  # Health check at root path
  healthProbePath: /
  
  # No sticky sessions
  sessionAffinity: SESSION_AFFINITY_NONE
```

**Expected Behavior**:
1. All traffic goes to `primary-us-east-1`
2. If primary fails 3 consecutive health checks (≈3 minutes), traffic moves to backup
3. When primary recovers, traffic automatically returns to primary

---

## Example 2: Geographic Routing (Multi-Region)

**Use Case**: Reduce latency by routing users to the origin closest to them.

**Features**:
- Three origins in different regions
- Geographic steering policy
- Cloudflare automatically routes based on user location

```yaml
apiVersion: cloudflare.openmcf.org/v1
kind: CloudflareLoadBalancer
metadata:
  name: global-app-lb
spec:
  hostname: app.example.com
  zoneId:
    value: "abc123def456"
  
  # Geographic steering - route to nearest origin
  steeringPolicy: STEERING_GEO
  
  origins:
    # North America origin
    - name: us-east-origin
      address: 203.0.113.10
      weight: 1
    
    # Europe origin
    - name: eu-west-origin
      address: 198.51.100.20
      weight: 1
    
    # Asia-Pacific origin
    - name: ap-south-origin
      address: 192.0.2.30
      weight: 1
  
  proxied: true
  healthProbePath: /health
  sessionAffinity: SESSION_AFFINITY_COOKIE
```

**Expected Behavior**:
- Users in North America → `us-east-origin`
- Users in Europe → `eu-west-origin`
- Users in Asia → `ap-south-origin`
- If the nearest origin is unhealthy, route to next nearest healthy origin

**Note**: Geographic steering requires Cloudflare Business plan or higher ($10/month add-on).

---

## Example 3: Weighted A/B Testing

**Use Case**: Gradually roll out a new application version by sending a percentage of traffic to it.

**Features**:
- Weighted traffic distribution
- 90% to stable version, 10% to canary
- Use for blue-green deployments or feature testing

```yaml
apiVersion: cloudflare.openmcf.org/v1
kind: CloudflareLoadBalancer
metadata:
  name: canary-deployment-lb
spec:
  hostname: app.example.com
  zoneId:
    value: "abc123def456"
  
  # Random steering with weights
  steeringPolicy: STEERING_RANDOM
  
  origins:
    # Stable version receives 90% of traffic
    - name: stable-v1
      address: 203.0.113.10
      weight: 90
    
    # Canary version receives 10% of traffic
    - name: canary-v2
      address: 198.51.100.20
      weight: 10
  
  proxied: true
  healthProbePath: /healthz
  sessionAffinity: SESSION_AFFINITY_COOKIE  # Ensure users stick to their version
```

**Expected Behavior**:
- 90% of requests go to stable-v1
- 10% of requests go to canary-v2
- Session affinity ensures the same user always hits the same version

**Rollout Strategy**:
1. Start with weight: 95 (stable) / 5 (canary)
2. Monitor error rates and performance
3. Gradually increase canary weight: 90/10, 80/20, 50/50
4. Once validated, route 100% to new version

---

## Example 4: Multi-Cloud Load Balancing

**Use Case**: Balance traffic across AWS, GCP, Azure, and on-premises data centers.

**Features**:
- Vendor-neutral load balancing
- Automatic failover across clouds
- Avoids cloud provider lock-in

```yaml
apiVersion: cloudflare.openmcf.org/v1
kind: CloudflareLoadBalancer
metadata:
  name: multi-cloud-lb
spec:
  hostname: api.example.com
  zoneId:
    value: "abc123def456"
  
  steeringPolicy: STEERING_OFF  # Priority-based failover
  
  origins:
    # Primary: AWS us-east-1
    - name: aws-us-east-1
      address: ec2-203-0-113-10.compute-1.amazonaws.com
      weight: 1
    
    # Secondary: GCP us-central1
    - name: gcp-us-central1
      address: 198.51.100.20
      weight: 1
    
    # Tertiary: Azure eastus
    - name: azure-eastus
      address: 192.0.2.30
      weight: 1
    
    # Last resort: On-premises
    - name: on-prem-datacenter
      address: 10.0.0.50
      weight: 1
  
  proxied: true
  healthProbePath: /health
  sessionAffinity: SESSION_AFFINITY_NONE
```

**Expected Behavior**:
1. All traffic goes to AWS (highest priority)
2. If AWS fails, traffic moves to GCP
3. If GCP fails, traffic moves to Azure
4. On-premises is the last resort

**Multi-Cloud Benefits**:
- Resilience against cloud provider outages
- Cost optimization (move traffic to cheaper regions)
- Data sovereignty compliance (keep EU traffic in EU)

---

## Example 5: Session Affinity for Stateful Apps

**Use Case**: E-commerce application with shopping cart state stored in-memory on origins.

**Features**:
- Cookie-based session stickiness
- Users always route to the same origin
- Prevents lost cart data

```yaml
apiVersion: cloudflare.openmcf.org/v1
kind: CloudflareLoadBalancer
metadata:
  name: ecommerce-lb
spec:
  hostname: shop.example.com
  zoneId:
    value: "abc123def456"
  
  steeringPolicy: STEERING_RANDOM  # Distribute new sessions evenly
  
  origins:
    - name: app-server-1
      address: 203.0.113.10
      weight: 1
    
    - name: app-server-2
      address: 198.51.100.20
      weight: 1
    
    - name: app-server-3
      address: 192.0.2.30
      weight: 1
  
  proxied: true  # Required for session affinity
  healthProbePath: /health
  
  # Enable cookie-based sticky sessions
  sessionAffinity: SESSION_AFFINITY_COOKIE
```

**Expected Behavior**:
1. First request from user → Cloudflare selects an origin randomly
2. Cloudflare sets `__cflb` cookie in response
3. All subsequent requests from that user route to the same origin
4. If that origin becomes unhealthy, user is routed to a new origin (cart may be lost)

**Important**: Session affinity requires `proxied: true`. It does NOT work in DNS-only mode.

**Best Practice**: For production, use a shared session store (Redis, Memcached) instead of in-memory sessions to avoid data loss during failover.

---

## Example 6: DNS-Only Mode (Non-Proxied)

**Use Case**: Load balance non-HTTP traffic (e.g., game servers, MQTT, custom protocols) or when you don't want Cloudflare to proxy traffic.

**Features**:
- DNS returns origin IP directly
- No Layer 7 features (no WAF, no caching, no session affinity)
- Lower latency (no Cloudflare proxy hop)

```yaml
apiVersion: cloudflare.openmcf.org/v1
kind: CloudflareLoadBalancer
metadata:
  name: game-server-lb
spec:
  hostname: game.example.com
  zoneId:
    value: "abc123def456"
  
  steeringPolicy: STEERING_GEO  # Route to nearest game server
  
  origins:
    - name: us-game-server
      address: 203.0.113.10
      weight: 1
    
    - name: eu-game-server
      address: 198.51.100.20
      weight: 1
  
  # DNS-only mode (gray cloud)
  proxied: false
  
  healthProbePath: /status
  
  # Session affinity is NOT supported in DNS-only mode
  sessionAffinity: SESSION_AFFINITY_NONE
```

**Expected Behavior**:
1. DNS query for `game.example.com` → Cloudflare returns IP of nearest healthy origin
2. Client connects directly to that origin (no Cloudflare proxy)
3. DNS response is cached by client/resolver (typically 5 minutes)

**Use Cases for DNS-Only**:
- Non-HTTP protocols (TCP, UDP, MQTT, game servers)
- When you need the client's real IP (not Cloudflare's edge IPs)
- When you manage SSL/TLS termination yourself

**Limitation**: Health check failures take longer to propagate (DNS caching).

---

## Example 7: Custom Health Check Path

**Use Case**: Health endpoint at a specific path that performs deep health validation.

**Features**:
- Custom health check path
- Validates database connectivity, cache availability, etc.
- More reliable than checking root path

```yaml
apiVersion: cloudflare.openmcf.org/v1
kind: CloudflareLoadBalancer
metadata:
  name: api-lb-deep-health
spec:
  hostname: api.example.com
  zoneId:
    value: "abc123def456"
  
  steeringPolicy: STEERING_OFF
  
  origins:
    - name: primary
      address: 203.0.113.10
      weight: 1
    
    - name: secondary
      address: 198.51.100.20
      weight: 1
  
  proxied: true
  
  # Custom health check endpoint
  # This endpoint should:
  # - Return HTTP 200 if healthy
  # - Return HTTP 5xx if unhealthy
  # - Check database, cache, external dependencies
  healthProbePath: /api/v1/healthz
  
  sessionAffinity: SESSION_AFFINITY_NONE
```

**Health Check Best Practices**:

1. **Return 200 only if truly healthy**:
   ```json
   {
     "status": "healthy",
     "database": "connected",
     "cache": "connected",
     "external_api": "reachable"
   }
   ```

2. **Return 5xx if dependencies are down**:
   ```json
   {
     "status": "unhealthy",
     "database": "disconnected"
   }
   ```

3. **Keep health checks lightweight**: Should respond in <1 second

4. **Don't check static assets**: Use an endpoint that validates service functionality

---

## Example 8: Three-Tier Failover

**Use Case**: Production system with primary, backup, and last-resort origins.

**Features**:
- Three levels of failover
- Health checks ensure each tier is validated before use
- Last resort origin can be a static "maintenance mode" page

```yaml
apiVersion: cloudflare.openmcf.org/v1
kind: CloudflareLoadBalancer
metadata:
  name: three-tier-failover-lb
spec:
  hostname: app.example.com
  zoneId:
    value: "abc123def456"
  
  # Priority-based failover
  steeringPolicy: STEERING_OFF
  
  origins:
    # Tier 1: Primary production origin
    - name: primary-prod
      address: 203.0.113.10
      weight: 1
    
    # Tier 2: Hot standby (ready to take traffic)
    - name: hot-standby
      address: 198.51.100.20
      weight: 1
    
    # Tier 3: Last resort - maintenance mode page
    - name: maintenance-mode
      address: 192.0.2.30
      weight: 1
  
  proxied: true
  healthProbePath: /health
  sessionAffinity: SESSION_AFFINITY_COOKIE
```

**Expected Behavior**:
1. All traffic → `primary-prod`
2. If primary fails → `hot-standby`
3. If hot-standby fails → `maintenance-mode` (static page saying "We'll be back soon")

**Last Resort Pattern**: The maintenance-mode origin can be:
- A static S3 bucket with a "Down for Maintenance" page
- A minimal nginx server serving a 503 page
- An always-healthy origin that shows degraded functionality

---

## Testing Your Load Balancer

### Verify Health Checks

1. Check Cloudflare dashboard → Traffic → Load Balancing → Monitors
2. Verify all origins show as "Healthy"
3. If an origin shows "Unhealthy", check:
   - Origin server is running
   - Health endpoint returns HTTP 200
   - Firewall allows Cloudflare health check IPs

### Test Failover

1. Stop the primary origin server
2. Wait for health checks to detect failure (60 seconds for Pro plan)
3. Verify traffic moves to secondary origin
4. Restart primary origin
5. Verify traffic fails back to primary

### Test Session Affinity

```bash
# Make multiple requests with same cookie
curl -c cookies.txt https://app.example.com/
curl -b cookies.txt https://app.example.com/
curl -b cookies.txt https://app.example.com/

# All requests should hit the same origin (check response headers or logs)
```

---

## Integration Examples

### Example: Load Balancer with CloudflareDnsZone

Reference a managed DNS zone:

```yaml
---
# Create the DNS zone first
apiVersion: cloudflare.openmcf.org/v1
kind: CloudflareDnsZone
metadata:
  name: example-zone
spec:
  domain_name: example.com
  account_id:
    value: "xyz789"

---
# Then create the load balancer
apiVersion: cloudflare.openmcf.org/v1
kind: CloudflareLoadBalancer
metadata:
  name: api-lb
spec:
  hostname: api.example.com
  zoneId:
    ref: example-zone  # References the CloudflareDnsZone
  origins:
    - name: origin-1
      address: 203.0.113.10
  proxied: true
  healthProbePath: /health
```

### Example: Load Balancer with Cloudflare Tunnel

Use Cloudflare Tunnel to load balance private origins:

```yaml
apiVersion: cloudflare.openmcf.org/v1
kind: CloudflareLoadBalancer
metadata:
  name: internal-api-lb
spec:
  hostname: internal-api.example.com
  zoneId:
    value: "abc123"
  
  origins:
    # Private origin accessible via Cloudflare Tunnel
    - name: private-origin-1
      address: <tunnel-uuid-1>.cfargotunnel.com
      weight: 1
    
    - name: private-origin-2
      address: <tunnel-uuid-2>.cfargotunnel.com
      weight: 1
  
  proxied: true
  healthProbePath: /health
  steeringPolicy: STEERING_OFF
```

**Use Case**: Load balance between multiple private data centers without public IPs.

---

## Common Issues and Solutions

### Issue: Origins always show as unhealthy

**Cause**: Firewall blocking Cloudflare health check IPs

**Solution**:
1. Allow Cloudflare IP ranges in your origin firewall
2. Verify health endpoint returns HTTP 200
3. Check health endpoint path is correct

### Issue: Failover not working

**Cause**: Health check path misconfigured or returning wrong status code

**Solution**:
1. Verify `health_probe_path` points to a real endpoint
2. Ensure endpoint returns HTTP 2xx (not 3xx redirect)
3. Check health check interval (60s minimum on Pro plan)

### Issue: Session affinity not working

**Cause**: DNS-only mode enabled

**Solution**:
- Session affinity requires `proxied: true` (orange cloud)
- Change `proxied: false` to `proxied: true`

### Issue: Traffic not balanced evenly

**Cause**: Using `steering_policy: STEERING_OFF` (failover mode)

**Solution**:
- For even distribution, use `steering_policy: STEERING_RANDOM` with equal weights

---

## Next Steps

- Read the [README.md](./README.md) for detailed field documentation
- Explore the [research documentation](./docs/README.md) for architectural deep dive
- Deploy using [Pulumi](./iac/pulumi/README.md) or [Terraform](./iac/tf/README.md)
- Review [Cloudflare Load Balancer documentation](https://developers.cloudflare.com/load-balancing/)

---

**Questions or Issues?** Refer to the [README.md](./README.md) or Cloudflare's official documentation.

