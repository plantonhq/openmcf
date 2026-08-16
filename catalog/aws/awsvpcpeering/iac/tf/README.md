# AwsVpcPeering — Terraform/OpenTofu module

Manages one side of a VPC peering connection: the request arm (`aws_vpc_peering_connection`) or the accept arm (`aws_vpc_peering_connection_accepter`) — the spec's CEL guarantees exactly one.

Module facts worth knowing before editing:

- **The peering topology is fixed for life** — `vpc_id`, `peer_vpc_id`, `peer_owner_id`, and `peer_region` all replace the connection on change.
- **`auto_accept` is same-account, same-region only** — the provider hard-errors on `peer_region` + `auto_accept`; the spec's CELs front-load both walls. Cross-account/cross-region peerings stay `pending-acceptance` until an accept-arm instance accepts.
- **DNS-resolution options are managed in-line as the single owner** — the standalone `aws_vpc_peering_connection_options` resource fights this form (the provider's own docs warn the two overwrite each other). Options need an ACTIVE connection; on a pending request AWS rejects the modification until accepted.
- **The accept arm's destroy is a no-op at AWS** — it abandons management without deleting the peering; only the requester side deletes. The request arm's delete also swallows the failed-state transition error.
- **The duplicate-connection trap** — AWS returns the EXISTING `pcx-` id for a second request on the same VPC pair; never declare the same pair twice.

Outputs mirror the Pulumi module key-for-key: `peering_connection_id` (import ID), `accept_status`.
