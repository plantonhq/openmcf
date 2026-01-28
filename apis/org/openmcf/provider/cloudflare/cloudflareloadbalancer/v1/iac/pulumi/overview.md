# Pulumi Module Overview: Cloudflare Load Balancer

## Architecture

This document explains the architectural decisions, resource dependency graph, and design rationale for the Cloudflare Load Balancer Pulumi module.

## Resource Model

The Cloudflare Load Balancer requires three distinct resources with specific parent-child relationships:

```
┌─────────────────────────────────────────────────────────────┐
│                    Cloudflare Account                        │
│                                                               │
│  ┌─────────────────────────┐                                │
│  │ LoadBalancerMonitor     │                                │
│  │ (Health Check)          │                                │
│  │ - Type: HTTPS           │                                │
│  │ - Path: /health         │                                │
│  │ - Interval: 60s         │                                │
│  └──────────┬──────────────┘                                │
│             │                                                 │
│             │ references                                      │
│             ↓                                                 │
│  ┌─────────────────────────┐                                │
│  │ LoadBalancerPool        │                                │
│  │ (Origin Group)          │                                │
│  │ - Origins: [...list]    │                                │
│  │ - Monitor ID: ↑         │                                │
│  └──────────┬──────────────┘                                │
│             │                                                 │
└─────────────┼─────────────────────────────────────────────────┘
              │ references
              ↓
┌─────────────────────────────────────────────────────────────┐
│                    Cloudflare Zone                           │
│  ┌─────────────────────────┐                                │
│  │ LoadBalancer            │                                │
│  │ (DNS Hostname)          │                                │
│  │ - Name: api.example.com │                                │
│  │ - Pool IDs: [↑]         │                                │
│  │ - Steering: off         │                                │
│  └─────────────────────────┘                                │
└─────────────────────────────────────────────────────────────┘
```

### Resource Hierarchy Explained

**1. LoadBalancerMonitor (Account-Level)**

- **Scope**: Account-wide resource
- **Purpose**: Defines health check parameters (path, interval, expected codes)
- **Reusability**: Can be shared across multiple pools and load balancers
- **API Endpoint**: `POST /accounts/{account_id}/load_balancers/monitors`

**Why account-level?** Cloudflare's design allows health check configurations to be reused across zones. For example, if you have 10 zones (domains), you don't need 10 separate health check configurations if they all use the same `/health` endpoint.

**2. LoadBalancerPool (Account-Level)**

- **Scope**: Account-wide resource
- **Purpose**: Groups origin servers and associates them with a health monitor
- **Dependency**: References a Monitor ID
- **Reusability**: Can be referenced by multiple load balancers across zones
- **API Endpoint**: `POST /accounts/{account_id}/load_balancers/pools`

**Why account-level?** Pools represent backend infrastructure (origin servers) that may serve multiple applications across different zones. For instance, a pool of Kubernetes nodes might back both `api.example.com` and `app.example.com`.

**3. LoadBalancer (Zone-Level)**

- **Scope**: Zone-specific resource (tied to a DNS zone)
- **Purpose**: Maps a DNS hostname to one or more pools with routing policies
- **Dependency**: References one or more Pool IDs
- **API Endpoint**: `POST /zones/{zone_id}/load_balancers`

**Why zone-level?** The load balancer creates a DNS record (`api.example.com`) within a specific zone (`example.com`). DNS records are inherently zone-scoped.

## Dependency Graph

The Pulumi module creates resources in this order:

```
1. Monitor → 2. Pool → 3. Load Balancer
   (uses Health Check ID) (uses Pool ID)
```

### Explicit Dependencies

```go
createdMonitor, err := cloudflare.NewLoadBalancerMonitor(...)
// ↓
createdPool, err := cloudflare.NewLoadBalancerPool(...,
    Monitor: createdMonitor.ID(),  // Explicit dependency
)
// ↓
createdLoadBalancer, err := cloudflare.NewLoadBalancer(...,
    DefaultPools: pulumi.StringArray{createdPool.ID()},  // Explicit dependency
)
```

