# Segmented Hub

The isolation-first posture: both default-table dials off, so nothing routes until an `AwsTransitGatewayRouteTable` explicitly associates and propagates attachments. Production cannot see non-production unless a route table says so.

## When to Use

- Prod/non-prod isolation with a shared-services domain both can reach
- Inspection topologies where an appliance VPC hair-pins inter-spoke traffic
- Any environment where "deny by default" is the network posture

## What It Configures

- **`defaultRouteTableAssociation: false`** — new attachments join no table until a custom table associates them
- **`defaultRouteTablePropagation: false`** — no CIDRs advertise anywhere until a custom table accepts them

## What to Customize

- Replace `<aws-region>` with your region
- Carve the domains with `AwsTransitGatewayRouteTable` resources (one per isolation zone) and attach spokes with their default-table membership pinned off
- AWS QUIRK: flipping either dial back to enabled later REPLACES the gateway — choose the posture up front
