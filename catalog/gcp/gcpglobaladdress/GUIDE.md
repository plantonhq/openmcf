# GcpGlobalAddress Guide

The judgment this guide protects: a global address is the anycast front
door of a global load balancer or the CIDR block private services live
in. Both roles are long-lived promises — recreation changes the IP or
tears the peering range out from under managed services.

## Two very different jobs, one resource

EXTERNAL global addresses are load-balancer VIPs: one IP, announced from
every Google edge. INTERNAL ones with `purpose: VPC_PEERING` are the CIDR
ranges Private Services Access carves out for Cloud SQL, AlloyDB, and
Memorystore. Confusing the two is the classic mistake — the VIP form
takes no network, the peering form REQUIRES a network and prefix.

## Size peering ranges for the decade

The `prefixLength` on a VPC_PEERING range is effectively permanent:
managed services allocate subnets inside it, and once instances exist the
range cannot shrink or move without rebuilding them. A /16 is the
conventional default; going smaller because "we only have two Cloud SQL
instances" is how the third instance fails to provision. Add a second
range later rather than resizing the first.

## The recreate trap

Everything except `labels` is ForceNew. For the VIP form, the replacement
IP breaks registrar DNS and every partner allow-list; for the peering
form, it strands the service networking connection built on top. A plan
that replaces this resource deserves the same scrutiny as deleting it.

## Teardown discipline

`PREVENT` is the right default for both roles once anything real points
at the address — destroy fails instead of releasing an anycast IP someone
else can claim, or tearing the floor out of a PSA range. `ABANDON` keeps
the reservation while dropping management. GCP refuses to release an
address in active use (a forwarding rule or service networking
connection), but PREVENT also covers the gap after those are destroyed.
