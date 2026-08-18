# AwsCloudMapNamespace — Operational Guide

Live-earned judgment lands here as proof runs and adopter operations teach it; the notes below are the forge-time seed.

## Declared instances and runtime registrations do not mix

ECS service discovery registers and deregisters task instances at runtime. Declare static instances only on services the manifest fully owns — a declared instance on an ECS-managed service invites two owners for one registry, and force_destroy on such a service deletes ECS's registrations wholesale.

## Deregistration has no forgiveness

Deregistering an already-gone instance ERRORS at the provider (no NotFound tolerance upstream). Destroy instances through the module, never out-of-band; an out-of-band deregistration leaves the next declarative destroy red until state is reconciled.

## The private namespace's VPC is write-only

AWS never returns which VPC a private namespace binds — the provider carries it in config and imports need the `{namespace_id}:{vpc_id}` composite. Losing the VPC id means reading it from the namespace's hosted zone associations, not from Cloud Map.

## Alias registrations are all-or-nothing

`alias_dns_name` (a Route 53 ALIAS to an ELB) cannot combine with ANY other instance attribute — AWS rejects the mix. One instance is an alias, or it is an address; never both.

## SRV records change the discovery contract

A service publishing SRV returns hostname+port tuples, and clients must ask for SRV (most naive resolvers ask A). Use SRV only when consumers understand it (gRPC clients, service meshes); A + a fixed port covers the common case with zero client changes.

## Custom health starts UNHEALTHY-free but needs a heartbeat

`health_check_custom_config` services mark instances healthy until told otherwise, but the workload must push UpdateInstanceCustomHealthStatus to ever mark one unhealthy — silence means permanently healthy, not monitored.
