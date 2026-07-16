---
title: "Transit Gateway VPC Attachment"
description: "Transit Gateway VPC Attachment deployment documentation"
icon: "package"
order: 100
componentName: "awstransitgatewayvpcattachment"
---

# AWS Transit Gateway VPC Attachment

Connects one VPC to a Transit Gateway hub through elastic network interfaces in the chosen subnets, making the VPC a spoke of the hub-and-spoke topology. The attachment ID it exports is the join key that Transit Gateway route tables associate, propagate, and route against.

## What Gets Created

When you deploy an AwsTransitGatewayVpcAttachment resource, Planton provisions:

- **VPC Attachment** — an `aws_ec2_transit_gateway_vpc_attachment` wiring the referenced VPC into the referenced gateway through one ENI per listed subnet, with per-attachment DNS, IPv6, appliance-mode, and security-group-referencing options and explicit default-route-table membership control

## Prerequisites

- **AWS credentials** configured via a Planton provider connection
- **An `AwsTransitGateway`** (or the ID of an existing gateway)
- **An `AwsVpc` with subnets** in each Availability Zone you want reachable through the gateway

## Quick Start

Create a file `attachment.yaml`:

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsTransitGatewayVpcAttachment
metadata:
  name: app-vpc-attachment
spec:
  region: us-east-1
  transitGatewayId:
    valueFrom:
      kind: AwsTransitGateway
      name: core-hub
      fieldPath: status.outputs.transit_gateway_id
  vpcId:
    valueFrom:
      kind: AwsVpc
      name: app-vpc
      fieldPath: status.outputs.vpc_id
  subnetIds:
    - valueFrom:
        kind: AwsSubnet
        name: app-tgw-subnet-az1
        fieldPath: status.outputs.subnet_id
```

Deploy it:

```bash
planton apply -f attachment.yaml
```

## Key Behaviors

- **Gateway and VPC are create-time immutable**; the subnet set updates in place, so AZs can be added or removed without replacement.
- **Default route table membership is inheritable**: left unset, the gateway's own dials decide; pin `defaultRouteTableAssociation: false` when a custom `AwsTransitGatewayRouteTable` owns the association (an attachment can be associated with at most one table).
- **Return-path reminder**: spoke subnets still need their own routes toward the gateway (`AwsSubnet.routes` with `targetType: transit_gateway`).

## Stack Outputs

| Output | Description |
|--------|-------------|
| `attachment_id` | Referenced by route table associations, propagations, and routes |
| `attachment_arn` | For IAM policies |
| `vpc_owner_id` | Account owning the attached VPC |