Pulumi automatically orders resource creation based on these dependencies.

## OpenMCF Abstraction

### The Core Problem

When using the Cloudflare API or Terraform directly, users must:
1. Manually create a Monitor and capture its ID
2. Create a Pool, passing the Monitor ID
3. Create a Load Balancer, passing the Pool ID(s)

This is verbose, error-prone, and requires deep knowledge of Cloudflare's resource model.

### The Solution: Denormalized API

OpenMCF provides a **flattened, denormalized API** that hides this complexity:

```yaml
# User's input (simple and intuitive)
spec:
  hostname: api.example.com
  origins:
    - name: primary
      address: 203.0.113.10
    - name: secondary
      address: 198.51.100.20
  health_probe_path: /health
```

Behind the scenes, the Pulumi module:
1. Creates a Monitor using `health_probe_path`
2. Creates a Pool with the list of origins
3. Creates a Load Balancer referencing the pool
4. Wires everything together by ID

**User sees**: Simple, inline origin definitions  
**Module handles**: Complex resource fan-out and ID management

## Implementation Details

### File Structure

```
module/
├── main.go             # Entry point - Resources() function
├── locals.go           # Local variables (metadata, spec shortcuts)
├── load_balancer.go    # Core logic (creates monitor, pool, LB)
└── outputs.go          # Output constant definitions
```

### Function Flow

**1. `main.go:Resources()`**

- Entry point called by Pulumi CLI
- Initializes locals (shortcuts for metadata and spec)
- Creates Cloudflare provider
- Calls `load_balancer()` to create resources

**2. `load_balancer.go:load_balancer()`**

This is where the magic happens:

```go
func load_balancer(ctx *pulumi.Context, locals *Locals, provider *cloudflare.Provider) {
    // Step 1: Create Monitor
    monitor, err := cloudflare.NewLoadBalancerMonitor(ctx, "monitor", &LoadBalancerMonitorArgs{
        Type: "http",
        Path: locals.CloudflareLoadBalancer.Spec.HealthProbePath,
        // ...
    })

    // Step 2: Create Pool (references Monitor)
    var poolOrigins []LoadBalancerPoolOriginArgs
    for _, origin := range locals.CloudflareLoadBalancer.Spec.Origins {
        poolOrigins = append(poolOrigins, LoadBalancerPoolOriginArgs{
            Name:    origin.Name,
            Address: origin.Address,
            Weight:  origin.Weight,
        })
    }

    pool, err := cloudflare.NewLoadBalancerPool(ctx, "pool", &LoadBalancerPoolArgs{
        Origins: poolOrigins,
        Monitor: monitor.ID(),  // Wire by ID
    })

    // Step 3: Create Load Balancer (references Pool)
    lb, err := cloudflare.NewLoadBalancer(ctx, "load_balancer", &LoadBalancerArgs{
        Name:         locals.CloudflareLoadBalancer.Spec.Hostname,
        DefaultPools: pulumi.StringArray{pool.ID()},  // Wire by ID
        // ...
    })

    // Step 4: Export outputs
    ctx.Export("load_balancer_id", lb.ID())
}
```

### Enum to String Mapping

Protobuf enums are integers, but Cloudflare expects string values. The module handles this translation:

```go
// Session Affinity: Proto enum → Cloudflare string
switch locals.CloudflareLoadBalancer.Spec.SessionAffinity {
case 0:  // SESSION_AFFINITY_NONE
    affinity = pulumi.StringPtr("none")
case 1:  // SESSION_AFFINITY_COOKIE
    affinity = pulumi.StringPtr("cookie")
}

// Steering Policy: Proto enum → Cloudflare string
switch locals.CloudflareLoadBalancer.Spec.SteeringPolicy {
case 0:  // STEERING_OFF (active-passive failover)
    steering = pulumi.StringPtr("off")
case 1:  // STEERING_GEO
    steering = pulumi.StringPtr("geo")
case 2:  // STEERING_RANDOM
    steering = pulumi.StringPtr("random")
}
```

