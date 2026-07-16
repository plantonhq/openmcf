# Private Service Connect Index Endpoint

The strongest network isolation for vector search: consumers connect
through a PSC service attachment, no VPC peering required.

## What this preset creates

An index endpoint named `PSC Search` in `us-central1` exposed via
Private Service Connect. After a `GcpVertexAiDeployedIndex` places an
index onto the endpoint, its `service_attachment` output is the address
allowlisted consumer projects target with PSC forwarding rules.

## When to use

- Multi-project architectures where peering every consumer VPC is
  unmanageable
- Strict tenant isolation: each consumer connects through its own
  forwarding rule
- Organizations standardizing on PSC as the private-connectivity model

## Remix ideas

- Trim `projectAllowlist` to exactly the consumer projects that should
  reach the indexes — the allowlist is the access-control boundary.
- Consumers create a forwarding rule per index endpoint; automate that
  side with a `GcpForwardingRule` targeting the deployed index's
  `service_attachment` output.
