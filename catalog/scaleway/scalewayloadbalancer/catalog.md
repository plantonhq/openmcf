# Scaleway Load Balancer

Deploys a managed Layer 4/7 Load Balancer on Scaleway as a composite resource that bundles a dedicated Flexible IP, the Load Balancer appliance, backend server pools with health checks, frontend listeners, and optional TLS certificates (Let's Encrypt or custom PEM) into a single declarative unit. Supports ValueFromRef for Private Network dependency wiring in InfraCharts.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Flexible IP** -- a dedicated public IPv4 address with independent lifecycle that survives Load Balancer replacement, preserving DNS records and firewall rules
- **Load Balancer** -- a zonal traffic distribution appliance with the configured type, optional Private Network attachment, and SSL compatibility level
- **Backends** -- one or more named server pools with configurable health checks (TCP, HTTP, or HTTPS), load-balancing algorithms, sticky sessions, and connection timeouts
- **Frontends** -- one or more named listeners on specific TCP ports that route incoming traffic to backends, with optional TLS certificate references and HTTP/3 support
- **TLS Certificates** -- created only when `certificates` entries are provided; either auto-provisioned via Let's Encrypt or supplied as custom PEM chains
- **Scaleway Tags** -- resource metadata tags (resource name, kind, organization, environment) applied automatically for tracking and governance

## Before You Deploy

### Scaleway Account

- **A Scaleway account** with an active project and API access key pair.
- **A Private Network** in the target region for backend server connectivity. Provide the Private Network UUID directly or reference a ScalewayPrivateNetwork Cloud Resource via ValueFromRef. Optional -- if omitted, the LB reaches backends via public IPs only.
- **Backend servers** with known IP addresses on the Private Network (or public IPs if no Private Network is used). At least one backend server IP is required.
- **A DNS A record** pointing to the LB's public IP before enabling Let's Encrypt certificates. The domain must resolve to the LB for ACME validation to succeed.

## Deploy

### Console

Open the deployment store, find **Scaleway Load Balancer**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **HTTPS Load Balancer with Let's Encrypt** preset in the [Presets](#presets) tab for a production configuration with automatic TLS.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: scaleway.planton.dev/v1
kind: ScalewayLoadBalancer
metadata:
  name: web-lb
  org: acme-corp
  env: prod
spec:
  zone: fr-par-1
  type: LB-S
  backends:
    - name: web
      serverIps:
        - "10.0.1.5"
      forwardPort: 80
      forwardProtocol: http
  frontends:
    - name: http
      inboundPort: 80
      backendName: web
```

```shell
planton apply -f scaleway-lb.yaml
```

This creates a small Load Balancer in the Paris zone with one HTTP backend and one HTTP frontend on port 80. No Private Network, TLS certificates, or health check customization is configured. A Stack Job tracks the provisioning in real time.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the Load Balancer to a Private Network deployed in the same InfraPipeline:

```yaml
spec:
  privateNetworkId:
    valueFrom:
      kind: ScalewayPrivateNetwork
      name: app-network
      fieldPath: status.outputs.private_network_id
```

The InfraPipeline resolves the dependency graph, deploys the Private Network first, then provisions the Load Balancer attached to it.

## Key Configuration

These are the most important decisions when configuring a Load Balancer. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Load Balancer type** -- The `type` field sets bandwidth and pricing tier. `LB-S` (up to 400 Mbps) is sufficient for most web applications. `LB-GP-M` (up to 4 Gbps) handles high-traffic production workloads. Type can be changed after creation.

**Backend health checks** -- Each backend's `healthCheck` controls how the LB detects and removes unhealthy servers. Use `tcp` for non-HTTP services, `http` with a `/health` endpoint for web applications, or `https` when backends require TLS. Default: TCP check with 5-second interval and 3 retries.

**TLS certificates** -- Add `certificates` entries with either `letsencrypt` (auto-provisioned and auto-renewed) or `customCertificate` (user-provided PEM chain). Frontends reference certificates by name in `certificateNames` for HTTPS termination. Let's Encrypt requires the domain to resolve to the LB's public IP.

**Sticky sessions** -- Set `stickySessions` on a backend to `cookie` for HTTP session affinity or `table` for TCP-level affinity. Required when backend servers maintain client session state.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **ScalewayPrivateNetwork** (optional) | `privateNetworkId` | `status.outputs.private_network_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `lb_id` | Unique identifier of the created Load Balancer | Scaleway API operations, monitoring dashboards |
| `lb_ip_address` | Public IPv4 address assigned to the Load Balancer's Flexible IP | ScalewayDnsRecord A records, firewall allowlists, external service whitelisting |
| `lb_ip_id` | Unique identifier of the Flexible IP resource | IP lifecycle management, reassignment to replacement Load Balancers |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**HTTPS with Let's Encrypt** -- A small Load Balancer on a Private Network with automatic TLS certificate provisioning, HTTP health checks, and two frontends (HTTPS on 443 and HTTP on 80). The standard production configuration for web applications. Start from the **HTTPS Load Balancer with Let's Encrypt** preset.

**Simple HTTP** -- A minimal Load Balancer with a single HTTP frontend and TCP health checks. No TLS or Private Network. Suitable for development, internal services, or scenarios where TLS is terminated elsewhere. Start from the **Simple HTTP Load Balancer** preset.

## Works With

- [**Scaleway Private Network**](/cloud-catalog/scaleway-private-network) -- provides the network for backend server connectivity via private IPs