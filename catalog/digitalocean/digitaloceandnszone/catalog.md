# DNS Zone on DigitalOcean

Deploys a DNS zone (domain) on DigitalOcean with optional inline DNS records covering every type the DigitalOcean API accepts -- A, AAAA, CNAME, MX, TXT, SRV, NS, CAA, and SOA -- plus the create-only apex-A convenience. DigitalOcean manages the authoritative nameservers for the zone, and record values can reference outputs from other Cloud Resources via ValueFromRef. Integrates with Planton's Provider Connections for DigitalOcean API token management.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **DigitalOcean Domain** -- a DNS zone served by DigitalOcean's nameservers (`ns1.digitalocean.com`, `ns2.digitalocean.com`, `ns3.digitalocean.com`); optionally seeded with an untracked apex A record via `ipAddress`
- **DNS Records** -- created only when `records` are provided; one DigitalOcean DNS record per value entry (multi-value records fan out), with type-specific fields for MX priority, SRV priority/weight/port, and CAA flags/tag enforced at validation time

## Before You Deploy

### Planton Setup

- **DigitalOcean Provider Connection** -- an active connection in the Connect module with a DigitalOcean API token. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline API token authentication.

### DigitalOcean Account

- **A domain you control** -- adding a zone does not require owning the domain (DigitalOcean hosts it immediately), but public resolution starts only after your registrar delegates to DigitalOcean's nameservers. Note domain names are unique across ALL DigitalOcean accounts.
- **IP addresses or hostnames for records** -- A records take IPv4 addresses, AAAA take IPv6, CNAME/MX/NS/SRV take target hostnames (author them fully qualified with a trailing dot -- that is how DigitalOcean stores them).

## Deploy

### Console

Open the deployment store, find **DNS Zone on DigitalOcean**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Simple Website** preset in the [Presets](#presets) tab to create a zone with apex and www records.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: digital-ocean.planton.dev/v1alpha1
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

This creates a DNS zone for `example.com` with a single A record pointing the apex at the specified IP address. A Stack Job tracks the provisioning in real time.

## Key Configuration

These are the most important decisions when configuring a DNS zone. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Domain name** -- The `domainName` field must be a valid fully-qualified domain name (e.g., `example.com`). After provisioning, update the domain's nameservers at your registrar to DigitalOcean's set (the `name_servers` output). DNS propagation can take up to 48 hours.

**Record types** -- Each record in the `records` list specifies a `type` (DigitalOcean accepts A, AAAA, CNAME, MX, TXT, SRV, NS, CAA, SOA; ALIAS and PTR are rejected at validation time), a `name` (use `@` for the apex), one or more `values` (each value becomes its own record -- two A values make round-robin), and an optional `ttlSeconds`. MX records require `priority`, SRV records require `priority`/`weight`/`port`, and CAA records require `flags`/`tag` -- all enforced before any provisioner runs.

**TTL** -- The `ttlSeconds` field controls how long resolvers cache the record; omit it to take DigitalOcean's default (1800 seconds). Use shorter TTLs (300) during migrations, longer (3600-86400) for stable records.

**ValueFromRef in record values** -- Record `values` accept ValueFromRef references, so records can point at outputs of other Cloud Resources (a Droplet's `ipv4_address`, a load balancer's IP) instead of hardcoded values.

**`ipAddress`** -- a create-only convenience that seeds an apex A record the platform never tracks afterwards. Prefer declaring the apex record in `records`; use `ipAddress` only when migrating a configuration that already relies on it.

## Outputs and Dependencies

### What This Component Consumes

This component has no foreign key dependencies (record values may optionally reference any resource's outputs).

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `zone_name` | Domain name of the DNS zone | DNS record `domain` references, App Platform custom domains |
| `zone_id` | The zone's resource identifier -- the domain name itself, not a UUID | API operations, imports |
| `name_servers` | DigitalOcean's fixed authoritative nameserver set | Domain registrar NS delegation |
| `urn` | The domain's uniform resource name (`do:domain:example.com`) | DigitalOcean project assignment, audit |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Simple website zone** -- apex A record plus a www CNAME. Minimal configuration for getting a domain live quickly. Start from the **Simple Website** preset.

**Production zone with email** -- website records plus MX pairs with priorities, an SPF policy, and CAA certificate authority pinning. Start from the **Production With Email** preset.

## Works With

- [**DNS Record on DigitalOcean**](/cloud-catalog/digital-ocean-dns-record) -- standalone records referencing this zone's `zone_name` output, for records owned by other teams or charts
