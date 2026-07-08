---
title: "Inspection Domain"
description: "The hair-pin pattern's spoke side: every flow leaving a spoke is default-routed through the inspection VPC's attachment, where stateful appliances decide what continues."
type: "preset"
rank: "02"
presetSlug: "02-inspection-domain"
componentSlug: "transit-gateway-route-table"
componentTitle: "Transit Gateway Route Table"
provider: "aws"
icon: "package"
order: 2
---

# Inspection Domain

The hair-pin pattern's spoke side: every flow leaving a spoke is default-routed through the inspection VPC's attachment, where stateful appliances decide what continues.

## When to Use

- Centralized east-west and north-south inspection with a firewall/IDS VPC
- Pairs with a second route table associated to the inspection attachment that routes traffic onward to its real destinations

## What It Configures

- **Spoke associations with NO propagations** — spokes cannot learn routes to each other; everything they send follows the static default
- **A `0.0.0.0/0` static route** to the inspection VPC attachment (which should have `applianceModeSupport: true` for AZ-symmetric flows)

## What to Customize

- Replace `<aws-region>`, the hub, and the attachment references with your resources
- Build the companion "post-inspection" route table: associated to the inspection attachment, propagating all spokes, so inspected traffic reaches its destination