This abstraction allows the proto API to use strongly-typed enums while the Pulumi module handles Cloudflare's string-based API.

## Design Decisions

### Decision 1: Single Pool Per Load Balancer

**Choice**: Create one pool containing all origins

**Rationale**: 
- Simplifies the API (users don't need to think about pools)
- Covers 80% use case (one pool with multiple origins for failover)
- Advanced use cases (multiple pools per LB for geo-routing) can be added later

**Trade-off**: 
- ✅ Simple API for common patterns
- ❌ Cannot represent complex multi-pool scenarios (e.g., separate US and EU pools)

**Future Enhancement**: Add optional `pools` field to support advanced multi-pool configurations

### Decision 2: Health Check Defaults

**Choice**: Default health probe path to `/`

**Rationale**:
- Most origin servers have a root path that returns 200
- Users can override with `health_probe_path` if needed
- Reduces required configuration fields

**Trade-off**:
- ✅ Fewer required fields
- ❌ Root path might not represent true service health

**Best Practice**: Users should set `health_probe_path` to a dedicated health endpoint (e.g., `/healthz`)

### Decision 3: Automatic Fallback Pool

**Choice**: Set `fallback_pool_id` to the same pool as `default_pool_ids`

**Rationale**:
- Cloudflare requires a fallback pool (pool of last resort)
- If all origins in the default pool fail, traffic still needs to go somewhere
- Using the same pool ensures traffic continues flowing (even if degraded)

**Trade-off**:
- ✅ Prevents total outage if all origins fail
- ❌ Sends traffic to unhealthy origins as a last resort

**Alternative**: Allow users to define a separate "maintenance mode" pool with a static error page

### Decision 4: Proxied Mode Default

**Choice**: Default `proxied` to `true` (orange cloud)

**Rationale**:
- 80% of load balancers benefit from Layer 7 features (WAF, caching, session affinity)
- Proxied mode is Cloudflare's recommended default
- DNS-only mode is the exception, not the rule

**Trade-off**:
- ✅ Enables advanced features by default
- ❌ Adds latency (one extra hop through Cloudflare proxy)

**When to override**: Use `proxied: false` for non-HTTP protocols or when you need the client's real IP

## Secret Management

### API Token Handling

Cloudflare API tokens are highly sensitive. The module supports two patterns:

**1. Environment Variables (Recommended)**

```bash
export CLOUDFLARE_API_TOKEN="your-token"
pulumi up
```

The Pulumi Cloudflare provider automatically reads `CLOUDFLARE_API_TOKEN`.

**2. Pulumi Config Secrets**

```bash
pulumi config set cloudflare:apiToken --secret "your-token"
pulumi up
```

Pulumi encrypts the secret before storing it in the state file.

**Why environment variables are preferred**:
- Secrets never touch the state file
- Compatible with CI/CD systems (GitHub Actions secrets, etc.)
- Follows 12-factor app principles

## Error Handling

The module uses Go's standard error handling pattern:

```go
monitor, err := cloudflare.NewLoadBalancerMonitor(...)
if err != nil {
    return fmt.Errorf("failed to create monitor: %w", err)
}
```

Pulumi captures these errors and displays them with context:

```
error: failed to create load balancer: failed to create pool: API Error: 
  zone_id 'invalid' is not a valid zone ID
```

## Outputs

The module exports three outputs:

| Output | Type | Description |
|--------|------|-------------|
| `load_balancer_id` | string | Cloudflare LB resource ID |
| `load_balancer_dns_record_name` | string | Hostname (e.g., `api.example.com`) |
| `load_balancer_cname_target` | string | CNAME target for external DNS |

These outputs can be used by other Pulumi stacks or external systems:

```bash
# Get load balancer ID for use in scripts
LB_ID=$(pulumi stack output load_balancer_id)

# Use in other Pulumi stacks via StackReference
lb_id = pulumi.StackReference("org/stack-name").get_output("load_balancer_id")
```

## Provider Configuration

The module uses the `pulumicloudflareprovider.Get()` helper to create a Cloudflare provider:

```go
cloudflareProvider, err := pulumicloudflareprovider.Get(ctx, stackInput.ProviderConfig)
```

This helper:
1. Extracts API token from `ProviderConfig` or environment variables
2. Creates a Cloudflare provider instance
3. Handles authentication errors gracefully

All resources are created using this provider:

```go
monitor, err := cloudflare.NewLoadBalancerMonitor(ctx, "monitor", &args,
    pulumi.Provider(cloudflareProvider),  // Use explicit provider
)
```

**Why explicit providers?** Allows multiple Cloudflare accounts in the same Pulumi program (e.g., dev and prod accounts).

## Testing Strategy

### Unit Tests

The module includes unit tests for:
- Enum-to-string conversion logic
- Local variable initialization
- Error handling paths

Run tests:

```bash
cd module
go test -v ./...
```

### Integration Tests

Integration tests deploy real Cloudflare resources (requires API token):

```bash
export CLOUDFLARE_API_TOKEN="..."
export CLOUDFLARE_ZONE_ID="..."
go test -v -tags=integration ./...
```

**Cost**: Each integration test run costs ~$0 (Cloudflare has no per-request charges, only monthly subscription)

## Performance Considerations

### Resource Creation Time

Typical deployment timeline:
- Monitor creation: ~1 second
- Pool creation: ~2 seconds (includes health check initialization)
- Load Balancer creation: ~3 seconds
- **Total**: 5-10 seconds

### State File Size

Each resource adds ~1-2 KB to the Pulumi state file:
- Monitor: ~1 KB
- Pool: ~1.5 KB (depends on origin count)
- Load Balancer: ~1.5 KB
- **Total**: ~4 KB per load balancer

## Future Enhancements

### 1. Multi-Pool Support

Allow users to define multiple pools for advanced geo-routing:

```yaml
spec:
  pools:
    - name: us-pool
      origins: [...]
    - name: eu-pool
      origins: [...]
  region_pools:
    US: us-pool
    EU: eu-pool
```

### 2. Custom Monitor Types

Support TCP and UDP health checks:

```yaml
spec:
  health_check:
    type: TCP  # or UDP
    port: 443
```

### 3. Advanced Steering Rules

Support Cloudflare's dynamic latency steering:

```yaml
spec:
  steering_policy: STEERING_DYNAMIC_LATENCY  # Enterprise only
```

### 4. Maintenance Mode Pool

Allow a separate fallback pool for maintenance pages:

```yaml
spec:
  fallback_pool:
    origins:
      - name: maintenance
        address: s3-maintenance-page.s3-website-us-east-1.amazonaws.com
```

## Comparison: Pulumi vs. Terraform

| Aspect | Pulumi (This Module) | Terraform |
|--------|----------------------|-----------|
| **Language** | Go (strongly typed) | HCL (declarative) |
| **Secret Management** | Built-in encryption | Requires external vault |
| **Type Safety** | Compile-time checks | Runtime errors |
| **Conditionals** | Go if/else | HCL count/for_each |
| **Reusability** | Go functions | Terraform modules |
| **Learning Curve** | Moderate (requires Go knowledge) | Low (HCL is simple) |

**Recommendation**: 
- Use Pulumi if you prefer type-safe, general-purpose languages
- Use Terraform if you prefer declarative DSLs and have an existing Terraform codebase

## References

- [Cloudflare Load Balancer API Docs](https://developers.cloudflare.com/api/operations/load-balancers-create-load-balancer)
- [Pulumi Cloudflare Provider](https://www.pulumi.com/registry/packages/cloudflare/)
- [Component README](../../README.md)
- [Research Documentation](../../docs/README.md)

---

**Questions?** Review the [README.md](./README.md) or consult Cloudflare's API documentation.

