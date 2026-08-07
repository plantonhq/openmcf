# DNS Zone on Cloudflare

Deploys a DNS zone on Cloudflare with configurable plan level, proxy defaults, and optional inline DNS records. The zone integrates with Planton's Provider Connections for Cloudflare credential management and serves as the foundation for DNS record management and Cloudflare proxy services.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Cloudflare Zone** -- a DNS zone for the specified domain, created under the target Cloudflare account with the selected plan level (Free, Pro, Business, or Enterprise)
- **DNS Records** -- created only when `records` entries are provided; supports A, AAAA, CNAME, MX, TXT, SRV, NS, and CAA record types with configurable TTL, proxy status, and priority
- **Cloudflare Labels** -- resource metadata applied to the zone for organization and environment tracking

## Before You Deploy

### Planton Setup

- **Cloudflare Provider Connection** -- an active connection in the Connect module with a Cloudflare API token. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline API token authentication.

### Cloudflare Account

- **A Cloudflare account** with the target domain's registrar nameservers pointed to Cloudflare, or the ability to update nameservers after zone creation. The `accountId` field identifies which Cloudflare account owns the zone.
- **Domain registration** -- the domain specified in `zoneName` must be registered with a domain registrar. Cloudflare will assign nameservers after zone creation that must be configured at the registrar.

## Deploy

### Console

Open the deployment store, find **DNS Zone on Cloudflare**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Free Plan Zone** preset in the [Presets](#presets) tab to pre-populate a working configuration.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: cloudflare.planton.dev/v1
kind: CloudflareDnsZone
metadata:
  name: example-zone
  org: acme-corp
  env: prod
spec:
  zoneName: example.com
  accountId: "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4"
  plan: FREE
```

```shell
planton apply -f cloudflare-dns-zone.yaml
```

This creates a Cloudflare DNS zone for example.com on the Free plan with no inline DNS records. Cloudflare assigns nameservers that must be configured at your domain registrar. A Stack Job tracks the provisioning in real time.

## Key Configuration

These are the most important decisions when configuring a DNS zone. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Plan level** -- Defaults to `FREE`, which includes basic DNS hosting, CDN, and DDoS protection for proxied records. Choose `PRO`, `BUSINESS`, or `ENTERPRISE` for advanced features like WAF rulesets, image optimization, and priority support. Plan changes affect billing immediately.

**Paused mode** -- Set `paused: true` to create the zone in DNS-only mode with no Cloudflare proxy, CDN, or security features active. Useful during initial migration when you want to verify DNS resolution before enabling Cloudflare's proxy layer.

**Default proxy behavior** -- Set `defaultProxied: true` to make new DNS records default to proxied (orange cloud) mode. When false (default), new records default to DNS-only (grey cloud). This affects only the default for new records created in the Cloudflare dashboard -- records in the `records` field have their own explicit `proxied` setting.

**Inline records vs. separate resources** -- The `records` field lets you manage DNS records as part of the zone configuration. For small, stable record sets this is convenient. For dynamic or individually managed records, use separate CloudflareDnsRecord resources that reference this zone via ValueFromRef.

## Outputs and Dependencies

### What This Component Consumes

This component has no foreign key dependencies.

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `zone_id` | The Cloudflare Zone ID of the created DNS zone | CloudflareDnsRecord zone references, CloudflareR2Bucket custom domain configuration, CloudflareWorker DNS routing |
| `nameservers` | The list of nameserver addresses assigned to this DNS zone | Domain registrar NS record configuration |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Free plan zone** -- A DNS zone on the Free plan with no inline records. Add DNS records separately as CloudflareDnsRecord resources for independent lifecycle management. Suitable for most domains that need DNS hosting with basic CDN and DDoS protection. Start from the **Free Plan Zone** preset.

## Works With

This component operates independently and does not reference other deployment components.