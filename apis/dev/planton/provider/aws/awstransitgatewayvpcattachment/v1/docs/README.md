# AwsTransitGatewayVpcAttachment -- Research Notes

## Provider Surface

Modeled 1:1 on `aws_ec2_transit_gateway_vpc_attachment`:

- `transit_gateway_id` (Required, ForceNew) -- spec ref -> `AwsTransitGateway.status.outputs.transit_gateway_id`.
- `vpc_id` (Required, ForceNew) -- spec ref -> `AwsVpc.status.outputs.vpc_id`.
- `subnet_ids` (Required, min 1, set; in-place updatable) -- spec `repeated StringValueOrRef` -> `AwsSubnet`.
- `dns_support` (Optional, default "enable") -- spec `optional bool`.
- `ipv6_support` / `appliance_mode_support` (Optional, default "disable") -- plain bools (spec default false == AWS default), always sent explicitly.
- `security_group_referencing_support` (Optional+Computed -- inherits the GATEWAY's dial) -- spec `optional bool`, presence-honest: unset is never sent.
- `transit_gateway_default_route_table_association` / `_propagation` (Optional+Computed bools -- AWS derives from the gateway's default dials) -- spec `optional bool` pair, presence-honest.
- Computed: `arn`, `vpc_owner_id`.

The attachment has no name attribute; the Name tag (from `metadata.name`) is the console identity in BOTH engines.

## Design Decisions

- **First-class kind, not a folded block on the gateway.** Many-per-gateway, independent lifecycle, and the reference target of route-table associations/propagations/routes -- the split test is met three ways (the App Runner / Batch / Cognito decomposition class).
- **Inheritance kept honest.** The three Optional+Computed attributes (SG referencing, default association, default propagation) stay tri-state: unset inherits the gateway posture exactly as the provider models it, instead of pinning a value the user never chose.
- **Registry prerequisites `[AwsTransitGateway, AwsSubnet]`** -- both are hard deploy requirements (required refs); the VPC arrives transitively through the subnet's own prerequisite chain.

## Deferral Ledger

See the family ledger in the `AwsTransitGateway` research notes (peering, Connect, multicast, policy/metering tables, cross-account accepters, default-table designations, flow logs).

## Verification

- Spec tests cover required-ness of region/gateway/VPC/subnets and the fully-dialed shape.
- E2E: chains the shared VPC + two-AZ subnet fixtures and a gateway prerequisite -- create -> DescribeTransitGatewayVpcAttachments (available) -> destroy -> verify deleted/absent (state-aware: deleting/deleted count as absent).
