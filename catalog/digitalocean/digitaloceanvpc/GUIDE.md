# DigitalOcean VPC -- Operational Guide

Judgment calls that matter when you run private networks on DigitalOcean.

## Plan the range once, or let DigitalOcean plan it forever

`ipRangeCidr` is immutable: the only way to change it is replacing the VPC, which means evacuating every Droplet, cluster, load balancer, and database inside it first. Decide up front which world you are in. If nothing will ever peer with this network or VPN into it, omit the range — DigitalOcean picks a non-conflicting block and the `ip_range` output tells you what it chose. The moment corporate ranges, site-to-site VPNs, or VPC peering enter the picture, choose the range yourself and record the allocation, because DigitalOcean's auto-assigned 10.x blocks will eventually collide with somebody's office network.

## Sizing inside /16–/24

DigitalOcean accepts prefixes from /16 (65,536 addresses) down to /24 (256). A /24 sounds roomy for a dozen Droplets, but managed resources quietly consume addresses too — every DOKS node, load balancer, and database cluster member takes one. Kubernetes clusters are the heavy consumer: autoscaling to twenty nodes eats twenty addresses. A /20 (4,096) is a comfortable default for an environment; reserve /24s for genuinely small, fixed-size networks.

## You cannot make a VPC the region's default — and should not want to

Each region has a default VPC that DigitalOcean manages; resources created without an explicit network land there. That flag is computed, not settable — no manifest, Terraform config, or API write flips it. Treat the default VPC as the untyped landing zone and this kind's VPCs as the deliberate ones: always wire the `vpc` reference on Droplets, clusters, load balancers, and databases explicitly, and membership never depends on which VPC happens to be the regional default.

## One region, no bridges — plan for it

A VPC exists in exactly one region, and members can only join from that region. Cross-region private connectivity does not come from this kind: it comes from VPC peering (a separate resource DigitalOcean offers) or from routing over the public network with TLS. If a workload spans regions, design the split now — one VPC per region with non-overlapping ranges, so future peering stays possible.

## Destroy members first; the VPC goes last

DigitalOcean refuses to delete a VPC that still contains resources, and the module's delete retries only paper over short races (a Droplet mid-destroy), not real membership. Tear environments down in dependency order — workloads, then load balancers and databases, then the VPC. The same applies in reverse for creation, which is why other kinds' E2E lanes install this VPC as their first fixture.
