---
title: "Weighted Canary"
description: "This preset creates an A record with a weighted round-robin routing policy — 95% of DNS answers point at the stable target, 5% at the canary. Shifting the weights progresses the rollout without..."
type: "preset"
rank: "03"
presetSlug: "03-weighted-canary"
componentSlug: "dns-record-on-google-cloud"
componentTitle: "DNS Record on Google Cloud"
provider: "gcp"
icon: "package"
order: 3
---

# Weighted Canary

This preset creates an A record with a weighted round-robin routing policy —
95% of DNS answers point at the stable target, 5% at the canary. Shifting
the weights progresses the rollout without touching consumers.

## When to Use

- Canary rollouts of a new backend or region behind one hostname
- A/B traffic splitting at the DNS layer
- Staging a target at weight 0 before shifting traffic onto it

## Key Configuration Choices

- **Weighted round robin (`routingPolicy.wrr`)** — traffic splits by weight
  ratio; the weights need not sum to 100.
- **60-second TTL** — weight changes take effect quickly; raise it after
  the rollout settles.
- **Static values per entry** — entries can instead use
  `healthCheckedTargets` to withdraw unhealthy targets automatically.

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<gcp-project-id>` | GCP project ID | `GcpProject` outputs |
| `<dns-zone-name>` | Name of the Cloud DNS managed zone | `GcpDnsZone` status outputs |
| `<stable-ipv4-address>` | Current production target IP | Load balancer or VM external IP |
| `<canary-ipv4-address>` | New target receiving the canary slice | Load balancer or VM external IP |

The sample FQDN `api.example.com.` is a realistic placeholder for the
pattern-validated `name` field — replace it with your hostname, keeping the
trailing dot.

## Related Presets

- **01-a-record** — plain static A record without traffic steering
- **02-cname-record** — alias to another hostname
