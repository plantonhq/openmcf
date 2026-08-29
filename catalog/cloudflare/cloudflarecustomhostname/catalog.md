# Custom Hostname on Cloudflare

Attaches a customer's own domain to a Cloudflare-for-SaaS zone. It extends Cloudflare's edge -- TLS termination, caching, WAF -- onto a hostname your customer owns (e.g. `support.acme.com`), with a per-customer certificate that Cloudflare provisions and auto-renews. The customer points their hostname at the SaaS zone via CNAME and proves control using the ownership-verification records exported in the stack outputs. Integrates with Planton's Provider Connections for Cloudflare credential management.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Custom Hostname** -- a customer hostname onboarded to the SaaS zone
- **Per-hostname Certificate** -- a managed DV certificate (by default), auto-renewed by Cloudflare

## Before You Deploy

### Planton Setup

- **Cloudflare Provider Connection** -- an active connection in the Connect module with a Cloudflare API token that has SSL and Certificates edit access. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline API token authentication.

### Cloudflare Account

- **A SaaS zone with a fallback origin** -- the zone needs a `CloudflareCustomHostnameFallbackOrigin` so traffic has a default backend (unless this hostname sets its own origin override).

## Deploy

### Console

Open the deployment store, find **Custom Hostname on Cloudflare**, and click **Deploy**. The creation wizard captures the SaaS zone and hostname, an optional per-hostname origin override (with a live connection diagram), and an optional TLS/SSL block.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: cloudflare.planton.dev/v1
kind: CloudflareCustomHostname
metadata:
  name: acme-support-hostname
  org: acme-corp
  env: prod
spec:
  zoneId:
    valueFrom:
      kind: CloudflareDnsZone
      name: saas-zone
      fieldPath: status.outputs.zone_id
  hostname: support.acme.com
```

```shell
planton apply -f cloudflare-custom-hostname.yaml
```

This onboards `support.acme.com` to the SaaS zone with a managed DV certificate. The customer then CNAMEs their hostname to the zone and adds the ownership-verification record from the outputs. A Stack Job tracks the provisioning in real time.

## Key Configuration

These are the most important decisions when configuring a custom hostname. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Zone (`zoneId`)** -- The SaaS zone this hostname is added to. Immutable -- changing it replaces the resource.

**Hostname (`hostname`)** -- The customer's own (external) hostname to onboard. Immutable.

**Custom Origin Server (`customOriginServer`)** -- Optional per-hostname origin override; leave empty to use the zone's fallback origin. A literal backend hostname or a reference to another resource's output. Editable.

**Custom Metadata (`customMetadata`)** -- Arbitrary key/value data (e.g. a tenant id) available to Workers and rules. Editable.

**SSL (`ssl`)** -- Optional per-hostname certificate configuration. Leave it unset for a managed DV certificate; open it to change issuance, upload your own certificate (Enterprise), or tune TLS termination. The private key fields are reference-only managed secrets.

## Outputs and Dependencies

### What This Component Consumes

The custom hostname references a **CloudflareDnsZone** (via `zoneId`) and, optionally, a backend endpoint (via `customOriginServer`). It relies on the zone's **CloudflareCustomHostnameFallbackOrigin** when no override is set.

### What This Component Provides

After provisioning, `status.outputs` contains:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `custom_hostname_id` | The Cloudflare-assigned identifier | Verification, dashboards |
| `ownership_verification_name` / `type` / `value` | The DNS record the customer adds to prove control | Customer onboarding hand-off |
| `ownership_verification_http_url` / `http_body` | The HTTP alternative for ownership verification | Customer onboarding hand-off |
| `verification_errors` | Any verification errors reported by Cloudflare | Troubleshooting |
| `created_at` | RFC3339 creation timestamp | Auditing |
| `zone_id` | The SaaS zone the hostname belongs to | Verification, imports, composition |

Activation status is deliberately not an output — it transitions asynchronously (`pending` → `pending_validation` → `active`); read it from the Cloudflare API or dashboard.

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Managed DV onboarding** -- onboard a customer hostname with an auto-provisioned certificate and hand them the ownership-verification record.

**Dedicated backend** -- set a custom origin server to route a specific customer to their own backend.

## Works With

- [**Custom Hostname Fallback Origin on Cloudflare**](/cloud-catalog/cloudflare-custom-hostname-fallback-origin) -- the zone-level default backend this hostname routes to
- [**DNS Zone on Cloudflare**](/cloud-catalog/cloudflare-dns-zone) -- the SaaS zone the hostname is added to
