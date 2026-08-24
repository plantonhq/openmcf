# Standby Failover IP

This preset reserves an IPv4 address WITHOUT assigning it -- capacity held for a failover or migration you have not executed yet. Be deliberate: this exact shape is the one that bills — a monthly charge, prorated hourly — until assigned or destroyed.

## When to Use

- Pre-allocating the public address for an environment before its droplets exist, so DNS can be published early
- Holding a known address through a droplet rebuild that spans manifests

## Key Configuration Choices

- **No `droplet`** -- the reservation stands alone; add the reference (or a literal droplet id) in a later apply to assign in place.
- **Explicit `ipVersion: ipv4`** -- stated even though it is the default, because the unassigned-billing warning is IPv4-specific (IPv6 reservations are free).

## What You Get

A held public address in the `reserved_ip_address` output -- and a running meter until you assign or destroy it. Do not let this preset outlive its purpose.
