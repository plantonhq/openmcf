---
title: "DNS Zone"
description: "DNS Zone deployment documentation"
icon: "package"
order: 100
componentName: "scalewaydnszone"
---

# Scaleway DNS Zone

Deploys a DNS zone on Scaleway with optional inline DNS records as a composite resource. A zone represents a delegated portion of the DNS namespace for a domain you own. Supports both root zones (e.g., `example.com`) and subdomain zones (e.g., `staging.example.com`) for environment-scoped DNS delegation. Inline records can reference other resources' outputs via ValueFromRef for dynamic values like Load Balancer IPs, while standalone ScalewayDnsRecord resources provide explicit DAG edges for cross-resource wiring in InfraCharts.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **DNS Zone** -- a `scaleway_domain_zone` resource for the configured domain and optional subdomain prefix, with Scaleway-assigned nameservers
- **Inline DNS Records** -- created only when `records` entries are provided; one `scaleway_domain_record` per entry for static records like MX, TXT (SPF/DKIM), CAA, and NS

Note: Scaleway DNS zones and records do not support tags. Unlike most other Scaleway resources, the DNS API does not accept tags or labels.

## Before You Deploy

### Scaleway Account

- **A Scaleway account** with an active project and API access key pair.
- **A registered domain name** at a domain registrar (Namecheap, Google Domains, etc.). Scaleway does not perform domain registration.
- **Registrar access** to configure nameservers. After zone creation, delegate the domain to the nameservers from `status.outputs.nameServers` at your registrar. DNS queries will not resolve through Scaleway until delegation is complete.

## Deploy

### Console

Open the deployment store, find **Scaleway DNS Zone**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Root Zone** preset in the [Presets](#presets) tab to create a zone for your domain with no inline records.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: scaleway.planton.dev/v1
kind: ScalewayDnsZone
metadata:
  name: main-zone
  org: acme-corp
  env: prod
spec:
  domain: example.com
```

```shell
planton apply -f scaleway-dns-zone.yaml
```

This creates a root DNS zone for `example.com` with no inline records. Records are managed separately as standalone ScalewayDnsRecord resources. A Stack Job tracks the provisioning in real time. After deployment, configure the nameservers from `status.outputs.nameServers` at your domain registrar.

## Key Configuration

These are the most important decisions when configuring a DNS zone. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Domain** -- The `domain` field specifies the registered parent domain (e.g., `example.com`). Cannot be changed after creation. The domain must already exist at a registrar.

**Subdomain** -- The `subdomain` field creates a delegated subzone (e.g., `subdomain: staging` creates `staging.example.com`). Leave empty for the root zone. Useful for environment-scoped DNS management where each environment has its own zone and nameservers.

**Inline vs standalone records** -- Use the `records` field for static records known at zone creation time (MX, SPF/TXT, CAA). Use standalone ScalewayDnsRecord resources for records that reference other infrastructure outputs (A records to Load Balancer IPs, CNAMEs to Kapsule endpoints) -- these create explicit dependency edges in InfraChart DAGs.

## Outputs and Dependencies

### What This Component Consumes

This component has no foreign key dependencies.

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `zone_name` | Computed zone name (e.g., `example.com` or `staging.example.com`) | ScalewayDnsRecord `zoneName` field for record placement |
| `name_servers` | Nameservers assigned by Scaleway for this zone | Domain registrar NS delegation configuration |
| `name_servers_default` | Default nameservers assigned by Scaleway DNS infrastructure | Verification and troubleshooting |
| `name_servers_master` | Primary nameserver(s) for the zone | Secondary/slave DNS configuration |
| `status` | Current zone state (e.g., `active`, `pending`, `error`) | Health monitoring, deployment verification |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Root zone with standalone records** -- An empty root zone for a domain with all records managed as separate ScalewayDnsRecord resources. The recommended production pattern because it creates explicit dependency edges in InfraChart DAGs and allows records to reference outputs from Load Balancers, Instances, and Kapsule clusters. Start from the **Root Zone** preset.

## Works With

This component operates independently and does not reference other components.