# AwsTransitGatewayVpcAttachment — Terraform Module

## Overview

This Terraform module provisions a Transit Gateway VPC attachment -- the connection that plugs one VPC into a Transit Gateway hub through ENIs in the chosen subnets.

## Module Structure

```
main.tf       — aws_ec2_transit_gateway_vpc_attachment
locals.tf     — identity tags + enable/disable dial mapping (tri-state aware)
outputs.tf    — attachment_id, attachment_arn, vpc_owner_id
variables.tf  — Generator-owned typed contract (metadata, spec)
provider.tf   — AWS provider configuration (>= 6.26.0)
```

## Behavior Notes

- The gateway and VPC references are create-time immutable; the subnet set
  updates in place.
- Tri-state options are omitted when unset so the provider/AWS default --
  or the gateway-inherited value, for SG referencing and the default-table
  membership pair -- applies.
- The Name tag (from `metadata.name`) is the attachment's console identity.

## Outputs

| Output | Description |
|--------|-------------|
| `attachment_id` | Join key for route table associations, propagations, routes |
| `attachment_arn` | For IAM policies |
| `vpc_owner_id` | Account owning the attached VPC |
