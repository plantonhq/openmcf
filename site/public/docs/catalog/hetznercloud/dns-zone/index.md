---
title: "DNS Zone"
description: "DNS Zone deployment documentation"
icon: "package"
order: 100
componentName: "hetznerclouddnszone"
---

# Hetzner Cloud DNS Zone

Deploys a DNS zone on Hetzner Cloud's authoritative nameservers with configurable record sets. Supports primary mode (records managed directly through this component) and secondary mode (records synchronized from an external primary nameserver via AXFR/IXFR zone transfer). Record sets group DNS records by (name, type) pair using `hcloud_zone_rrset`, allowing in-place value updates without destroying the record set. Record values support ValueFromRef for cross-resource wiring in InfraChart deployments.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **DNS Zone** -- an `hcloud_zone` resource representing the domain on Hetzner Cloud's nameservers
- **Record Sets** -- one `hcloud_zone_rrset` per entry in the record sets list, each managing all records for a unique (name, type) pair

## Before You Deploy

### Hetzner Cloud Account

- **A Hetzner Cloud account** with an active project and DNS API access.
- **A domain name** you own or control -- you will need to update the domain's NS records at your registrar to point to Hetzner's nameservers after zone creation.

### For Secondary Mode

- **An external primary nameserver** with zone transfer (AXFR/IXFR) enabled for Hetzner Cloud's nameserver IPs.
- **TSIG credentials** (optional) if the primary nameserver requires authenticated zone transfers.

## Deploy

### Console

Open the deployment store, find **Hetzner Cloud DNS Zone**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields including domain name, mode, and record sets.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: hetzner-cloud.planton.dev/v1
kind: HetznerCloudDnsZone
metadata:
  name: example-zone
  org: acme-corp
  env: prod
spec:
  domainName: "example.com"
  mode: primary
  ttl: 3600
  recordSets:
    - name: "@"
      type: A
      records:
        - value:
            value: "93.184.216.34"
    - name: "www"
      type: CNAME
      ttl: 300
      records:
        - value:
            value: "example.com."
```

```shell
planton apply -f hetznercloud-dns-zone.yaml
```

This creates a DNS zone for example.com with an A record at the apex and a CNAME for www. A Stack Job tracks the provisioning in real time. Update your domain's NS records at your registrar to the nameservers returned in the outputs.

### InfraChart

When deploying as part of a web application, use ValueFromRef to wire DNS records to a load balancer's public IP:

```yaml
spec:
  domainName: "example.com"
  mode: primary
  recordSets:
    - name: "@"
      type: A
      records:
        - value:
            valueFrom:
              kind: HetznerCloudLoadBalancer
              name: web-lb
              fieldPath: status.outputs.ipv4_address
    - name: "www"
      type: CNAME
      records:
        - value:
            value: "example.com."
```

The InfraPipeline resolves the dependency graph, provisions the load balancer first, then creates DNS records pointing to its IP.

## Key Configuration

These are the most important decisions when configuring a DNS zone. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Domain name** -- The `domainName` field specifies the DNS domain (e.g., "example.com"). Changing forces replacement.

**Mode** -- The `mode` field selects `primary` (records managed here) or `secondary` (records pulled from an external primary). In secondary mode, `recordSets` must be empty and `primaryNameservers` is required.

**Default TTL** -- The `ttl` field sets the default Time To Live for records (default: 3600 seconds). Individual record sets can override with their own TTL.

**Record sets** -- The `recordSets` field defines DNS records grouped by (name, type). Each record set specifies a `name` (relative to the zone, "@" for apex), `type` (A, AAAA, CNAME, MX, TXT, etc.), optional `ttl` override, and one or more record `values`. Values support ValueFromRef for cross-resource wiring.

## Outputs and Dependencies

### What This Component Consumes

Record values can reference any component output via ValueFromRef. Common patterns:

| Dependency | Record Type | ValueFromRef Path |
|------------|------------|-------------------|
| **HetznerCloudServer** | A | `status.outputs.ipv4_address` |
| **HetznerCloudLoadBalancer** | A | `status.outputs.ipv4_address` |
| **HetznerCloudFloatingIp** | A | `status.outputs.ip_address` |
| **HetznerCloudPrimaryIp** | A | `status.outputs.ip_address` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `zone_id` | Hetzner Cloud numeric ID of the zone | API operations, resource tracking |
| `nameservers` | Authoritative nameservers assigned to the zone | Domain registrar NS record configuration |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Primary zone with A records** -- A primary zone with apex and www records pointing to server or load balancer IPs. The standard starting point for hosting a domain on Hetzner Cloud.

**Secondary zone** -- A secondary zone that synchronizes records from an external primary nameserver, providing redundancy and geographic distribution.

## Works With

- [**Hetzner Cloud Server**](/cloud-catalog/hetznercloud-server) -- A records point to server IPv4 addresses
- [**Hetzner Cloud Load Balancer**](/cloud-catalog/hetznercloud-load-balancer) -- A records point to LB IPv4 addresses
- [**Hetzner Cloud Primary IP**](/cloud-catalog/hetznercloud-primary-ip) -- A records point to managed IP addresses
- [**Hetzner Cloud Floating IP**](/cloud-catalog/hetznercloud-floating-ip) -- A records point to failover IP addresses
- [**Hetzner Cloud Certificate**](/cloud-catalog/hetznercloud-certificate) -- managed certificates require DNS records pointing to a load balancer
