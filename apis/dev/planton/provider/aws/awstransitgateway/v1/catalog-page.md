# AWS Transit Gateway

Deploys an AWS Transit Gateway -- the regional hub-and-spoke networking core that replaces VPC peering meshes. The gateway is a pure hub: VPC attachments and custom route tables are their own composable resources (`AwsTransitGatewayVpcAttachment`, `AwsTransitGatewayRouteTable`) that reference the gateway's outputs, so spokes and routing domains come and go without touching the hub.

## What Gets Created

When you deploy an AwsTransitGateway resource, Planton provisions:

- **Transit Gateway** — an `aws_ec2_transit_gateway` with the configured ASN, routing defaults, DNS/ECMP/multicast/encryption dials, and optional TGW Connect CIDR blocks
- **Default Route Tables** — created by AWS when default association/propagation are enabled; their IDs are exported as stack outputs

## Prerequisites

- **AWS credentials** configured via a Planton provider connection

## Quick Start

Create a file `tgw.yaml`:

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsTransitGateway
metadata:
  name: core-hub
spec:
  region: us-east-1
  description: Full-mesh hub for the platform VPCs
```

Deploy it:

```bash
planton apply -f tgw.yaml
```

Then attach VPCs with `AwsTransitGatewayVpcAttachment` resources referencing the gateway, and (for segmented topologies) carve routing domains with `AwsTransitGatewayRouteTable`.

## Key Behaviors

- **Full mesh by default**: with `defaultRouteTableAssociation` and `defaultRouteTablePropagation` enabled (the defaults), every attachment can reach every other attachment with zero routing configuration.
- **Segmentation posture is a create-time call**: flipping either default-table dial from disabled back to enabled REPLACES the gateway (AWS's asymmetric rule); disabling updates in place.
- **The Name tag is the console identity** -- gateways have no name attribute; this component sets it from `metadata.name` in both engines.

## Stack Outputs

| Output | Description |
|--------|-------------|
| `transit_gateway_id` | Referenced by attachments, route tables, and subnet routes |
| `transit_gateway_arn` | For IAM policies and AWS RAM sharing |
| `owner_id` | Owning account ID |
| `association_default_route_table_id` | Default association table (empty when disabled) |
| `propagation_default_route_table_id` | Default propagation table (empty when disabled) |
