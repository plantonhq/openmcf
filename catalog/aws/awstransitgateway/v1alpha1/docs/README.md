# AwsTransitGateway -- Research Notes

## Provider Surface

Modeled 1:1 on `aws_ec2_transit_gateway`:

- `description` (Optional) -- spec `description`.
- `amazon_side_asn` (Optional, default 64512; replaces on change) -- spec `amazon_side_asn` with the provider's private ranges as CEL (64512-65534, 4200000000-4294967294); 0 falls through to the AWS default.
- `default_route_table_association` / `default_route_table_propagation` (Optional, default "enable"; asymmetric ForceNew -- disable -> enable replaces, enable -> disable updates in place) -- spec `optional bool` pair; unset falls through to the provider default.
- `dns_support` / `vpn_ecmp_support` (Optional, default "enable") -- spec `optional bool` pair.
- `auto_accept_shared_attachments` / `security_group_referencing_support` / `multicast_support` (Optional, default "disable"; multicast is ForceNew) -- plain bools (spec default false == AWS default), always sent explicitly.
- `encryption_support` (Optional+Computed) -- spec `optional bool`; unset is never sent so AWS computes the effective posture. Provider floor >= 6.26.0: the argument landed in 6.25.0 with a crash regression fixed in 6.26.0.
- `transit_gateway_cidr_blocks` (Optional, max 5; IPv4 prefix <= /24, IPv6 <= /64, never 169.254.0.0/16) -- the provider's validator bounds as CEL.
- Computed: `arn`, `owner_id`, `association_default_route_table_id`, `propagation_default_route_table_id`.

The gateway has no name attribute; the Name tag (from `metadata.name`) is the console identity in BOTH engines.

## Design Decisions

- **Pure hub -- attachments split OUT.** The previous shape embedded `vpc_attachments` (min_items 1) and exported a `vpc_attachment_ids` map. Attachments are many-per-gateway, independently lifecycled, and are the reference target of route-table associations/propagations/routes -- the split test is met three ways. A gateway with zero attachments is also legitimate (it is how hubs are pre-staged before spokes migrate in), which the embedded min_items 1 made unrepresentable.
- **Tri-state dials.** The enable/disable dials whose AWS default is "enable" are `optional bool` so an omitted dial inherits the AWS default while an explicit false is expressible -- the plain-bool shape could not distinguish "unset" from "disable".
- **encryption_support presence-honest.** Optional+Computed in the provider; only sent when the spec pins it.

## Deferral Ledger (family demand-check, recorded once for the three kinds)

- `aws_ec2_transit_gateway_peering_attachment` + `_accepter` -- DEFER: cross-region/cross-account topology; no chart composes cross-region today. A future kind composes against the exported `transit_gateway_id` with zero rework (the RDS global-cluster precedent).
- `aws_ec2_transit_gateway_vpc_attachment_accepter` -- DEFER: the cross-account RAM-shared side of the attachment handshake (the Route 53 cross-account association-plane precedent); same-account attachments are the shipped kind.
- `aws_ec2_transit_gateway_connect` + `_connect_peer` -- DEFER: the SD-WAN/GRE appliance integration surface; composes against the exported gateway ID and an attachment ID. The gateway's `transit_gateway_cidr_blocks` (its prerequisite) IS modeled.
- `aws_ec2_transit_gateway_multicast_domain` + `_association` + `_group_member` + `_group_source` -- DEFER: a separate multicast product surface; the gateway's `multicast_support` dial (its prerequisite) IS modeled.
- `aws_ec2_transit_gateway_policy_table` + `_association` -- DEFER: Cloud WAN / network-manager peering policy surface, not part of the core VPC routing story.
- `aws_ec2_transit_gateway_metering_policy` + `_entry` + `_list` -- DEFER: brand-new usage-metering surface; revisit on pull.
- `aws_ec2_transit_gateway_default_route_table_association` / `_propagation` -- DEFER: swapping the gateway's DEFAULT table designation onto a custom table is a niche migration maneuver; modeling it on the gateway would create a circular reference (gateway -> route table -> gateway), and on the route table it invites two tables claiming the default. The segmented topology (explicit associations on custom tables, default dials off) covers the real need without it.
- Transit Gateway flow logs (`aws_flow_log` with `transit_gateway_id`) -- DEFER: the whole flow-log surface (VPC/subnet/ENI/TGW) is one future kind, not per-carrier fields.

## Verification

- Spec tests cover the ASN ranges, the CIDR-block bounds (count, prefix length per family, link-local), and required region.
- E2E: dependency-free leaf lane -- create -> DescribeTransitGateways (available) -> destroy -> verify deleted/absent (state-aware: deleting/deleted count as absent).
