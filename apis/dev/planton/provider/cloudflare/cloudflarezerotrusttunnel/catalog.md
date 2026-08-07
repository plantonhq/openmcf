# Zero Trust Tunnel on Cloudflare

Provisions a Cloudflare Tunnel (cloudflared): a secure, outbound-only connection from a private network to Cloudflare's edge. A tunnel exposes private HTTP/TCP/SSH/RDP services via public hostnames (ingress rules) and makes private IP ranges reachable to WARP clients (via `CloudflareZeroTrustTunnelRoute`) -- without opening a single inbound firewall port. The connector authenticates with the run token exported in the stack outputs. Integrates with Planton's Provider Connections for Cloudflare credential management.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Tunnel** -- the cloudflared tunnel object and its run token
- **Ingress configuration** -- public-hostname rules and origin-connection settings, when the tunnel is Cloudflare-managed
- **Cloudflare Labels** -- resource metadata applied for organization and environment tracking

## Before You Deploy

### Planton Setup

- **Cloudflare Provider Connection** -- an active connection in the Connect module with a Cloudflare API token that has Cloudflare Tunnel edit access. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline API token authentication.

### Cloudflare Account

- **Zero Trust enabled** -- the account must have Cloudflare Zero Trust (Cloudflare One) set up.
- **A connector to run** -- after deploy, run `cloudflared` somewhere on the private network with the run token from the outputs.

## Deploy

### Console

Open the deployment store, find **Zero Trust Tunnel on Cloudflare**, and click **Deploy**. The creation wizard captures the tunnel identity (account, name, configuration source, and an optional managed-secret run secret), the public-hostname ingress rules, and the tunnel-level origin connection defaults (including a Cloudflare Access block).

### CLI

Create a manifest and apply it:

```yaml
apiVersion: cloudflare.planton.dev/v1
kind: CloudflareZeroTrustTunnel
metadata:
  name: prod-connector
  org: acme-corp
  env: prod
spec:
  accountId: a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4
  name: prod-connector
  configSrc: cloudflare
  ingress:
    - hostname: app.example.com
      service: http://localhost:8080
    - service: http_status:404
```

```shell
planton apply -f cloudflare-zero-trust-tunnel.yaml
```

This exposes `app.example.com` through the tunnel to a local service on port 8080, with a catch-all 404 for everything else. A Stack Job tracks the provisioning in real time.

## Key Configuration

These are the most important decisions when configuring a tunnel. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Account (`accountId`)** -- The account that owns the tunnel. Immutable -- changing it replaces the tunnel.

**Configuration Source (`configSrc`)** -- `cloudflare` (default) manages ingress here as desired state; `local` leaves ingress to a cloudflared YAML file on the origin. Immutable -- changing it replaces the tunnel.

**Ingress (`ingress`)** -- Public-hostname rules, evaluated top to bottom. The final rule must be a catch-all (a service with no hostname, e.g. `http_status:404`). A DNS CNAME for each hostname must point at the tunnel's CNAME target.

**Origin Defaults (`originRequest`)** -- Connection settings applied to every ingress rule unless a rule overrides them. The Cloudflare Access sub-block (`originRequest.access`) requires an Access JWT on the matched hostnames, referencing `CloudflareZeroTrustAccessApplication` resources.

**Tunnel Secret (`tunnelSecret`)** -- An optional, reference-only managed secret for a locally-managed connector. Leave empty to let Cloudflare generate one.

## Outputs and Dependencies

### What This Component Consumes

A tunnel's Access block references **CloudflareZeroTrustAccessApplication** resources (via `originRequest.access.audTag`).

### What This Component Provides

After provisioning, `status.outputs` contains:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `tunnel_id` | The Cloudflare-assigned UUID of the tunnel | Referenced by `CloudflareZeroTrustTunnelRoute` |
| `tunnel_cname` | The CNAME target for public hostnames | Point a `CloudflareDnsRecord` CNAME at it |
| `tunnel_token` | The connector run token (sensitive) | `cloudflared tunnel run --token <token>` |
| `tunnel_status` | The tunnel health (inactive/degraded/healthy/down) | Dashboards, alerting |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Public hostname** -- an ingress rule maps a public hostname to a local HTTP service, with a catch-all 404.

**Access-protected** -- the origin Access block requires an Access application's JWT so only authorized users reach a hostname.

**Private network connector** -- pair the tunnel with `CloudflareZeroTrustTunnelRoute` resources to make private CIDRs reachable to WARP clients.

## Works With

- [**Zero Trust Tunnel Route on Cloudflare**](/cloud-catalog/cloudflare-zero-trust-tunnel-route) -- advertises private CIDRs through this tunnel
- [**Zero Trust Tunnel Virtual Network on Cloudflare**](/cloud-catalog/cloudflare-zero-trust-tunnel-virtual-network) -- the routing segment a route belongs to
- [**Zero Trust Access Application on Cloudflare**](/cloud-catalog/cloudflare-zero-trust-access-application) -- the Access apps an ingress Access block references
- [**DNS Record on Cloudflare**](/cloud-catalog/cloudflare-dns-record) -- the CNAME that points a public hostname at the tunnel
