# AWS Transit Gateway Route Table

A Transit Gateway route table is one isolated routing domain inside a Transit Gateway hub. It owns its domain's membership and routes:

- **Associations** -- which attachments USE this table for their outbound lookups (at most one association per attachment, gateway-wide).
- **Propagations** -- which attachments ADVERTISE their routes into this table (any number of tables per attachment).
- **Static routes** -- a destination CIDR forwarded through one attachment, or blackholed. Statics win longest-prefix ties over propagated routes.
- **Prefix list references** -- route a managed prefix list's whole CIDR set via one attachment (or blackhole it), tracking list membership automatically.

The canonical segmented topology -- production spokes that reach shared services but not each other, an inspection VPC that hair-pins inter-spoke traffic -- is built from exactly this resource plus attachments created with their default-table membership turned off.

## When to Use

- Any topology beyond the default full mesh: prod/non-prod isolation, inspection hair-pinning, egress-VPC default routes, blackholed quarantine CIDRs.
- Keep the gateway's `defaultRouteTableAssociation`/`defaultRouteTablePropagation` disabled (or pin the relevant attachments' membership off) so custom tables own the routing.

## When NOT to Use

- A pure full-mesh topology needs no custom tables -- the gateway's default route table handles everything.

## Prerequisites

- An `AwsTransitGateway`.
- `AwsTransitGatewayVpcAttachment` resources (or literal attachment IDs for VPN / Direct Connect / peering attachments created outside the Planton graph).

## Spec Fields

| Field | Type | Default | Description |
|---|---|---|---|
| `region` | string | (required) | AWS region; must match the gateway |
| `transitGatewayId` | ref | (required) | The gateway this table belongs to (create-time immutable) |
| `associations` | ref[] | [] | Attachments associated with this table. Each must have default association off and appear in no other table's associations |
| `propagations` | ref[] | [] | Attachments advertising their routes into this table |
| `routes` | object[] | [] | Static routes: `destinationCidrBlock` (unique) + exactly one of `attachmentId` / `blackhole` |
| `prefixListReferences` | object[] | [] | Prefix list routes: `prefixListId` (unique) + exactly one of `attachmentId` / `blackhole` |

## Stack Outputs

| Output | Description |
|---|---|
| `route_table_id` | The table ID (tgw-rtb-...) |
| `route_table_arn` | ARN for IAM policies |

## Composition

- References an `AwsTransitGateway` and any number of `AwsTransitGatewayVpcAttachment`s.
- Association/propagation/route entries accept literal attachment IDs, so VPN, Direct Connect, and peering attachments created outside the graph participate in the domain too.

## Example: isolated production domain

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsTransitGatewayRouteTable
metadata:
  name: prod-domain
spec:
  region: us-east-1
  transitGatewayId:
    valueFrom:
      kind: AwsTransitGateway
      name: segmented-hub
      fieldPath: status.outputs.transit_gateway_id
  associations:
    - valueFrom:
        kind: AwsTransitGatewayVpcAttachment
        name: prod-vpc-attachment
        fieldPath: status.outputs.attachment_id
  propagations:
    - valueFrom:
        kind: AwsTransitGatewayVpcAttachment
        name: shared-services-attachment
        fieldPath: status.outputs.attachment_id
  routes:
    # Everything else leaves through the egress VPC.
    - destinationCidrBlock: 0.0.0.0/0
      attachmentId:
        valueFrom:
          kind: AwsTransitGatewayVpcAttachment
          name: egress-vpc-attachment
          fieldPath: status.outputs.attachment_id
    # The non-prod CIDR can never be reached from this domain.
    - destinationCidrBlock: 10.200.0.0/16
      blackhole: true
```

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
