# AWS Private Network Hub

A Transit Gateway hub for the VPCs you already run. Two spokes attach into
their own routing domains; by default the domains are segmented — production
and non-production cannot exchange a packet, and blackhole guardrails keep it
that way even as routes grow — or flip one toggle for plain any-to-any
connectivity that replaces a VPC-peering mesh.

Peering does not scale: every new VPC means N-1 new peering connections,
each with route-table edits on both sides, and no central place to see or
control who reaches whom. A hub makes connectivity a routing-domain decision
made once, in one place — and it is also where VPN, Direct Connect, and
inspection VPCs attach later without touching the spokes.

## Architecture

```
        production VPC (existing)         non-production VPC (existing)
               |                                     |
   [AwsTransitGatewayVpcAttachment]     [AwsTransitGatewayVpcAttachment]
     pinned out of default tables         pinned out of default tables
               |  associated with                    |  associated with
               v                                     v
   [AwsTransitGatewayRouteTable]        [AwsTransitGatewayRouteTable]
        prod-domain                          nonprod-domain
     segmented:  blackhole nonprod CIDR   segmented:  blackhole prod CIDR
     toggle on:  propagate nonprod        toggle on:  propagate prod
               \                                     /
                \___________ [AwsTransitGateway] ___/
                  default association/propagation OFF,
                  BGP ASN pinned, DNS support on
```

## Included Cloud Resources

| Resource | Kind | Purpose |
|---|---|---|
| Hub | `AwsTransitGateway` | The regional router; default-table dials off so membership is always explicit. |
| Prod / non-prod attachments | `AwsTransitGatewayVpcAttachment` | Connect the existing VPCs; ENIs anchored in the subnets you pick. |
| Prod / non-prod domains | `AwsTransitGatewayRouteTable` | Per-spoke routing views: association + (blackhole guardrail XOR cross-propagation). |

## Parameters

| Name | Description | Default | Required |
|---|---|---|---|
| `aws_region` | Region of the hub and the spoke VPCs. | `us-east-1` | yes |
| `hub_name` | Gateway name and companion-resource prefix. | `my-network-hub` | yes |
| `amazon_side_asn` | AWS-side BGP ASN (immutable; pick from 64512–65534). | `64512` | yes |
| `prod_vpc_id` / `nonprod_vpc_id` | The existing spoke VPCs. | example ids | yes |
| `prod_vpc_cidr` / `nonprod_vpc_cidr` | Spoke CIDRs — the blackhole guardrails while segmented. | `10.0.0.0/16` / `10.1.0.0/16` | yes |
| `prod_subnet_ids` / `nonprod_subnet_ids` | Attachment ENI subnets, one per AZ. | example ids | yes |
| `spokes_interconnected_enabled` | Any-to-any between the spokes instead of segmentation. | `false` | no |

## Completing the path (post-deploy)

The hub owns the gateway side of routing; the **return routes inside each
spoke VPC** belong to the spoke's owner — this chart never modifies
resources it does not own. For each subnet that should send traffic through
the hub, add a route targeting the gateway. On a Planton-managed subnet
(e.g. from network-foundation) that is one inline route — `targetId` is
deliberately polymorphic, so reference the hub's output explicitly:

```yaml
routes:
  - destinationCidrBlock: 10.1.0.0/16     # the peer spoke (or wider, e.g. 10.0.0.0/8)
    targetType: transit_gateway
    targetId:
      valueFrom:
        kind: AwsTransitGateway
        name: my-network-hub
        fieldPath: status.outputs.transit_gateway_id
```

For subnets managed outside Planton: `aws ec2 create-route
--route-table-id rtb-... --destination-cidr-block <peer-cidr>
--transit-gateway-id <status.outputs.transit_gateway_id>`.

While the hub is segmented, cross-spoke packets die at the gateway's
blackhole even after these routes exist — the spoke routes decide what
*enters* the hub; the domains decide what *crosses* it. That layering is
what lets you pre-wire spoke routes broadly (`10.0.0.0/8` once) and govern
reachability purely at the hub.

Also mind AWS's one-association-per-attachment rule: an attachment joins
exactly one routing domain. Re-homing a spoke means removing it from one
table's `associations` and adding it to another's in the same change.

## Choosing the posture

- **Segmented (default)** — environment isolation as infrastructure:
  nothing in non-prod can reach prod even when someone fat-fingers a route,
  because the blackhole wins on longest-prefix match. Shared capabilities
  (below) still reach both.
- **Interconnected** — the peering-mesh replacement: both spokes learn each
  other's routes through propagation. Flip `spokes_interconnected_enabled`
  and the guardrails are replaced by propagations in one apply.

## Day-2 guidance

- **Shared-services domain**: attach the shared VPC (DNS, CI, tooling) with
  its own `AwsTransitGatewayRouteTable`, propagate BOTH spoke attachments
  into it, and propagate ITS attachment into both spoke domains — every
  spoke reaches the shared services, spokes still cannot reach each other.
- **Centralized egress / inspection**: point a static default route in each
  spoke domain at an inspection-VPC attachment (the route-table spec's
  `routes[].attachmentId` arm); enable `applianceModeSupport` on that
  attachment so stateful firewalls see both directions of a flow.
- **Hybrid connectivity**: VPN and Direct Connect attach to the same hub
  and slot into domains the same way — the pinned `amazon_side_asn` is
  already reserved for their BGP sessions. The route-table spec accepts
  literal attachment ids for attachment types that are not yet first-class
  kinds.
- **More spokes**: copy the attachment + domain pair per VPC (or re-deploy
  the chart per environment boundary). The hub's per-attachment cost
  (~$36/month + per-GB data processing) replaces per-peering route
  sprawl, not per-peering dollars — peering itself is free; what you buy
  is a governable center.
- **Cross-region**: peer two hubs with a transit gateway peering
  attachment (static routes only across the peering; literal attachment
  ids work in the `routes` arm today).
