# AwsTransitGatewayRouteTable -- Research Notes

## Provider Surface

The kind folds five provider resources that share the table's lifecycle:

- `aws_ec2_transit_gateway_route_table`: `transit_gateway_id` (Required, ForceNew) -- spec ref -> `AwsTransitGateway`. Computed: `arn`, `default_association_route_table`, `default_propagation_route_table` (the two computed booleans are introspection, not join keys -- not exported).
- `aws_ec2_transit_gateway_route_table_association` (per `associations[]` entry): `transit_gateway_attachment_id` + `transit_gateway_route_table_id`, both ForceNew. `replace_existing_association` NOT exposed -- silently stealing an attachment from another table would hide a cross-document conflict the user should resolve in the manifests; the AWS error is the honest signal.
- `aws_ec2_transit_gateway_route_table_propagation` (per `propagations[]` entry): same key pair, ForceNew.
- `aws_ec2_transit_gateway_route` (per `routes[]` entry): `destination_cidr_block` (ForceNew; IPv4 or IPv6), `transit_gateway_attachment_id` XOR `blackhole` (the exclusivity is a spec CEL; AWS expects the attachment ABSENT for blackholes).
- `aws_ec2_transit_gateway_prefix_list_reference` (per `prefix_list_references[]` entry): `prefix_list_id` (ForceNew), attachment XOR blackhole (same CEL). `prefix_list_owner_id` computed, not exported (cross-account prefix lists compose by ID).

The table has no name attribute; the Name tag (from `metadata.name`) is the console identity in BOTH engines.

## Design Decisions

- **First-class kind.** Many-per-gateway with independent lifecycles; the routing-domain unit segmentation policy hangs off. Folding tables into the gateway would bury attachment references inside the hub and churn the hub on every domain edit.
- **Associations/propagations folded ON the table** (the ElastiCache user-group membership precedent), not modeled on the attachment. Deciding argument: routing domains must include attachments created OUTSIDE the Planton graph -- VPN, Direct Connect, and peering attachments -- and only a table-side `StringValueOrRef` list accepts those as literal IDs. The table is the routing-domain definition; membership IS routing policy.
- **Gateway-wide association uniqueness is AWS-enforced.** One association per attachment across all tables cannot be validated across documents (the two-listeners-one-port class); stated loudly in the field comments and this doc instead.
- **Members keyed by stable identifiers** (attachment ID / destination CIDR / prefix list ID; uniqueness CELs on the spec) so a membership edit adds/removes exactly one provider resource in both engines.
- **Registry prerequisites `[AwsTransitGateway]` only** -- an empty table is deployable; attachment references are optional composition (scenarios declare them via the e2e-prerequisites annotation).

## Deferral Ledger

See the family ledger in the `AwsTransitGateway` research notes. Route-table-specific: the default-table DESIGNATION resources (`aws_ec2_transit_gateway_default_route_table_association`/`_propagation`) are deferred -- a niche default-swap maneuver that would either create a circular gateway<->table reference or invite two tables claiming the default; the segmented pattern (explicit associations, default dials off) covers the need.

## Verification

- Spec tests cover required region/gateway, the route and prefix-list target XORs (neither/both/empty), and destination/prefix-list uniqueness.
- E2E: the composed lane declares a VPC attachment via the e2e-prerequisites annotation and proves the isolated-routing-domain story live -- custom table + association + static route + blackhole -> DescribeTransitGatewayRouteTables (available) -> destroy -> verify deleted/absent.
