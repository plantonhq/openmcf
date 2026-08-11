---
title: "DNS Zone"
description: "DNS Zone deployment documentation"
icon: "package"
order: 100
componentName: "digitaloceandnszone"
---

# DNS Zone on DigitalOcean

Deploys a DNS zone (domain) on DigitalOcean with optional inline DNS records including A, AAAA, CNAME, MX, TXT, SRV, and CAA types. DigitalOcean manages the authoritative nameservers for the zone, and record values can reference outputs from other Cloud Resources via ValueFromRef. Integrates with Planton's Provider Connections for DigitalOcean API token management.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **DigitalOcean Domain** -- a DNS zone registered with DigitalOcean's nameservers (`ns1.digitalocean.com`, `ns2.digitalocean.com`, `ns3.digitalocean.com`)
- **DNS Records** -- created only when `records` are provided; one DigitalOcean DNS record per value entry, with type-specific fields for MX priority, SRV weight/port, and CAA flags/tag

## Before You Deploy

### Planton Setup

- **DigitalOcean Provider Connection** -- an active connection in the Connect module with a DigitalOcean API token. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline API token authentication.

### DigitalOcean Account

- **A registered domain name** -- you must own the domain. After creating the zone, update your domain registrar's nameservers to `ns1.digitalocean.com`, `ns2.digitalocean.com`, and `ns3.digitalocean.com`.
- **IP addresses or hostnames for records** -- A records require IPv4 addresses, AAAA records require IPv6 addresses, CNAME records require target hostnames, and MX records require mail server hostnames.

## Deploy

### Console

Open the deployment store, find **DNS Zone on DigitalOcean**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Simple Website** preset in the [Presets](#presets) tab to create a zone with a root A record.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: digitalocean.planton.dev/v1
kind: DigitalOceanDnsZone
metadata:
  name: example-com
  org: acme-corp
  env: prod
spec:
  domainName: example.com
  records:
    - name: "@"
      type: A
      values:
        - value: "203.0.113.10"
      ttlSeconds: 3600
```

```shell
planton apply -f do-dns-zone.yaml
```

This creates a DNS zone for `example.com` with a single A record pointing the root domain to the specified IP address. No MX, TXT, or other records are configured. A Stack Job tracks the provisioning in real time.

## Key Configuration

These are the most important decisions when configuring a DNS zone. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Domain name** -- The `domainName` field must be a valid fully-qualified domain name (e.g., `example.com`). After provisioning, update the domain's nameservers at your registrar to point to DigitalOcean's nameservers. DNS propagation can take up to 48 hours.

**Record types** -- Each record in the `records` list specifies a `type` (A, AAAA, CNAME, MX, TXT, SRV, CAA), a `name` (use `@` for the apex domain), one or more `values`, and a `ttlSeconds`. MX records require `priority`, SRV records require `priority`, `weight`, and `port`, and CAA records require `flags` and `tag`.

**TTL** -- The `ttlSeconds` field controls how long resolvers cache the record. Default is 3600 seconds (1 hour). Use shorter TTLs (300 seconds) during migrations for faster propagation, and longer TTLs (86400 seconds) for stable records to reduce lookup latency.

**ValueFromRef in record values** -- Record `values` accept ValueFromRef references, allowing DNS records to point to IP addresses or hostnames output by other Cloud Resources (e.g., a Droplet's `ipv4_address` or a load balancer's IP).

## Outputs and Dependencies

### What This Component Consumes

This component has no foreign key dependencies.

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `zone_name` | Domain name of the DNS zone | App Platform custom domain configuration |
| `zone_id` | Unique identifier of the DNS zone | API operations, record management |
| `name_servers` | Nameserver addresses for the zone | Domain registrar NS delegation |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Simple website zone** -- a DNS zone with a single root A record pointing to a web server or load balancer IP. Minimal configuration for getting a domain live quickly. Start from the **Simple Website** preset.

**Production zone with email** -- a DNS zone with A, MX, and TXT (SPF) records for a production website that also receives email. MX directs mail to your provider; SPF authorizes the provider to send on behalf of the domain. Start from the **Production With Email** preset.

## Works With

This component operates independently and does not reference other components.