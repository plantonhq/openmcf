---
title: "Managed Default Domain"
description: "Adopt a flat hub into a managed routing domain without a two-step migration: this table becomes the gateway's default association and propagation table, and attachments still parked on the old..."
type: "preset"
rank: "03"
presetSlug: "03-managed-default-domain"
componentSlug: "transit-gateway-route-table"
componentTitle: "Transit Gateway Route Table"
provider: "aws"
icon: "package"
order: 3
---

# Managed Default Domain

Adopt a flat hub into a managed routing domain without a two-step migration: this table becomes the gateway's default association and propagation table, and attachments still parked on the old default table are taken over in the same apply.

## When to Use

- Migrating a default-enabled ("flat mesh") hub toward governed routing: every attachment that keeps its default-table dials on now lands in THIS table, where your statics and blackholes apply
- Taking over attachments already associated elsewhere — including attachments on a RAM-shared gateway whose default membership the sharing account controls

## What It Configures

- **`setAsDefaultAssociationTable` / `setAsDefaultPropagationTable`** — repoint the GATEWAY's default-table pointers at this table; meaningful only on a gateway whose default dials are enabled, at most one claimant table per gateway, and removing a flag restores the gateway's original default table
- **`associations[].replaceExistingAssociation`** — disassociate the attachment's current association (the old default table) and associate it here, in one apply
- **A blackhole route** — governance takes effect the moment attachments land: the quarantined CIDR is unreachable regardless of propagations

## What to Customize

- Replace `<aws-region>`, the hub, and the attachment references with your resources
- Drop the designation flags to run this as an ordinary opt-in domain instead of the gateway-wide default
