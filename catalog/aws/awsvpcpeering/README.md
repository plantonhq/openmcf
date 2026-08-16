# AwsVpcPeering

One side of a VPC peering connection, as a request-XOR-accept mode union: the request arm creates the peering from its VPC toward a peer (same-account auto-accept supported); the accept arm adopts and accepts a pending connection by id from the accepter side. DNS-resolution options fold into both arms — the standalone options resource records composed.

## Highlights

- **Cross-account and cross-region topologies are two instances of one kind**: the requester runs the request arm, the peer account/region runs the accept arm against the shared `pcx-` id (chart-wired to this kind's own `peering_connection_id` output for same-account chains).
- **The provider's walls are CELs**: auto-accept forbids `peer_region` AND `peer_owner_id` (AWS activates cross-boundary peerings only from the accepter), and the accepter-side DNS option from the request arm requires auto-accept.
- **Contracts taught in place**: the duplicate-connection trap (AWS returns the existing id for a repeated VPC pair), options requiring an ACTIVE connection, the accept arm's no-op destroy (abandon, never delete — only the requester side deletes).

## Both Engines

Both modules render whichever arm the spec configured and export the same outputs: `peering_connection_id` (import ID for both arms), `accept_status`.

## Chart Wiring

`request.vpc_id` / `request.peer_vpc_id` → AwsVpc `vpc_id`; `accept.vpc_peering_connection_id` → another AwsVpcPeering's `peering_connection_id`. After activation, each side adds routes toward the peer CIDR (AwsVpc route surfaces) — peering carries no routes itself.
