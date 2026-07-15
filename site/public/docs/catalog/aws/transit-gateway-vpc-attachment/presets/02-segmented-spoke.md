---
title: "Segmented Spoke"
description: "A spoke for segmented topologies: default-table membership pinned off, so the attachment routes nothing until an `AwsTransitGatewayRouteTable` explicitly associates it and accepts its propagation."
type: "preset"
rank: "02"
presetSlug: "02-segmented-spoke"
componentSlug: "transit-gateway-vpc-attachment"
componentTitle: "Transit Gateway VPC Attachment"
provider: "aws"
icon: "package"
order: 2
---

# Segmented Spoke

A spoke for segmented topologies: default-table membership pinned off, so the attachment routes nothing until an `AwsTransitGatewayRouteTable` explicitly associates it and accepts its propagation.

## When to Use

- Prod/non-prod isolation domains
- Any spoke whose reachability is defined by a custom routing domain rather than the gateway's default mesh

## What It Configures

- **`defaultRouteTableAssociation: false`** — the association belongs to exactly one custom route table (AWS allows one association per attachment, gateway-wide)
- **`defaultRouteTablePropagation: false`** — the VPC's CIDRs advertise only into the tables that list this attachment in `propagations`

## What to Customize

- Replace `<aws-region>` and the referenced hub/VPC/subnet names with your resources
- Add this attachment to the right `AwsTransitGatewayRouteTable`'s `associations` (its own domain) and other tables' `propagations` (domains that may reach it)
