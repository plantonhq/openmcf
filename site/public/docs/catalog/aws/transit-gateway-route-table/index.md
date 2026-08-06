---
title: "Transit Gateway Route Table"
description: "Transit Gateway Route Table deployment documentation"
icon: "package"
order: 100
componentName: "awstransitgatewayroutetable"
---

# AWS Transit Gateway Route Table

Defines one isolated routing domain inside a Transit Gateway hub: which attachments use the table, which advertise into it, and its static and prefix-list routes. Segmented topologies (prod/non-prod isolation, inspection hair-pinning, egress VPCs, blackholed quarantine ranges) are built from this resource.

## What Gets Created

When you deploy an AwsTransitGatewayRouteTable resource, Planton provisions:

- **Route Table** — an `aws_ec2_transit_gateway_route_table` on the referenced gateway
- **Associations** — one `aws_ec2_transit_gateway_route_table_association` per associated attachment
- **Propagations** — one `aws_ec2_transit_gateway_route_table_propagation` per propagating attachment
- **Static Routes** — one `aws_ec2_transit_gateway_route` per entry (attachment-forwarded or blackhole)
- **Prefix List References** — one `aws_ec2_transit_gateway_prefix_list_reference` per entry

## Prerequisites

- **AWS credentials** configured via a Planton provider connection
- **An `AwsTransitGateway`**, and attachments to route (Planton-managed or literal IDs for VPN / Direct Connect / peering attachments)

## Quick Start

Create a file `route-table.yaml`:

```yaml
apiVersion: aws.planton.dev/v1alpha1
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
  routes:
    - destinationCidrBlock: 10.200.0.0/16
      blackhole: true
```

Deploy it:

```bash
planton apply -f route-table.yaml
```

## Key Behaviors

- **One association per attachment, gateway-wide** -- AWS enforces this at apply time; an associated attachment must have its default-table association turned off.
- **Statics beat propagations** on longest-prefix ties -- use them for inspection detours, egress default routes, and blackholes.
- **Members are keyed by stable identifiers** (attachment ID, destination CIDR, prefix list ID), so membership changes add or remove exactly one underlying resource.

## Stack Outputs

| Output | Description |
|--------|-------------|
| `route_table_id` | The table ID (tgw-rtb-...) |
| `route_table_arn` | For IAM policies |
