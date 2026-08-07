# VPC-Peered Index Endpoint

A private serving surface reachable only inside a peered VPC — vector
search that never touches the public internet.

## What this preset creates

An index endpoint named `Private Search` in `us-central1`, peered into
the `ml-vpc` network via Private Services Access. Deployed indexes are
queried over private gRPC addresses (the deployed index's
`match_grpc_address` output) from inside the VPC.

## When to use

- Compliance regimes that forbid public data-plane endpoints
- Latency-sensitive querying from workloads already inside the VPC
- Reserved-IP-range deployments (pair with the deployed index's
  `reservedIpRanges` field)

## Prerequisites

The VPC must have Private Services Access configured — a
`GcpServiceNetworkingConnection` with an allocated peering range —
before the endpoint is created, or the create fails.

## Remix ideas

- Replace the `valueFrom` reference with a literal
  `projects/{projectNumber}/global/networks/{name}` value when the
  network is not Planton-managed.
- Pin deployed indexes to reserved address ranges by creating
  `GcpGlobalAddress` resources (purpose `VPC_PEERING`) and naming them in
  the deployed index's `reservedIpRanges`.
