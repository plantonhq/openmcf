# Private VPC-Peered Endpoint

Prediction serving reachable only from inside a peered VPC, encrypted
under a customer-managed key, with an isolated dedicated DNS name.

## What this preset creates

An endpoint named `Internal Scoring` in `us-central1`, peered with the
referenced `GcpVpcNetwork` (`scoring-vpc`) and CMEK-encrypted under the
referenced `GcpKmsKey` (`scoring-key`). Prediction traffic never leaves
private address space, and the dedicated DNS name isolates this
endpoint's traffic from other users of the shared regional endpoint.

## Prerequisites

The referenced VPC must already carry Private Services Access — compose
a `GcpGlobalAddress` (purpose VPC_PEERING) and a
`GcpServiceNetworkingConnection` on the same network. The KMS key must
live in the endpoint's region.

## When to use

- Production inference on regulated or sensitive data
- Serving that must be unreachable from the public internet
- Environments with CMEK mandates

## Remix ideas

- Swap VPC peering for Private Service Connect (see the private-psc
  preset) when you want per-project access control instead of
  network-level reachability — note PSC excludes the dedicated DNS.
- Add `requestResponseLoggingConfig` to sample predictions into BigQuery
  for drift monitoring.
