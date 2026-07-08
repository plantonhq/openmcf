---
title: "Presets"
description: "Ready-to-deploy configuration presets for Transit Gateway VPC Attachment"
type: "preset-list"
componentSlug: "transit-gateway-vpc-attachment"
componentTitle: "Transit Gateway VPC Attachment"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-mesh-spoke"
    rank: "01"
    title: "Mesh Spoke"
    excerpt: "The default spoke shape: two-AZ subnet coverage on a full-mesh hub. Default route table membership is inherited from the gateway, so the spoke lands in the mesh with zero routing configuration."
  - slug: "02-segmented-spoke"
    rank: "02"
    title: "Segmented Spoke"
    excerpt: "A spoke for segmented topologies: default-table membership pinned off, so the attachment routes nothing until an `AwsTransitGatewayRouteTable` explicitly associates it and accepts its propagation."
  - slug: "03-inspection-vpc"
    rank: "03"
    title: "Inspection VPC Attachment"
    excerpt: "The attachment shape for a shared-services VPC hosting stateful inspection appliances (firewall, IDS/IPS): appliance mode keeps each flow's return traffic in the Availability Zone it entered through,..."
---

# Transit Gateway VPC Attachment Presets

Ready-to-deploy configuration presets for Transit Gateway VPC Attachment. Each preset is a complete manifest you can copy, customize, and deploy.
