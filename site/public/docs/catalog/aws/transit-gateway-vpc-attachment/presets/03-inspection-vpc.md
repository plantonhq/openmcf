---
title: "Inspection VPC Attachment"
description: "The attachment shape for a shared-services VPC hosting stateful inspection appliances (firewall, IDS/IPS): appliance mode keeps each flow's return traffic in the Availability Zone it entered through,..."
type: "preset"
rank: "03"
presetSlug: "03-inspection-vpc"
componentSlug: "transit-gateway-vpc-attachment"
componentTitle: "Transit Gateway VPC Attachment"
provider: "aws"
icon: "package"
order: 3
---

# Inspection VPC Attachment

The attachment shape for a shared-services VPC hosting stateful inspection appliances (firewall, IDS/IPS): appliance mode keeps each flow's return traffic in the Availability Zone it entered through, preserving the symmetric routing stateful appliances require.

## When to Use

- ONLY on the attachment of the VPC that hosts the appliances — spoke attachments never need appliance mode
- Hair-pin topologies where spoke route tables default-route through the inspection VPC

## What It Configures

- **`applianceModeSupport: true`** — AZ-symmetric flow placement for stateful inspection
- **Default-table membership off** — inspection VPCs live in custom routing domains by definition

## What to Customize

- Replace `<aws-region>` and the referenced hub/VPC/subnet names with your resources
- Steer traffic through this attachment with static routes (`0.0.0.0/0` or specific prefixes) in the spoke domains' `AwsTransitGatewayRouteTable` resources
