---
title: "Committed Production Capacity"
description: "This preset buys two model units on a one-month commitment — the steady-state production shape: discounted rates in exchange for a term that bills in full and cannot be canceled or destroyed early."
type: "preset"
rank: "02"
presetSlug: "02-committed-production"
componentSlug: "bedrock-provisioned-throughput"
componentTitle: "Bedrock Provisioned Throughput"
provider: "aws"
icon: "package"
order: 2
---

# Committed Production Capacity

This preset buys two model units on a one-month commitment — the
steady-state production shape: discounted rates in exchange for a term
that bills in full and cannot be canceled or destroyed early.

## When to Use

- Workloads whose demand is proven by a no-commitment validation period
- Cost optimization once utilization is consistently high

## Key Configuration Choices

- **OneMonth before SixMonths.** The month term is the cheapest way to
  validate a commitment posture; move to SixMonths when demand is
  certain.
- **Treat this manifest as procurement.** The purchase is a financial
  commitment — route changes through the same review as any reserved
  spend.

## After Deployment

Note the term end date: the resource refuses destroy/replace until then,
so capacity changes must be planned as overlapping purchases.
