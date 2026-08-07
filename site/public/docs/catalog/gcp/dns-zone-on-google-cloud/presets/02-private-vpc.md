---
title: "Private VPC DNS Zone"
description: "Creates a private Cloud DNS managed zone visible only to resources on a VPC network. Use this for internal service discovery (`*.internal.example.com`) without exposing names to the public internet."
type: "preset"
rank: "02"
presetSlug: "02-private-vpc"
componentSlug: "dns-zone-on-google-cloud"
componentTitle: "DNS Zone on Google Cloud"
provider: "gcp"
icon: "package"
order: 2
---

# Private VPC DNS Zone

Creates a private Cloud DNS managed zone visible only to resources on a VPC network. Use this for internal service discovery (`*.internal.example.com`) without exposing names to the public internet.

## When to Use

- Microservice DNS inside a VPC
- GKE workloads that resolve internal hostnames via Cloud DNS private zones

## Key Configuration Choices

- **visibility: private** — required for private zones
- **privateVisibilityConfig.networks** — binds the zone to a GcpVpcNetwork via `network_self_link`
- **dnsName explicit** — set when the zone domain differs from metadata.name

## Placeholders to Replace

| Placeholder | Description |
|---|---|
| `my-gcp-project-123` | GCP project ID |
| `app-vpc` | Your GcpVpcNetwork resource name |
| `internal.example.com` | Private domain suffix |

## Related Presets

- **01-public-zone** — internet-facing authoritative zone
- **03-private-dnssec** — public zone with DNSSEC

## Related Components

- [GcpVpcNetwork](/docs/catalog/gcp/gcpvpcnetwork) — VPC that can see this zone
- [GcpDnsRecord](/docs/catalog/gcp/gcpdnsrecord) — records within the private zone
