---
title: "Growing Producer Capacity with Multiple Ranges"
description: "The expansion pattern: a connection carrying its original reserved range plus a second one appended later, for when managed-service instances exhausted the first allocation."
type: "preset"
rank: "02"
presetSlug: "02-multi-range-growth"
componentSlug: "service-networking-connection"
componentTitle: "Service Networking Connection"
provider: "gcp"
icon: "package"
order: 2
---

# Growing Producer Capacity with Multiple Ranges

The expansion pattern: a connection carrying its original reserved range plus a second one appended later, for when managed-service instances exhausted the first allocation.

## Why this shape

One connection exists per (network, service) pair — GCP rejects a second connection, so more capacity never means another connection. Instead, reserve a new `VPC_PEERING` global address range and append it to `reservedPeeringRanges` on the existing resource. The update applies in place: every service subnet the producer already carved from the original range keeps working.

## When you need it

- Cloud SQL / AlloyDB instance creation starts failing with range-exhaustion errors.
- A new service family (e.g. Filestore) needs an allocation the original /20 cannot fit.

## Remix ideas

- Prevent the exhaustion entirely: size the first range at /16 (the producer default request is /24 per service subnet, and instances multiply).
- Keep expansion ranges adjacent in your IP plan so network ACL reasoning stays simple.
