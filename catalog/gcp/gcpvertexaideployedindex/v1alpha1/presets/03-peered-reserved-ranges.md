# Peered Endpoint with Reserved Ranges and JWT Auth

Deploys onto a VPC-peered endpoint with the deployment pinned to
reserved IP ranges, access logging on, and JWT authentication on the
private query endpoint — the full private-serving posture.

## What this preset creates

A DeployedIndex on the referenced peered endpoint, its serving nodes
drawing IPs from the named reserved `VPC_PEERING` global address,
accepting only JWTs signed by the allowed issuer for the listed
audience. Private-endpoint access logs flow to Cloud Logging. The
private query address lands in the `match_grpc_address` output; there
is no public domain on a peered endpoint.

## When to use

- Serving inside a private VPC where the IP space must be predictable
  (firewall rules, on-prem routing over the peering)
- Deployments that need caller authentication beyond network
  reachability
- Environments with audit requirements on query-endpoint access

## Remix ideas

- Reserved ranges are `VPC_PEERING` global addresses referenced by
  name — the same reservations a `GcpServiceNetworkingConnection`
  composes; the FK targets `GcpGlobalAddress.status.outputs.name`.
- Plan deployment groups and range sets together: once `prod` ships
  with `[vertex-ai-range-a]`, the API permanently pairs group `prod`
  with exactly that set.
- JWT issuers are service-account emails — clients present JWTs signed
  by that account with an `aud` matching `audiences`.
