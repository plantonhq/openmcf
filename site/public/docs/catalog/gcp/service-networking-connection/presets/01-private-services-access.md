---
title: "Private Services Access for Managed Services"
description: "The canonical private services access setup: peer a VPC with Google's managed-services producer (`servicenetworking.googleapis.com`) so Cloud SQL, AlloyDB, Memorystore (PRIVATE_SERVICE_ACCESS mode),..."
type: "preset"
rank: "01"
presetSlug: "01-private-services-access"
componentSlug: "service-networking-connection"
componentTitle: "Service Networking Connection"
provider: "gcp"
icon: "package"
order: 1
---

# Private Services Access for Managed Services

The canonical private services access setup: peer a VPC with Google's managed-services producer (`servicenetworking.googleapis.com`) so Cloud SQL, AlloyDB, Memorystore (PRIVATE_SERVICE_ACCESS mode), and Filestore can hand out private IPs inside the network.

## What this creates

One connection between the referenced VPC and the Google producer, backed by one reserved `VPC_PEERING` address range. After it deploys, private-IP managed services in this network start working — without it, a Cloud SQL instance with `privateIpEnabled` fails to create.

## The composition

This preset is the middle node of a three-resource chain:

1. **GcpVpcNetwork** — the network (`prod-vpc` here).
2. **GcpGlobalAddress** — an INTERNAL address with purpose `VPC_PEERING` and a prefix length (16 reserves a /16, the common default; producers carve every service subnet out of this space, so size generously).
3. **GcpServiceNetworkingConnection** (this preset) — peers the two.

## Remix ideas

- Add a second range to `reservedPeeringRanges` later when the producer runs out of space — an in-place update that never disturbs existing service subnets.
- Point `network` at a literal self-link to attach to a VPC managed outside Planton.
