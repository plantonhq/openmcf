---
title: "Static-IP Internet-Facing NLB"
description: "This preset creates the headline NLB use case: an internet-facing load balancer whose public IPs never change. Each subnet mapping pins one NLB node to an Elastic IP, referenced from an..."
type: "preset"
rank: "02"
presetSlug: "02-static-ip-internet-facing"
componentSlug: "nlb"
componentTitle: "NLB"
provider: "aws"
icon: "package"
order: 2
---

# Static-IP Internet-Facing NLB

This preset creates the headline NLB use case: an internet-facing load
balancer whose public IPs never change. Each subnet mapping pins one NLB
node to an Elastic IP, referenced from an `AwsElasticIp` resource via
`valueFrom` — the addresses survive scaling events, node replacements, and
AWS maintenance, so partners and firewalls can allowlist them permanently.
Attach `AwsLbListener` resources against its `load_balancer_arn` output for
ports; the destinations are `AwsLbTargetGroup` resources.

## When to Use

- Partner integrations, corporate firewalls, or legacy systems that
  allowlist by IP address rather than hostname
- DNS pinning: A records pointing at addresses that never move
- Any internet-facing Layer-4 entry point (TCP/UDP/TLS) — including fronting
  an ALB via a target group with `targetType: alb` when you need static IPs
  *and* Layer-7 routing

## Key Configuration Choices

- **Elastic IPs by reference** — `allocationId` resolves from each
  `AwsElasticIp` resource's `allocation_id` output at deploy time, keeping
  the IP's lifecycle independent of (and longer than) the NLB's
- **Two mappings across AZs** — one node and one static IP per Availability
  Zone; subnet mappings are add-only in AWS, so this pair is the durable
  minimum for HA
- **Deletion protection** (`deleteProtectionEnabled: true`) — deleting an
  NLB orphans its listeners and leaves the Elastic IPs billing as
  unattached; production entry points should refuse casual deletion
- **No security groups** — attaching any is a one-way door (the last one can
  never be removed); add them only when committed to operating them

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<nlb-name>` | Unique name for the NLB (AWS caps it at 32 characters) | Choose a descriptive name (e.g., `partner-ingress`) |
| `<aws-region>` | AWS region code (e.g., `us-east-1`) | Your deployment region |
| `<public-subnet-id-az1>` | Public subnet in the first Availability Zone | AWS VPC console or `AwsSubnet` status outputs |
| `<public-subnet-id-az2>` | Public subnet in the second Availability Zone | AWS VPC console or `AwsSubnet` status outputs |
| `<elastic-ip-resource-name-az1>` | `AwsElasticIp` resource for the first AZ's static IP | Your AwsElasticIp manifest's `metadata.name` |
| `<elastic-ip-resource-name-az2>` | `AwsElasticIp` resource for the second AZ's static IP | Your AwsElasticIp manifest's `metadata.name` |

## Common Additions

- Add `dns` with a Route53 zone and hostnames — alias records give clients a
  name while the static IPs remain the allowlisting contract
- Add `crossZoneLoadBalancingEnabled: true` if targets are unevenly
  distributed across the two AZs
- Add `zonalShiftEnabled: true` to let Amazon Application Recovery Controller
  drain an impaired AZ
- Set `ipAddressType: dualstack` to serve IPv6 clients alongside the static
  IPv4 addresses

## Related Presets

- **01-internal** — the plain VPC-internal variant
- **03-private-link-hardened** — security groups, PrivateLink enforcement,
  and access logs
