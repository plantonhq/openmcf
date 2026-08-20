# Consolidated Billing Only

This preset creates a billing-only organization: shared invoices and
volume discounts across accounts, nothing else.

## When to Use

- Estates that want one bill across accounts but deliberately no
  central governance (no SCPs, no trusted access, no delegation)
- Organizations migrating in from separate accounts, before the
  governance decision is made

## What You Get

- A CONSOLIDATED_BILLING organization — member accounts share the
  management account's bill
- Validation that keeps the advanced arms out: service access, policy
  types, delegated admins, and the resource policy all require ALL
  features (the spec rejects them here)

## Customize

- Upgrading later is IN-PLACE: switch `featureSet` to `ALL` (AWS's
  EnableAllFeatures) and the advanced arms open up
- The reverse — ALL back to billing-only — REPLACES the entire
  organization; treat the upgrade as one-way in practice
