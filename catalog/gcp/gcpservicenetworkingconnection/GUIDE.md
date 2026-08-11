# GcpServiceNetworkingConnection Guide

The judgment this guide protects: this peering is a singleton with
gravity. One connection exists per (network, service) pair, everything
private-IP (Cloud SQL, AlloyDB, Memorystore PRIVATE_SERVICE_ACCESS,
Filestore) hangs off it, and both its capacity and its teardown follow
rules that surprise people.

## One connection per pair — widen, never duplicate

GCP rejects a second connection for the same (network, service) pair.
When the producer runs out of address space, the answer is ALWAYS to
append ranges to THIS resource's `reservedPeeringRanges` — growth is
additive and safe (already-provisioned producer subnets are untouched).
A second connection resource for the same pair fails its create.

## Size ranges generously up front

Producers carve service subnets out of your reserved ranges and cannot
use space that is too fragmented. A /16 is the standard default; a
too-small range surfaces LATER, as instance-creation failures when the
producer cannot allocate — not at connection time. Reserve big once
rather than appending slivers under incident pressure.

## The composition

The connection composes two first-class kinds: the `GcpVpcNetwork` being
peered and one or more `GcpGlobalAddress` resources (INTERNAL,
VPC_PEERING purpose) referenced BY NAME — not by self-link or CIDR. In a
chart: network → address range(s) → this connection → private-IP
instances, in that order.

## Adopting a pre-existing connection

When a connection for the pair already exists outside management (a
gcloud or console flow), the create fails with "Cannot modify allocated
ranges". `updateOnCreationFail: true` converts that failure into an
in-place update of the existing connection's ranges — a deliberate
adoption lever. Leave it false otherwise: silently adopting someone
else's peering is worse than failing loudly.

## Teardown ordering is the trap

GCP refuses to delete the connection while any producer still holds
subnets in the reserved ranges — destroy the Cloud SQL / AlloyDB /
Memorystore instances FIRST, then the connection, then the address
ranges. In chart teardown this ordering falls out of the dependency
graph; in manual cleanup it is the step people miss, and the destroy
failure ("producer still has allocations") is the system enforcing it.

`deletionPolicy` completes the picture: `PREVENT` for the connection
under a production database fleet; `ABANDON` removes it from management
while the peering keeps serving — the historical escape hatch for stuck
teardowns, at the price of an unmanaged peering that a later cleanup
will not know about.

## No cloud-side name

The connection is addressed by (network, service); `metadata.name` names
only the Planton resource, and the peering appears on the VPC named
after the service (e.g. `servicenetworking-googleapis-com`).
