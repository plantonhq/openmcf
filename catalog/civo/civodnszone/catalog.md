# DNS Zone on Civo

Creates a DNS zone (domain) on Civo's managed DNS service with optional inline DNS records for A, AAAA, CNAME, MX, and TXT record types. The zone and its records are provisioned as a single unit, and Civo assigns nameservers that must be configured at your domain registrar. Integrates with Planton's Provider Connections for Civo credential management.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Civo DNS Domain** -- a managed DNS zone for the specified domain name, assigned Civo nameservers (ns0.civo.com, ns1.civo.com, ns2.civo.com)
- **Civo DNS Domain Records** -- created only when `records` entries are specified in the spec; each record is provisioned with the configured type, name, value(s), and TTL

## Before You Deploy

### Planton Setup

- **Civo Provider Connection** -- an active connection in the Connect module with a Civo API token. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline API token authentication.

### Civo Account

- **A registered domain name** -- you must own the domain and have access to its registrar to update nameserver records to point to Civo's nameservers after zone creation.
- **Target IP addresses or hostnames** for any inline DNS records -- these can be literal values or ValueFromRef references to other Cloud Resources (e.g., a CivoComputeInstance or CivoIpAddress).

## Deploy

### Console

Open the deployment store, find **DNS Zone on Civo**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Simple Website** preset in the [Presets](#presets) tab for a domain with apex A record and www CNAME.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: civo.planton.dev/v1
kind: CivoDnsZone
metadata:
  name: app-domain
  org: acme-corp
  env: prod
spec:
  domainName: example.com
  records:
    - name: "@"
      type: A
      values:
        - value: "192.0.2.1"
      ttlSeconds: 3600
    - name: www
      type: CNAME
      values:
        - value: "example.com"
      ttlSeconds: 3600
```

```shell
planton apply -f civo-dns-zone.yaml
```

This creates a DNS zone for `example.com` with an apex A record and a www CNAME pointing to the apex. After provisioning, update your domain registrar's nameservers to point to the Civo nameservers returned in the outputs. A Stack Job tracks the provisioning in real time.

## Key Configuration

These are the most important decisions when configuring a Civo DNS zone. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Domain name** -- The `domainName` field must be a valid fully-qualified domain name (e.g., `"example.com"`). This becomes the zone apex. Subdomains are managed as records within the zone.

**Inline records** -- The `records` list defines DNS records created alongside the zone. Each record requires a `name` (use `"@"` for the apex), `type` (A, AAAA, CNAME, MX, TXT), and at least one entry in `values`. Records can use ValueFromRef to reference outputs from other Cloud Resources, such as IP addresses from a CivoIpAddress.

**TTL** -- Each record's `ttlSeconds` controls resolver caching duration. Defaults to 3600 seconds (1 hour) when not specified. Use shorter TTLs during DNS migrations and longer TTLs for stable records to reduce query volume.

**Nameserver delegation** -- After creating the zone, Civo assigns nameservers (ns0.civo.com, ns1.civo.com, ns2.civo.com). Update your domain registrar's NS records to point to these nameservers. DNS propagation typically completes within 24--48 hours.

## Outputs and Dependencies

### What This Component Consumes

This component has no foreign key dependencies.

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `zone_name` | Domain name of the DNS zone | Application configuration, certificate requests |
| `zone_id` | Unique identifier of the DNS zone on Civo | CivoDnsRecord zone references, API operations |
| `name_servers` | Nameserver addresses assigned to the zone | Domain registrar NS record configuration |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Simple website** -- an apex A record pointing to a server IP and a www CNAME aliasing to the apex. Covers the standard website configuration where both `example.com` and `www.example.com` resolve to the same server. Start from the **Simple Website** preset.

**Website with email** -- extends the simple website pattern with an MX record for mail routing and a TXT record for SPF email authentication. Covers the majority of small business and SaaS domain configurations. Start from the **With Email** preset.

## Works With

This component operates independently and does not reference other deployment components.