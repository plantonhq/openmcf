---
title: "Proxy-Only Subnet (Regional Managed Proxy)"
description: "The address space GCP's Envoy-based regional load balancers allocate their proxies from. Every VPC region that hosts a regional (internal or external) Application Load Balancer must have exactly one..."
type: "preset"
rank: "03"
presetSlug: "03-proxy-only"
componentSlug: "subnetwork"
componentTitle: "Subnetwork"
provider: "gcp"
icon: "package"
order: 3
---

# Proxy-Only Subnet (Regional Managed Proxy)

The address space GCP's Envoy-based regional load balancers allocate their proxies from. Every VPC region that hosts a regional (internal or external) Application Load Balancer must have exactly one ACTIVE proxy-only subnet — creating it is the prerequisite step load-balancer sessions trip over.

## When to Use

- Before creating any regional internal or regional external Application Load Balancer in a VPC region
- When staging a proxy address-space migration (create a `role: BACKUP` twin, then promote it)

## Remix Notes

- Size at least /23 (Google's recommendation) — the proxy fleet scales with load-balancer traffic, and this range cannot be shared with workloads.
- Workloads never live here: no secondary ranges, no Private Google Access needed — the purpose locks the subnet to proxy use.
- One ACTIVE proxy-only subnet per region per VPC; a second one must be `role: BACKUP` until promoted.
- `GLOBAL_MANAGED_PROXY` is the cross-region equivalent for cross-region internal ALBs — change `purpose` if that is the target.
