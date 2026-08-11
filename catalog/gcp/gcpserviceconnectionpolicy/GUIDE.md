# GcpServiceConnectionPolicy Guide

The judgment this guide protects: the policy is invisible until it is
missing. PSC-first managed services (Memorystore for Valkey, Redis
Cluster, and the producers that follow them) cannot place endpoints in
your network without one — and instance creation fails with a
connectivity error that says nothing about policies.

## Deploy the policy before the first instance

One policy per (network, service class, region) triple, created BEFORE
any instance of that class in that region. In charts, make the producer
instance depend on the policy; the error you avoid ("connectivity error"
at instance create) is among GCP's least self-explanatory.

## The service class is the producer's published name

`serviceClass` is not free text: Google publishes one identifier per
service (`gcp-memorystore` for Memorystore for Valkey,
`gcp-memorystore-redis` for Redis Cluster); third-party producers
publish their own. A typo creates a valid policy that authorizes
nothing — the instance-create failure looks identical to the
missing-policy case.

## Subnets are ordinary; the limit is a guardrail

`pscConfig.subnetworks` accepts regular-purpose subnets in the policy's
region — no special PSC purpose (unlike PSC NAT subnets for published
services). `pscConfig.limit` caps how many endpoints the automation may
create: in shared networks set it deliberately, so a runaway
instance-creation loop hits the cap instead of exhausting subnet space.

## The producer allowlist moves as a pair

`producerInstanceLocation: CUSTOM_RESOURCE_HIERARCHY_LEVELS` and
`allowedGoogleProducersResourceHierarchyLevels` only work together —
validation enforces both-or-neither, because an allowlist under the
default location mode is silently ignored by GCP.

## Mutation map

`location`, `network`, `serviceClass`, and the name are ForceNew.
Everything inside `pscConfig` (plus description and labels) updates in
place — growing subnets and raising limits never recreates the policy or
disturbs existing endpoints.

## Teardown discipline

Deleting the policy strands every PSC endpoint created under it and
blocks new instances of the class in that region — while the producer
instances keep running and failing. `deletionPolicy: PREVENT` is the
right posture wherever managed instances depend on the policy; `ABANDON`
leaves it authorizing connections unmanaged.
