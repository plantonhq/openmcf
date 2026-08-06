# AwsTransitGatewayRouteTable — Terraform Module

## Overview

This Terraform module provisions a Transit Gateway route table -- one isolated routing domain -- together with its folded associations, propagations, static routes, and prefix list references.

## Module Structure

```
main.tf       — route table + per-member association/propagation/route/prefix-list resources
locals.tf     — identity tags + stable-key maps for the folded members
outputs.tf    — route_table_id, route_table_arn
variables.tf  — Generator-owned typed contract (metadata, spec)
provider.tf   — AWS provider configuration (>= 6.26.0)
```

## Behavior Notes

- Folded members are keyed by stable identifiers (attachment ID,
  destination CIDR, prefix list ID), so a membership edit touches exactly
  one provider resource.
- An attachment can be ASSOCIATED with at most one route table across the
  whole gateway; AWS enforces that at apply time. Propagations carry no
  such limit.
- Blackhole routes send the attachment argument as null -- AWS expects it
  ABSENT, not empty.

## Outputs

| Output | Description |
|--------|-------------|
| `route_table_id` | The table ID (tgw-rtb-...) |
| `route_table_arn` | For IAM policies |
