---
title: "Reserved Capacity for a Traffic Event"
description: "This preset pre-provisions an internet-facing ALB for a known traffic surge -- a product launch, a ticket sale, a marketing spike -- instead of waiting for the load balancer to scale organically. It..."
type: "preset"
rank: "03"
presetSlug: "03-reserved-capacity"
componentSlug: "alb"
componentTitle: "ALB"
provider: "aws"
icon: "package"
order: 3
---

# Reserved Capacity for a Traffic Event

This preset pre-provisions an internet-facing ALB for a known traffic surge
-- a product launch, a ticket sale, a marketing spike -- instead of waiting
for the load balancer to scale organically. It combines the two
launch-planning knobs: a Load Balancer Capacity Unit (LCU) reservation and
BYOIP addressing from a VPC IPAM pool, so the event traffic arrives on
pre-warmed capacity at addresses your organization owns and has allowlisted.

## When to Use

- A dated traffic event whose peak you can estimate (launches, sales,
  broadcasts) -- organic ALB scaling reacts in minutes; a reservation is
  ready at second zero
- Workloads whose clients allowlist source ranges: the IPAM pool keeps the
  ALB's public addresses inside your own BYOIP range across rebuilds
- Replacing the old "pre-warming ticket to AWS support" workflow with a
  declarative setting

## Key Configuration Choices

- **`minimumLoadBalancerCapacityUnits: 500`** -- the reserved LCU floor.
  Size it from the event estimate (one LCU is roughly 25 new conn/s, 3,000
  active conn, 2.22 Mbps processed bytes -- whichever dimension peaks
  first). **The reservation BILLS while set, on top of normal usage: set it
  for the event window, then remove the field to release it.** Minimums and
  maximums depend on the account's service quotas.
- **`ipv4IpamPoolId`** -- the ALB draws its public IPv4 addresses from your
  VPC IPAM pool (BYOIP) instead of AWS's ranges. Supply a literal pool id;
  removing the field moves the ALB back to AWS-assigned addresses in place.
  Drop this field if you only need the capacity reservation.
- **`deleteProtectionEnabled: true`** -- an event ALB is exactly the one
  nobody should delete mid-event

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<alb-name>` | Name for the ALB (max 32 chars after truncation) | Your naming convention |
| `<aws-region>` | AWS region code (e.g., `us-east-1`) | Your deployment region |
| `<public-subnet-id-az1>` / `<public-subnet-id-az2>` | Public subnets in two AZs | Your VPC's subnet list |
| `<alb-security-group-id>` | Security group opening the listener ports | Your VPC's security groups |
| `<ipam-pool-id>` | The public-scope VPC IPAM pool holding your BYOIP range | VPC console → IPAM → Pools |

## Common Additions

- Add `dns` (zone + hostnames) to point your event domain at the ALB, as in
  **01-internet-facing**
- Add `accessLogs` to capture the event's request log in S3
- Set `zonalShiftEnabled: true` so Application Recovery Controller can pull
  the event's traffic out of an impaired AZ

## Related Presets

- **01-internet-facing** -- the everyday public ALB with DNS
- **02-internal-hardened** -- the private, locked-down shape
