---
title: "DNS Record"
description: "DNS Record deployment documentation"
icon: "package"
order: 100
componentName: "aliclouddnsrecord"
---

# AliCloud DNS Record

Deploys a DNS record on Alibaba Cloud's Alidns service within an existing hosted zone. The record maps a host record (subdomain) to a value such as an IP address, CNAME alias, or mail server. The parent domain must already exist in Alidns -- either managed by the AliCloudDnsZone component or added manually.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **DNS Record** -- an `alicloud_dns_record` resource in the specified domain with configurable type, value, TTL, priority, and resolution line

## Before You Deploy

### Planton Setup

- **AliCloud Provider Connection** -- an active connection in the Connect module with credentials for the target Alibaba Cloud account. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery.

### Alibaba Cloud Account

- **An existing DNS zone** -- the `domainName` must already be registered in Alidns (create one with AliCloudDnsZone).
- **NS delegation** -- ensure your registrar's NS records point to the Alidns nameservers for the record to resolve publicly.

## Deploy

### Console

Open the deployment store, find **AliCloud DNS Record**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields including domain, host record, type, and value.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: alicloud.planton.dev/v1
kind: AliCloudDnsRecord
metadata:
  name: api-record
  org: acme-corp
  env: prod
spec:
  region: cn-hangzhou
  domainName: example.com
  rr: api
  type: A
  value: "47.100.1.1"
  ttl: 600
```

```shell
planton apply -f alicloud-dns-record.yaml
```

This creates an A record for api.example.com pointing to 47.100.1.1. A Stack Job tracks the provisioning in real time.

### InfraChart

DNS records are typically the last resource in a dependency chain. Use ValueFromRef to wire the record value to an upstream resource's output (e.g., an ALB's DNS name):

```yaml
spec:
  region: cn-hangzhou
  domainName: example.com
  rr: www
  type: CNAME
  value: "alb-abc123.cn-hangzhou.alb.aliyuncs.com"
```

## Key Configuration

These are the most important decisions when configuring a DNS record. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Record type** -- The `type` field determines what the record maps to: A (IPv4), AAAA (IPv6), CNAME (alias), MX (mail), TXT (verification/SPF), NS (delegation), SRV (service), CAA (certificate authority).

**TTL** -- The `ttl` field controls how long resolvers cache this record in seconds. Lower values (60-600) allow faster propagation of changes. Higher values (3600-86400) reduce DNS query load.

**Priority** -- The `priority` field is required for MX records (1 = highest). Ignored for other types.

**Resolution line** -- The `line` field enables intelligent DNS routing by ISP or geography. Use "default" for standard resolution.

## Outputs and Dependencies

### What This Component Provides

After provisioning, `status.outputs` contains:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `record_id` | Record ID assigned by Alibaba Cloud | Record management |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**A record** -- Maps a subdomain to an IPv4 address. Start from the **A Record** preset.

**CNAME record** -- Aliases a subdomain to another domain (e.g., ALB or CDN endpoint). Start from the **CNAME Record** preset.

**MX record** -- Configures mail delivery with priority routing. Start from the **MX Record** preset.

## Works With

- [**AliCloud DNS Zone**](/cloud-catalog/ali-cloud-dns-zone) -- the parent domain this record belongs to
- [**AliCloud Application Load Balancer**](/cloud-catalog/ali-cloud-application-load-balancer) -- common CNAME target (ALB dns_name)
- [**AliCloud Network Load Balancer**](/cloud-catalog/ali-cloud-network-load-balancer) -- common CNAME target (NLB dns_name)
- [**AliCloud EIP Address**](/cloud-catalog/ali-cloud-eip-address) -- common A record target (EIP ip_address)
