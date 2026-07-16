# AwsTransitGateway — Terraform Module

## Overview

This Terraform module provisions an AWS Transit Gateway -- the regional hub that VPC attachments, route tables, VPN connections, and Direct Connect gateways compose onto.

## Module Structure

```
main.tf       — aws_ec2_transit_gateway
locals.tf     — identity tags + enable/disable dial mapping (tri-state aware)
outputs.tf    — transit_gateway_id/arn, owner_id, default route table pair
variables.tf  — Generator-owned typed contract (metadata, spec)
provider.tf   — AWS provider configuration (>= 6.26.0)
```

## Behavior Notes

- Tri-state dials (proto `optional bool`) map null -> null so an omitted
  dial falls through to the provider/AWS default instead of being pinned.
- The default-table association/propagation dials carry AWS's asymmetric
  replacement rule: disable -> enable replaces the gateway; enable ->
  disable updates in place.
- The Name tag (from `metadata.name`) is the gateway's console identity --
  Transit Gateways have no name attribute.

## Outputs

| Output | Description |
|--------|-------------|
| `transit_gateway_id` | Join key for attachments, route tables, subnet routes |
| `transit_gateway_arn` | For IAM policies and RAM sharing |
| `owner_id` | Owning account ID |
| `association_default_route_table_id` | Default association table (empty when disabled) |
| `propagation_default_route_table_id` | Default propagation table (empty when disabled) |
