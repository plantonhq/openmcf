---
title: "Private Backend Access"
description: "The production connector shape: two private subnets in different availability zones and a dedicated egress security group, composed from first-class networking resources in the same resource graph."
type: "preset"
rank: "01"
presetSlug: "01-private-backend-access"
componentSlug: "app-runner-vpc-connector"
componentTitle: "App Runner VPC Connector"
provider: "aws"
icon: "package"
order: 1
---

# Private Backend Access

The production connector shape: two private subnets in different availability zones and a dedicated egress security group, composed from first-class networking resources in the same resource graph.

## When to Use

- App Runner services that read from RDS, ElastiCache, or internal APIs inside a VPC
- Any fleet of services sharing one VPC egress path

## What It Configures

- **Two-AZ subnet spread** — App Runner routes egress only through AZs the connector has an ENI in, so a single-AZ connector is a single point of failure
- **A dedicated egress security group** — grant this group egress to your databases' ports, and admit ingress FROM it on the databases' own groups

## What to Customize

- Replace `<aws-region>` and the referenced subnet/security-group names with your networking resources
- Reference from services via `vpcConnectorArn` — one connector serves any number of services
- Remember everything here is create-time immutable: changing subnets or groups replaces the connector (a new revision under the same name)
